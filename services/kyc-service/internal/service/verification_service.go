package service

import (
	"fmt"
	"time"

	"sparkfund/services/kyc-service/internal/model"
	"sparkfund/services/kyc-service/internal/repository"

	"github.com/google/uuid"
)

// VerificationService handles business logic for document verification
type VerificationService struct {
	verificationRepo *repository.VerificationRepository
	documentRepo     *repository.DocumentRepository
}

// NewVerificationService creates a new verification service
func NewVerificationService(verificationRepo *repository.VerificationRepository, documentRepo *repository.DocumentRepository) *VerificationService {
	return &VerificationService{
		verificationRepo: verificationRepo,
		documentRepo:     documentRepo,
	}
}

// CreateVerification creates a new verification record
func (s *VerificationService) CreateVerification(documentID uuid.UUID, method model.VerificationMethod) (*model.Verification, error) {
	// Check if document exists
	document, err := s.documentRepo.GetByID(documentID)
	if err != nil {
		return nil, err
	}

	// Create verification record
	verification := &model.Verification{
		DocumentID: documentID,
		Status:     model.VerificationStatusPending,
		Method:     method,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ExpiresAt:  &time.Time{},
	}

	// Set expiration time based on document type
	switch document.Type {
	case model.DocumentTypeID:
		verification.ExpiresAt = &time.Time{}
		*verification.ExpiresAt = time.Now().Add(24 * time.Hour)
	case model.DocumentTypeProofOfAddress:
		verification.ExpiresAt = &time.Time{}
		*verification.ExpiresAt = time.Now().Add(48 * time.Hour)
	default:
		verification.ExpiresAt = &time.Time{}
		*verification.ExpiresAt = time.Now().Add(72 * time.Hour)
	}

	err = s.verificationRepo.Create(verification)
	if err != nil {
		return nil, err
	}

	return verification, nil
}

// GetVerification retrieves verification details by ID
func (s *VerificationService) GetVerification(id uuid.UUID) (*model.Verification, error) {
	return s.verificationRepo.GetByID(id)
}

// GetDocumentVerifications retrieves verification records for a document
func (s *VerificationService) GetDocumentVerifications(documentID uuid.UUID) ([]*model.Verification, error) {
	return s.verificationRepo.GetByDocumentID(documentID)
}

// UpdateVerificationStatus updates an existing verification record
func (s *VerificationService) UpdateVerificationStatus(id uuid.UUID, status model.VerificationStatus, verifierID uuid.UUID, notes string) error {
	verification, err := s.verificationRepo.GetByID(id)
	if err != nil {
		return err
	}

	verification.Status = status
	verification.VerifierID = verifierID
	verification.Notes = notes
	verification.UpdatedAt = time.Now()

	if status == model.VerificationStatusCompleted {
		now := time.Now()
		verification.CompletedAt = &now
	}

	err = s.verificationRepo.Update(verification)
	if err != nil {
		return err
	}

	// Add history entry
	historyEntry := &model.VerificationHistory{
		VerificationID: id,
		Status:         status,
		Notes:          notes,
		CreatedBy:      verifierID,
		CreatedAt:      time.Now(),
	}

	return s.verificationRepo.AddHistoryEntry(historyEntry)
}

// DeleteVerification soft deletes a verification record
func (s *VerificationService) DeleteVerification(id uuid.UUID) error {
	return s.verificationRepo.Delete(id)
}

// GetVerificationHistory retrieves the complete verification history for a document
func (s *VerificationService) GetVerificationHistory(id uuid.UUID) ([]*model.VerificationHistory, error) {
	return s.verificationRepo.GetHistory(id)
}

// GetVerificationStats retrieves statistics about verifications
func (s *VerificationService) GetVerificationStats() (*model.VerificationStats, error) {
	return s.verificationRepo.GetStats()
}

// GetVerificationSummary retrieves a summary of verifications for a document
func (s *VerificationService) GetVerificationSummary(documentID uuid.UUID) (*model.VerificationSummary, error) {
	return s.verificationRepo.GetSummary(documentID)
}

// GetExpiredVerifications retrieves expired verifications
func (s *VerificationService) GetExpiredVerifications() ([]*model.Verification, error) {
	return s.verificationRepo.GetExpired()
}

// GetPendingVerifications retrieves pending verifications
func (s *VerificationService) GetPendingVerifications() ([]*model.Verification, error) {
	return s.verificationRepo.GetPending()
}

// GetFailedVerifications retrieves failed verifications
func (s *VerificationService) GetFailedVerifications() ([]*model.Verification, error) {
	return s.verificationRepo.GetFailed()
}

// GetVerificationsByVerifier retrieves verifications by verifier
func (s *VerificationService) GetVerificationsByVerifier(verifierID uuid.UUID, page, pageSize int) ([]*model.Verification, int64, error) {
	return s.verificationRepo.GetByVerifier(verifierID, page, pageSize)
}

// GetVerificationsByMethod retrieves verifications by method
func (s *VerificationService) GetVerificationsByMethod(method model.VerificationMethod, page, pageSize int) ([]*model.Verification, int64, error) {
	return s.verificationRepo.GetByMethod(method, page, pageSize)
}

// GetVerificationsByDateRange retrieves verifications by date range
func (s *VerificationService) GetVerificationsByDateRange(startDate, endDate time.Time, page, pageSize int) ([]*model.Verification, int64, error) {
	return s.verificationRepo.GetByDateRange(startDate, endDate, page, pageSize)
}

// UpdateConfidenceScore updates the confidence score of a verification
func (s *VerificationService) UpdateConfidenceScore(id uuid.UUID, score float64) error {
	verification, err := s.verificationRepo.GetByID(id)
	if err != nil {
		return err
	}

	verification.ConfidenceScore = score
	verification.UpdatedAt = time.Now()

	return s.verificationRepo.Update(verification)
}

// UpdateVerificationMetadata updates the metadata of a verification
func (s *VerificationService) UpdateVerificationMetadata(id uuid.UUID, metadata map[string]interface{}) error {
	verification, err := s.verificationRepo.GetByID(id)
	if err != nil {
		return err
	}

	verification.Metadata = metadata
	verification.UpdatedAt = time.Now()

	return s.verificationRepo.Update(verification)
}

// ValidateVerification validates verification details
func (s *VerificationService) ValidateVerification(data *model.VerificationDetails) error {
	if data.DocumentID == uuid.Nil {
		return fmt.Errorf("document ID is required")
	}
	if data.VerifiedBy == uuid.Nil {
		return fmt.Errorf("verifier ID is required")
	}
	if data.VerificationMethod == "" {
		return fmt.Errorf("verification method is required")
	}
	if data.ConfidenceScore < 0 || data.ConfidenceScore > 1 {
		return fmt.Errorf("confidence score must be between 0 and 1")
	}

	return nil
}
