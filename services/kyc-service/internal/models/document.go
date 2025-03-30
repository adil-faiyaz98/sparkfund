package models

import (
	"time"

	"github.com/google/uuid"
)

// DocumentType represents the type of KYC document
type DocumentType string

const (
	DocumentTypePassport           DocumentType = "passport"
	DocumentTypeNationalID         DocumentType = "national_id"
	DocumentTypeDrivingLicense     DocumentType = "driving_license"
	DocumentTypeProofOfAddress     DocumentType = "proof_of_address"
	DocumentTypeBankStatement      DocumentType = "bank_statement"
	DocumentTypeTaxReturn          DocumentType = "tax_return"
	DocumentTypeEmploymentContract DocumentType = "employment_contract"
	DocumentTypeUtilityBill        DocumentType = "utility_bill"
)

// DocumentStatus represents the verification status of a document
type DocumentStatus string

const (
	DocumentStatusPending  DocumentStatus = "pending"
	DocumentStatusVerified DocumentStatus = "verified"
	DocumentStatusRejected DocumentStatus = "rejected"
	DocumentStatusExpired  DocumentStatus = "expired"
)

// VerificationMethod represents the method used to verify a document
type VerificationMethod string

const (
	VerificationMethodManual     VerificationMethod = "manual"
	VerificationMethodAutomated  VerificationMethod = "automated"
	VerificationMethodThirdParty VerificationMethod = "third_party"
)

// Document represents a KYC document in the system
type Document struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Type      DocumentType   `json:"document_type" gorm:"type:varchar(50);not null"`
	Status    DocumentStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	FileData  []byte         `json:"-" gorm:"type:bytea;not null"`
	FileHash  string         `json:"file_hash" gorm:"type:varchar(64);not null"`
	MimeType  string         `json:"mime_type" gorm:"type:varchar(100);not null"`
	FileSize  int64          `json:"file_size" gorm:"type:bigint;not null"`
	Metadata  JSONMap        `json:"metadata" gorm:"type:jsonb"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp with time zone;not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp with time zone;not null"`
	DeletedAt *time.Time     `json:"-" gorm:"type:timestamp with time zone"`
}

// VerificationDetails represents the verification information for a document
type VerificationDetails struct {
	ID                 uuid.UUID          `json:"id" gorm:"type:uuid;primary_key"`
	DocumentID         uuid.UUID          `json:"document_id" gorm:"type:uuid;not null;unique"`
	VerifiedBy         uuid.UUID          `json:"verified_by" gorm:"type:uuid;not null"`
	VerifiedAt         time.Time          `json:"verified_at" gorm:"type:timestamp with time zone;not null"`
	VerificationMethod VerificationMethod `json:"verification_method" gorm:"type:varchar(20);not null"`
	ConfidenceScore    float64            `json:"confidence_score" gorm:"type:float;not null"`
	RejectionReason    *string            `json:"rejection_reason,omitempty" gorm:"type:text"`
	Notes              *string            `json:"notes,omitempty" gorm:"type:text"`
	CreatedAt          time.Time          `json:"created_at" gorm:"type:timestamp with time zone;not null"`
	UpdatedAt          time.Time          `json:"updated_at" gorm:"type:timestamp with time zone;not null"`
}

// JSONMap is a type alias for map[string]interface{} to handle JSON data
type JSONMap map[string]interface{}

// TableName specifies the table name for the Document model
func (Document) TableName() string {
	return "documents"
}

// TableName specifies the table name for the VerificationDetails model
func (VerificationDetails) TableName() string {
	return "verification_details"
}
