package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID            uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email         string     `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash  string     `json:"-" gorm:"type:varchar(255);not null"`
	FirstName     string     `json:"first_name" gorm:"type:varchar(100)"`
	LastName      string     `json:"last_name" gorm:"type:varchar(100)"`
	Role          string     `json:"role" gorm:"type:varchar(20);not null;default:'user'"`
	MFAEnabled    bool       `json:"mfa_enabled" gorm:"not null;default:false"`
	MFASecret     string     `json:"-" gorm:"type:varchar(32)"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	LastLoginIP   string     `json:"last_login_ip,omitempty" gorm:"type:varchar(45)"`
	LoginAttempts int        `json:"-" gorm:"not null;default:0"`
	LockedUntil   *time.Time `json:"-"`
	CreatedAt     time.Time  `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"not null;default:now()"`
}

// Session represents a user session
type Session struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	RefreshToken string    `json:"-" gorm:"type:varchar(255);not null"`
	UserAgent    string    `json:"user_agent" gorm:"type:varchar(255)"`
	IPAddress    string    `json:"ip_address" gorm:"type:varchar(45)"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// DocumentType represents the type of document
type DocumentType string

const (
	DocumentTypePassport       DocumentType = "PASSPORT"
	DocumentTypeDriversLicense DocumentType = "DRIVERS_LICENSE"
	DocumentTypeIDCard         DocumentType = "ID_CARD"
	DocumentTypeBankStatement  DocumentType = "BANK_STATEMENT"
	DocumentTypeUtilityBill    DocumentType = "UTILITY_BILL"
	DocumentTypeSelfie         DocumentType = "SELFIE"
)

// Document represents a document in the system
type Document struct {
	ID          uuid.UUID    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID      uuid.UUID    `json:"user_id" gorm:"type:uuid;not null"`
	Type        DocumentType `json:"type" gorm:"type:varchar(50);not null"`
	Name        string       `json:"name" gorm:"type:varchar(255);not null"`
	Path        string       `json:"path" gorm:"type:varchar(255);not null"`
	Size        int64        `json:"size" gorm:"not null"`
	ContentType string       `json:"content_type" gorm:"type:varchar(100);not null"`
	Metadata    []byte       `json:"metadata,omitempty" gorm:"type:jsonb"`
	UploadedAt  time.Time    `json:"uploaded_at" gorm:"not null;default:now()"`
	ExpiresAt   *time.Time   `json:"expires_at,omitempty"`
	CreatedAt   time.Time    `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt   time.Time    `json:"updated_at" gorm:"not null;default:now()"`
}

// VerificationStatus represents the status of a verification
type VerificationStatus string

const (
	VerificationStatusPending   VerificationStatus = "PENDING"
	VerificationStatusInProcess VerificationStatus = "IN_PROCESS"
	VerificationStatusCompleted VerificationStatus = "COMPLETED"
	VerificationStatusFailed    VerificationStatus = "FAILED"
	VerificationStatusExpired   VerificationStatus = "EXPIRED"
)

// VerificationMethod represents the method of verification
type VerificationMethod string

const (
	VerificationMethodDocument  VerificationMethod = "DOCUMENT"
	VerificationMethodBiometric VerificationMethod = "BIOMETRIC"
	VerificationMethodAI        VerificationMethod = "AI"
)

// Verification represents a verification in the system
type Verification struct {
	ID          uuid.UUID          `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID      uuid.UUID          `json:"user_id" gorm:"type:uuid;not null"`
	KYCID       uuid.UUID          `json:"kyc_id" gorm:"type:uuid;not null"`
	DocumentID  *uuid.UUID         `json:"document_id" gorm:"type:uuid"`
	SelfieID    *uuid.UUID         `json:"selfie_id" gorm:"type:uuid"`
	VerifierID  *uuid.UUID         `json:"verifier_id" gorm:"type:uuid"`
	Method      VerificationMethod `json:"method" gorm:"type:varchar(20);not null"`
	Status      VerificationStatus `json:"status" gorm:"type:varchar(20);not null"`
	Notes       string             `json:"notes" gorm:"type:text"`
	CompletedAt *time.Time         `json:"completed_at"`
	CreatedAt   time.Time          `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt   time.Time          `json:"updated_at" gorm:"not null;default:now()"`
}

// DocumentAnalysisResult represents the result of AI document analysis
type DocumentAnalysisResult struct {
	ID             uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	VerificationID uuid.UUID `json:"verification_id" gorm:"type:uuid;not null"`
	DocumentID     uuid.UUID `json:"document_id" gorm:"type:uuid;not null"`
	DocumentType   string    `json:"document_type" gorm:"type:varchar(50);not null"`
	IsAuthentic    bool      `json:"is_authentic" gorm:"not null"`
	Confidence     float64   `json:"confidence" gorm:"type:float;not null"`
	ExtractedData  []byte    `json:"extracted_data" gorm:"type:jsonb"`
	Issues         []byte    `json:"issues" gorm:"type:jsonb"`
	CreatedAt      time.Time `json:"created_at" gorm:"not null;default:now()"`
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
	ID             uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	VerificationID uuid.UUID `json:"verification_id" gorm:"type:uuid;not null"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	RiskScore      float64   `json:"risk_score" gorm:"type:float;not null"`
	RiskLevel      string    `json:"risk_level" gorm:"type:varchar(20);not null"`
	RiskFactors    []byte    `json:"risk_factors" gorm:"type:jsonb"`
	DeviceInfo     []byte    `json:"device_info" gorm:"type:jsonb"`
	IPAddress      string    `json:"ip_address" gorm:"type:varchar(45)"`
	Location       string    `json:"location" gorm:"type:varchar(100)"`
	CreatedAt      time.Time `json:"created_at" gorm:"not null;default:now()"`
}

// AnomalyDetectionResult represents the result of anomaly detection
type AnomalyDetectionResult struct {
	ID             uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	VerificationID uuid.UUID `json:"verification_id" gorm:"type:uuid;not null"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	IsAnomaly      bool      `json:"is_anomaly" gorm:"not null"`
	AnomalyScore   float64   `json:"anomaly_score" gorm:"type:float;not null"`
	AnomalyType    string    `json:"anomaly_type" gorm:"type:varchar(50)"`
	Reasons        []byte    `json:"reasons" gorm:"type:jsonb"`
	DeviceInfo     []byte    `json:"device_info" gorm:"type:jsonb"`
	CreatedAt      time.Time `json:"created_at" gorm:"not null;default:now()"`
}

// AIModelInfo represents information about an AI model
type AIModelInfo struct {
	ID            uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name          string    `json:"name" gorm:"type:varchar(100);not null"`
	Version       string    `json:"version" gorm:"type:varchar(20);not null"`
	Type          string    `json:"type" gorm:"type:varchar(50);not null"`
	Accuracy      float64   `json:"accuracy" gorm:"type:float;not null"`
	LastTrainedAt time.Time `json:"last_trained_at" gorm:"not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"not null;default:now()"`
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

// AutoMigrate migrates all models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Session{},
		&Document{},
		&Verification{},
		&DocumentAnalysisResult{},
		&FaceMatchResult{},
		&RiskAnalysisResult{},
		&AnomalyDetectionResult{},
		&AIModelInfo{},
	)
}
