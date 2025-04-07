package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DocumentType represents the type of document
type DocumentType string

const (
	// Standard document types
	DocumentTypeID              DocumentType = "id" // Keep for backward compatibility
	DocumentTypePassport        DocumentType = "passport"
	DocumentTypeDriversLicense  DocumentType = "drivers_license"
	DocumentTypeNationalID      DocumentType = "national_id"
	DocumentTypeResidencePermit DocumentType = "residence_permit"

	// Additional document types
	DocumentTypeProofOfAddress  DocumentType = "proof_of_address"
	DocumentTypeTaxDocument     DocumentType = "tax_document"
	DocumentTypeBankStatement   DocumentType = "bank_statement"
	DocumentTypeEmploymentProof DocumentType = "employment_proof"
	DocumentTypeUtilityBill     DocumentType = "utility_bill"
	DocumentTypeOther           DocumentType = "other"
)

// DocumentStatus represents the status of a document
type DocumentStatus string

const (
	DocumentStatusPending    DocumentStatus = "pending"
	DocumentStatusInReview   DocumentStatus = "in_review"
	DocumentStatusVerified   DocumentStatus = "verified"
	DocumentStatusRejected   DocumentStatus = "rejected"
	DocumentStatusExpired    DocumentStatus = "expired"
	DocumentStatusIncomplete DocumentStatus = "incomplete"
)

// Document represents a KYC document
type Document struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID            uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	KYCVerificationID *uuid.UUID     `gorm:"type:uuid;index" json:"kyc_verification_id,omitempty"`
	Type              DocumentType   `gorm:"type:varchar(50);not null;index" json:"type"`
	Status            DocumentStatus `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`

	// File information
	FileName string `gorm:"type:varchar(255);not null" json:"file_name"`
	FileSize int64  `gorm:"not null" json:"file_size"`
	MimeType string `gorm:"type:varchar(100);not null" json:"mime_type"`
	FileHash string `gorm:"type:varchar(64);not null" json:"file_hash"`
	FilePath string `gorm:"type:varchar(255);not null" json:"file_path"`
	FileURL  string `gorm:"type:varchar(255)" json:"file_url,omitempty"`

	// Document details
	DocumentNumber   string     `gorm:"type:varchar(100)" json:"document_number,omitempty"`
	IssueDate        *time.Time `json:"issue_date,omitempty"`
	ExpiryDate       *time.Time `json:"expiry_date,omitempty"`
	IssuingCountry   string     `gorm:"type:varchar(100)" json:"issuing_country,omitempty"`
	IssuingAuthority string     `gorm:"type:varchar(255)" json:"issuing_authority,omitempty"`

	// Verification information
	VerificationID  *uuid.UUID `gorm:"type:uuid" json:"verification_id,omitempty"`
	ConfidenceScore float64    `gorm:"type:float;default:0" json:"confidence_score,omitempty"`
	IsValid         bool       `gorm:"default:false" json:"is_valid"`

	// Additional data
	Metadata map[string]interface{} `gorm:"type:jsonb" json:"metadata,omitempty"`

	// Timestamps
	CreatedAt  time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"not null" json:"updated_at"`
	ExpiresAt  *time.Time     `json:"expires_at,omitempty"`
	VerifiedAt *time.Time     `json:"verified_at,omitempty"`
	RejectedAt *time.Time     `json:"rejected_at,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Rejection information
	RejectionReason string     `gorm:"type:text" json:"rejection_reason,omitempty"`
	RejectedBy      *uuid.UUID `gorm:"type:uuid" json:"rejected_by,omitempty"`
}

// TableName specifies the table name for the Document model
func (Document) TableName() string {
	return "documents"
}

// IsExpired checks if the document is expired
func (d *Document) IsExpired() bool {
	if d.ExpiryDate == nil {
		return false
	}
	return time.Now().After(*d.ExpiryDate)
}

// DocumentStats represents statistics for documents
type DocumentStats struct {
	TotalCount         int64
	PendingCount       int64
	VerifiedCount      int64
	RejectedCount      int64
	ExpiredCount       int64
	AverageFileSize    int64
	TotalFileSize      int64
	DocumentTypeCounts map[DocumentType]int64
	ProcessingTimeAvg  time.Duration
	VerificationRate   float64 // Percentage of documents that were verified
}

// DocumentSummary represents a summary of a document
type DocumentSummary struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	Type            DocumentType
	Status          DocumentStatus
	FileName        string
	FileSize        int64
	DocumentNumber  string
	ExpiryDate      *time.Time
	CreatedAt       time.Time
	VerifiedAt      *time.Time
	ProcessingTime  time.Duration
	ConfidenceScore float64
}

// DocumentHistory represents a history entry for a document
type DocumentHistory struct {
	ID         uuid.UUID              `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	DocumentID uuid.UUID              `gorm:"type:uuid;not null;index" json:"document_id"`
	Status     DocumentStatus         `gorm:"type:varchar(20);not null" json:"status"`
	Notes      string                 `gorm:"type:text" json:"notes"`
	CreatedBy  uuid.UUID              `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt  time.Time              `gorm:"not null" json:"created_at"`
	Metadata   map[string]interface{} `gorm:"type:jsonb" json:"metadata,omitempty"`
}

// TableName specifies the table name for the DocumentHistory model
func (DocumentHistory) TableName() string {
	return "document_history"
}

// DocumentMetadata represents structured metadata for a document
type DocumentMetadata struct {
	OriginalFileName string                 `json:"original_file_name,omitempty"`
	Resolution       string                 `json:"resolution,omitempty"`
	PageCount        int                    `json:"page_count,omitempty"`
	ExtractedFields  map[string]string      `json:"extracted_fields,omitempty"`
	SecurityFeatures []string               `json:"security_features,omitempty"`
	VerificationData map[string]interface{} `json:"verification_data,omitempty"`
}
