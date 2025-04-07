package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// DocumentStatus represents the status of a document
type DocumentStatus string

const (
	DocStatusPending    DocumentStatus = "PENDING"
	DocStatusInReview   DocumentStatus = "IN_REVIEW"
	DocStatusVerified   DocumentStatus = "VERIFIED"
	DocStatusRejected   DocumentStatus = "REJECTED"
	DocStatusExpired    DocumentStatus = "EXPIRED"
	DocStatusIncomplete DocumentStatus = "INCOMPLETE"
)

// DocumentMetadata represents structured metadata for a document
type DocumentMetadata struct {
	OriginalFileName string                 `json:"original_file_name,omitempty"`
	Resolution       string                 `json:"resolution,omitempty"`
	PageCount        int                    `json:"page_count,omitempty"`
	ExtractedFields  map[string]string      `json:"extracted_fields,omitempty"`
	SecurityFeatures []string               `json:"security_features,omitempty"`
	VerificationData map[string]interface{} `json:"verification_data,omitempty"`
	ExtractionData   map[string]string      `json:"extraction_data,omitempty"`
}

// DocumentSummary represents a summary of a document
type DocumentSummary struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	Type            DocumentType
	Status          DocumentStatus
	FileName        string
	DocumentNumber  string
	ExpiryDate      *time.Time
	CreatedAt       time.Time
	VerifiedAt      *time.Time
	ProcessingTime  time.Duration
	ConfidenceScore float64
}

// EnhancedDocument represents a KYC document in the domain model with additional fields
type EnhancedDocument struct {
	ID                uuid.UUID      `json:"id"`
	UserID            uuid.UUID      `json:"user_id"`
	KYCID             *uuid.UUID     `json:"kyc_id,omitempty"`
	Type              DocumentType   `json:"type"`
	Status            DocumentStatus `json:"status"`
	
	// File information
	FileName          string         `json:"file_name"`
	FileSize          int64          `json:"file_size"`
	MimeType          string         `json:"mime_type"`
	FileHash          string         `json:"file_hash"`
	FilePath          string         `json:"file_path,omitempty"`
	FileURL           string         `json:"file_url,omitempty"`
	
	// Document details
	DocumentNumber    string         `json:"document_number,omitempty"`
	IssueDate         *time.Time     `json:"issue_date,omitempty"`
	ExpiryDate        *time.Time     `json:"expiry_date,omitempty"`
	IssuingCountry    string         `json:"issuing_country,omitempty"`
	IssuingAuthority  string         `json:"issuing_authority,omitempty"`
	
	// Verification information
	VerificationID    *uuid.UUID     `json:"verification_id,omitempty"`
	ConfidenceScore   float64        `json:"confidence_score,omitempty"`
	IsValid           bool           `json:"is_valid"`
	
	// Additional data
	Metadata          Metadata       `json:"metadata,omitempty"`
	
	// Timestamps
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	ExpiresAt         *time.Time     `json:"expires_at,omitempty"`
	VerifiedAt        *time.Time     `json:"verified_at,omitempty"`
	RejectedAt        *time.Time     `json:"rejected_at,omitempty"`
	
	// Rejection information
	RejectionReason   string         `json:"rejection_reason,omitempty"`
	RejectedBy        *uuid.UUID     `json:"rejected_by,omitempty"`
}

// Validate performs basic validation on the document
func (d *EnhancedDocument) Validate() error {
	if d.ID == uuid.Nil {
		return errors.New("document ID is required")
	}
	
	if d.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}
	
	if d.Type == "" {
		return errors.New("document type is required")
	}
	
	if d.FileName == "" {
		return errors.New("file name is required")
	}
	
	if d.FileSize <= 0 {
		return errors.New("file size must be greater than zero")
	}
	
	if d.FileHash == "" {
		return errors.New("file hash is required")
	}
	
	return nil
}

// IsExpired checks if the document is expired
func (d *EnhancedDocument) IsExpired() bool {
	if d.ExpiryDate == nil {
		return false
	}
	return time.Now().After(*d.ExpiryDate)
}

// IsVerified checks if the document is verified
func (d *EnhancedDocument) IsVerified() bool {
	return d.Status == DocStatusVerified && d.VerifiedAt != nil
}

// ToDocument converts an EnhancedDocument to a Document
func (d *EnhancedDocument) ToDocument() Document {
	return Document{
		ID:          d.ID,
		KYCID:       d.KYCID,
		Type:        d.Type,
		Hash:        d.FileHash,
		Status:      string(d.Status),
		ContentType: d.MimeType,
		Size:        d.FileSize,
		Metadata:    d.Metadata,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

// FromDocument creates an EnhancedDocument from a Document
func FromDocument(doc Document) EnhancedDocument {
	return EnhancedDocument{
		ID:          doc.ID,
		KYCID:       &doc.KYCID,
		Type:        doc.Type,
		Status:      DocumentStatus(doc.Status),
		FileHash:    doc.Hash,
		MimeType:    doc.ContentType,
		FileSize:    doc.Size,
		Metadata:    doc.Metadata,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}
}
