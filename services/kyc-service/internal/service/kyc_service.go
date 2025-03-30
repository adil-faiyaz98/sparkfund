package service

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"

	"sparkfund/services/kyc-service/internal/ai"
	"sparkfund/services/kyc-service/internal/models"
	"sparkfund/services/kyc-service/internal/repository"
)

// KYCService handles business logic for KYC operations
type KYCService struct {
	kycRepo    *repository.KYCRepository
	docRepo    *repository.DocumentRepository
	fraudModel ai.FraudModel
}

// NewKYCService creates a new KYC service instance
func NewKYCService(
	kycRepo *repository.KYCRepository,
	docRepo *repository.DocumentRepository,
	fraudModel ai.FraudModel,
) *KYCService {
	return &KYCService{
		kycRepo:    kycRepo,
		docRepo:    docRepo,
		fraudModel: fraudModel,
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

// CreateKYC creates a new KYC verification request
func (s *KYCService) CreateKYC(ctx context.Context, userID uuid.UUID, data *models.KYCVerification) (*models.KYCVerification, error) {
	// Get user's document
	doc, err := s.docRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user document: %w", err)
	}
	if doc == nil {
		return nil, fmt.Errorf("no document found for user %s", userID)
	}

	// Prepare fraud detection features
	features := ai.FraudFeatures{
		TransactionAmount:    data.TransactionAmount,
		TransactionTime:      time.Now(),
		TransactionFrequency: 1, // This should be calculated from user's history
		UserAge:              calculateAge(data.DateOfBirth),
		AccountAge:           0,   // This should be calculated from user's account creation date
		PreviousFraudReports: 0,   // This should be fetched from user's history
		DocumentQuality:      0.8, // This should be calculated from document verification
		DocumentAge:          calculateDocumentAge(doc.IssueDate),
		DocumentType:         string(doc.Type), // Convert DocumentType to string
		CountryRiskScore:     0.5,              // This should be fetched from a risk database
		IPRiskScore:          0.5,              // This should be calculated from IP analysis
		LoginAttempts:        1,                // This should be fetched from auth service
		FailedLogins:         0,                // This should be fetched from auth service
		LastLoginTime:        time.Now(),       // This should be fetched from auth service
		PEPStatus:            false,            // This should be fetched from PEP database
		SanctionStatus:       false,            // This should be fetched from sanctions database
		WatchlistStatus:      false,            // This should be fetched from watchlist database
	}

	// Get fraud prediction
	prediction, err := s.fraudModel.Predict(features)
	if err != nil {
		return nil, fmt.Errorf("failed to get fraud prediction: %w", err)
	}

	// Update KYC data with fraud detection results
	data.RiskLevel = string(prediction.RiskLevel)
	data.RiskScore = int(prediction.RiskScore)
	data.Notes = prediction.Explanation

	// Set initial status based on risk level
	switch prediction.RiskLevel {
	case ai.FraudRiskLow:
		data.Status = "PENDING"
	case ai.FraudRiskMedium:
		data.Status = "IN_REVIEW"
	case ai.FraudRiskHigh:
		data.Status = "FLAGGED"
	}

	// Create KYC record
	if err := s.kycRepo.Create(ctx, data); err != nil {
		return nil, fmt.Errorf("failed to create KYC record: %w", err)
	}

	return data, nil
}

// GetKYC retrieves a KYC verification by user ID
func (s *KYCService) GetKYC(ctx context.Context, userID uuid.UUID) (*models.KYCVerification, error) {
	return s.kycRepo.GetByUserID(ctx, userID)
}

// UpdateKYC updates an existing KYC verification
func (s *KYCService) UpdateKYC(ctx context.Context, userID uuid.UUID, data *models.KYCVerification) error {
	// Get existing KYC record
	existing, err := s.kycRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get existing KYC record: %w", err)
	}

	// Update fields
	existing.FirstName = data.FirstName
	existing.LastName = data.LastName
	existing.DateOfBirth = data.DateOfBirth
	existing.Address = data.Address
	existing.City = data.City
	existing.Country = data.Country
	existing.PostalCode = data.PostalCode
	existing.PhoneNumber = data.PhoneNumber
	existing.Email = data.Email

	// Update KYC record
	if err := s.kycRepo.Update(ctx, existing); err != nil {
		return fmt.Errorf("failed to update KYC record: %w", err)
	}

	return nil
}

// UpdateKYCStatus updates the status of a KYC verification
func (s *KYCService) UpdateKYCStatus(ctx context.Context, userID uuid.UUID, status string, notes string) error {
	return s.kycRepo.UpdateStatus(ctx, userID, status, notes)
}

// ListKYC retrieves KYC verifications with pagination and filtering
func (s *KYCService) ListKYC(ctx context.Context, status string, page, pageSize int) ([]models.KYCVerification, int64, error) {
	return s.kycRepo.List(ctx, status, page, pageSize)
}

// DeleteKYC soft deletes a KYC verification
func (s *KYCService) DeleteKYC(ctx context.Context, userID uuid.UUID) error {
	return s.kycRepo.Delete(ctx, userID)
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

// Helper functions
func calculateAge(dateOfBirth time.Time) int {
	now := time.Now()
	years := now.Year() - dateOfBirth.Year()
	if now.Month() < dateOfBirth.Month() || (now.Month() == dateOfBirth.Month() && now.Day() < dateOfBirth.Day()) {
		years--
	}
	return years
}

func calculateDocumentAge(issueDate time.Time) int {
	now := time.Now()
	days := int(now.Sub(issueDate).Hours() / 24)
	return days
}
