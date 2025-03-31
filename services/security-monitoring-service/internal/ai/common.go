package ai

import (
	"errors"
	"time"
)

// Common errors
var (
	ErrInsufficientData = errors.New("insufficient data for training")
	ErrModelNotFound    = errors.New("model not found")
	ErrInvalidConfig    = errors.New("invalid configuration")
	ErrTrainingFailed   = errors.New("model training failed")
	ErrPredictionFailed = errors.New("model prediction failed")
)

// Common types
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

// Model loading functions
func loadModel(path string) (AnomalyModel, error) {
	// Implement model loading logic
	return nil, nil
}

func loadBehaviorModel(path string) (BehaviorModel, error) {
	// Implement behavior model loading logic
	return nil, nil
}

func loadThreatModel(path string) (ThreatModel, error) {
	// Implement threat model loading logic
	return nil, nil
}
