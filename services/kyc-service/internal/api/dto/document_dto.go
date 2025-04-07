package dto

import (
	"time"

	"github.com/google/uuid"
	"sparkfund/services/kyc-service/internal/domain"
)

// DocumentUploadRequest represents a request to upload a document
type DocumentUploadRequest struct {
	Type     string                 `json:"type" binding:"required"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// DocumentResponse represents a document response
type DocumentResponse struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	Type            string    `json:"type"`
	Status          string    `json:"status"`
	FileName        string    `json:"file_name"`
	FileSize        int64     `json:"file_size"`
	MimeType        string    `json:"mime_type"`
	FileURL         string    `json:"file_url,omitempty"`
	DocumentNumber  string    `json:"document_number,omitempty"`
	IssueDate       string    `json:"issue_date,omitempty"`
	ExpiryDate      string    `json:"expiry_date,omitempty"`
	IssuingCountry  string    `json:"issuing_country,omitempty"`
	ConfidenceScore float64   `json:"confidence_score,omitempty"`
	IsValid         bool      `json:"is_valid"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	VerifiedAt      string    `json:"verified_at,omitempty"`
}

// DocumentListResponse represents a paginated list of documents
type DocumentListResponse struct {
	Documents []DocumentResponse `json:"documents"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
}

// DocumentStatusUpdateRequest represents a request to update document status
type DocumentStatusUpdateRequest struct {
	Status  string `json:"status" binding:"required"`
	Notes   string `json:"notes,omitempty"`
	VerifierID uuid.UUID `json:"verifier_id" binding:"required"`
}

// DocumentStatsResponse represents document statistics
type DocumentStatsResponse struct {
	TotalCount         int64                  `json:"total_count"`
	PendingCount       int64                  `json:"pending_count"`
	VerifiedCount      int64                  `json:"verified_count"`
	RejectedCount      int64                  `json:"rejected_count"`
	ExpiredCount       int64                  `json:"expired_count"`
	AverageFileSize    int64                  `json:"average_file_size"`
	TotalFileSize      int64                  `json:"total_file_size"`
	ProcessingTimeAvg  string                 `json:"processing_time_avg"`
	VerificationRate   float64                `json:"verification_rate"`
	DocumentTypeCounts map[string]int64       `json:"document_type_counts"`
}

// FromDomainDocument converts a domain document to a document response
func FromDomainDocument(doc *domain.EnhancedDocument) DocumentResponse {
	response := DocumentResponse{
		ID:              doc.ID,
		UserID:          doc.UserID,
		Type:            string(doc.Type),
		Status:          string(doc.Status),
		FileName:        doc.FileName,
		FileSize:        doc.FileSize,
		MimeType:        doc.MimeType,
		FileURL:         doc.FileURL,
		DocumentNumber:  doc.DocumentNumber,
		ConfidenceScore: doc.ConfidenceScore,
		IsValid:         doc.IsValid,
		CreatedAt:       doc.CreatedAt,
		UpdatedAt:       doc.UpdatedAt,
	}

	if doc.IssueDate != nil {
		response.IssueDate = doc.IssueDate.Format("2006-01-02")
	}

	if doc.ExpiryDate != nil {
		response.ExpiryDate = doc.ExpiryDate.Format("2006-01-02")
	}

	if doc.VerifiedAt != nil {
		response.VerifiedAt = doc.VerifiedAt.Format("2006-01-02T15:04:05Z")
	}

	return response
}

// FromDomainDocuments converts a slice of domain documents to document responses
func FromDomainDocuments(docs []*domain.EnhancedDocument) []DocumentResponse {
	responses := make([]DocumentResponse, len(docs))
	for i, doc := range docs {
		responses[i] = FromDomainDocument(doc)
	}
	return responses
}

// FromDomainDocumentStats converts domain document stats to a document stats response
func FromDomainDocumentStats(stats *domain.DocumentStats) DocumentStatsResponse {
	response := DocumentStatsResponse{
		TotalCount:        stats.TotalCount,
		PendingCount:      stats.PendingCount,
		VerifiedCount:     stats.VerifiedCount,
		RejectedCount:     stats.RejectedCount,
		ExpiredCount:      stats.ExpiredCount,
		AverageFileSize:   stats.AverageFileSize,
		TotalFileSize:     stats.TotalFileSize,
		ProcessingTimeAvg: stats.ProcessingTimeAvg.String(),
		VerificationRate:  stats.VerificationRate,
		DocumentTypeCounts: make(map[string]int64),
	}

	for docType, count := range stats.DocumentTypeCounts {
		response.DocumentTypeCounts[string(docType)] = count
	}

	return response
}
