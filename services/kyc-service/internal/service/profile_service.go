package service

import (
	"context"
	"fmt"
	"time"

	"sparkfund/services/kyc-service/internal/models"
	"sparkfund/services/kyc-service/internal/repository"

	"github.com/google/uuid"
)

// ProfileService handles business logic for KYC profile operations
type ProfileService struct {
	profileRepo *repository.ProfileRepository
	docRepo     *repository.DocumentRepository
}

// NewProfileService creates a new profile service
func NewProfileService(profileRepo *repository.ProfileRepository, docRepo *repository.DocumentRepository) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
		docRepo:     docRepo,
	}
}

// CreateProfile creates a new KYC profile
func (s *ProfileService) CreateProfile(ctx context.Context, userID uuid.UUID, data *models.KYCProfile) (*models.KYCProfile, error) {
	// Set initial status
	data.Status = models.ProfileStatusPending
	data.RiskLevel = models.RiskLevelMedium // Default risk level
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	// Create profile
	if err := s.profileRepo.Create(ctx, data); err != nil {
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	return data, nil
}

// GetProfile retrieves a KYC profile by user ID
func (s *ProfileService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.KYCProfile, error) {
	return s.profileRepo.GetByUserID(ctx, userID)
}

// UpdateProfile updates an existing KYC profile
func (s *ProfileService) UpdateProfile(ctx context.Context, userID uuid.UUID, data *models.KYCProfile) error {
	// Get existing profile
	existing, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get existing profile: %w", err)
	}

	// Update fields
	existing.FirstName = data.FirstName
	existing.LastName = data.LastName
	existing.DateOfBirth = data.DateOfBirth
	existing.Nationality = data.Nationality
	existing.Address = data.Address
	existing.PhoneNumber = data.PhoneNumber
	existing.Email = data.Email
	existing.TaxID = data.TaxID
	existing.Occupation = data.Occupation
	existing.SourceOfFunds = data.SourceOfFunds
	existing.ExpectedTransactionVolume = data.ExpectedTransactionVolume
	existing.UpdatedAt = time.Now()

	// Update profile
	if err := s.profileRepo.Update(ctx, existing); err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

// UpdateProfileStatus updates the status of a KYC profile
func (s *ProfileService) UpdateProfileStatus(ctx context.Context, userID uuid.UUID, status models.ProfileStatus, notes string) error {
	// Get existing profile
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get profile: %w", err)
	}

	// Update status
	profile.Status = status
	profile.Notes = &notes
	profile.UpdatedAt = time.Now()

	// Update profile
	if err := s.profileRepo.Update(ctx, profile); err != nil {
		return fmt.Errorf("failed to update profile status: %w", err)
	}

	return nil
}

// UpdateRiskLevel updates the risk level of a KYC profile
func (s *ProfileService) UpdateRiskLevel(ctx context.Context, userID uuid.UUID, riskLevel models.RiskLevel, notes string) error {
	// Get existing profile
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get profile: %w", err)
	}

	// Update risk level
	profile.RiskLevel = riskLevel
	profile.RiskNotes = &notes
	profile.UpdatedAt = time.Now()

	// Update profile
	if err := s.profileRepo.Update(ctx, profile); err != nil {
		return fmt.Errorf("failed to update risk level: %w", err)
	}

	return nil
}

// ListProfiles retrieves KYC profiles with pagination and filtering
func (s *ProfileService) ListProfiles(ctx context.Context, status *models.ProfileStatus, riskLevel *models.RiskLevel, page, pageSize int) ([]models.KYCProfile, int64, error) {
	return s.profileRepo.List(ctx, status, riskLevel, page, pageSize)
}

// GetProfileStats retrieves statistics about KYC profiles
func (s *ProfileService) GetProfileStats(ctx context.Context) (*models.ProfileStats, error) {
	return s.profileRepo.GetStats(ctx)
}

// DeleteProfile soft deletes a KYC profile
func (s *ProfileService) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	return s.profileRepo.Delete(ctx, userID)
}

// ValidateProfile validates a KYC profile
func (s *ProfileService) ValidateProfile(profile *models.KYCProfile) error {
	if profile.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if profile.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	if profile.DateOfBirth.IsZero() {
		return fmt.Errorf("date of birth is required")
	}
	if profile.Nationality == "" {
		return fmt.Errorf("nationality is required")
	}
	if profile.Address == "" {
		return fmt.Errorf("address is required")
	}
	if profile.PhoneNumber == "" {
		return fmt.Errorf("phone number is required")
	}
	if profile.Email == "" {
		return fmt.Errorf("email is required")
	}
	if profile.TaxID == "" {
		return fmt.Errorf("tax ID is required")
	}
	if profile.Occupation == "" {
		return fmt.Errorf("occupation is required")
	}
	if profile.SourceOfFunds == "" {
		return fmt.Errorf("source of funds is required")
	}
	if profile.ExpectedTransactionVolume <= 0 {
		return fmt.Errorf("expected transaction volume must be greater than 0")
	}

	return nil
}
