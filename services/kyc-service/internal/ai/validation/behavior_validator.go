package validation

import (
    "context"
    "fmt"
    "time"

    "github.com/sparkfund/kyc-service/internal/models"
)

// BehaviorValidator handles behavioral pattern validation
type BehaviorValidator struct {
    *BaseValidator
    aiClient AIClient
}

// NewBehaviorValidator creates a new behavior validator
func NewBehaviorValidator(config *Config) *BehaviorValidator {
    return &BehaviorValidator{
        BaseValidator: NewBaseValidator(config),
        aiClient:      NewAIClient(config.AIConfig),
    }
}

// Validate performs behavioral validation
func (bv *BehaviorValidator) Validate(ctx context.Context, data *models.BehaviorData) (*ValidationResult, error) {
    startTime := time.Now()
    defer bv.metrics.RecordValidationDuration("behavior", time.Since(startTime))

    result := &ValidationResult{
        ValidatedAt: time.Now(),
        ModelInfo:   bv.modelInfo,
    }

    // Analyze user patterns
    patternScore, patternWarnings, err := bv.analyzePatterns(ctx, data)
    if err != nil {
        return nil, fmt.Errorf("pattern analysis failed: %w", err)
    }
    result.Warnings = append(result.Warnings, patternWarnings...)

    // Analyze anomalies
    anomalyScore, anomalyWarnings, err := bv.detectAnomalies(ctx, data)
    if err != nil {
        return nil, fmt.Errorf("anomaly detection failed: %w", err)
    }
    result.Warnings = append(result.Warnings, anomalyWarnings...)

    // Analyze risk patterns
    riskScore, riskWarnings, err := bv.assessRiskPatterns(ctx, data)
    if err != nil {
        return nil, fmt.Errorf("risk pattern assessment failed: %w", err)
    }
    result.Warnings = append(result.Warnings, riskWarnings...)

    // Calculate final score
    result.Score = bv.calculateFinalScore(patternScore, anomalyScore, riskScore)
    result.IsValid = result.Score >= bv.config.MinBehaviorScore
    result.Confidence = bv.calculateConfidence(patternScore, anomalyScore, riskScore)

    return result, nil
}

func (bv *BehaviorValidator) analyzePatterns(ctx context.Context, data *models.BehaviorData) (float64, []string, error) {
    resp, err := bv.aiClient.AnalyzeBehaviorPatterns(ctx, data)
    if err != nil {
        return 0, nil, err
    }

    var warnings []string
    if resp.UnusualPatternScore > bv.config.UnusualPatternThreshold {
        warnings = append(warnings, "Unusual behavior patterns detected")
    }
    if resp.VelocityScore > bv.config.VelocityThreshold {
        warnings = append(warnings, "High velocity behavior detected")
    }

    return resp.PatternScore, warnings, nil
}

func (bv *BehaviorValidator) detectAnomalies(ctx context.Context, data *models.BehaviorData) (float64, []string, error) {
    resp, err := bv.aiClient.DetectBehaviorAnomalies(ctx, data)
    if err != nil {
        return 0, nil, err
    }

    var warnings []string
    for _, anomaly := range resp.Anomalies {
        if anomaly.Score > bv.config.AnomalyThreshold {
            warnings = append(warnings, fmt.Sprintf("Anomaly detected: %s", anomaly.Description))
        }
    }

    return resp.AnomalyScore, warnings, nil
}

func (bv *BehaviorValidator) assessRiskPatterns(ctx context.Context, data *models.BehaviorData) (float64, []string, error) {
    resp, err := bv.aiClient.AssessRiskPatterns(ctx, data)
    if err != nil {
        return 0, nil, err
    }

    var warnings []string
    if resp.FraudRiskScore > bv.config.FraudPatternThreshold {
        warnings = append(warnings, "Suspicious behavior patterns detected")
    }
    if resp.ComplianceRiskScore > bv.config.CompliancePatternThreshold {
        warnings = append(warnings, "Compliance risk patterns detected")
    }

    return resp.RiskScore, warnings, nil
}

func (bv *BehaviorValidator) calculateFinalScore(patternScore, anomalyScore, riskScore float64) float64 {
    weights := bv.config.BehaviorWeights
    return (patternScore*weights.Pattern + 
            anomalyScore*weights.Anomaly + 
            riskScore*weights.Risk) / 
            (weights.Pattern + weights.Anomaly + weights.Risk)
}

func (bv *BehaviorValidator) calculateConfidence(patternScore, anomalyScore, riskScore float64) float64 {
    // Confidence is based on the consistency of our behavioral analysis
    return (patternScore + anomalyScore + riskScore) / 3
}