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

// Model loading functions
func loadKYCAnomalyModel(path string) (KYCAnomalyModel, error) {
	// Implement KYC anomaly model loading logic
	return nil, nil
}

func loadKYCBehaviorModel(path string) (KYCBehaviorModel, error) {
	// Implement KYC behavior model loading logic
	return nil, nil
}

func loadKYCThreatModel(path string) (KYCThreatModel, error) {
	// Implement KYC threat model loading logic
	return nil, nil
}
