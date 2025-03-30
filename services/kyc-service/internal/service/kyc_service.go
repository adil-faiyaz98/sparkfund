package service

import (
	"context"
	"fmt"
	"time"
	"math/rand"

	"kyc-service/internal/models"
	"kyc-service/internal/repository"
)

// KYCService handles business logic for KYC operations
type KYCService struct {
	repo *repository.KYCRepository
}

// NewKYCService creates a new KYC service instance
func NewKYCService() *KYCService {
	return &KYCService{
		repo: repository.NewKYCRepository(),
	}
}

// Placeholder function for fraud detection
func (s *KYCService) checkFraud(kyc *models.KYC) (float64, error) {
	// TODO: Replace with actual AI model integration
	// This is a placeholder that returns a random fraud score
	return rand.Float64(), nil
}

// Placeholder function for sanctions list screening
func (s *KYCService) checkSanctions(kyc *models.KYC) (bool, error) {
	// TODO: Replace with actual sanctions list screening API integration
	// This is a placeholder that always returns false (no match)
	return false, nil

}

// CreateKYC creates a new KYC record
func (s *KYCService) CreateKYC(ctx context.Context, kyc *models.KYC) error {
	// Validate required fields
	if err := s.validateKYC(kyc); err != nil {
		return err
	}

	// Check if KYC record already exists for the user
	existing, err := s.repo.GetByUserID(ctx, kyc.UserID)
	if err != nil {
		return fmt.Errorf("failed to check existing KYC record: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("KYC record already exists for user %s", kyc.UserID)
	}

	// AI-powered checks for fraud and sanctions
	fraudScore, err := s.checkFraud(kyc)
	if err != nil {
		return fmt.Errorf("failed to check fraud: %w", err)
	}

	sanctionsMatch, err := s.checkSanctions(kyc)
	if err != nil {
		return fmt.Errorf("failed to check sanctions: %w", err)
	}

	// Determine initial status based on AI checks
	if fraudScore > 0.7 || sanctionsMatch {
		kyc.Status = "reviewing"
	} else {
		kyc.Status = "pending"
	}

	return s.repo.Create(ctx, kyc)
}

// GetKYC retrieves a KYC record by user ID
func (s *KYCService) GetKYC(ctx context.Context, userID string) (*models.KYC, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// UpdateKYC updates an existing KYC record
func (s *KYCService) UpdateKYC(ctx context.Context, kyc *models.KYC) error {
	// Validate required fields
	if err := s.validateKYC(kyc); err != nil {
		return err
	}

	// Check if KYC record exists
	existing, err := s.repo.GetByUserID(ctx, kyc.UserID)
	if err != nil {
		return fmt.Errorf("failed to check existing KYC record: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("KYC record not found for user %s", kyc.UserID)
	}

	return s.repo.Update(ctx, kyc)
}

// UpdateKYCStatus updates the status of a KYC record
func (s *KYCService) UpdateKYCStatus(ctx context.Context, userID, status, rejectionReason string, reviewerID string) error {
	// Validate status
	if !isValidStatus(status) {
		return fmt.Errorf("invalid status: %s", status)
	}

	// Check if KYC record exists
	existing, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check existing KYC record: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("KYC record not found for user %s", userID)
	}

	// Update status and reviewer information
	updates := map[string]interface{}{
		"status":       status,
		"reviewed_by":  reviewerID,
		"reviewed_at":  time.Now(),
	}
	if rejectionReason != "" {
		updates["rejection_reason"] = rejectionReason
	}

	return s.repo.UpdateStatus(ctx, userID, status, rejectionReason)
}

// ListKYC retrieves a list of KYC records with optional filtering
func (s *KYCService) ListKYC(ctx context.Context, status string, page, pageSize int) ([]models.KYC, int64, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Validate status if provided
	if status != "" && !isValidStatus(status) {
		return nil, 0, fmt.Errorf("invalid status: %s", status)
	}

	return s.repo.List(ctx, status, page, pageSize)
}

// DeleteKYC deletes a KYC record
func (s *KYCService) DeleteKYC(ctx context.Context, userID string) error {
	// Check if KYC record exists
	existing, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check existing KYC record: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("KYC record not found for user %s", userID)
	}

	return s.repo.Delete(ctx, userID)
}

// validateKYC validates the required fields of a KYC record
func (s *KYCService) validateKYC(kyc *models.KYC) error {
	if kyc.UserID == "" {
		return fmt.Errorf("user ID is required")
	}
	if kyc.DocumentType == "" {
		return fmt.Errorf("document type is required")
	}
	if kyc.DocumentNumber == "" {
		return fmt.Errorf("document number is required")
	}
	if kyc.DocumentURL == "" {
		return fmt.Errorf("document URL is required")
	}
	if kyc.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if kyc.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	if kyc.DateOfBirth.IsZero() {
		return fmt.Errorf("date of birth is required")
	}
	if kyc.Address == "" {
		return fmt.Errorf("address is required")
	}
	if kyc.City == "" {
		return fmt.Errorf("city is required")
	}
	if kyc.Country == "" {
		return fmt.Errorf("country is required")
	}
	if kyc.PostalCode == "" {
		return fmt.Errorf("postal code is required")
	}
	if kyc.PhoneNumber == "" {
		return fmt.Errorf("phone number is required")
	}
	if kyc.Email == "" {
		return fmt.Errorf("email is required")
	}
	return nil
}

// isValidStatus checks if a given status is valid
func isValidStatus(status string) bool {
	validStatuses := map[string]bool{
		"pending":   true,
		"approved":  true,
		"rejected":  true,
		"reviewing": true,
	}
	return validStatuses[status]
}
