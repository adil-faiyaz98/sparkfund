package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// KYCStatus represents the status of a KYC verification
type KYCStatus string

const (
	KYCStatusPending    KYCStatus = "pending"
	KYCStatusInReview   KYCStatus = "in_review"
	KYCStatusVerified   KYCStatus = "verified"
	KYCStatusApproved   KYCStatus = "approved"
	KYCStatusRejected   KYCStatus = "rejected"
	KYCStatusExpired    KYCStatus = "expired"
	KYCStatusFlagged    KYCStatus = "flagged"
	KYCStatusIncomplete KYCStatus = "incomplete"
)

// RiskLevel represents the risk level of a KYC verification
type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "low"
	RiskLevelMedium RiskLevel = "medium"
	RiskLevelHigh   RiskLevel = "high"
)

// KYC represents a KYC verification record
type KYC struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	// Personal information
	FirstName   string `gorm:"not null" json:"first_name"`
	LastName    string `gorm:"not null" json:"last_name"`
	DateOfBirth string `gorm:"not null" json:"date_of_birth"` // Format: YYYY-MM-DD
	Nationality string `gorm:"type:varchar(100)" json:"nationality,omitempty"`
	Email       string `gorm:"type:varchar(255)" json:"email,omitempty"`
	PhoneNumber string `gorm:"type:varchar(50)" json:"phone_number,omitempty"`

	// Address information
	Address    string `gorm:"not null" json:"address"`
	City       string `gorm:"not null" json:"city"`
	State      string `gorm:"type:varchar(100)" json:"state,omitempty"`
	Country    string `gorm:"not null" json:"country"`
	PostalCode string `gorm:"not null" json:"postal_code"`

	// Document information
	DocumentType   string     `gorm:"not null" json:"document_type"`
	DocumentNumber string     `gorm:"not null" json:"document_number"`
	DocumentFront  string     `gorm:"not null" json:"document_front"`
	DocumentBack   string     `gorm:"not null" json:"document_back"`
	SelfieImage    string     `gorm:"not null" json:"selfie_image"`
	DocumentExpiry *time.Time `json:"document_expiry,omitempty"`

	// Risk assessment
	RiskLevel         RiskLevel `gorm:"type:varchar(20);default:'medium'" json:"risk_level"`
	RiskScore         float64   `gorm:"type:float;default:50" json:"risk_score"`
	TransactionAmount float64   `gorm:"type:decimal(15,2)" json:"transaction_amount,omitempty"`

	// Status information
	Status          KYCStatus `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	RejectionReason string    `gorm:"type:text" json:"rejection_reason,omitempty"`
	Notes           string    `gorm:"type:text" json:"notes,omitempty"`

	// Verification information
	VerifiedBy *uuid.UUID `gorm:"type:uuid" json:"verified_by,omitempty"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	ReviewedBy *uuid.UUID `gorm:"type:uuid" json:"reviewed_by,omitempty"`
	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`

	// Timestamps
	CreatedAt   time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"not null" json:"updated_at"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Additional data
	Metadata map[string]interface{} `gorm:"type:jsonb" json:"metadata,omitempty"`
}

// TableName specifies the table name for the KYC model
func (KYC) TableName() string {
	return "kyc_verifications"
}

// KYCRequest represents a request to create a KYC verification
type KYCRequest struct {
	FirstName         string  `json:"first_name" binding:"required"`
	LastName          string  `json:"last_name" binding:"required"`
	DateOfBirth       string  `json:"date_of_birth" binding:"required"` // Format: YYYY-MM-DD
	Nationality       string  `json:"nationality,omitempty"`
	Email             string  `json:"email,omitempty"`
	PhoneNumber       string  `json:"phone_number,omitempty"`
	Address           string  `json:"address" binding:"required"`
	City              string  `json:"city" binding:"required"`
	State             string  `json:"state,omitempty"`
	Country           string  `json:"country" binding:"required"`
	PostalCode        string  `json:"postal_code" binding:"required"`
	DocumentType      string  `json:"document_type" binding:"required"`
	DocumentNumber    string  `json:"document_number" binding:"required"`
	DocumentFront     string  `json:"document_front" binding:"required"`
	DocumentBack      string  `json:"document_back" binding:"required"`
	SelfieImage       string  `json:"selfie_image" binding:"required"`
	DocumentExpiry    string  `json:"document_expiry,omitempty"` // Format: YYYY-MM-DD
	TransactionAmount float64 `json:"transaction_amount,omitempty"`
}

// KYCResponse represents a response for a KYC verification
type KYCResponse struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	Status          KYCStatus  `json:"status"`
	RiskLevel       RiskLevel  `json:"risk_level"`
	RiskScore       float64    `json:"risk_score"`
	RejectionReason string     `json:"rejection_reason,omitempty"`
	VerifiedAt      *time.Time `json:"verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
}

// KYCStats represents statistics for KYC verifications
type KYCStats struct {
	TotalCount            int64
	PendingCount          int64
	VerifiedCount         int64
	RejectedCount         int64
	FlaggedCount          int64
	AverageRiskScore      float64
	HighRiskCount         int64
	MediumRiskCount       int64
	LowRiskCount          int64
	AverageProcessingTime time.Duration
}

// KYCReview represents a review of a KYC verification
type KYCReview struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	KYCID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"kyc_id"`
	ReviewerID     uuid.UUID      `gorm:"type:uuid;not null" json:"reviewer_id"`
	Status         string         `gorm:"type:varchar(20);not null" json:"status"` // APPROVED, REJECTED
	Reason         string         `gorm:"type:text" json:"reason,omitempty"`
	RiskAssessment string         `gorm:"type:text" json:"risk_assessment,omitempty"`
	Notes          string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt      time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for the KYCReview model
func (KYCReview) TableName() string {
	return "kyc_reviews"
}

// KYCWatchlist represents a watchlist entry for a user
type KYCWatchlist struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Reason     string         `gorm:"type:text;not null" json:"reason"`
	RiskLevel  RiskLevel      `gorm:"type:varchar(20);not null" json:"risk_level"`
	Source     string         `gorm:"type:varchar(100)" json:"source,omitempty"`
	ExpiryDate *time.Time     `json:"expiry_date,omitempty"`
	CreatedBy  uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt  time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for the KYCWatchlist model
func (KYCWatchlist) TableName() string {
	return "kyc_watchlist"
}
