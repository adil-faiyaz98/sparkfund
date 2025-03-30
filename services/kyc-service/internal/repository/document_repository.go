package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sparkfund/services/kyc-service/internal/model"
	"sparkfund/services/kyc-service/internal/models"
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
func (r *DocumentRepository) Create(ctx context.Context, doc *models.Document) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

// GetByID retrieves a document by ID
func (r *DocumentRepository) GetByID(id uuid.UUID) (*model.Document, error) {
	var document model.Document
	err := r.db.First(&document, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}

// GetByUserID retrieves a document by user ID
func (r *DocumentRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Document, error) {
	var doc models.Document
	err := r.db.WithContext(ctx).First(&doc, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

// Update updates an existing document
func (r *DocumentRepository) Update(ctx context.Context, doc *models.Document) error {
	return r.db.WithContext(ctx).Save(doc).Error
}

// Delete soft deletes a document
func (r *DocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Document{}, "id = ?", id).Error
}

// GetHistory retrieves the history of a document
func (r *DocumentRepository) GetHistory(documentID uuid.UUID) ([]*model.DocumentHistory, error) {
	var history []*model.DocumentHistory
	err := r.db.Where("document_id = ?", documentID).Order("created_at DESC").Find(&history).Error
	if err != nil {
		return nil, err
	}
	return history, nil
}

// AddHistoryEntry adds a new history entry to a document
func (r *DocumentRepository) AddHistoryEntry(entry *model.DocumentHistory) error {
	return r.db.Create(entry).Error
}

// GetStats retrieves statistics about documents
func (r *DocumentRepository) GetStats() (*model.DocumentStats, error) {
	var stats model.DocumentStats

	// Get total count
	err := r.db.Model(&model.Document{}).Count(&stats.TotalCount).Error
	if err != nil {
		return nil, err
	}

	// Get counts by status
	err = r.db.Model(&model.Document{}).Where("status = ?", model.DocumentStatusPending).Count(&stats.PendingCount).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&model.Document{}).Where("status = ?", model.DocumentStatusVerified).Count(&stats.VerifiedCount).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&model.Document{}).Where("status = ?", model.DocumentStatusRejected).Count(&stats.RejectedCount).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&model.Document{}).Where("status = ?", model.DocumentStatusExpired).Count(&stats.ExpiredCount).Error
	if err != nil {
		return nil, err
	}

	// Calculate average file size
	err = r.db.Model(&model.Document{}).Select("AVG(file_size)").Row().Scan(&stats.AverageFileSize)
	if err != nil {
		return nil, err
	}

	// Calculate total file size
	err = r.db.Model(&model.Document{}).Select("SUM(file_size)").Row().Scan(&stats.TotalFileSize)
	if err != nil {
		return nil, err
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
func (r *DocumentRepository) GetByType(documentType model.DocumentType, page, pageSize int) ([]*model.Document, int64, error) {
	var documents []*model.Document
	var total int64

	query := r.db.Model(&model.Document{}).Where("type = ?", documentType)

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
func (r *DocumentRepository) GetByDateRange(startDate, endDate time.Time, page, pageSize int) ([]*model.Document, int64, error) {
	var documents []*model.Document
	var total int64

	query := r.db.Model(&model.Document{}).
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
func (r *DocumentRepository) UpdateStatus(id uuid.UUID, status model.DocumentStatus, notes string, updatedBy uuid.UUID) error {
	document, err := r.GetByID(id)
	if err != nil {
		return err
	}

	document.Status = status
	document.UpdatedAt = time.Now()

	switch status {
	case model.DocumentStatusVerified:
		now := time.Now()
		document.VerifiedAt = &now
	case model.DocumentStatusRejected:
		now := time.Now()
		document.RejectedAt = &now
		document.RejectionReason = notes
	}

	err = r.Update(document)
	if err != nil {
		return err
	}

	// Add history entry
	historyEntry := &model.DocumentHistory{
		DocumentID: id,
		Status:     status,
		Notes:      notes,
		CreatedBy:  updatedBy,
		CreatedAt:  time.Now(),
	}

	return r.AddHistoryEntry(historyEntry)
}
