package security

import (
	"time"
)

// Location represents geographic location information
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

// RiskLevel represents a risk level with associated score
type RiskLevel struct {
	Level  string
	Score  float64
	Reason string
}

// Helper function to convert risk level string to score
func riskLevelToScore(level string) float64 {
	switch level {
	case "high":
		return 0.8
	case "medium":
		return 0.5
	case "low":
		return 0.2
	default:
		return 0.0
	}
}
