package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"sparkfund/services/kyc-service/internal/model"
)

// AIClient is a client for the AI service
type AIClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logger     *logrus.Logger
}

// NewAIClient creates a new AI client
func NewAIClient(logger *logrus.Logger) *AIClient {
	baseURL := os.Getenv("AI_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://ai-service:8000" // Default URL for Docker Compose
	}

	apiKey := os.Getenv("AI_SERVICE_API_KEY")
	if apiKey == "" {
		apiKey = "your-api-key" // Default API key
	}

	return &AIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// AnalyzeDocument analyzes a document using the AI service
func (c *AIClient) AnalyzeDocument(ctx context.Context, documentID uuid.UUID, verificationID uuid.UUID, documentPath string) (*model.DocumentAnalysisResult, error) {
	c.logger.WithFields(logrus.Fields{
		"document_id":     documentID,
		"verification_id": verificationID,
	}).Info("Analyzing document")

	// Create request body
	reqBody := map[string]interface{}{
		"document_id":     documentID.String(),
		"verification_id": verificationID.String(),
	}

	// Convert request body to JSON
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/v1/document/analyze", c.baseURL),
		bytes.NewBuffer(reqBodyJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI service returned non-OK status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var result model.DocumentAnalysisResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// MatchFaces matches a selfie with a document photo using the AI service
func (c *AIClient) MatchFaces(ctx context.Context, documentID uuid.UUID, selfieID uuid.UUID, verificationID uuid.UUID) (*model.FaceMatchResult, error) {
	c.logger.WithFields(logrus.Fields{
		"document_id":     documentID,
		"selfie_id":       selfieID,
		"verification_id": verificationID,
	}).Info("Matching faces")

	// Create request body
	reqBody := map[string]interface{}{
		"document_id":     documentID.String(),
		"selfie_id":       selfieID.String(),
		"verification_id": verificationID.String(),
	}

	// Convert request body to JSON
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/v1/face/match", c.baseURL),
		bytes.NewBuffer(reqBodyJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI service returned non-OK status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var result model.FaceMatchResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// AnalyzeRisk analyzes risk based on user data and device information using the AI service
func (c *AIClient) AnalyzeRisk(ctx context.Context, userID uuid.UUID, verificationID uuid.UUID, deviceInfo model.DeviceInfo) (*model.RiskAnalysisResult, error) {
	c.logger.WithFields(logrus.Fields{
		"user_id":         userID,
		"verification_id": verificationID,
	}).Info("Analyzing risk")

	// Create request body
	reqBody := map[string]interface{}{
		"user_id":         userID.String(),
		"verification_id": verificationID.String(),
		"device_info":     deviceInfo,
	}

	// Convert request body to JSON
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/v1/risk/analyze", c.baseURL),
		bytes.NewBuffer(reqBodyJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI service returned non-OK status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var result model.RiskAnalysisResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DetectAnomalies detects anomalies in user behavior using the AI service
func (c *AIClient) DetectAnomalies(ctx context.Context, userID uuid.UUID, verificationID uuid.UUID, deviceInfo model.DeviceInfo) (*model.AnomalyDetectionResult, error) {
	c.logger.WithFields(logrus.Fields{
		"user_id":         userID,
		"verification_id": verificationID,
	}).Info("Detecting anomalies")

	// Create request body
	reqBody := map[string]interface{}{
		"user_id":         userID.String(),
		"verification_id": verificationID.String(),
		"device_info":     deviceInfo,
	}

	// Convert request body to JSON
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/v1/anomaly/detect", c.baseURL),
		bytes.NewBuffer(reqBodyJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI service returned non-OK status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var result model.AnomalyDetectionResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ListAIModels lists all AI models using the AI service
func (c *AIClient) ListAIModels(ctx context.Context) ([]*model.AIModelInfo, error) {
	c.logger.Info("Listing AI models")

	// Create request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/api/v1/models", c.baseURL),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-API-Key", c.apiKey)

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI service returned non-OK status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var result struct {
		Models []*model.AIModelInfo `json:"models"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result.Models, nil
}
