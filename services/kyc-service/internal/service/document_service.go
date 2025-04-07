package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"sparkfund/services/kyc-service/internal/domain"
	"sparkfund/services/kyc-service/internal/mapper"
	"sparkfund/services/kyc-service/internal/model"
	"sparkfund/services/kyc-service/internal/repository"

	"github.com/google/uuid"
)

// DocumentService handles business logic for document operations
type DocumentService struct {
	docRepo   *repository.DocumentRepository
	verRepo   *repository.VerificationRepository
	uploadDir string
}

// NewDocumentService creates a new document service
func NewDocumentService(docRepo *repository.DocumentRepository, verRepo *repository.VerificationRepository, uploadDir string) *DocumentService {
	return &DocumentService{
		docRepo:   docRepo,
		verRepo:   verRepo,
		uploadDir: uploadDir,
	}
}

// UploadDocument handles document upload and processing
func (s *DocumentService) UploadDocument(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader, docType string, metadata map[string]interface{}) (*domain.EnhancedDocument, error) {
	// Validate file
	if err := s.validateFile(file); err != nil {
		return nil, fmt.Errorf("invalid file: %w", err)
	}

	// Read file data
	fileData, err := s.readFileData(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Calculate file hash
	fileHash := s.calculateFileHash(fileData)

	// Create document record
	doc := &model.Document{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      model.DocumentType(docType),
		Status:    model.DocumentStatusPending,
		FileName:  file.Filename,
		FileSize:  file.Size,
		MimeType:  file.Header.Get("Content-Type"),
		FileHash:  fileHash,
		FilePath:  filepath.Join(s.uploadDir, fileHash),
		Metadata:  metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save document
	if err := s.docRepo.Create(ctx, doc); err != nil {
		return nil, fmt.Errorf("failed to save document: %w", err)
	}

	// Convert to domain model
	domainDoc := mapper.DocumentModelToDomain(doc)

	return domainDoc, nil
}

// GetDocument retrieves a document by ID
func (s *DocumentService) GetDocument(ctx context.Context, id uuid.UUID) (*domain.EnhancedDocument, error) {
	doc, err := s.docRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return mapper.DocumentModelToDomain(doc), nil
}

// ListDocuments retrieves documents for a user with pagination
func (s *DocumentService) ListDocuments(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*domain.EnhancedDocument, int64, error) {
	docs, total, err := s.docRepo.GetByUserIDPaginated(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	
	return mapper.DocumentModelsToDomains(docs), total, nil
}

// UpdateDocumentStatus updates the status of a document
func (s *DocumentService) UpdateDocumentStatus(ctx context.Context, id uuid.UUID, status domain.DocumentStatus, verifierID uuid.UUID, notes string) error {
	// Update document status
	if err := s.docRepo.UpdateStatus(ctx, id, model.DocumentStatus(status), notes, verifierID); err != nil {
		return fmt.Errorf("failed to update document status: %w", err)
	}

	return nil
}

// DeleteDocument soft deletes a document
func (s *DocumentService) DeleteDocument(ctx context.Context, id uuid.UUID) error {
	return s.docRepo.Delete(ctx, id)
}

// GetDocumentsByStatus retrieves documents by status with pagination
func (s *DocumentService) GetDocumentsByStatus(ctx context.Context, status domain.DocumentStatus, page, pageSize int) ([]*domain.EnhancedDocument, int64, error) {
	docs, total, err := s.docRepo.GetByType(ctx, model.DocumentType(status), page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	
	return mapper.DocumentModelsToDomains(docs), total, nil
}

// GetDocumentsByDateRange retrieves documents by date range with pagination
func (s *DocumentService) GetDocumentsByDateRange(ctx context.Context, startDate, endDate time.Time, page, pageSize int) ([]*domain.EnhancedDocument, int64, error) {
	docs, total, err := s.docRepo.GetByDateRange(ctx, startDate, endDate, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	
	return mapper.DocumentModelsToDomains(docs), total, nil
}

// GetDocumentStats retrieves statistics about documents
func (s *DocumentService) GetDocumentStats(ctx context.Context) (*domain.DocumentStats, error) {
	stats, err := s.docRepo.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	
	// Convert to domain model
	domainStats := &domain.DocumentStats{
		TotalCount:         stats.TotalCount,
		PendingCount:       stats.PendingCount,
		VerifiedCount:      stats.VerifiedCount,
		RejectedCount:      stats.RejectedCount,
		ExpiredCount:       stats.ExpiredCount,
		AverageFileSize:    stats.AverageFileSize,
		TotalFileSize:      stats.TotalFileSize,
		ProcessingTimeAvg:  stats.ProcessingTimeAvg,
		VerificationRate:   stats.VerificationRate,
		DocumentTypeCounts: make(map[domain.DocumentType]int64),
	}
	
	// Convert document type counts
	for docType, count := range stats.DocumentTypeCounts {
		domainStats.DocumentTypeCounts[domain.DocumentType(docType)] = count
	}
	
	return domainStats, nil
}

// validateFile validates the uploaded file
func (s *DocumentService) validateFile(file *multipart.FileHeader) error {
	// Check file size (max 10MB)
	if file.Size > 10*1024*1024 {
		return errors.New("file size exceeds 10MB limit")
	}

	// Check file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".pdf":  true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	if !allowedExts[ext] {
		return errors.New("unsupported file type")
	}

	return nil
}

// readFileData reads the file data from the multipart file
func (s *DocumentService) readFileData(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	return io.ReadAll(src)
}

// calculateFileHash calculates SHA-256 hash of the file data
func (s *DocumentService) calculateFileHash(data []byte) string {
	hash := sha256.Sum256(data)
	return base64.URLEncoding.EncodeToString(hash[:])
}
