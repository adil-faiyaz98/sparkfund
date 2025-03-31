package ai

import (
	"context"
)

// ModelMetadata contains information about the model execution
type ModelMetadata struct {
	ModelVersion    string
	LatencyMs       float64
	ConfidenceScore float64
}

// MLClient defines the interface for ML model inference
type MLClient interface {
	// Predict makes a prediction using the specified model
	Predict(ctx context.Context, modelName string, features map[string]interface{}) (map[string]interface{}, ModelMetadata, error)

	// GetExplanations gets explanations for a prediction (for regulatory compliance)
	GetExplanations(ctx context.Context, modelName string, features map[string]interface{}) (map[string]interface{}, error)
}

// RemoteMLClient implements MLClient for a remote model server
type RemoteMLClient struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
	logger     *logrus.Logger
}

// NewRemoteMLClient creates a new client for the remote ML service
func NewRemoteMLClient(baseURL string, apiKey string, logger *logrus.Logger) *RemoteMLClient {
	return &RemoteMLClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey: apiKey,
		logger: logger,
	}
}

// Predict makes a prediction using the specified model
func (c *RemoteMLClient) Predict(ctx context.Context, modelName string, features map[string]interface{}) (map[string]interface{}, ModelMetadata, error) {
	startTime := time.Now()

	requestBody := map[string]interface{}{
		"model_name": modelName,
		"features":   features,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, ModelMetadata{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/predict", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, ModelMetadata{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, ModelMetadata{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ModelMetadata{}, fmt.Errorf("model server returned status: %s", resp.Status)
	}

	var response struct {
		Predictions map[string]interface{} `json:"predictions"`
		Metadata    struct {
			ModelVersion    string  `json:"model_version"`
			ConfidenceScore float64 `json:"confidence_score"`
		} `json:"metadata"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, ModelMetadata{}, err
	}

	metadata := ModelMetadata{
		ModelVersion:    response.Metadata.ModelVersion,
		LatencyMs:       float64(time.Since(startTime).Milliseconds()),
		ConfidenceScore: response.Metadata.ConfidenceScore,
	}

	return response.Predictions, metadata, nil
}
