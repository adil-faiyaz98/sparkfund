package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sparkfund/services/kyc-service/internal/model"
)

// KYCRepository handles database operations for KYC records
type KYCRepository struct {
	db *gorm.DB
}

// NewKYCRepository creates a new KYC repository
func NewKYCRepository(db *gorm.DB) *KYCRepository {
	return &KYCRepository{
		db: db,
	}
}

// Create creates a new KYC record
func (r *KYCRepository) Create(ctx context.Context, kyc *model.KYC) error {
	return r.db.WithContext(ctx).Create(kyc).Error
}

// GetByID retrieves a KYC record by ID
func (r *KYCRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.KYC, error) {
	var kyc model.KYC
	err := r.db.WithContext(ctx).First(&kyc, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &kyc, nil
}

// GetByUserID retrieves a KYC record by user ID
func (r *KYCRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.KYC, error) {
	var kyc model.KYC
	err := r.db.WithContext(ctx).First(&kyc, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &kyc, nil
}

// Update updates an existing KYC record
func (r *KYCRepository) Update(ctx context.Context, kyc *model.KYC) error {
	return r.db.WithContext(ctx).Save(kyc).Error
}

// Delete soft deletes a KYC record
func (r *KYCRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.KYC{}, "id = ?", id).Error
}

// List retrieves KYC records with pagination
func (r *KYCRepository) List(ctx context.Context, page, pageSize int) ([]*model.KYC, int64, error) {
	var kycs []*model.KYC
	var total int64

	query := r.db.WithContext(ctx).Model(&model.KYC{})

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&kycs).Error
	if err != nil {
		return nil, 0, err
	}

	return kycs, total, nil
}

// GetByStatus retrieves KYC records by status with pagination
func (r *KYCRepository) GetByStatus(ctx context.Context, status model.KYCStatus, page, pageSize int) ([]*model.KYC, int64, error) {
	var kycs []*model.KYC
	var total int64

	query := r.db.WithContext(ctx).Model(&model.KYC{}).Where("status = ?", status)

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&kycs).Error
	if err != nil {
		return nil, 0, err
	}

	return kycs, total, nil
}

// GetByRiskLevel retrieves KYC records by risk level with pagination
func (r *KYCRepository) GetByRiskLevel(ctx context.Context, riskLevel model.RiskLevel, page, pageSize int) ([]*model.KYC, int64, error) {
	var kycs []*model.KYC
	var total int64

	query := r.db.WithContext(ctx).Model(&model.KYC{}).Where("risk_level = ?", riskLevel)

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&kycs).Error
	if err != nil {
		return nil, 0, err
	}

	return kycs, total, nil
}

// CreateReview creates a new KYC review record
func (r *KYCRepository) CreateReview(ctx context.Context, review *model.KYCReview) error {
	return r.db.WithContext(ctx).Create(review).Error
}

// GetReviews retrieves KYC review records for a KYC record
func (r *KYCRepository) GetReviews(ctx context.Context, kycID uuid.UUID) ([]*model.KYCReview, error) {
	var reviews []*model.KYCReview
	err := r.db.WithContext(ctx).Where("kyc_id = ?", kycID).Order("created_at DESC").Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

// GetStats retrieves KYC statistics
func (r *KYCRepository) GetStats(ctx context.Context) (*model.KYCStats, error) {
	var stats model.KYCStats

	// Get total count
	err := r.db.WithContext(ctx).Model(&model.KYC{}).Count(&stats.TotalCount).Error
	if err != nil {
		return nil, err
	}

	// Get counts by status
	err = r.db.WithContext(ctx).Model(&model.KYC{}).Where("status = ?", model.KYCStatusPending).Count(&stats.PendingCount).Error
	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).Model(&model.KYC{}).Where("status = ?", model.KYCStatusVerified).Count(&stats.VerifiedCount).Error
	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).Model(&model.KYC{}).Where("status = ?", model.KYCStatusRejected).Count(&stats.RejectedCount).Error
	if err != nil {
		return nil, err
	}

	// Get counts by risk level
	err = r.db.WithContext(ctx).Model(&model.KYC{}).Where("risk_level = ?", model.RiskLevelLow).Count(&stats.LowRiskCount).Error
	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).Model(&model.KYC{}).Where("risk_level = ?", model.RiskLevelMedium).Count(&stats.MediumRiskCount).Error
	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).Model(&model.KYC{}).Where("risk_level = ?", model.RiskLevelHigh).Count(&stats.HighRiskCount).Error
	if err != nil {
		return nil, err
	}

	// Calculate average processing time for verified KYCs
	var avgProcessingTime float64
	err = r.db.WithContext(ctx).Model(&model.KYC{}).
		Where("status = ? AND verified_at IS NOT NULL", model.KYCStatusVerified).
		Select("AVG(EXTRACT(EPOCH FROM (verified_at - created_at)))").
		Row().Scan(&avgProcessingTime)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	stats.AverageProcessingTime = time.Duration(avgProcessingTime) * time.Second

	// Calculate verification rate
	if stats.TotalCount > 0 {
		stats.VerificationRate = float64(stats.VerifiedCount) / float64(stats.TotalCount) * 100
	}

	return &stats, nil
}
