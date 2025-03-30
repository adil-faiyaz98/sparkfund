package handlers

import (
	"net/http"
	"strconv"

	"sparkfund/services/kyc-service/internal/models"
	"sparkfund/services/kyc-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DocumentHandlers handles HTTP requests for document operations
type DocumentHandlers struct {
	docService *service.DocumentService
}

// NewDocumentHandlers creates new document handlers
func NewDocumentHandlers(docService *service.DocumentService) *DocumentHandlers {
	return &DocumentHandlers{
		docService: docService,
	}
}

// UploadDocument handles document upload requests
func (h *DocumentHandlers) UploadDocument(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Get document type from query
	docType := models.DocumentType(c.Query("type"))
	if !docType.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document type"})
		return
	}

	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
		return
	}

	// Get metadata from form
	metadata := make(map[string]interface{})
	for key, value := range c.Request.Form {
		if key != "file" && key != "type" {
			metadata[key] = value[0]
		}
	}

	// Upload document
	doc, err := h.docService.UploadDocument(c.Request.Context(), userID, file, docType, metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, doc)
}

// GetDocument handles document retrieval requests
func (h *DocumentHandlers) GetDocument(c *gin.Context) {
	// Parse document ID
	docID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	// Get document
	doc, err := h.docService.GetDocument(c.Request.Context(), docID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	c.JSON(http.StatusOK, doc)
}

// ListDocuments handles document listing requests
func (h *DocumentHandlers) ListDocuments(c *gin.Context) {
	// Get user ID from context
	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// List documents
	docs, total, err := h.docService.ListDocuments(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": docs,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
	})
}

// UpdateDocumentStatus handles document status update requests
func (h *DocumentHandlers) UpdateDocumentStatus(c *gin.Context) {
	// Parse document ID
	docID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	// Parse request body
	var req struct {
		Status          models.DocumentStatus     `json:"status" binding:"required"`
		VerifierID      uuid.UUID                 `json:"verifier_id" binding:"required"`
		Method          models.VerificationMethod `json:"method" binding:"required"`
		ConfidenceScore float64                   `json:"confidence_score" binding:"required,min=0,max=1"`
		Notes           string                    `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update document status
	err = h.docService.UpdateDocumentStatus(
		c.Request.Context(),
		docID,
		req.Status,
		req.VerifierID,
		req.Method,
		req.ConfidenceScore,
		req.Notes,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "document status updated successfully"})
}

// DeleteDocument handles document deletion requests
func (h *DocumentHandlers) DeleteDocument(c *gin.Context) {
	// Parse document ID
	docID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	// Delete document
	err = h.docService.DeleteDocument(c.Request.Context(), docID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "document deleted successfully"})
}

// GetPendingDocuments handles pending document listing requests
func (h *DocumentHandlers) GetPendingDocuments(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// List pending documents
	docs, total, err := h.docService.GetPendingDocuments(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": docs,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
	})
}

// GetExpiredDocuments handles expired document listing requests
func (h *DocumentHandlers) GetExpiredDocuments(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// List expired documents
	docs, total, err := h.docService.GetExpiredDocuments(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": docs,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
	})
}
