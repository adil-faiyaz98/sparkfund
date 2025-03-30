package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sparkfund/services/kyc-service/internal/model"
)

// VerificationRepository handles database operations for verification details
type VerificationRepository struct {
	db *gorm.DB
}

// NewVerificationRepository creates a new verification repository
func NewVerificationRepository(db *gorm.DB) *VerificationRepository {
	return &VerificationRepository{db: db}
}

// Create creates new verification details
func (r *VerificationRepository) Create(verification *model.Verification) error {
	return r.db.Create(verification).Error
}

// GetByID retrieves verification details by ID
func (r *VerificationRepository) GetByID(id uuid.UUID) (*model.Verification, error) {
	var verification model.Verification
	err := r.db.First(&verification, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// GetByDocumentID retrieves verification details by document ID
func (r *VerificationRepository) GetByDocumentID(documentID uuid.UUID) ([]*model.Verification, error) {
	var verifications []*model.Verification
	err := r.db.Where("document_id = ?", documentID).Find(&verifications).Error
	if err != nil {
		return nil, err
	}
	return verifications, nil
}

// Update updates verification details
func (r *VerificationRepository) Update(verification *model.Verification) error {
	return r.db.Save(verification).Error
}

// Delete deletes a verification
func (r *VerificationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Verification{}, "id = ?", id).Error
}

// GetHistory retrieves verification history
func (r *VerificationRepository) GetHistory(verificationID uuid.UUID) ([]*model.VerificationHistory, error) {
	var history []*model.VerificationHistory
	err := r.db.Where("verification_id = ?", verificationID).Order("created_at DESC").Find(&history).Error
	if err != nil {
		return nil, err
	}
	return history, nil
}

// AddHistoryEntry adds a new history entry
func (r *VerificationRepository) AddHistoryEntry(entry *model.VerificationHistory) error {
	return r.db.Create(entry).Error
}

// GetStats retrieves verification statistics
func (r *VerificationRepository) GetStats() (*model.VerificationStats, error) {
	var stats model.VerificationStats

	// Get total count
	err := r.db.Model(&model.Verification{}).Count(&stats.TotalCount).Error
	if err != nil {
		return nil, err
	}

	// Get counts by status
	err = r.db.Model(&model.Verification{}).Where("status = ?", model.VerificationStatusPending).Count(&stats.PendingCount).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&model.Verification{}).Where("status = ?", model.VerificationStatusCompleted).Count(&stats.CompletedCount).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&model.Verification{}).Where("status = ?", model.VerificationStatusFailed).Count(&stats.FailedCount).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&model.Verification{}).Where("status = ?", model.VerificationStatusExpired).Count(&stats.ExpiredCount).Error
	if err != nil {
		return nil, err
	}

	// Calculate average confidence score
	err = r.db.Model(&model.Verification{}).Where("status = ?", model.VerificationStatusCompleted).Select("AVG(confidence_score)").Row().Scan(&stats.AverageConfidence)
	if err != nil {
		return nil, err
	}

	// Calculate completion rate
	if stats.TotalCount > 0 {
		stats.CompletionRate = float64(stats.CompletedCount) / float64(stats.TotalCount) * 100
	}

	// Calculate average processing time
	var avgTime time.Duration
	err = r.db.Model(&model.Verification{}).
		Where("status = ? AND completed_at IS NOT NULL", model.VerificationStatusCompleted).
		Select("AVG(EXTRACT(EPOCH FROM (completed_at - created_at)) * interval '1 second')").
		Row().Scan(&avgTime)
	if err != nil {
		return nil, err
	}
	stats.AverageProcessingTime = avgTime

	return &stats, nil
}

// GetSummary retrieves verification summary
func (r *VerificationRepository) GetSummary(documentID uuid.UUID) (*model.VerificationSummary, error) {
	var verification model.Verification
	err := r.db.Where("document_id = ?", documentID).Order("created_at DESC").First(&verification).Error
	if err != nil {
		return nil, err
	}

	summary := &model.VerificationSummary{
		ID:              verification.ID,
		DocumentID:      verification.DocumentID,
		Status:          verification.Status,
		Method:          verification.Method,
		ConfidenceScore: verification.ConfidenceScore,
		CreatedAt:       verification.CreatedAt,
		CompletedAt:     verification.CompletedAt,
	}

	if verification.CompletedAt != nil {
		summary.ProcessingTime = verification.CompletedAt.Sub(verification.CreatedAt)
	}

	return summary, nil
}

