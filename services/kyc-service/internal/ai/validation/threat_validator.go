package validation

import (
	"context"
	"fmt"
	"time"

	"github.com/sparkfund/kyc-service/internal/models"
)

// ThreatValidator handles threat and risk validation
type ThreatValidator struct {
	*BaseValidator
	aiClient AIClient
}

// NewThreatValidator creates a new threat validator
func NewThreatValidator(config *Config) *ThreatValidator {
	return &ThreatValidator{
		BaseValidator: NewBaseValidator(config),
		aiClient:      NewAIClient(config.AIConfig),
	}
}

// Validate performs threat validation
func (tv *ThreatValidator) Validate(ctx context.Context, input *models.KYCData) (*ValidationResult, error) {
	startTime := time.Now()
	defer tv.metrics.RecordValidationDuration("threat", time.Since(startTime))

	result := &ValidationResult{
		ValidatedAt: time.Now(),
		ModelInfo:   tv.modelInfo,
	}

	// Check for identity threats
	identityScore, identityWarnings, err := tv.checkIdentityThreats(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("identity threat check failed: %w", err)
	}
	result.Warnings = append(result.Warnings, identityWarnings...)

	// Check for financial risks
	financialScore, financialWarnings, err := tv.checkFinancialRisks(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("financial risk check failed: %w", err)
	}
	result.Warnings = append(result.Warnings, financialWarnings...)

	// Check for compliance risks
	complianceScore, complianceWarnings, err := tv.checkComplianceRisks(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("compliance risk check failed: %w", err)
	}
	result.Warnings = append(result.Warnings, complianceWarnings...)

	// Calculate final score
	result.Score = tv.calculateFinalScore(identityScore, financialScore, complianceScore)
	result.IsValid = result.Score >= tv.config.MinThreatScore
	result.Confidence = tv.calculateConfidence(identityScore, financialScore, complianceScore)

	return result, nil
}

func (tv *ThreatValidator) checkIdentityThreats(ctx context.Context, input *models.KYCData) (float64, []string, error) {
	resp, err := tv.aiClient.AnalyzeIdentityThreats(ctx, input)
	if err != nil {
		return 0, nil, err
	}

	var warnings []string
	if resp.SyntheticIdentityRisk > tv.config.SyntheticIdentityThreshold {
		warnings = append(warnings, "High synthetic identity risk detected")
	}
	if resp.IdentityTheftRisk > tv.config.IdentityTheftThreshold {
		warnings = append(warnings, "Elevated identity theft risk detected")
	}

	return resp.ThreatScore, warnings, nil
}

func (tv *ThreatValidator) checkFinancialRisks(ctx context.Context, input *models.KYCData) (float64, []string, error) {
	resp, err := tv.aiClient.AnalyzeFinancialRisks(ctx, input)
	if err != nil {
		return 0, nil, err
	}

	var warnings []string
	if resp.FraudRisk > tv.config.FraudRiskThreshold {
		warnings = append(warnings, "High fraud risk detected")
	}
	if resp.MoneyLaunderingRisk > tv.config.AMLRiskThreshold {
		warnings = append(warnings, "Elevated AML risk detected")
	}

	return resp.RiskScore, warnings, nil
}

func (tv *ThreatValidator) checkComplianceRisks(ctx context.Context, input *models.KYCData) (float64, []string, error) {
	resp, err := tv.aiClient.AnalyzeComplianceRisks(ctx, input)
	if err != nil {
		return 0, nil, err
	}

	var warnings []string
	if resp.SanctionsRisk > tv.config.SanctionsThreshold {
		warnings = append(warnings, "Potential sanctions list match detected")
	}
	if resp.PEPRisk > tv.config.PEPThreshold {
		warnings = append(warnings, "Potential PEP status detected")
	}

	return resp.ComplianceScore, warnings, nil
}

func (tv *ThreatValidator) calculateFinalScore(identityScore, financialScore, complianceScore float64) float64 {
	weights := tv.config.ThreatWeights
	return (identityScore*weights.Identity +
		financialScore*weights.Financial +
		complianceScore*weights.Compliance) /
		(weights.Identity + weights.Financial + weights.Compliance)
}

func (tv *ThreatValidator) calculateConfidence(identityScore, financialScore, complianceScore float64) float64 {
	// Confidence is based on the reliability of our risk assessments
	return (identityScore + financialScore + complianceScore) / 3
}
