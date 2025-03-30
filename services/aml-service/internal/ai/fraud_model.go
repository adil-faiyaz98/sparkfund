package ai

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
)

// FraudModelImpl implements the FraudModel interface using a Python-trained model
type FraudModelImpl struct {
	modelPath string
}

// NewFraudModel creates a new instance of FraudModelImpl
func NewFraudModel(modelPath string) (*FraudModelImpl, error) {
	return &FraudModelImpl{
		modelPath: modelPath,
	}, nil
}

// Predict returns a fraud prediction for the given features
func (m *FraudModelImpl) Predict(features FraudFeatures) (*FraudPrediction, error) {
	// Convert features to JSON
	featuresJSON, err := json.Marshal(features)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal features: %w", err)
	}

	// Call Python script for prediction
	cmd := exec.Command("python3", filepath.Join(m.modelPath, "fraud_detection.py"), "predict", string(featuresJSON))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run prediction: %w, output: %s", err, string(output))
	}

	// Parse prediction result
	var prediction FraudPrediction
	if err := json.Unmarshal(output, &prediction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal prediction: %w", err)
	}

	return &prediction, nil
}

// Update updates the model with new training data
func (m *FraudModelImpl) Update(features []FraudFeatures, labels []bool) error {
	// Convert training data to JSON
	trainingData := struct {
		Features []FraudFeatures `json:"features"`
		Labels   []bool          `json:"labels"`
	}{
		Features: features,
		Labels:   labels,
	}

	dataJSON, err := json.Marshal(trainingData)
	if err != nil {
		return fmt.Errorf("failed to marshal training data: %w", err)
	}

	// Call Python script for model update
	cmd := exec.Command("python3", filepath.Join(m.modelPath, "fraud_detection.py"), "update", string(dataJSON))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update model: %w, output: %s", err, string(output))
	}

	return nil
}

// GetModelInfo returns information about the current model
func (m *FraudModelImpl) GetModelInfo() map[string]interface{} {
	// Call Python script to get model info
	cmd := exec.Command("python3", filepath.Join(m.modelPath, "fraud_detection.py"), "info")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("failed to get model info: %v", err),
		}
	}

	// Parse model info
	var info map[string]interface{}
	if err := json.Unmarshal(output, &info); err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("failed to parse model info: %v", err),
		}
	}

	return info
} 