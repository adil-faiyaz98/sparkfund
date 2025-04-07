package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sparkfund/services/kyc-service/internal/domain"
	"sparkfund/services/kyc-service/internal/mapper"
	"sparkfund/services/kyc-service/internal/model"
	"sparkfund/services/kyc-service/internal/repository"

	"github.com/google/uuid"
)

// VerificationService handles business logic for verification operations
type VerificationService struct {
	verRepo     *repository.VerificationRepository
	docRepo     *repository.DocumentRepository
	kycRepo     *repository.KYCRepository
}

// NewVerificationService creates a new verification service
func NewVerificationService(verRepo *repository.VerificationRepository, docRepo *repository.DocumentRepository, kycRepo *repository.KYCRepository) *VerificationService {
	return &VerificationService{
		verRepo: verRepo,
		docRepo: docRepo,
		kycRepo: kycRepo,
	}
}

// CreateVerification creates a new verification record
func (s *VerificationService) CreateVerification(ctx context.Context, documentID uuid.UUID, method domain.VerificationMethod) (*domain.EnhancedVerification, error) {
	// Check if document exists
	document, err := s.docRepo.GetByID(ctx, documentID)
	if err != nil {
		return nil, err
	}

	// Create verification record
	verification := &model.Verification{
		ID:              uuid.New(),
		DocumentID:      &documentID,
		Type:            model.VerificationType(domain.VerificationTypeDocument),
		Status:          model.VerificationStatusPending,
		Method:          model.VerificationMethod(method),
		ConfidenceScore: 0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Save verification
	if err := s.verRepo.Create(ctx, verification); err != nil {
		return nil, fmt.Errorf("failed to create verification: %w", err)
	}

	// Convert to domain model
	domainVerification := mapper.VerificationModelToDomain(verification)

	return domainVerification, nil
}

// GetVerification retrieves a verification by ID
func (s *VerificationService) GetVerification(ctx context.Context, id uuid.UUID) (*domain.EnhancedVerification, error) {
	verification, err := s.verRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapper.VerificationModelToDomain(verification), nil
}

// ListVerifications retrieves verifications with pagination
func (s *VerificationService) ListVerifications(ctx context.Context, page, pageSize int) ([]*domain.EnhancedVerification, int64, error) {
	verifications, total, err := s.verRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return mapper.VerificationModelsToDomains(verifications), total, nil
}

// UpdateVerificationStatus updates the status of a verification
func (s *VerificationService) UpdateVerificationStatus(ctx context.Context, id uuid.UUID, status domain.VerificationStatus, confidenceScore float64, notes string) error {
	// Get existing verification
	verification, err := s.verRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Update verification
	verification.Status = model.VerificationStatus(status)
	verification.ConfidenceScore = confidenceScore
	verification.Notes = notes
	verification.UpdatedAt = time.Now()

	if status == domain.VerStatusCompleted || status == domain.VerStatusApproved {
		now := time.Now()
		verification.CompletedAt = &now
	}

	// Save verification
	if err := s.verRepo.Update(ctx, verification); err != nil {
		return fmt.Errorf("failed to update verification: %w", err)
	}

	// If document verification, update document status
	if verification.DocumentID != nil && (status == domain.VerStatusApproved || status == domain.VerStatusRejected) {
		var docStatus model.DocumentStatus
		if status == domain.VerStatusApproved {
			docStatus = model.DocumentStatusVerified
		} else {
			docStatus = model.DocumentStatusRejected
		}

		if err := s.docRepo.UpdateStatus(ctx, *verification.DocumentID, docStatus, notes, id); err != nil {
			return fmt.Errorf("failed to update document status: %w", err)
		}
	}

	return nil
}

// CreateVerificationResult creates a result for a verification
func (s *VerificationService) CreateVerificationResult(ctx context.Context, verificationID uuid.UUID, result *domain.VerificationResult) error {
	// Check if verification exists
	verification, err := s.verRepo.GetByID(ctx, verificationID)
	if err != nil {
		return err
	}

	// Create model result
	modelResult := mapper.VerificationResultDomainToModel(result)
	modelResult.VerificationID = verificationID
	modelResult.ID = uuid.New()
	modelResult.CreatedAt = time.Now()
	modelResult.UpdatedAt = time.Now()

	// Save result
	if err := s.verRepo.CreateResult(ctx, modelResult); err != nil {
		return fmt.Errorf("failed to create verification result: %w", err)
	}

	// Update verification status based on result
	var status model.VerificationStatus
	if result.Success {
		status = model.VerificationStatusApproved
	} else {
		status = model.VerificationStatusRejected
	}

	verification.Status = status
	verification.ConfidenceScore = result.Score
	verification.UpdatedAt = time.Now()
	now := time.Now()
	verification.CompletedAt = &now

	// Save verification
	if err := s.verRepo.Update(ctx, verification); err != nil {
		return fmt.Errorf("failed to update verification: %w", err)
	}

	return nil
}

// GetVerificationsByDocument retrieves verifications for a document
func (s *VerificationService) GetVerificationsByDocument(ctx context.Context, documentID uuid.UUID) ([]*domain.EnhancedVerification, error) {
	verifications, err := s.verRepo.GetByDocumentID(ctx, documentID)
	if err != nil {
		return nil, err
	}

	return mapper.VerificationModelsToDomains(verifications), nil
}

// GetVerificationsByKYC retrieves verifications for a KYC record
func (s *VerificationService) GetVerificationsByKYC(ctx context.Context, kycID uuid.UUID) ([]*domain.EnhancedVerification, error) {
	verifications, err := s.verRepo.GetByKYCID(ctx, kycID)
	if err != nil {
		return nil, err
	}

	return mapper.VerificationModelsToDomains(verifications), nil
}

// ValidateVerification validates verification details
func (s *VerificationService) ValidateVerification(data *domain.EnhancedVerification) error {
	if data.ID == uuid.Nil {
		return errors.New("verification ID is required")
	}

	if data.DocumentID == nil && data.KYCID == nil {
		return errors.New("either document ID or KYC ID is required")
	}

	if data.Type == "" {
		return errors.New("verification type is required")
	}

	if data.Method == "" {
		return errors.New("verification method is required")
	}

	return nil
}
