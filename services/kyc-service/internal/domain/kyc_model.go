package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// KYCStatus represents the status of a KYC verification
type KYCStatus string

const (
	KYCStatusPending    KYCStatus = "PENDING"
	KYCStatusInReview   KYCStatus = "IN_REVIEW"
	KYCStatusVerified   KYCStatus = "VERIFIED"
	KYCStatusApproved   KYCStatus = "APPROVED"
	KYCStatusRejected   KYCStatus = "REJECTED"
	KYCStatusExpired    KYCStatus = "EXPIRED"
	KYCStatusFlagged    KYCStatus = "FLAGGED"
	KYCStatusIncomplete KYCStatus = "INCOMPLETE"
)

// RiskLevel represents the risk level of a KYC verification
type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "LOW"
	RiskLevelMedium RiskLevel = "MEDIUM"
	RiskLevelHigh   RiskLevel = "HIGH"
)

// EnhancedKYC represents a KYC verification in the domain model
type EnhancedKYC struct {
	ID              uuid.UUID      `json:"id"`
	UserID          uuid.UUID      `json:"user_id"`
	
	// Personal information
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	DateOfBirth     string         `json:"date_of_birth"` // Format: YYYY-MM-DD
	Nationality     string         `json:"nationality,omitempty"`
	Email           string         `json:"email,omitempty"`
	PhoneNumber     string         `json:"phone_number,omitempty"`
	
	// Address information
	Address         string         `json:"address"`
	City            string         `json:"city"`
	State           string         `json:"state,omitempty"`
	Country         string         `json:"country"`
	PostalCode      string         `json:"postal_code"`
	
	// Document information
	DocumentType    string         `json:"document_type"`
	DocumentNumber  string         `json:"document_number"`
	DocumentFront   string         `json:"document_front"`
	DocumentBack    string         `json:"document_back"`
	SelfieImage     string         `json:"selfie_image"`
	DocumentExpiry  *time.Time     `json:"document_expiry,omitempty"`
	
	// Risk assessment
	RiskLevel       RiskLevel      `json:"risk_level"`
	RiskScore       float64        `json:"risk_score"`
	TransactionAmount float64      `json:"transaction_amount,omitempty"`
	
	// Status information
	Status          KYCStatus      `json:"status"`
	RejectionReason string         `json:"rejection_reason,omitempty"`
	Notes           string         `json:"notes,omitempty"`
	
	// Verification information
	VerifiedBy      *uuid.UUID     `json:"verified_by,omitempty"`
	VerifiedAt      *time.Time     `json:"verified_at,omitempty"`
	ReviewedBy      *uuid.UUID     `json:"reviewed_by,omitempty"`
	ReviewedAt      *time.Time     `json:"reviewed_at,omitempty"`
	
	// Related entities
	Documents       []EnhancedDocument    `json:"documents,omitempty"`
	Verifications   []EnhancedVerification `json:"verifications,omitempty"`
	
	// Timestamps
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	CompletedAt     *time.Time     `json:"completed_at,omitempty"`
	ExpiresAt       *time.Time     `json:"expires_at,omitempty"`
	
	// Additional data
	Metadata        Metadata       `json:"metadata,omitempty"`
}

// Validate performs basic validation on the KYC verification
func (k *EnhancedKYC) Validate() error {
	if k.ID == uuid.Nil {
		return errors.New("KYC ID is required")
	}
	
	if k.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}
	
	if k.FirstName == "" {
		return errors.New("first name is required")
	}
	
	if k.LastName == "" {
		return errors.New("last name is required")
	}
	
	if k.DateOfBirth == "" {
		return errors.New("date of birth is required")
	}
	
	if k.Address == "" {
		return errors.New("address is required")
	}
	
	if k.City == "" {
		return errors.New("city is required")
	}
	
	if k.Country == "" {
		return errors.New("country is required")
	}
	
	if k.PostalCode == "" {
		return errors.New("postal code is required")
	}
	
	if k.DocumentType == "" {
		return errors.New("document type is required")
	}
	
	if k.DocumentNumber == "" {
		return errors.New("document number is required")
	}
	
	if k.DocumentFront == "" {
		return errors.New("document front is required")
	}
	
	if k.DocumentBack == "" {
		return errors.New("document back is required")
	}
	
	if k.SelfieImage == "" {
		return errors.New("selfie image is required")
	}
	
	return nil
}

// IsExpired checks if the KYC verification is expired
func (k *EnhancedKYC) IsExpired() bool {
	if k.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*k.ExpiresAt)
}

// IsVerified checks if the KYC verification is verified
func (k *EnhancedKYC) IsVerified() bool {
	return (k.Status == KYCStatusVerified || k.Status == KYCStatusApproved) && k.VerifiedAt != nil
}

// IsHighRisk checks if the KYC verification is high risk
func (k *EnhancedKYC) IsHighRisk() bool {
	return k.RiskLevel == RiskLevelHigh || k.RiskScore >= 75
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
	ID             uuid.UUID  `json:"id"`
	KYCID          uuid.UUID  `json:"kyc_id"`
	ReviewerID     uuid.UUID  `json:"reviewer_id"`
	Status         string     `json:"status"` // APPROVED, REJECTED
	Reason         string     `json:"reason,omitempty"`
	RiskAssessment string     `json:"risk_assessment,omitempty"`
	Notes          string     `json:"notes,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ToKYC converts an EnhancedKYC to a KYC
func (k *EnhancedKYC) ToKYC() KYC {
	return KYC{
		ID:            k.ID,
		UserID:        k.UserID,
		Status:        k.Status,
		RiskLevel:     string(k.RiskLevel),
		CreatedAt:     k.CreatedAt,
		UpdatedAt:     k.UpdatedAt,
		CompletedAt:   k.CompletedAt,
	}
}

// FromKYC creates an EnhancedKYC from a KYC
func FromKYC(kyc KYC) EnhancedKYC {
	return EnhancedKYC{
		ID:            kyc.ID,
		UserID:        kyc.UserID,
		Status:        kyc.Status,
		RiskLevel:     RiskLevel(kyc.RiskLevel),
		CreatedAt:     kyc.CreatedAt,
		UpdatedAt:     kyc.UpdatedAt,
		CompletedAt:   kyc.CompletedAt,
	}
}
