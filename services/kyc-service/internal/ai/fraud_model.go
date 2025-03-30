package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// PythonFraudModel implements the FraudModel interface using a Python-trained model
type PythonFraudModel struct {
	modelPath    string
	pythonScript string
}

// NewPythonFraudModel creates a new instance of PythonFraudModel
func NewPythonFraudModel(modelPath string) (*PythonFraudModel, error) {
	// Ensure the model file exists
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("model file not found: %s", modelPath)
	}

	// Get the directory of the current Go file
	scriptDir := filepath.Dir(os.Args[0])
	pythonScript := filepath.Join(scriptDir, "ai", "fraud_detection.py")

	return &PythonFraudModel{
		modelPath:    modelPath,
		pythonScript: pythonScript,
	}, nil
}

// Predict implements the FraudModel interface
func (m *PythonFraudModel) Predict(features FraudFeatures) (*FraudPrediction, error) {
	// Convert features to JSON
	featuresJSON, err := json.Marshal(features)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal features: %w", err)
	}

	// Prepare the Python command
	cmd := exec.Command("python", m.pythonScript, "predict", m.modelPath, string(featuresJSON))

	// Execute the command and get output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute Python script: %w\nOutput: %s", err, output)
	}

	// Parse the prediction
	var prediction FraudPrediction
	if err := json.Unmarshal(output, &prediction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal prediction: %w", err)
	}

	return &prediction, nil
}

// Update implements the FraudModel interface
func (m *PythonFraudModel) Update(trainingData []FraudFeatures, labels []bool) error {
	// Convert training data to JSON
	data := struct {
		Features []FraudFeatures `json:"features"`
		Labels   []bool          `json:"labels"`
	}{
		Features: trainingData,
		Labels:   labels,
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal training data: %w", err)
	}

	// Prepare the Python command
	cmd := exec.Command("python", m.pythonScript, "update", m.modelPath, string(dataJSON))

	// Execute the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute Python script: %w\nOutput: %s", err, output)
	}

	return nil
}

// GetModelInfo implements the FraudModel interface
func (m *PythonFraudModel) GetModelInfo() map[string]interface{} {
	return map[string]interface{}{
		"type":        "python",
		"model_path":  m.modelPath,
		"script_path": m.pythonScript,
		"last_update": time.Now(),
	}
}
