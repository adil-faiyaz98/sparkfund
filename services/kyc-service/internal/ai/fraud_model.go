package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type FraudRiskLevel string

const (
	FraudRiskLow    FraudRiskLevel = "LOW"
	FraudRiskMedium FraudRiskLevel = "MEDIUM"
	FraudRiskHigh   FraudRiskLevel = "HIGH"
)

type FraudFeatures struct {
	TransactionAmount float64   `json:"transaction_amount"`
	TransactionTime   time.Time `json:"transaction_time"`
	UserAge           int       `json:"user_age"`
	DocumentType      string    `json:"document_type"`
	CountryRiskScore  float64   `json:"country_risk_score"`
}

type FraudPrediction struct {
	RiskLevel   FraudRiskLevel `json:"risk_level"`
	RiskScore   float64        `json:"risk_score"`
	Explanation string         `json:"explanation"`
}

type FraudModel interface {
	Predict(features FraudFeatures) (*FraudPrediction, error)
}

type MLServiceFraudModel struct {
	baseURL    string
	httpClient *http.Client
}

func NewMLServiceFraudModel(baseURL string, timeout time.Duration) *MLServiceFraudModel {
	return &MLServiceFraudModel{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (m *MLServiceFraudModel) Predict(features FraudFeatures) (*FraudPrediction, error) {
	reqBody, err := json.Marshal(features)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal features: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, m.baseURL+"/predict", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service returned status %d", resp.StatusCode)
	}

	var prediction FraudPrediction
	if err := json.NewDecoder(resp.Body).Decode(&prediction); err != nil {
		return nil, fmt.Errorf("failed to decode prediction: %w", err)
	}

	return &prediction, nil
}
