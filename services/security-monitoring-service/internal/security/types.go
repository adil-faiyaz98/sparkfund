package security

import (
	"time"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID            string
	UserID        string
	Amount        float64
	Currency      string
	Timestamp     time.Time
	Status        string
	RiskLevel     string
	RiskScore     float64
	Location      Location
	Device        DeviceInfo
	RecipientID   string
	RecipientName string
	Description   string
	Category      string
	Flags         []string
}

// Location represents transaction location information
type Location struct {
	Country    string
	City       string
	IP         string
	Latitude   float64
	Longitude  float64
	IsVPN      bool
	IsProxy    bool
	Confidence float64
}

// DeviceInfo represents device information
type DeviceInfo struct {
	DeviceID      string
	DeviceType    string
	OS            string
	Browser       string
	UserAgent     string
	IsKnownDevice bool
	FirstSeen     time.Time
	LastSeen      time.Time
}

// Session represents a user session
type Session struct {
	ID        string
	UserID    string
	Timestamp time.Time
	Location  Location
	Device    DeviceInfo
	RiskScore float64
	Token     string
	ExpiresAt time.Time
}

// RiskAssessment represents the result of a risk assessment
type RiskAssessment struct {
	RiskScore       float64
	RiskLevel       string
	Factors         []RiskFactor
	Recommendations []string
	Timestamp       time.Time
}

// RiskFactor represents a specific risk factor
type RiskFactor struct {
	Name        string
	Score       float64
	Weight      float64
	Description string
	Details     map[string]interface{}
}

// RiskProfile represents a user's risk profile
type RiskProfile struct {
	UserID             string
	BaseRiskScore      float64
	RiskFactors        []RiskFactor
	LastAssessment     time.Time
	TransactionHistory []TransactionSummary
	BehaviorPatterns   map[string]interface{}
}

// TransactionSummary represents a summary of a transaction for risk assessment
type TransactionSummary struct {
	ID        string
	Amount    float64
	Timestamp time.Time
	RiskScore float64
	Status    string
}

// MFAFactor represents a multi-factor authentication factor
type MFAFactor struct {
	ID          string
	Type        string
	Secret      string
	Status      string
	CreatedAt   time.Time
	LastUsed    time.Time
	RetryCount  int
	LockedUntil time.Time
	Metadata    map[string]interface{}
}

// MFAStatus represents the status of MFA verification
type MFAStatus struct {
	Verified         bool
	FactorsUsed      []string
	RemainingFactors int
	Locked           bool
	LockedUntil      time.Time
	Error            string
}

// Permission represents a system permission
type Permission struct {
	Name        string
	Description string
	Category    string
	RiskLevel   string
}

// AccessDecision represents the result of an access control decision
type AccessDecision struct {
	Allowed     bool
	Roles       []string
	Permissions []string
	Reason      string
	RiskLevel   string
	ExpiresAt   time.Time
}

// DynamicRule represents a dynamic access control rule
type DynamicRule struct {
	Name        string
	Condition   string
	Action      string
	Priority    int
	Description string
}
