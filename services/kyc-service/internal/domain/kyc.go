package domain

import (
	"time"

	"github.com/google/uuid"
)

type KYCStatus string

const (
	StatusPending    KYCStatus = "PENDING"
	StatusVerified   KYCStatus = "VERIFIED"
	StatusRejected   KYCStatus = "REJECTED"
	StatusIncomplete KYCStatus = "INCOMPLETE"
)

type KYC struct {
	ID            uuid.UUID      `json:"id"`
	UserID        uuid.UUID      `json:"user_id"`
	Status        KYCStatus      `json:"status"`
	RiskLevel     string         `json:"risk_level"`
	Documents     []Document     `json:"documents"`
	Verifications []Verification `json:"verifications"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	CompletedAt   *time.Time     `json:"completed_at,omitempty"`
}

// Document represents a KYC document in the domain model
type Document struct {
	ID     uuid.UUID      `json:"id"`
	UserID uuid.UUID      `json:"user_id"`
	KYCID  *uuid.UUID     `json:"kyc_id,omitempty"`
	Type   DocumentType   `json:"type"`
	Status DocumentStatus `json:"status"`

	// File information
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	MimeType string `json:"mime_type"`
	FileHash string `json:"file_hash"`
	FilePath string `json:"file_path,omitempty"`
	FileURL  string `json:"file_url,omitempty"`

	// Document details
	DocumentNumber   string     `json:"document_number,omitempty"`
	IssueDate        *time.Time `json:"issue_date,omitempty"`
	ExpiryDate       *time.Time `json:"expiry_date,omitempty"`
	IssuingCountry   string     `json:"issuing_country,omitempty"`
	IssuingAuthority string     `json:"issuing_authority,omitempty"`

	// Verification information
	VerificationID  *uuid.UUID `json:"verification_id,omitempty"`
	ConfidenceScore float64    `json:"confidence_score,omitempty"`
	IsValid         bool       `json:"is_valid"`

	// Additional data
	Metadata Metadata `json:"metadata,omitempty"`

	// Timestamps
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	RejectedAt *time.Time `json:"rejected_at,omitempty"`

	// Rejection information
	RejectionReason string     `json:"rejection_reason,omitempty"`
	RejectedBy      *uuid.UUID `json:"rejected_by,omitempty"`
}

type DocumentType string

const (
	DocTypePassport      DocumentType = "PASSPORT"
	DocTypeDriverLicense DocumentType = "DRIVERS_LICENSE"
	DocTypeIDCard        DocumentType = "ID_CARD"
	DocTypeUtilityBill   DocumentType = "UTILITY_BILL"
	DocTypeBankStatement DocumentType = "BANK_STATEMENT"
)

type Metadata map[string]interface{}

type Verification struct {
	ID        uuid.UUID           `json:"id"`
	KYCID     uuid.UUID           `json:"kyc_id"`
	Type      VerificationType    `json:"type"`
	Status    VerificationStatus  `json:"status"`
	Result    *VerificationResult `json:"result,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type VerificationType string

const (
	VerificationDocument  VerificationType = "DOCUMENT"
	VerificationBiometric VerificationType = "BIOMETRIC"
	VerificationAddress   VerificationType = "ADDRESS"
)

type VerificationStatus string

const (
	VerificationPending  VerificationStatus = "PENDING"
	VerificationApproved VerificationStatus = "APPROVED"
	VerificationRejected VerificationStatus = "REJECTED"
	VerificationFailed   VerificationStatus = "FAILED"
)

type VerificationResult struct {
	ID             uuid.UUID   `json:"id"`
	VerificationID uuid.UUID   `json:"verification_id"`
	Score          float64     `json:"score"`
	Details        interface{} `json:"details"`
	CreatedAt      time.Time   `json:"created_at"`
}
