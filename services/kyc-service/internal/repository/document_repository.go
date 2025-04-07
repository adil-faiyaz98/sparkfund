package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sparkfund/services/kyc-service/internal/model"
)

// DocumentRepository handles database operations for documents
type DocumentRepository struct {
	db *gorm.DB
}

// NewDocumentRepository creates a new document repository
func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

// Create creates a new document
func (r *DocumentRepository) Create(ctx context.Context, doc *model.Document) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

// GetByID retrieves a document by ID
func (r *DocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Document, error) {
	var document model.Document
	err := r.db.WithContext(ctx).First(&document, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}

// GetByUserID retrieves a document by user ID
func (r *DocumentRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Document, error) {
	var doc model.Document
	err := r.db.WithContext(ctx).First(&doc, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

// GetByUserIDPaginated retrieves documents for a user with pagination
func (r *DocumentRepository) GetByUserIDPaginated(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*model.Document, int64, error) {
	var documents []*model.Document
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Document{}).Where("user_id = ?", userID)

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&documents).Error
	if err != nil {
		return nil, 0, err
	}

	return documents, total, nil
}

// Update updates an existing document
func (r *DocumentRepository) Update(ctx context.Context, doc *model.Document) error {
	return r.db.WithContext(ctx).Save(doc).Error
}

// Delete soft deletes a document
func (r *DocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Document{}, "id = ?", id).Error
}

// GetHistory retrieves the history of a document
func (r *DocumentRepository) GetHistory(ctx context.Context, documentID uuid.UUID) ([]*model.DocumentHistory, error) {
	var history []*model.DocumentHistory
	err := r.db.WithContext(ctx).Where("document_id = ?", documentID).Order("created_at DESC").Find(&history).Error
	if err != nil {
		return nil, err
	}
	return history, nil
}

// AddHistoryEntry adds a new history entry to a document
func (r *DocumentRepository) AddHistoryEntry(ctx context.Context, entry *model.DocumentHistory) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

// GetStats retrieves statistics about documents
func (r *DocumentRepository) GetStats(ctx context.Context) (*model.DocumentStats, error) {
	var stats model.DocumentStats

	// Get total count
	err := r.db.WithContext(ctx).Model(&model.Document{}).Count(&stats.TotalCount).Error
	if err != nil {
		return nil, err
	}

	// Get counts by status
	err = r.db.WithContext(ctx).Model(&model.Document{}).Where("status = ?", model.DocumentStatusPending).Count(&stats.PendingCount).Error
	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).Model(&model.Document{}).Where("status = ?", model.DocumentStatusVerified).Count(&stats.VerifiedCount).Error
	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).Model(&model.Document{}).Where("status = ?", model.DocumentStatusRejected).Count(&stats.RejectedCount).Error
	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).Model(&model.Document{}).Where("status = ?", model.DocumentStatusExpired).Count(&stats.ExpiredCount).Error
	if err != nil {
		return nil, err
	}

	// Calculate average file size
	err = r.db.WithContext(ctx).Model(&model.Document{}).Select("AVG(file_size)").Row().Scan(&stats.AverageFileSize)
	if err != nil {
		return nil, err
	}

	// Calculate total file size
	err = r.db.WithContext(ctx).Model(&model.Document{}).Select("SUM(file_size)").Row().Scan(&stats.TotalFileSize)
	if err != nil {
		return nil, err
	}

	// Calculate average processing time for verified documents
	var avgProcessingTime float64
	err = r.db.WithContext(ctx).Model(&model.Document{}).
		Where("status = ? AND verified_at IS NOT NULL", model.DocumentStatusVerified).
		Select("AVG(EXTRACT(EPOCH FROM (verified_at - created_at)))").
		Row().Scan(&avgProcessingTime)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	stats.ProcessingTimeAvg = time.Duration(avgProcessingTime) * time.Second

	// Calculate verification rate
	if stats.TotalCount > 0 {
		stats.VerificationRate = float64(stats.VerifiedCount) / float64(stats.TotalCount) * 100
	}

	// Get counts by document type
	stats.DocumentTypeCounts = make(map[model.DocumentType]int64)
	var typeCounts []struct {
		Type  model.DocumentType
		Count int64
	}
	err = r.db.Model(&model.Document{}).Select("type, COUNT(*) as count").Group("type").Find(&typeCounts).Error
	if err != nil {
		return nil, err
	}

	for _, tc := range typeCounts {
		stats.DocumentTypeCounts[tc.Type] = tc.Count
	}

	return &stats, nil
}

