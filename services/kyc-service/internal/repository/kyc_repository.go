package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sparkfund/services/kyc-service/internal/models"
)

// KYCRepository handles database operations for KYC verifications
type KYCRepository struct {
	db *gorm.DB
}

// NewKYCRepository creates a new KYC repository
func NewKYCRepository(db *gorm.DB) *KYCRepository {
	return &KYCRepository{db: db}
}

// Create creates a new KYC verification
func (r *KYCRepository) Create(ctx context.Context, kyc *models.KYCVerification) error {
	return r.db.WithContext(ctx).Create(kyc).Error
}

// GetByUserID retrieves a KYC verification by user ID
func (r *KYCRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.KYCVerification, error) {
	var kyc models.KYCVerification
	err := r.db.WithContext(ctx).First(&kyc, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &kyc, nil
}

// Update updates an existing KYC verification
func (r *KYCRepository) Update(ctx context.Context, kyc *models.KYCVerification) error {
	return r.db.WithContext(ctx).Save(kyc).Error
}

// UpdateStatus updates the status of a KYC verification
func (r *KYCRepository) UpdateStatus(ctx context.Context, userID uuid.UUID, status string, notes string) error {
	return r.db.WithContext(ctx).Model(&models.KYCVerification{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"status":      status,
			"notes":       notes,
			"reviewed_at": time.Now(),
		}).Error
}

// List retrieves KYC verifications with pagination and filtering
func (r *KYCRepository) List(ctx context.Context, status string, page, pageSize int) ([]models.KYCVerification, int64, error) {
	var verifications []models.KYCVerification
	var total int64

	query := r.db.WithContext(ctx).Model(&models.KYCVerification{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&verifications).Error
	if err != nil {
		return nil, 0, err
	}

	return verifications, total, nil
}

// Delete soft deletes a KYC verification
func (r *KYCRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.KYCVerification{}, "user_id = ?", userID).Error
}
