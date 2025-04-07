package model

import (
	"time"

	"github.com/google/uuid"
)

// DocumentAnalysisResult represents the result of AI document analysis
type DocumentAnalysisResult struct {
	ID             uuid.UUID         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	VerificationID uuid.UUID         `json:"verification_id" gorm:"type:uuid;not null"`
	DocumentID     uuid.UUID         `json:"document_id" gorm:"type:uuid;not null"`
	DocumentType   string            `json:"document_type" gorm:"type:varchar(50);not null"`
	IsAuthentic    bool              `json:"is_authentic" gorm:"not null"`
	Confidence     float64           `json:"confidence" gorm:"type:float;not null"`
	ExtractedData  map[string]string `json:"extracted_data" gorm:"type:jsonb"`
	Issues         []string          `json:"issues" gorm:"type:jsonb"`
	CreatedAt      time.Time         `json:"created_at" gorm:"not null;default:now()"`
}

// FaceMatchResult represents the result of face matching analysis
type FaceMatchResult struct {
	ID             uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	VerificationID uuid.UUID `json:"verification_id" gorm:"type:uuid;not null"`
	DocumentID     uuid.UUID `json:"document_id" gorm:"type:uuid;not null"`
	SelfieID       uuid.UUID `json:"selfie_id" gorm:"type:uuid;not null"`
	IsMatch        bool      `json:"is_match" gorm:"not null"`
	Confidence     float64   `json:"confidence" gorm:"type:float;not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"not null;default:now()"`
}

// RiskAnalysisResult represents the result of risk analysis
type RiskAnalysisResult struct {
	ID             uuid.UUID         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	VerificationID uuid.UUID         `json:"verification_id" gorm:"type:uuid;not null"`
	UserID         uuid.UUID         `json:"user_id" gorm:"type:uuid;not null"`
	RiskScore      float64           `json:"risk_score" gorm:"type:float;not null"`
	RiskLevel      string            `json:"risk_level" gorm:"type:varchar(20);not null"`
	RiskFactors    []string          `json:"risk_factors" gorm:"type:jsonb"`
	DeviceInfo     map[string]string `json:"device_info" gorm:"type:jsonb"`
	IPAddress      string            `json:"ip_address" gorm:"type:varchar(45)"`
	Location       string            `json:"location" gorm:"type:varchar(100)"`
	CreatedAt      time.Time         `json:"created_at" gorm:"not null;default:now()"`
}

// AnomalyDetectionResult represents the result of anomaly detection
type AnomalyDetectionResult struct {
	ID             uuid.UUID         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	VerificationID uuid.UUID         `json:"verification_id" gorm:"type:uuid;not null"`
	UserID         uuid.UUID         `json:"user_id" gorm:"type:uuid;not null"`
	IsAnomaly      bool              `json:"is_anomaly" gorm:"not null"`
	AnomalyScore   float64           `json:"anomaly_score" gorm:"type:float;not null"`
	AnomalyType    string            `json:"anomaly_type" gorm:"type:varchar(50)"`
	Reasons        []string          `json:"reasons" gorm:"type:jsonb"`
	DeviceInfo     map[string]string `json:"device_info" gorm:"type:jsonb"`
	CreatedAt      time.Time         `json:"created_at" gorm:"not null;default:now()"`
}

// DeviceInfo represents information about a user's device
type DeviceInfo struct {
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	DeviceType   string    `json:"device_type"`
	OS           string    `json:"os"`
	Browser      string    `json:"browser"`
	MacAddress   string    `json:"mac_address,omitempty"`
	Location     string    `json:"location,omitempty"`
	Coordinates  string    `json:"coordinates,omitempty"`
	ISP          string    `json:"isp,omitempty"`
	CountryCode  string    `json:"country_code,omitempty"`
	CapturedTime time.Time `json:"captured_time"`
}

// AIModelInfo represents information about an AI model
type AIModelInfo struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name         string    `json:"name" gorm:"type:varchar(100);not null"`
	Version      string    `json:"version" gorm:"type:varchar(20);not null"`
	Type         string    `json:"type" gorm:"type:varchar(50);not null"`
	Accuracy     float64   `json:"accuracy" gorm:"type:float;not null"`
	LastTrainedAt time.Time `json:"last_trained_at" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"not null;default:now()"`
}
