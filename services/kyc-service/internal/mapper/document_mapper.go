package mapper

import (
	"time"

	"github.com/google/uuid"

	"sparkfund/services/kyc-service/internal/domain"
	"sparkfund/services/kyc-service/internal/model"
)

// DocumentModelToDomain converts a model.Document to a domain.EnhancedDocument
func DocumentModelToDomain(doc *model.Document) *domain.EnhancedDocument {
	if doc == nil {
		return nil
	}

	var metadata domain.Metadata
	if doc.Metadata != nil {
		metadata = domain.Metadata(doc.Metadata)
	}

	return &domain.EnhancedDocument{
		ID:                doc.ID,
		UserID:            doc.UserID,
		KYCID:             doc.KYCVerificationID,
		Type:              domain.DocumentType(doc.Type),
		Status:            domain.DocumentStatus(doc.Status),
		FileName:          doc.FileName,
		FileSize:          doc.FileSize,
		MimeType:          doc.MimeType,
		FileHash:          doc.FileHash,
		FilePath:          doc.FilePath,
		FileURL:           doc.FileURL,
		DocumentNumber:    doc.DocumentNumber,
		IssueDate:         doc.IssueDate,
		ExpiryDate:        doc.ExpiryDate,
		IssuingCountry:    doc.IssuingCountry,
		IssuingAuthority:  doc.IssuingAuthority,
		VerificationID:    doc.VerificationID,
		ConfidenceScore:   doc.ConfidenceScore,
		IsValid:           doc.IsValid,
		Metadata:          metadata,
		CreatedAt:         doc.CreatedAt,
		UpdatedAt:         doc.UpdatedAt,
		ExpiresAt:         doc.ExpiresAt,
		VerifiedAt:        doc.VerifiedAt,
		RejectedAt:        doc.RejectedAt,
		RejectionReason:   doc.RejectionReason,
		RejectedBy:        doc.RejectedBy,
	}
}

// DocumentDomainToModel converts a domain.EnhancedDocument to a model.Document
func DocumentDomainToModel(doc *domain.EnhancedDocument) *model.Document {
	if doc == nil {
		return nil
	}

	var metadata map[string]interface{}
	if doc.Metadata != nil {
		metadata = map[string]interface{}(doc.Metadata)
	}

	return &model.Document{
		ID:                doc.ID,
		UserID:            doc.UserID,
		KYCVerificationID: doc.KYCID,
		Type:              model.DocumentType(doc.Type),
		Status:            model.DocumentStatus(doc.Status),
		FileName:          doc.FileName,
		FileSize:          doc.FileSize,
		MimeType:          doc.MimeType,
		FileHash:          doc.FileHash,
		FilePath:          doc.FilePath,
		FileURL:           doc.FileURL,
		DocumentNumber:    doc.DocumentNumber,
		IssueDate:         doc.IssueDate,
		ExpiryDate:        doc.ExpiryDate,
		IssuingCountry:    doc.IssuingCountry,
		IssuingAuthority:  doc.IssuingAuthority,
		VerificationID:    doc.VerificationID,
		ConfidenceScore:   doc.ConfidenceScore,
		IsValid:           doc.IsValid,
		Metadata:          metadata,
		CreatedAt:         doc.CreatedAt,
		UpdatedAt:         doc.UpdatedAt,
		ExpiresAt:         doc.ExpiresAt,
		VerifiedAt:        doc.VerifiedAt,
		RejectedAt:        doc.RejectedAt,
		RejectionReason:   doc.RejectionReason,
		RejectedBy:        doc.RejectedBy,
	}
}

// DocumentSummaryModelToDomain converts a model.DocumentSummary to a domain.DocumentSummary
func DocumentSummaryModelToDomain(summary *model.DocumentSummary) *domain.DocumentSummary {
	if summary == nil {
		return nil
	}

	return &domain.DocumentSummary{
		ID:              summary.ID,
		UserID:          summary.UserID,
		Type:            domain.DocumentType(summary.Type),
		Status:          domain.DocumentStatus(summary.Status),
		FileName:        summary.FileName,
		DocumentNumber:  summary.DocumentNumber,
		ExpiryDate:      summary.ExpiryDate,
		CreatedAt:       summary.CreatedAt,
		VerifiedAt:      summary.VerifiedAt,
		ProcessingTime:  summary.ProcessingTime,
		ConfidenceScore: summary.ConfidenceScore,
	}
}

// DocumentSummaryDomainToModel converts a domain.DocumentSummary to a model.DocumentSummary
func DocumentSummaryDomainToModel(summary *domain.DocumentSummary) *model.DocumentSummary {
	if summary == nil {
		return nil
	}

	return &model.DocumentSummary{
		ID:              summary.ID,
		UserID:          summary.UserID,
		Type:            model.DocumentType(summary.Type),
		Status:          model.DocumentStatus(summary.Status),
		FileName:        summary.FileName,
		DocumentNumber:  summary.DocumentNumber,
		ExpiryDate:      summary.ExpiryDate,
		CreatedAt:       summary.CreatedAt,
		VerifiedAt:      summary.VerifiedAt,
		ProcessingTime:  summary.ProcessingTime,
		ConfidenceScore: summary.ConfidenceScore,
	}
}

// DocumentModelsToDomains converts a slice of model.Document to a slice of domain.EnhancedDocument
func DocumentModelsToDomains(docs []*model.Document) []*domain.EnhancedDocument {
	result := make([]*domain.EnhancedDocument, len(docs))
	for i, doc := range docs {
		result[i] = DocumentModelToDomain(doc)
	}
	return result
}

// DocumentDomainsToModels converts a slice of domain.EnhancedDocument to a slice of model.Document
func DocumentDomainsToModels(docs []*domain.EnhancedDocument) []*model.Document {
	result := make([]*model.Document, len(docs))
	for i, doc := range docs {
		result[i] = DocumentDomainToModel(doc)
	}
	return result
}
