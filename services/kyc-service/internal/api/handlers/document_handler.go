package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"sparkfund/services/kyc-service/internal/api/dto"
	"sparkfund/services/kyc-service/internal/domain"
	"sparkfund/services/kyc-service/internal/service"
)

// DocumentHandler handles document-related HTTP requests
type DocumentHandler struct {
	documentService *service.DocumentService
}

// NewDocumentHandler creates a new document handler
func NewDocumentHandler(documentService *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
	}
}

// RegisterRoutes registers the document routes
func (h *DocumentHandler) RegisterRoutes(router *gin.RouterGroup) {
	documents := router.Group("/documents")
	{
		documents.POST("", h.UploadDocument)
		documents.GET("/:id", h.GetDocument)
		documents.GET("", h.ListDocuments)
		documents.PUT("/:id/status", h.UpdateDocumentStatus)
		documents.DELETE("/:id", h.DeleteDocument)
		documents.GET("/stats", h.GetDocumentStats)
		documents.GET("/by-status/:status", h.GetDocumentsByStatus)
		documents.GET("/by-date-range", h.GetDocumentsByDateRange)
	}
}

// UploadDocument handles document upload
// @Summary Upload a document
// @Description Upload a new document for KYC verification
// @Tags documents
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Document file"
// @Param type formData string true "Document type"
// @Param metadata formData string false "Document metadata (JSON)"
// @Param user_id formData string true "User ID"
// @Success 201 {object} dto.DocumentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /documents [post]
func (h *DocumentHandler) UploadDocument(c *gin.Context) {
	// Parse user ID
	userIDStr := c.PostForm("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	// Get file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "File is required",
		})
		return
	}

	// Get document type
	docType := c.PostForm("type")
	if docType == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Document type is required",
		})
		return
	}

	// Parse metadata (optional)
	var metadata map[string]interface{}
	metadataStr := c.PostForm("metadata")
	if metadataStr != "" {
		if err := c.ShouldBindJSON(&metadata); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: "Invalid metadata format",
			})
			return
		}
	}

	// Upload document
	document, err := h.documentService.UploadDocument(c.Request.Context(), userID, file, docType, metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to upload document",
		})
		return
	}

	// Return response
	c.JSON(http.StatusCreated, dto.FromDomainDocument(document))
}

// GetDocument handles document retrieval
// @Summary Get a document
// @Description Get a document by ID
// @Tags documents
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} dto.DocumentResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /documents/{id} [get]
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	// Parse document ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid document ID",
		})
		return
	}

	// Get document
	document, err := h.documentService.GetDocument(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "Document not found",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.FromDomainDocument(document))
}

// ListDocuments handles document listing
// @Summary List documents
// @Description List documents for a user with pagination
// @Tags documents
// @Produce json
// @Param user_id query string true "User ID"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} dto.DocumentListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /documents [get]
func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	// Parse user ID
	userIDStr := c.Query("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	// Parse pagination parameters
	page, pageSize := getPaginationParams(c)

	// Get documents
	documents, total, err := h.documentService.ListDocuments(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to list documents",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.DocumentListResponse{
		Documents: dto.FromDomainDocuments(documents),
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	})
}

// UpdateDocumentStatus handles document status update
// @Summary Update document status
// @Description Update the status of a document
// @Tags documents
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Param request body dto.DocumentStatusUpdateRequest true "Status update request"
// @Success 200 {object} dto.DocumentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /documents/{id}/status [put]
func (h *DocumentHandler) UpdateDocumentStatus(c *gin.Context) {
	// Parse document ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid document ID",
		})
		return
	}

	// Parse request
	var req dto.DocumentStatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Update document status
	err = h.documentService.UpdateDocumentStatus(
		c.Request.Context(),
		id,
		domain.DocumentStatus(req.Status),
		req.VerifierID,
		req.Notes,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to update document status",
		})
		return
	}

	// Get updated document
	document, err := h.documentService.GetDocument(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "Document not found",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.FromDomainDocument(document))
}

// DeleteDocument handles document deletion
// @Summary Delete a document
// @Description Delete a document by ID
// @Tags documents
// @Produce json
// @Param id path string true "Document ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /documents/{id} [delete]
func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	// Parse document ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid document ID",
		})
		return
	}

	// Delete document
	err = h.documentService.DeleteDocument(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to delete document",
		})
		return
	}

	// Return response
	c.Status(http.StatusNoContent)
}

// GetDocumentStats handles document statistics retrieval
// @Summary Get document statistics
// @Description Get statistics about documents
// @Tags documents
// @Produce json
// @Success 200 {object} dto.DocumentStatsResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /documents/stats [get]
func (h *DocumentHandler) GetDocumentStats(c *gin.Context) {
	// Get document stats
	stats, err := h.documentService.GetDocumentStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get document statistics",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.FromDomainDocumentStats(stats))
}

// GetDocumentsByStatus handles document retrieval by status
// @Summary Get documents by status
// @Description Get documents by status with pagination
// @Tags documents
// @Produce json
// @Param status path string true "Document status"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} dto.DocumentListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /documents/by-status/{status} [get]
func (h *DocumentHandler) GetDocumentsByStatus(c *gin.Context) {
	// Parse status
	status := domain.DocumentStatus(c.Param("status"))

	// Parse pagination parameters
	page, pageSize := getPaginationParams(c)

	// Get documents
	documents, total, err := h.documentService.GetDocumentsByStatus(c.Request.Context(), status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get documents by status",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.DocumentListResponse{
		Documents: dto.FromDomainDocuments(documents),
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	})
}

// GetDocumentsByDateRange handles document retrieval by date range
// @Summary Get documents by date range
// @Description Get documents by date range with pagination
// @Tags documents
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} dto.DocumentListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /documents/by-date-range [get]
func (h *DocumentHandler) GetDocumentsByDateRange(c *gin.Context) {
	// Parse date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Start date and end date are required",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid start date format (YYYY-MM-DD)",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid end date format (YYYY-MM-DD)",
		})
		return
	}

	// Add one day to end date to include the end date in the range
	endDate = endDate.Add(24 * time.Hour)

	// Parse pagination parameters
	page, pageSize := getPaginationParams(c)

	// Get documents
	documents, total, err := h.documentService.GetDocumentsByDateRange(c.Request.Context(), startDate, endDate, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get documents by date range",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.DocumentListResponse{
		Documents: dto.FromDomainDocuments(documents),
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	})
}

// Helper function to get pagination parameters
func getPaginationParams(c *gin.Context) (int, int) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return page, pageSize
}