// GetSummary retrieves a summary of a document
func (r *DocumentRepository) GetSummary(documentID uuid.UUID) (*model.DocumentSummary, error) {
	var document model.Document
	err := r.db.First(&document, "id = ?", documentID).Error
	if err != nil {
		return nil, err
	}

	summary := &model.DocumentSummary{
		ID:         document.ID,
		Type:       document.Type,
		Status:     document.Status,
		FileName:   document.FileName,
		FileSize:   document.FileSize,
		CreatedAt:  document.CreatedAt,
		VerifiedAt: document.VerifiedAt,
	}

	if document.VerifiedAt != nil {
		summary.ProcessingTime = document.VerifiedAt.Sub(document.CreatedAt)
	}

	return summary, nil
}

// GetExpired retrieves all expired documents
func (r *DocumentRepository) GetExpired() ([]*model.Document, error) {
	var documents []*model.Document
	err := r.db.Where("expires_at <= ? AND status != ?", time.Now(), model.DocumentStatusExpired).Find(&documents).Error
	if err != nil {
		return nil, err
	}
	return documents, nil
}

// GetPending retrieves all pending documents
func (r *DocumentRepository) GetPending() ([]*model.Document, error) {
	var documents []*model.Document
	err := r.db.Where("status = ?", model.DocumentStatusPending).Find(&documents).Error
	if err != nil {
		return nil, err
	}
	return documents, nil
}

// GetRejected retrieves all rejected documents
func (r *DocumentRepository) GetRejected() ([]*model.Document, error) {
	var documents []*model.Document
	err := r.db.Where("status = ?", model.DocumentStatusRejected).Find(&documents).Error
	if err != nil {
		return nil, err
	}
	return documents, nil
}

// GetByType retrieves documents by type
func (r *DocumentRepository) GetByType(ctx context.Context, documentType model.DocumentType, page, pageSize int) ([]*model.Document, int64, error) {
	var documents []*model.Document
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Document{}).Where("type = ?", documentType)

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&documents).Error
	if err != nil {
		return nil, 0, err
	}

	return documents, total, nil
}

// GetByDateRange retrieves documents by date range
func (r *DocumentRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time, page, pageSize int) ([]*model.Document, int64, error) {
	var documents []*model.Document
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Document{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate)

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&documents).Error
	if err != nil {
		return nil, 0, err
	}

	return documents, total, nil
}

// UpdateStatus updates the status of a document
func (r *DocumentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.DocumentStatus, notes string, updatedBy uuid.UUID) error {
	document, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	document.Status = status
	document.UpdatedAt = time.Now()

	switch status {
	case model.DocumentStatusVerified:
		now := time.Now()
		document.VerifiedAt = &now
		document.RejectedAt = nil
		document.RejectionReason = ""
	case model.DocumentStatusRejected:
		now := time.Now()
		document.RejectedAt = &now
		document.RejectionReason = notes
	}

	err = r.Update(ctx, document)
	if err != nil {
		return err
	}

	// Add history entry
	historyEntry := &model.DocumentHistory{
		ID:         uuid.New(),
		DocumentID: id,
		Status:     status,
		Notes:      notes,
		CreatedBy:  updatedBy,
		CreatedAt:  time.Now(),
	}

	return r.AddHistoryEntry(ctx, historyEntry)
}
