package repository

import (
	"context"
	"errors"

	"kyc-service/internal/database"
	"kyc-service/internal/models"

	"gorm.io/gorm"
)

// KYCRepository handles database operations for KYC records
type KYCRepository struct {
	db *gorm.DB
}

// NewKYCRepository creates a new KYC repository instance
func NewKYCRepository() *KYCRepository {
	return &KYCRepository{
		db: database.DB,
	}
}

// Create creates a new KYC record
func (r *KYCRepository) Create(ctx context.Context, kyc *models.KYC) error {
	return r.db.WithContext(ctx).Create(kyc).Error
}

// GetByUserID retrieves a KYC record by user ID
func (r *KYCRepository) GetByUserID(ctx context.Context, userID string) (*models.KYC, error) {
	var kyc models.KYC
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&kyc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &kyc, nil
}

// Update updates an existing KYC record
func (r *KYCRepository) Update(ctx context.Context, kyc *models.KYC) error {
	return r.db.WithContext(ctx).Save(kyc).Error
}

// UpdateStatus updates the status of a KYC record
func (r *KYCRepository) UpdateStatus(ctx context.Context, userID, status string, rejectionReason string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if rejectionReason != "" {
		updates["rejection_reason"] = rejectionReason
	}
	return r.db.WithContext(ctx).Model(&models.KYC{}).Where("user_id = ?", userID).Updates(updates).Error
}

// List retrieves a list of KYC records with optional filtering
func (r *KYCRepository) List(ctx context.Context, status string, page, pageSize int) ([]models.KYC, int64, error) {
	var kycs []models.KYC
	var total int64

	query := r.db.WithContext(ctx).Model(&models.KYC{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&kycs).Error
	if err != nil {
		return nil, 0, err
	}

	return kycs, total, nil
}

// Delete deletes a KYC record
func (r *KYCRepository) Delete(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.KYC{}).Error
}
