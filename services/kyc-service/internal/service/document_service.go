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

	"sparkfund/services/kyc-service/internal/models"
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
func (s *DocumentService) UploadDocument(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader, docType models.DocumentType, metadata map[string]interface{}) (*models.Document, error) {
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
	doc := &models.Document{
		ID:       uuid.New(),
		UserID:   userID,
		Type:     docType,
		Status:   models.DocumentStatusPending,
		FileData: fileData,
		FileHash: fileHash,
		MimeType: file.Header.Get("Content-Type"),
		FileSize: file.Size,
		Metadata: metadata,
	}

	// Save document
	if err := s.docRepo.Create(ctx, doc); err != nil {
		return nil, fmt.Errorf("failed to save document: %w", err)
	}

	return doc, nil
}

// GetDocument retrieves a document by ID
func (s *DocumentService) GetDocument(ctx context.Context, id uuid.UUID) (*models.Document, error) {
	return s.docRepo.GetByID(ctx, id)
}

// ListDocuments retrieves documents for a user with pagination
func (s *DocumentService) ListDocuments(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.Document, int64, error) {
	return s.docRepo.GetByUserID(ctx, userID, page, pageSize)
}

// UpdateDocumentStatus updates the status of a document
func (s *DocumentService) UpdateDocumentStatus(ctx context.Context, id uuid.UUID, status models.DocumentStatus, verifierID uuid.UUID, method models.VerificationMethod, confidenceScore float64, notes string) error {
	// Update document status
	if err := s.docRepo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update document status: %w", err)
	}

	// Create verification details
	details := &models.VerificationDetails{
		ID:                 uuid.New(),
		DocumentID:         id,
		VerifiedBy:         verifierID,
		VerifiedAt:         time.Now(),
		VerificationMethod: method,
		ConfidenceScore:    confidenceScore,
		Notes:              &notes,
	}

	if err := s.verRepo.Create(ctx, details); err != nil {
		return fmt.Errorf("failed to create verification details: %w", err)
	}

	return nil
}

// DeleteDocument soft deletes a document
func (s *DocumentService) DeleteDocument(ctx context.Context, id uuid.UUID) error {
	return s.docRepo.Delete(ctx, id)
}

// GetPendingDocuments retrieves all pending documents
func (s *DocumentService) GetPendingDocuments(ctx context.Context, page, pageSize int) ([]models.Document, int64, error) {
	return s.docRepo.GetPendingDocuments(ctx, page, pageSize)
}

// GetExpiredDocuments retrieves all expired documents
func (s *DocumentService) GetExpiredDocuments(ctx context.Context, page, pageSize int) ([]models.Document, int64, error) {
	return s.docRepo.GetExpiredDocuments(ctx, page, pageSize)
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
