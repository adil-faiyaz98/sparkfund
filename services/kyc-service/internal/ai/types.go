package ai

import (
	"time"
)

// FraudRiskLevel represents the risk level assessment from the AI model
type FraudRiskLevel string

const (
	FraudRiskLow    FraudRiskLevel = "low"
	FraudRiskMedium FraudRiskLevel = "medium"
	FraudRiskHigh   FraudRiskLevel = "high"
)

// FraudScore represents the numerical risk score from the AI model
type FraudScore float64

// FraudFeatures represents the input features for the fraud detection model
type FraudFeatures struct {
	// Transaction Features
	TransactionAmount     float64   `json:"transaction_amount"`
	TransactionTime      time.Time `json:"transaction_time"`
	TransactionFrequency int       `json:"transaction_frequency"`
	
	// User Features
	UserAge              int       `json:"user_age"`
	AccountAge           int       `json:"account_age"` // Days since account creation
	PreviousFraudReports int       `json:"previous_fraud_reports"`
	
	// Document Features
	DocumentQuality      float64   `json:"document_quality"` // 0-1 score
	DocumentAge          int       `json:"document_age"`     // Days since document issue
	DocumentType         string    `json:"document_type"`
	
	// Location Features
	CountryRiskScore     float64   `json:"country_risk_score"` // 0-1 score
	IPRiskScore          float64   `json:"ip_risk_score"`      // 0-1 score
	
	// Behavioral Features
	LoginAttempts        int       `json:"login_attempts"`
	FailedLogins         int       `json:"failed_logins"`
	LastLoginTime        time.Time `json:"last_login_time"`
	
	// Additional Risk Factors
	PEPStatus            bool      `json:"pep_status"`            // Politically Exposed Person
	SanctionStatus       bool      `json:"sanction_status"`       // Sanctioned individual
	WatchlistStatus      bool      `json:"watchlist_status"`      // On watchlist
}

// FraudPrediction represents the output from the fraud detection model
type FraudPrediction struct {
	RiskLevel     FraudRiskLevel `json:"risk_level"`
	RiskScore     FraudScore     `json:"risk_score"`
	Confidence    float64        `json:"confidence"`     // Model's confidence in the prediction
	Explanation   string         `json:"explanation"`    // Human-readable explanation of the risk factors
	RiskFactors   []string       `json:"risk_factors"`   // List of identified risk factors
	RecommendedAction string     `json:"recommended_action"` // Suggested action based on risk level
}

// FraudModel represents the interface for fraud detection models
type FraudModel interface {
	// Predict returns a fraud prediction for the given features
	Predict(features FraudFeatures) (*FraudPrediction, error)
	
	// Update updates the model with new training data
	Update(trainingData []FraudFeatures, labels []bool) error
	
	// GetModelInfo returns information about the model
	GetModelInfo() map[string]interface{}
} 