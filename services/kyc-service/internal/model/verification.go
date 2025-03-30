package model

import (
	"time"

	"github.com/google/uuid"
)

// VerificationStatus represents the status of a verification
type VerificationStatus string

const (
	VerificationStatusPending    VerificationStatus = "pending"
	VerificationStatusInProgress VerificationStatus = "in_progress"
	VerificationStatusCompleted  VerificationStatus = "completed"
	VerificationStatusFailed     VerificationStatus = "failed"
	VerificationStatusExpired    VerificationStatus = "expired"
)

// VerificationMethod represents the method used for verification
type VerificationMethod string

const (
	VerificationMethodManual     VerificationMethod = "manual"
	VerificationMethodAutomated  VerificationMethod = "automated"
	VerificationMethodThirdParty VerificationMethod = "third_party"
)

// Verification represents a document verification record
type Verification struct {
	ID              uuid.UUID              `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	DocumentID      uuid.UUID              `gorm:"type:uuid;not null"`
	Status          VerificationStatus     `gorm:"type:varchar(20);not null;default:'pending'"`
	Method          VerificationMethod     `gorm:"type:varchar(20);not null;default:'manual'"`
	VerifierID      uuid.UUID              `gorm:"type:uuid"`
	ConfidenceScore float64                `gorm:"type:float;not null;default:0"`
	Notes           string                 `gorm:"type:text"`
	Metadata        map[string]interface{} `gorm:"type:jsonb"`
	CreatedAt       time.Time              `gorm:"not null"`
	UpdatedAt       time.Time              `gorm:"not null"`
	CompletedAt     *time.Time
	ExpiresAt       *time.Time
}

// VerificationStats represents statistics for verifications
type VerificationStats struct {
	TotalCount            int64
	PendingCount          int64
	CompletedCount        int64
	FailedCount           int64
	ExpiredCount          int64
	AverageConfidence     float64
	CompletionRate        float64
	AverageProcessingTime time.Duration
}

// VerificationSummary represents a summary of a verification
type VerificationSummary struct {
	ID              uuid.UUID
	DocumentID      uuid.UUID
	Status          VerificationStatus
	Method          VerificationMethod
	ConfidenceScore float64
	CreatedAt       time.Time
	CompletedAt     *time.Time
	ProcessingTime  time.Duration
}

// VerificationHistory represents a history entry for a verification
type VerificationHistory struct {
	ID             uuid.UUID          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	VerificationID uuid.UUID          `gorm:"type:uuid;not null"`
	Status         VerificationStatus `gorm:"type:varchar(20);not null"`
	Notes          string             `gorm:"type:text"`
	CreatedBy      uuid.UUID          `gorm:"type:uuid;not null"`
	CreatedAt      time.Time          `gorm:"not null"`
}
