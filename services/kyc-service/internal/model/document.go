package model

import (
	"time"

	"github.com/google/uuid"
)

// DocumentType represents the type of document
type DocumentType string

const (
	DocumentTypeID              DocumentType = "id"
	DocumentTypeProofOfAddress  DocumentType = "proof_of_address"
	DocumentTypeTaxDocument     DocumentType = "tax_document"
	DocumentTypeBankStatement   DocumentType = "bank_statement"
	DocumentTypeEmploymentProof DocumentType = "employment_proof"
	DocumentTypeOther           DocumentType = "other"
)

// DocumentStatus represents the status of a document
type DocumentStatus string

const (
	DocumentStatusPending   DocumentStatus = "pending"
	DocumentStatusVerified  DocumentStatus = "verified"
	DocumentStatusRejected  DocumentStatus = "rejected"
	DocumentStatusExpired   DocumentStatus = "expired"
)

// Document represents a KYC document
type Document struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null"`
	Type        DocumentType   `gorm:"type:varchar(50);not null"`
	Status      DocumentStatus `gorm:"type:varchar(20);not null;default:'pending'"`
	FileName    string         `gorm:"type:varchar(255);not null"`
	FileSize    int64          `gorm:"not null"`
	MimeType    string         `gorm:"type:varchar(100);not null"`
	FileHash    string         `gorm:"type:varchar(64);not null"`
	FilePath    string         `gorm:"type:varchar(255);not null"`
	Metadata    map[string]interface{} `gorm:"type:jsonb"`
	CreatedAt   time.Time      `gorm:"not null"`
	UpdatedAt   time.Time      `gorm:"not null"`
	ExpiresAt   *time.Time
	VerifiedAt  *time.Time
	RejectedAt  *time.Time
	RejectionReason string     `gorm:"type:text"`
}

// DocumentStats represents statistics for documents
type DocumentStats struct {
	TotalCount          int64
	PendingCount        int64
	VerifiedCount       int64
	RejectedCount       int64
	ExpiredCount        int64
	AverageFileSize     int64
	TotalFileSize       int64
	DocumentTypeCounts  map[DocumentType]int64
}

// DocumentSummary represents a summary of a document
type DocumentSummary struct {
	ID              uuid.UUID
	Type            DocumentType
	Status          DocumentStatus
	FileName        string
	FileSize        int64
	CreatedAt       time.Time
	VerifiedAt      *time.Time
	ProcessingTime  time.Duration
}

// DocumentHistory represents a history entry for a document
type DocumentHistory struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	DocumentID      uuid.UUID      `gorm:"type:uuid;not null"`
	Status          DocumentStatus `gorm:"type:varchar(20);not null"`
	Notes           string         `gorm:"type:text"`
	CreatedBy       uuid.UUID      `gorm:"type:uuid;not null"`
	CreatedAt       time.Time      `gorm:"not null"`
} 