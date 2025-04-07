package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// VerificationStatus represents the status of a verification
type VerificationStatus string

const (
	VerificationStatusPending    VerificationStatus = "pending"
	VerificationStatusInProgress VerificationStatus = "in_progress"
	VerificationStatusCompleted  VerificationStatus = "completed"
	VerificationStatusApproved   VerificationStatus = "approved"
	VerificationStatusRejected   VerificationStatus = "rejected"
	VerificationStatusFailed     VerificationStatus = "failed"
	VerificationStatusExpired    VerificationStatus = "expired"
)

// VerificationMethod represents the method used for verification
type VerificationMethod string

const (
	VerificationMethodManual     VerificationMethod = "manual"
	VerificationMethodAutomated  VerificationMethod = "automated"
	VerificationMethodThirdParty VerificationMethod = "third_party"
	VerificationMethodAI         VerificationMethod = "ai"
	VerificationMethodBiometric  VerificationMethod = "biometric"
	VerificationMethodDocument   VerificationMethod = "document"
	VerificationMethodFacial     VerificationMethod = "facial"
)

// VerificationType represents the type of verification
type VerificationType string

const (
	VerificationTypeDocument  VerificationType = "document"
	VerificationTypeBiometric VerificationType = "biometric"
	VerificationTypeAddress   VerificationType = "address"
	VerificationTypeIdentity  VerificationType = "identity"
	VerificationTypeLiveness  VerificationType = "liveness"
	VerificationTypeFacial    VerificationType = "facial"
)

// Verification represents a document verification record
type Verification struct {
	ID              uuid.UUID              `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	KYCID           *uuid.UUID             `gorm:"type:uuid;index" json:"kyc_id,omitempty"`
	DocumentID      *uuid.UUID             `gorm:"type:uuid;index" json:"document_id,omitempty"`
	Type            VerificationType       `gorm:"type:varchar(50);not null" json:"type"`
	Status          VerificationStatus     `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	Method          VerificationMethod     `gorm:"type:varchar(20);not null;default:'manual'" json:"method"`
	VerifierID      *uuid.UUID             `gorm:"type:uuid" json:"verifier_id,omitempty"`
	ConfidenceScore float64                `gorm:"type:float;not null;default:0" json:"confidence_score"`
	MatchScore      float64                `gorm:"type:float;default:0" json:"match_score,omitempty"`
	FraudScore      float64                `gorm:"type:float;default:0" json:"fraud_score,omitempty"`
	Notes           string                 `gorm:"type:text" json:"notes,omitempty"`
	Metadata        map[string]interface{} `gorm:"type:jsonb" json:"metadata,omitempty"`
	Result          map[string]interface{} `gorm:"type:jsonb" json:"result,omitempty"`
	ErrorMessage    string                 `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt       time.Time              `gorm:"not null" json:"created_at"`
	UpdatedAt       time.Time              `gorm:"not null" json:"updated_at"`
	CompletedAt     *time.Time             `json:"completed_at,omitempty"`
	ExpiresAt       *time.Time             `json:"expires_at,omitempty"`
	DeletedAt       gorm.DeletedAt         `gorm:"index" json:"-"`
}

// TableName specifies the table name for the Verification model
func (Verification) TableName() string {
	return "verifications"
}

// VerificationStats represents statistics for verifications
type VerificationStats struct {
	TotalCount            int64
	PendingCount          int64
	CompletedCount        int64
	ApprovedCount         int64
	RejectedCount         int64
	FailedCount           int64
	ExpiredCount          int64
	AverageConfidence     float64
	AverageMatchScore     float64
	AverageFraudScore     float64
	CompletionRate        float64
	SuccessRate           float64
	AverageProcessingTime time.Duration
	VerificationsByType   map[VerificationType]int64
	VerificationsByMethod map[VerificationMethod]int64
}

// VerificationSummary represents a summary of a verification
type VerificationSummary struct {
	ID              uuid.UUID
	KYCID           *uuid.UUID
	DocumentID      *uuid.UUID
	Type            VerificationType
	Status          VerificationStatus
	Method          VerificationMethod
	ConfidenceScore float64
	MatchScore      float64
	FraudScore      float64
	CreatedAt       time.Time
	CompletedAt     *time.Time
	ProcessingTime  time.Duration
	Success         bool
}

// VerificationHistory represents a history entry for a verification
type VerificationHistory struct {
	ID             uuid.UUID              `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	VerificationID uuid.UUID              `gorm:"type:uuid;not null;index" json:"verification_id"`
	Status         VerificationStatus     `gorm:"type:varchar(20);not null" json:"status"`
	Notes          string                 `gorm:"type:text" json:"notes,omitempty"`
	CreatedBy      uuid.UUID              `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt      time.Time              `gorm:"not null" json:"created_at"`
	Metadata       map[string]interface{} `gorm:"type:jsonb" json:"metadata,omitempty"`
}

// TableName specifies the table name for the VerificationHistory model
func (VerificationHistory) TableName() string {
	return "verification_history"
}

// VerificationResult represents the result of a verification process
type VerificationResult struct {
	ID             uuid.UUID              `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	VerificationID uuid.UUID              `gorm:"type:uuid;not null;index" json:"verification_id"`
	Score          float64                `gorm:"type:float;not null" json:"score"`
	MatchScore     float64                `gorm:"type:float" json:"match_score,omitempty"`
	FraudScore     float64                `gorm:"type:float" json:"fraud_score,omitempty"`
	Success        bool                   `gorm:"not null" json:"success"`
	Details        map[string]interface{} `gorm:"type:jsonb" json:"details,omitempty"`
	ErrorMessage   string                 `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt      time.Time              `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time              `gorm:"not null" json:"updated_at"`
	DeletedAt      gorm.DeletedAt         `gorm:"index" json:"-"`
}

// TableName specifies the table name for the VerificationResult model
func (VerificationResult) TableName() string {
	return "verification_results"
}
