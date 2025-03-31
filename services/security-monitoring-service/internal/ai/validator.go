package ai

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sparkfund/security-monitoring/internal/models"
)

// ModelValidator handles model validation and performance metrics
type ModelValidator struct {
	mu sync.RWMutex
}

// NewModelValidator creates a new model validator
func NewModelValidator() *ModelValidator {
	return &ModelValidator{}
}

// ValidateModel validates a model's performance
func (v *ModelValidator) ValidateModel(ctx context.Context, model Model, testData []models.SecurityData) (ValidationResult, error) {
	if len(testData) == 0 {
		return ValidationResult{}, fmt.Errorf("no test data provided")
	}

	var truePositives, falsePositives, falseNegatives int
	var totalPredictions int

	for _, data := range testData {
		prediction, err := model.Predict(ctx, data)
		if err != nil {
			return ValidationResult{}, fmt.Errorf("prediction failed: %v", err)
		}

		// Evaluate prediction against ground truth
		result := v.evaluatePrediction(prediction, data)
		truePositives += result.truePositives
		falsePositives += result.falsePositives
		falseNegatives += result.falseNegatives
		totalPredictions++
	}

	// Calculate metrics
	accuracy := float64(truePositives) / float64(totalPredictions)
	precision := v.calculatePrecision(truePositives, falsePositives)
	recall := v.calculateRecall(truePositives, falseNegatives)
	f1Score := v.calculateF1Score(precision, recall)

	return ValidationResult{
		Accuracy:    accuracy,
		Precision:   precision,
		Recall:      recall,
		F1Score:     f1Score,
		LastUpdated: time.Now(),
	}, nil
}

// evaluatePrediction evaluates a single prediction
func (v *ModelValidator) evaluatePrediction(prediction interface{}, data models.SecurityData) predictionResult {
	var result predictionResult

	// Extract ground truth from data
	groundTruth := v.extractGroundTruth(data)

	// Compare prediction with ground truth
	switch pred := prediction.(type) {
	case *models.Threat:
		result = v.evaluateThreatPrediction(pred, groundTruth)
	case *models.Intrusion:
		result = v.evaluateIntrusionPrediction(pred, groundTruth)
	case *models.Malware:
		result = v.evaluateMalwarePrediction(pred, groundTruth)
	case *models.Pattern:
		result = v.evaluatePatternPrediction(pred, groundTruth)
	}

	return result
}

// predictionResult stores the evaluation results for a single prediction
type predictionResult struct {
	truePositives  int
	falsePositives int
	falseNegatives int
}

// extractGroundTruth extracts ground truth from security data
func (v *ModelValidator) extractGroundTruth(data models.SecurityData) map[string]interface{} {
	groundTruth := make(map[string]interface{})

	// Extract known threats
	if threats, ok := data.Metadata["known_threats"].([]interface{}); ok {
		groundTruth["threats"] = threats
	}

	// Extract known intrusions
	if intrusions, ok := data.Metadata["known_intrusions"].([]interface{}); ok {
		groundTruth["intrusions"] = intrusions
	}

	// Extract known malware
	if malware, ok := data.Metadata["known_malware"].([]interface{}); ok {
		groundTruth["malware"] = malware
	}

	// Extract known patterns
	if patterns, ok := data.Metadata["known_patterns"].([]interface{}); ok {
		groundTruth["patterns"] = patterns
	}

	return groundTruth
}

// Model-specific evaluation functions
func (v *ModelValidator) evaluateThreatPrediction(pred *models.Threat, groundTruth map[string]interface{}) predictionResult {
	var result predictionResult
	knownThreats, ok := groundTruth["threats"].([]interface{})
	if !ok {
		return result
	}

	// Check if prediction matches any known threat
	matched := false
	for _, threat := range knownThreats {
		if threatMap, ok := threat.(map[string]interface{}); ok {
			if v.matchThreat(pred, threatMap) {
				result.truePositives++
				matched = true
				break
			}
		}
	}

	if !matched {
		if pred.Severity >= models.SeverityHigh {
			result.falsePositives++
		} else {
			result.falseNegatives++
		}
	}

	return result
}

func (v *ModelValidator) evaluateIntrusionPrediction(pred *models.Intrusion, groundTruth map[string]interface{}) predictionResult {
	var result predictionResult
	knownIntrusions, ok := groundTruth["intrusions"].([]interface{})
	if !ok {
		return result
	}

	// Check if prediction matches any known intrusion
	matched := false
	for _, intrusion := range knownIntrusions {
		if intrusionMap, ok := intrusion.(map[string]interface{}); ok {
			if v.matchIntrusion(pred, intrusionMap) {
				result.truePositives++
				matched = true
				break
			}
		}
	}

	if !matched {
		if pred.Severity >= models.SeverityHigh {
			result.falsePositives++
		} else {
			result.falseNegatives++
		}
	}

	return result
}

func (v *ModelValidator) evaluateMalwarePrediction(pred *models.Malware, groundTruth map[string]interface{}) predictionResult {
	var result predictionResult
	knownMalware, ok := groundTruth["malware"].([]interface{})
	if !ok {
		return result
	}

	// Check if prediction matches any known malware
	matched := false
	for _, malware := range knownMalware {
		if malwareMap, ok := malware.(map[string]interface{}); ok {
			if v.matchMalware(pred, malwareMap) {
				result.truePositives++
				matched = true
				break
			}
		}
	}

	if !matched {
		if pred.Severity >= models.SeverityHigh {
			result.falsePositives++
		} else {
			result.falseNegatives++
		}
	}

	return result
}

func (v *ModelValidator) evaluatePatternPrediction(pred *models.Pattern, groundTruth map[string]interface{}) predictionResult {
	var result predictionResult
	knownPatterns, ok := groundTruth["patterns"].([]interface{})
	if !ok {
		return result
	}

	// Check if prediction matches any known pattern
	matched := false
	for _, pattern := range knownPatterns {
		if patternMap, ok := pattern.(map[string]interface{}); ok {
			if v.matchPattern(pred, patternMap) {
				result.truePositives++
				matched = true
				break
			}
		}
	}

	if !matched {
		if pred.Severity >= models.SeverityHigh {
			result.falsePositives++
		} else {
			result.falseNegatives++
		}
	}

	return result
}

// Matching functions
func (v *ModelValidator) matchThreat(pred *models.Threat, known map[string]interface{}) bool {
	// Implement threat matching logic
	// This could include:
	// - Signature matching
	// - Behavior matching
	// - Context matching
	return false
}

func (v *ModelValidator) matchIntrusion(pred *models.Intrusion, known map[string]interface{}) bool {
	// Implement intrusion matching logic
	return false
}

func (v *ModelValidator) matchMalware(pred *models.Malware, known map[string]interface{}) bool {
	// Implement malware matching logic
	return false
}

func (v *ModelValidator) matchPattern(pred *models.Pattern, known map[string]interface{}) bool {
	// Implement pattern matching logic
	return false
}

// Metric calculation functions
func (v *ModelValidator) calculatePrecision(truePositives, falsePositives int) float64 {
	if truePositives+falsePositives == 0 {
		return 0
	}
	return float64(truePositives) / float64(truePositives+falsePositives)
}

func (v *ModelValidator) calculateRecall(truePositives, falseNegatives int) float64 {
	if truePositives+falseNegatives == 0 {
		return 0
	}
	return float64(truePositives) / float64(truePositives+falseNegatives)
}

func (v *ModelValidator) calculateF1Score(precision, recall float64) float64 {
	if precision+recall == 0 {
		return 0
	}
	return 2 * (precision * recall) / (precision + recall)
}
