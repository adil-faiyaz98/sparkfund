package handlers

import (
	"net/http"
	"strconv"

	"sparkfund/services/kyc-service/internal/models"
	"sparkfund/services/kyc-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// VerificationHandlers handles HTTP requests for document verification operations
type VerificationHandlers struct {
	verService *service.VerificationService
}

// NewVerificationHandlers creates new verification handlers
func NewVerificationHandlers(verService *service.VerificationService) *VerificationHandlers {
	return &VerificationHandlers{
		verService: verService,
	}
}

// CreateVerification handles verification creation requests
func (h *VerificationHandlers) CreateVerification(c *gin.Context) {
	// Parse request body
	var details models.VerificationDetails
	if err := c.ShouldBindJSON(&details); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate verification details
	if err := h.verService.ValidateVerification(&details); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create verification record
	err := h.verService.CreateVerification(c.Request.Context(), &details)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, details)
}

// GetVerification handles verification retrieval requests
func (h *VerificationHandlers) GetVerification(c *gin.Context) {
	// Parse verification ID
	verID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification ID"})
		return
	}

	// Get verification details
	details, err := h.verService.GetVerification(c.Request.Context(), verID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "verification not found"})
		return
	}

	c.JSON(http.StatusOK, details)
}

// ListVerifications handles verification listing requests
func (h *VerificationHandlers) ListVerifications(c *gin.Context) {
	// Parse document ID
	docID, err := uuid.Parse(c.Param("document_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// List verifications
	verifications, total, err := h.verService.ListVerifications(c.Request.Context(), docID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"verifications": verifications,
		"total":         total,
		"page":          page,
		"pageSize":      pageSize,
	})
}

// UpdateVerification handles verification update requests
func (h *VerificationHandlers) UpdateVerification(c *gin.Context) {
	// Parse verification ID
	verID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification ID"})
		return
	}

	// Parse request body
	var details models.VerificationDetails
	if err := c.ShouldBindJSON(&details); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update verification record
	err = h.verService.UpdateVerification(c.Request.Context(), verID, &details)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification updated successfully"})
}

// DeleteVerification handles verification deletion requests
func (h *VerificationHandlers) DeleteVerification(c *gin.Context) {
	// Parse verification ID
	verID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification ID"})
		return
	}

	// Delete verification record
	err = h.verService.DeleteVerification(c.Request.Context(), verID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification deleted successfully"})
}

// GetVerificationStats handles verification statistics requests
func (h *VerificationHandlers) GetVerificationStats(c *gin.Context) {
	// Get verification statistics
	stats, err := h.verService.GetVerificationStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetVerificationHistory handles verification history requests
func (h *VerificationHandlers) GetVerificationHistory(c *gin.Context) {
	// Parse document ID
	docID, err := uuid.Parse(c.Param("document_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	// Get verification history
	history, err := h.verService.GetVerificationHistory(c.Request.Context(), docID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// GetVerificationSummary handles verification summary requests
func (h *VerificationHandlers) GetVerificationSummary(c *gin.Context) {
	// Parse document ID
	docID, err := uuid.Parse(c.Param("document_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	// Get verification summary
	summary, err := h.verService.GetVerificationSummary(c.Request.Context(), docID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
