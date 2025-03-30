package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sparkfund/services/kyc-service/internal/models"
)

// ProfileRepository handles database operations for KYC profiles
type ProfileRepository struct {
	db *gorm.DB
}

// NewProfileRepository creates a new profile repository
func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// Create creates a new KYC profile
func (r *ProfileRepository) Create(ctx context.Context, profile *models.KYCProfile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

// GetByID retrieves a KYC profile by ID
func (r *ProfileRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.KYCProfile, error) {
	var profile models.KYCProfile
	err := r.db.WithContext(ctx).First(&profile, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

// GetByUserID retrieves a KYC profile by user ID
func (r *ProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.KYCProfile, error) {
	var profile models.KYCProfile
	err := r.db.WithContext(ctx).First(&profile, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

// Update updates a KYC profile
func (r *ProfileRepository) Update(ctx context.Context, profile *models.KYCProfile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

// UpdateStatus updates the status of a KYC profile
func (r *ProfileRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.KYCStatus) error {
	return r.db.WithContext(ctx).Model(&models.KYCProfile{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateRiskLevel updates the risk level of a KYC profile
func (r *ProfileRepository) UpdateRiskLevel(ctx context.Context, id uuid.UUID, riskLevel models.RiskLevel, riskScore float64) error {
	return r.db.WithContext(ctx).Model(&models.KYCProfile{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"risk_level": riskLevel,
			"risk_score": riskScore,
		}).Error
}

// Delete soft deletes a KYC profile
func (r *ProfileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.KYCProfile{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// GetByStatus retrieves KYC profiles by status
func (r *ProfileRepository) GetByStatus(ctx context.Context, status models.KYCStatus, page, pageSize int) ([]models.KYCProfile, int64, error) {
	var profiles []models.KYCProfile
	var total int64

	query := r.db.WithContext(ctx).Model(&models.KYCProfile{}).Where("status = ?", status)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&profiles).Error
	if err != nil {
		return nil, 0, err
	}

	return profiles, total, nil
}

// GetByRiskLevel retrieves KYC profiles by risk level
func (r *ProfileRepository) GetByRiskLevel(ctx context.Context, riskLevel models.RiskLevel, page, pageSize int) ([]models.KYCProfile, int64, error) {
	var profiles []models.KYCProfile
	var total int64

	query := r.db.WithContext(ctx).Model(&models.KYCProfile{}).Where("risk_level = ?", riskLevel)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&profiles).Error
	if err != nil {
		return nil, 0, err
	}

	return profiles, total, nil
}

// GetPendingProfiles retrieves all pending KYC profiles
func (r *ProfileRepository) GetPendingProfiles(ctx context.Context, page, pageSize int) ([]models.KYCProfile, int64, error) {
	return r.GetByStatus(ctx, models.KYCStatusPending, page, pageSize)
}

// GetInReviewProfiles retrieves all in-review KYC profiles
func (r *ProfileRepository) GetInReviewProfiles(ctx context.Context, page, pageSize int) ([]models.KYCProfile, int64, error) {
	return r.GetByStatus(ctx, models.KYCStatusInReview, page, pageSize)
}

// GetHighRiskProfiles retrieves all high-risk KYC profiles
func (r *ProfileRepository) GetHighRiskProfiles(ctx context.Context, page, pageSize int) ([]models.KYCProfile, int64, error) {
	return r.GetByRiskLevel(ctx, models.RiskLevelHigh, page, pageSize)
}

// GetProfileStats retrieves KYC profile statistics
func (r *ProfileRepository) GetProfileStats(ctx context.Context) (map[string]interface{}, error) {
	var stats struct {
		TotalProfiles    int64 `json:"total_profiles"`
		ProfilesByStatus map[string]int64
		ProfilesByRisk   map[string]int64
		AverageRiskScore float64 `json:"average_risk_score"`
		HighRiskRate     float64 `json:"high_risk_rate"`
	}

	// Get total profiles
	if err := r.db.WithContext(ctx).Model(&models.KYCProfile{}).Count(&stats.TotalProfiles).Error; err != nil {
		return nil, err
	}

	// Get profiles by status
	var statusCounts []struct {
		Status string
		Count  int64
	}
	if err := r.db.WithContext(ctx).Model(&models.KYCProfile{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		return nil, err
	}

	stats.ProfilesByStatus = make(map[string]int64)
	for _, sc := range statusCounts {
		stats.ProfilesByStatus[sc.Status] = sc.Count
	}

	// Get profiles by risk level
	var riskCounts []struct {
		RiskLevel string
		Count     int64
	}
	if err := r.db.WithContext(ctx).Model(&models.KYCProfile{}).
		Select("risk_level, COUNT(*) as count").
		Group("risk_level").
		Scan(&riskCounts).Error; err != nil {
		return nil, err
	}

	stats.ProfilesByRisk = make(map[string]int64)
	for _, rc := range riskCounts {
		stats.ProfilesByRisk[rc.RiskLevel] = rc.Count
	}

	// Get average risk score
	if err := r.db.WithContext(ctx).Model(&models.KYCProfile{}).
		Select("AVG(risk_score) as average_risk_score").
		Scan(&stats.AverageRiskScore).Error; err != nil {
		return nil, err
	}

	// Calculate high risk rate
	if stats.TotalProfiles > 0 {
		highRiskCount := stats.ProfilesByRisk[string(models.RiskLevelHigh)]
		stats.HighRiskRate = float64(highRiskCount) / float64(stats.TotalProfiles) * 100
	}

	return map[string]interface{}{
		"total_profiles":     stats.TotalProfiles,
		"profiles_by_status": stats.ProfilesByStatus,
		"profiles_by_risk":   stats.ProfilesByRisk,
		"average_risk_score": stats.AverageRiskScore,
		"high_risk_rate":     stats.HighRiskRate,
	}, nil
}
