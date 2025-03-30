package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"sparkfund/services/kyc-service/internal/model"
	"sparkfund/services/kyc-service/internal/service"
)

type DocumentHandler struct {
	documentService *service.DocumentService
}

func NewDocumentHandler(documentService *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
	}
}

type UploadDocumentRequest struct {
	Type model.DocumentType `form:"type" binding:"required"`
}

// UploadDocument handles document upload
func (h *DocumentHandler) UploadDocument(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req UploadDocumentRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}

	// Validate file
	if err := h.documentService.ValidateDocument(file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	document, err := h.documentService.UploadDocument(userID, req.Type, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, document)
}

// GetDocument retrieves a document by ID
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	document, err := h.documentService.GetDocument(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	c.JSON(http.StatusOK, document)
}

// GetUserDocuments retrieves all documents for a user
func (h *DocumentHandler) GetUserDocuments(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	documents, err := h.documentService.GetUserDocuments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, documents)
}

// UpdateDocumentStatus updates the status of a document
func (h *DocumentHandler) UpdateDocumentStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	updatedBy, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req struct {
		Status model.DocumentStatus `json:"status" binding:"required"`
		Notes  string               `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.documentService.UpdateDocumentStatus(id, req.Status, req.Notes, updatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "document status updated successfully"})
}

// DeleteDocument deletes a document
func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	err = h.documentService.DeleteDocument(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "document deleted successfully"})
}

// GetDocumentHistory retrieves the history of a document
func (h *DocumentHandler) GetDocumentHistory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	history, err := h.documentService.GetDocumentHistory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// GetDocumentStats retrieves document statistics
func (h *DocumentHandler) GetDocumentStats(c *gin.Context) {
	stats, err := h.documentService.GetDocumentStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetDocumentSummary retrieves a summary of a document
func (h *DocumentHandler) GetDocumentSummary(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	summary, err := h.documentService.GetDocumentSummary(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetExpiredDocuments retrieves all expired documents
func (h *DocumentHandler) GetExpiredDocuments(c *gin.Context) {
	documents, err := h.documentService.GetExpiredDocuments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, documents)
}

// GetPendingDocuments retrieves all pending documents
func (h *DocumentHandler) GetPendingDocuments(c *gin.Context) {
	documents, err := h.documentService.GetPendingDocuments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, documents)
}

// GetRejectedDocuments retrieves all rejected documents
func (h *DocumentHandler) GetRejectedDocuments(c *gin.Context) {
	documents, err := h.documentService.GetRejectedDocuments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, documents)
}

// GetDocumentsByType retrieves documents by type
func (h *DocumentHandler) GetDocumentsByType(c *gin.Context) {
	docType := model.DocumentType(c.Param("type"))
	if !docType.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document type"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	documents, total, err := h.documentService.GetDocumentsByType(docType, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": documents,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetDocumentsByDateRange retrieves documents within a date range
func (h *DocumentHandler) GetDocumentsByDateRange(c *gin.Context) {
	startDate, err := time.Parse(time.RFC3339, c.Query("start_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, c.Query("end_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	documents, total, err := h.documentService.GetDocumentsByDateRange(startDate, endDate, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": documents,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateDocumentMetadata updates the metadata of a document
func (h *DocumentHandler) UpdateDocumentMetadata(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	var req struct {
		Metadata map[string]interface{} `json:"metadata" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.documentService.UpdateDocumentMetadata(id, req.Metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "document metadata updated successfully"})
}
