package ai

import "time"

// FraudRiskLevel represents the level of fraud risk
type FraudRiskLevel string

const (
	FraudRiskLow    FraudRiskLevel = "LOW"
	FraudRiskMedium FraudRiskLevel = "MEDIUM"
	FraudRiskHigh   FraudRiskLevel = "HIGH"
)

// FraudScore represents a numerical risk score
type FraudScore float64

// FraudFeatures represents the input features for fraud detection
type FraudFeatures struct {
	// Transaction Features
	Amount            float64   `json:"amount"`
	Currency          string    `json:"currency"`
	TransactionTime   time.Time `json:"transaction_time"`
	TransactionType   string    `json:"transaction_type"`
	TransactionCount  int       `json:"transaction_count"`
	AverageAmount     float64   `json:"average_amount"`
	AmountDeviation   float64   `json:"amount_deviation"`
	
	// User Features
	UserAge           int       `json:"user_age"`
	AccountAge        int       `json:"account_age"`
	PreviousFraud     bool      `json:"previous_fraud"`
	RiskProfile       float64   `json:"risk_profile"`
	
	// Location Features
	CountryRisk       float64   `json:"country_risk"`
	IPRisk            float64   `json:"ip_risk"`
	LocationMismatch  bool      `json:"location_mismatch"`
	
	// Behavioral Features
	LoginAttempts     int       `json:"login_attempts"`
	FailedLogins      int       `json:"failed_logins"`
	DeviceChanges     int       `json:"device_changes"`
	
	// Additional Risk Factors
	PEPStatus         bool      `json:"pep_status"`
	SanctionList      bool      `json:"sanction_list"`
	WatchList         bool      `json:"watch_list"`
}

// FraudPrediction represents the output of fraud detection
type FraudPrediction struct {
	RiskLevel    FraudRiskLevel `json:"risk_level"`
	RiskScore    FraudScore     `json:"risk_score"`
	Explanation  string         `json:"explanation"`
	RiskFactors  []string       `json:"risk_factors"`
	Confidence   float64        `json:"confidence"`
}

// FraudModel interface defines the contract for fraud detection models
type FraudModel interface {
	// Predict returns a fraud prediction for the given features
	Predict(features FraudFeatures) (*FraudPrediction, error)
	
	// Update updates the model with new training data
	Update(features []FraudFeatures, labels []bool) error
	
	// GetModelInfo returns information about the current model
	GetModelInfo() map[string]interface{}
} 