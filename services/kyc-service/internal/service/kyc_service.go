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

// KYCService handles business logic for KYC operations
type KYCService struct {
	kycRepo         *repository.KYCRepository
	docRepo         *repository.DocumentRepository
	verRepo         *repository.VerificationRepository
	eventPublisher  EventPublisher
}

// EventPublisher defines the interface for publishing events
type EventPublisher interface {
	Publish(ctx context.Context, eventType string, payload interface{}) error
}

// NewKYCService creates a new KYC service
func NewKYCService(kycRepo *repository.KYCRepository, docRepo *repository.DocumentRepository, verRepo *repository.VerificationRepository, eventPublisher EventPublisher) *KYCService {
	return &KYCService{
		kycRepo:        kycRepo,
		docRepo:        docRepo,
		verRepo:        verRepo,
		eventPublisher: eventPublisher,
	}
}

// CreateKYC creates a new KYC verification
func (s *KYCService) CreateKYC(ctx context.Context, userID uuid.UUID, request *model.KYCRequest) (*domain.EnhancedKYC, error) {
	// Create KYC record
	kyc := &model.KYC{
		ID:                uuid.New(),
		UserID:            userID,
		FirstName:         request.FirstName,
		LastName:          request.LastName,
		DateOfBirth:       request.DateOfBirth,
		Nationality:       request.Nationality,
		Email:             request.Email,
		PhoneNumber:       request.PhoneNumber,
		Address:           request.Address,
		City:              request.City,
		State:             request.State,
		Country:           request.Country,
		PostalCode:        request.PostalCode,
		DocumentType:      request.DocumentType,
		DocumentNumber:    request.DocumentNumber,
		DocumentFront:     request.DocumentFront,
		DocumentBack:      request.DocumentBack,
		SelfieImage:       request.SelfieImage,
		Status:            model.KYCStatusPending,
		RiskLevel:         model.RiskLevelMedium,
		RiskScore:         50.0,
		TransactionAmount: request.TransactionAmount,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Parse document expiry if provided
	if request.DocumentExpiry != "" {
		expiry, err := time.Parse("2006-01-02", request.DocumentExpiry)
		if err == nil {
			kyc.DocumentExpiry = &expiry
		}
	}

	// Save KYC record
	if err := s.kycRepo.Create(ctx, kyc); err != nil {
		return nil, fmt.Errorf("failed to create KYC record: %w", err)
	}

	// Convert to domain model
	domainKYC := mapper.KYCModelToDomain(kyc)

	// Publish event
	if err := s.eventPublisher.Publish(ctx, "kyc.created", domainKYC); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to publish kyc.created event: %v\n", err)
	}

	return domainKYC, nil
}

// GetKYC retrieves a KYC verification by ID
func (s *KYCService) GetKYC(ctx context.Context, id uuid.UUID) (*domain.EnhancedKYC, error) {
	kyc, err := s.kycRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get documents for this KYC
	documents, err := s.docRepo.GetByUserID(ctx, kyc.UserID)
	if err != nil {
		return nil, err
	}

	// Get verifications for this KYC
	verifications, err := s.verRepo.GetByKYCID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Convert to domain model
	domainKYC := mapper.KYCModelToDomain(kyc)

	// Add documents and verifications
	if documents != nil {
		domainKYC.Documents = mapper.DocumentModelsToDomains([]*model.Document{documents})
	}

	if verifications != nil {
		domainKYC.Verifications = mapper.VerificationModelsToDomains(verifications)
	}

	return domainKYC, nil
}

// GetKYCByUserID retrieves a KYC verification by user ID
func (s *KYCService) GetKYCByUserID(ctx context.Context, userID uuid.UUID) (*domain.EnhancedKYC, error) {
	kyc, err := s.kycRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get documents for this user
	documents, err := s.docRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get verifications for this KYC
	verifications, err := s.verRepo.GetByKYCID(ctx, kyc.ID)
	if err != nil {
		return nil, err
	}

	// Convert to domain model
	domainKYC := mapper.KYCModelToDomain(kyc)

	// Add documents and verifications
	if documents != nil {
		domainKYC.Documents = mapper.DocumentModelsToDomains([]*model.Document{documents})
	}

	if verifications != nil {
		domainKYC.Verifications = mapper.VerificationModelsToDomains(verifications)
	}

	return domainKYC, nil
}

