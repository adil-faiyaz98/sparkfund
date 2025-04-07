package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// VerificationStatus represents the status of a verification
type VerificationStatus string

const (
	VerStatusPending    VerificationStatus = "PENDING"
	VerStatusInProgress VerificationStatus = "IN_PROGRESS"
	VerStatusCompleted  VerificationStatus = "COMPLETED"
	VerStatusApproved   VerificationStatus = "APPROVED"
	VerStatusRejected   VerificationStatus = "REJECTED"
	VerStatusFailed     VerificationStatus = "FAILED"
	VerStatusExpired    VerificationStatus = "EXPIRED"
)

// VerificationMethod represents the method used for verification
type VerificationMethod string

const (
	VerMethodManual     VerificationMethod = "MANUAL"
	VerMethodAutomated  VerificationMethod = "AUTOMATED"
	VerMethodThirdParty VerificationMethod = "THIRD_PARTY"
	VerMethodAI         VerificationMethod = "AI"
	VerMethodBiometric  VerificationMethod = "BIOMETRIC"
	VerMethodDocument   VerificationMethod = "DOCUMENT"
	VerMethodFacial     VerificationMethod = "FACIAL"
)

// EnhancedVerification represents a verification in the domain model
type EnhancedVerification struct {
	ID              uuid.UUID              `json:"id"`
	KYCID           *uuid.UUID             `json:"kyc_id,omitempty"`
	DocumentID      *uuid.UUID             `json:"document_id,omitempty"`
	Type            VerificationType       `json:"type"`
	Status          VerificationStatus     `json:"status"`
	Method          VerificationMethod     `json:"method"`
	VerifierID      *uuid.UUID             `json:"verifier_id,omitempty"`
	ConfidenceScore float64                `json:"confidence_score"`
	MatchScore      float64                `json:"match_score,omitempty"`
	FraudScore      float64                `json:"fraud_score,omitempty"`
	Notes           string                 `json:"notes,omitempty"`
	Metadata        Metadata               `json:"metadata,omitempty"`
	Result          Metadata               `json:"result,omitempty"`
	ErrorMessage    string                 `json:"error_message,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	CompletedAt     *time.Time             `json:"completed_at,omitempty"`
	ExpiresAt       *time.Time             `json:"expires_at,omitempty"`
}

// Validate performs basic validation on the verification
func (v *EnhancedVerification) Validate() error {
	if v.ID == uuid.Nil {
		return errors.New("verification ID is required")
	}
	
	if v.KYCID == nil && v.DocumentID == nil {
		return errors.New("either KYC ID or document ID is required")
	}
	
	if v.Type == "" {
		return errors.New("verification type is required")
	}
	
	if v.Method == "" {
		return errors.New("verification method is required")
	}
	
	return nil
}

// IsExpired checks if the verification is expired
func (v *EnhancedVerification) IsExpired() bool {
	if v.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*v.ExpiresAt)
}

// IsCompleted checks if the verification is completed
func (v *EnhancedVerification) IsCompleted() bool {
	return v.Status == VerStatusCompleted || v.Status == VerStatusApproved || v.CompletedAt != nil
}

// IsSuccessful checks if the verification was successful
func (v *EnhancedVerification) IsSuccessful() bool {
	return v.Status == VerStatusApproved || v.Status == VerStatusCompleted
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

// VerificationResult represents the result of a verification process
type VerificationResult struct {
	ID             uuid.UUID  `json:"id"`
	VerificationID uuid.UUID  `json:"verification_id"`
	Score          float64    `json:"score"`
	MatchScore     float64    `json:"match_score,omitempty"`
	FraudScore     float64    `json:"fraud_score,omitempty"`
	Success        bool       `json:"success"`
	Details        Metadata   `json:"details,omitempty"`
	ErrorMessage   string     `json:"error_message,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// ToVerification converts an EnhancedVerification to a Verification
func (v *EnhancedVerification) ToVerification() Verification {
	return Verification{
		ID:        v.ID,
		KYCID:     v.KYCID,
		Type:      v.Type,
		Status:    VerificationStatus(v.Status),
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}

// FromVerification creates an EnhancedVerification from a Verification
func FromVerification(ver Verification) EnhancedVerification {
	return EnhancedVerification{
		ID:        ver.ID,
		KYCID:     ver.KYCID,
		Type:      ver.Type,
		Status:    VerificationStatus(ver.Status),
		CreatedAt: ver.CreatedAt,
		UpdatedAt: ver.UpdatedAt,
	}
}