// GetExpired retrieves expired verifications
func (r *VerificationRepository) GetExpired() ([]*model.Verification, error) {
	var verifications []*model.Verification
	err := r.db.Where("expires_at <= ? AND status != ?", time.Now(), model.VerificationStatusExpired).Find(&verifications).Error
	if err != nil {
		return nil, err
	}
	return verifications, nil
}

// GetPending retrieves pending verifications
func (r *VerificationRepository) GetPending() ([]*model.Verification, error) {
	var verifications []*model.Verification
	err := r.db.Where("status = ?", model.VerificationStatusPending).Find(&verifications).Error
	if err != nil {
		return nil, err
	}
	return verifications, nil
}

// GetByVerifier retrieves all verifications done by a specific verifier
func (r *VerificationRepository) GetByVerifier(verifierID uuid.UUID, page, pageSize int) ([]*model.Verification, int64, error) {
	var verifications []*model.Verification
	var total int64

	query := r.db.Model(&model.Verification{}).Where("verifier_id = ?", verifierID)

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&verifications).Error
	if err != nil {
		return nil, 0, err
	}

	return verifications, total, nil
}

// GetByMethod retrieves verifications by method
func (r *VerificationRepository) GetByMethod(method model.VerificationMethod, page, pageSize int) ([]*model.Verification, int64, error) {
	var verifications []*model.Verification
	var total int64

	query := r.db.Model(&model.Verification{}).Where("method = ?", method)

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&verifications).Error
	if err != nil {
		return nil, 0, err
	}

	return verifications, total, nil
}

// GetByDateRange retrieves verifications within a date range
func (r *VerificationRepository) GetByDateRange(startDate, endDate time.Time, page, pageSize int) ([]*model.Verification, int64, error) {
	var verifications []*model.Verification
	var total int64

	query := r.db.Model(&model.Verification{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate)

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&verifications).Error
	if err != nil {
		return nil, 0, err
	}

	return verifications, total, nil
}

// GetFailed retrieves all failed verifications
func (r *VerificationRepository) GetFailed() ([]*model.Verification, error) {
	var verifications []*model.Verification
	err := r.db.Where("status = ?", model.VerificationStatusFailed).Find(&verifications).Error
	if err != nil {
		return nil, err
	}
	return verifications, nil
}

// GetVerificationStats retrieves verification statistics
func (r *VerificationRepository) GetVerificationStats(ctx context.Context) (map[string]interface{}, error) {
	var stats struct {
		TotalVerifications     int64   `json:"total_verifications"`
		AverageConfidenceScore float64 `json:"average_confidence_score"`
		RejectionRate          float64 `json:"rejection_rate"`
		VerificationsByMethod  map[string]int64
	}

	// Get total verifications
	if err := r.db.WithContext(ctx).Model(&model.VerificationDetails{}).Count(&stats.TotalVerifications).Error; err != nil {
		return nil, err
	}

	// Get average confidence score
	if err := r.db.WithContext(ctx).Model(&model.VerificationDetails{}).
		Select("AVG(confidence_score) as average_confidence_score").
		Scan(&stats.AverageConfidenceScore).Error; err != nil {
		return nil, err
	}

	// Get rejection rate
	var rejectedCount int64
	if err := r.db.WithContext(ctx).Model(&model.VerificationDetails{}).
		Where("rejection_reason IS NOT NULL").
		Count(&rejectedCount).Error; err != nil {
		return nil, err
	}
	if stats.TotalVerifications > 0 {
		stats.RejectionRate = float64(rejectedCount) / float64(stats.TotalVerifications) * 100
	}

	// Get verifications by method
	var methodCounts []struct {
		Method string
		Count  int64
	}
	if err := r.db.WithContext(ctx).Model(&model.VerificationDetails{}).
		Select("verification_method as method, COUNT(*) as count").
		Group("verification_method").
		Scan(&methodCounts).Error; err != nil {
		return nil, err
	}

	stats.VerificationsByMethod = make(map[string]int64)
	for _, mc := range methodCounts {
		stats.VerificationsByMethod[mc.Method] = mc.Count
	}

	return map[string]interface{}{
		"total_verifications":      stats.TotalVerifications,
		"average_confidence_score": stats.AverageConfidenceScore,
		"rejection_rate":           stats.RejectionRate,
		"verifications_by_method":  stats.VerificationsByMethod,
	}, nil
}