// UpdateKYCStatus updates the status of a KYC verification
func (s *KYCService) UpdateKYCStatus(ctx context.Context, id uuid.UUID, status domain.KYCStatus, notes string, reviewerID uuid.UUID) error {
	// Get existing KYC
	kyc, err := s.kycRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Update status
	kyc.Status = model.KYCStatus(status)
	kyc.Notes = notes
	kyc.ReviewedBy = &reviewerID
	kyc.ReviewedAt = &time.Time{}
	*kyc.ReviewedAt = time.Now()
	kyc.UpdatedAt = time.Now()

	if status == domain.KYCStatusVerified || status == domain.KYCStatusApproved {
		now := time.Now()
		kyc.CompletedAt = &now
	}

	// Save KYC
	if err := s.kycRepo.Update(ctx, kyc); err != nil {
		return fmt.Errorf("failed to update KYC status: %w", err)
	}

	// Create review record
	review := &model.KYCReview{
		ID:         uuid.New(),
		KYCID:      id,
		ReviewerID: reviewerID,
		Status:     string(status),
		Reason:     notes,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.kycRepo.CreateReview(ctx, review); err != nil {
		return fmt.Errorf("failed to create KYC review: %w", err)
	}

	// Publish event
	domainKYC := mapper.KYCModelToDomain(kyc)
	if err := s.eventPublisher.Publish(ctx, "kyc.status_updated", domainKYC); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to publish kyc.status_updated event: %v\n", err)
	}

	return nil
}

// UpdateKYCRiskLevel updates the risk level of a KYC verification
func (s *KYCService) UpdateKYCRiskLevel(ctx context.Context, id uuid.UUID, riskLevel domain.RiskLevel, riskScore float64, notes string, reviewerID uuid.UUID) error {
	// Get existing KYC
	kyc, err := s.kycRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Update risk level
	kyc.RiskLevel = model.RiskLevel(riskLevel)
	kyc.RiskScore = riskScore
	kyc.Notes = notes
	kyc.ReviewedBy = &reviewerID
	kyc.ReviewedAt = &time.Time{}
	*kyc.ReviewedAt = time.Now()
	kyc.UpdatedAt = time.Now()

	// Save KYC
	if err := s.kycRepo.Update(ctx, kyc); err != nil {
		return fmt.Errorf("failed to update KYC risk level: %w", err)
	}

	// Create review record
	review := &model.KYCReview{
		ID:             uuid.New(),
		KYCID:          id,
		ReviewerID:     reviewerID,
		Status:         "RISK_UPDATED",
		RiskAssessment: fmt.Sprintf("Risk level updated to %s with score %.2f", riskLevel, riskScore),
		Notes:          notes,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.kycRepo.CreateReview(ctx, review); err != nil {
		return fmt.Errorf("failed to create KYC review: %w", err)
	}

	// Publish event
	domainKYC := mapper.KYCModelToDomain(kyc)
	if err := s.eventPublisher.Publish(ctx, "kyc.risk_updated", domainKYC); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to publish kyc.risk_updated event: %v\n", err)
	}

	return nil
}

// ListKYCs retrieves KYC verifications with pagination
func (s *KYCService) ListKYCs(ctx context.Context, page, pageSize int) ([]*domain.EnhancedKYC, int64, error) {
	kycs, total, err := s.kycRepo.List(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return mapper.KYCModelsToDomains(kycs), total, nil
}

// GetKYCsByStatus retrieves KYC verifications by status with pagination
func (s *KYCService) GetKYCsByStatus(ctx context.Context, status domain.KYCStatus, page, pageSize int) ([]*domain.EnhancedKYC, int64, error) {
	kycs, total, err := s.kycRepo.GetByStatus(ctx, model.KYCStatus(status), page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return mapper.KYCModelsToDomains(kycs), total, nil
}

// GetKYCsByRiskLevel retrieves KYC verifications by risk level with pagination
func (s *KYCService) GetKYCsByRiskLevel(ctx context.Context, riskLevel domain.RiskLevel, page, pageSize int) ([]*domain.EnhancedKYC, int64, error) {
	kycs, total, err := s.kycRepo.GetByRiskLevel(ctx, model.RiskLevel(riskLevel), page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return mapper.KYCModelsToDomains(kycs), total, nil
}

// ValidateKYC validates a KYC verification
func (s *KYCService) ValidateKYC(kyc *domain.EnhancedKYC) error {
	if kyc.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}

	if kyc.FirstName == "" {
		return errors.New("first name is required")
	}

	if kyc.LastName == "" {
		return errors.New("last name is required")
	}

	if kyc.DateOfBirth == "" {
		return errors.New("date of birth is required")
	}

	if kyc.Address == "" {
		return errors.New("address is required")
	}

	if kyc.City == "" {
		return errors.New("city is required")
	}

	if kyc.Country == "" {
		return errors.New("country is required")
	}

	if kyc.PostalCode == "" {
		return errors.New("postal code is required")
	}

	if kyc.DocumentType == "" {
		return errors.New("document type is required")
	}

	if kyc.DocumentNumber == "" {
		return errors.New("document number is required")
	}

	if kyc.DocumentFront == "" {
		return errors.New("document front is required")
	}

	if kyc.DocumentBack == "" {
		return errors.New("document back is required")
	}

	if kyc.SelfieImage == "" {
		return errors.New("selfie image is required")
	}

	return nil
}
