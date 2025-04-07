package validation

import (
	"context"
	"fmt"
	"time"

	"github.com/sparkfund/kyc-service/internal/models"
)

// BiometricValidator handles biometric validation
type BiometricValidator struct {
	*BaseValidator
	aiClient AIClient
}

// NewBiometricValidator creates a new biometric validator
func NewBiometricValidator(config *Config) *BiometricValidator {
	return &BiometricValidator{
		BaseValidator: NewBaseValidator(config),
		aiClient:      NewAIClient(config.AIConfig),
	}
}

// Validate performs biometric validation
func (bv *BiometricValidator) Validate(ctx context.Context, data *models.BiometricData) (*ValidationResult, error) {
	startTime := time.Now()
	defer bv.metrics.RecordValidationDuration("biometric", time.Since(startTime))

	result := &ValidationResult{
		ValidatedAt: time.Now(),
		ModelInfo:   bv.modelInfo,
	}

	// Validate face matching
	faceScore, err := bv.validateFaceMatch(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("face matching failed: %w", err)
	}

	// Validate liveness
	livenessScore, err := bv.validateLiveness(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("liveness check failed: %w", err)
	}

	// Check for potential fraud indicators
	fraudScore, warnings := bv.checkFraudIndicators(ctx, data)
	result.Warnings = warnings

	// Calculate final score
	result.Score = bv.calculateFinalScore(faceScore, livenessScore, fraudScore)
	result.IsValid = result.Score >= bv.config.MinBiometricScore
	result.Confidence = bv.calculateConfidence(faceScore, livenessScore)

	return result, nil
}

func (bv *BiometricValidator) validateFaceMatch(ctx context.Context, data *models.BiometricData) (float64, error) {
	resp, err := bv.aiClient.CompareFaces(ctx, data.SelfieImage, data.DocumentImage)
	if err != nil {
		return 0, fmt.Errorf("face comparison failed: %w", err)
	}
	return resp.MatchScore, nil
}

func (bv *BiometricValidator) validateLiveness(ctx context.Context, data *models.BiometricData) (float64, error) {
	resp, err := bv.aiClient.CheckLiveness(ctx, data.SelfieImage)
	if err != nil {
		return 0, fmt.Errorf("liveness check failed: %w", err)
	}
	return resp.LivenessScore, nil
}

func (bv *BiometricValidator) checkFraudIndicators(ctx context.Context, data *models.BiometricData) (float64, []string) {
	var warnings []string
	var fraudScore float64 = 1.0

	// Check for digital manipulation
	if manipulationScore, err := bv.aiClient.CheckManipulation(ctx, data.SelfieImage); err == nil {
		if manipulationScore > bv.config.ManipulationThreshold {
			warnings = append(warnings, "Possible digital manipulation detected")
			fraudScore *= (1 - manipulationScore)
		}
	}

	// Check for presentation attacks
	if attackScore, err := bv.aiClient.CheckPresentationAttack(ctx, data.SelfieImage); err == nil {
		if attackScore > bv.config.PresentationAttackThreshold {
			warnings = append(warnings, "Possible presentation attack detected")
			fraudScore *= (1 - attackScore)
		}
	}

	return fraudScore, warnings
}

func (bv *BiometricValidator) calculateFinalScore(faceScore, livenessScore, fraudScore float64) float64 {
	weights := bv.config.BiometricWeights
	return (faceScore*weights.FaceMatch +
		livenessScore*weights.Liveness +
		fraudScore*weights.FraudCheck) /
		(weights.FaceMatch + weights.Liveness + weights.FraudCheck)
}

func (bv *BiometricValidator) calculateConfidence(faceScore, livenessScore float64) float64 {
	// Confidence is based on the reliability of our primary metrics
	return (faceScore + livenessScore) / 2
}
