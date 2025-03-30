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

type VerificationHandler struct {
	verificationService *service.VerificationService
}

func NewVerificationHandler(verificationService *service.VerificationService) *VerificationHandler {
	return &VerificationHandler{
		verificationService: verificationService,
	}
}

type CreateVerificationRequest struct {
	DocumentID uuid.UUID                `json:"document_id" binding:"required"`
	Method     model.VerificationMethod `json:"method" binding:"required"`
}

type UpdateVerificationStatusRequest struct {
	Status     model.VerificationStatus `json:"status" binding:"required"`
	VerifierID uuid.UUID                `json:"verifier_id" binding:"required"`
	Notes      string                   `json:"notes"`
}

type UpdateConfidenceScoreRequest struct {
	Score float64 `json:"score" binding:"required,min=0,max=100"`
}

type UpdateMetadataRequest struct {
	Metadata map[string]interface{} `json:"metadata" binding:"required"`
}

// CreateVerification handles the creation of a new verification
func (h *VerificationHandler) CreateVerification(c *gin.Context) {
	var req CreateVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	verification, err := h.verificationService.CreateVerification(req.DocumentID, req.Method)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, verification)
}

// GetVerification retrieves a verification by ID
func (h *VerificationHandler) GetVerification(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification ID"})
		return
	}

	verification, err := h.verificationService.GetVerification(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "verification not found"})
		return
	}

	c.JSON(http.StatusOK, verification)
}

// GetDocumentVerifications retrieves all verifications for a document
func (h *VerificationHandler) GetDocumentVerifications(c *gin.Context) {
	documentID, err := uuid.Parse(c.Param("document_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	verifications, err := h.verificationService.GetDocumentVerifications(documentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, verifications)
}

// UpdateVerificationStatus updates the status of a verification
func (h *VerificationHandler) UpdateVerificationStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification ID"})
		return
	}

	var req UpdateVerificationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.verificationService.UpdateVerificationStatus(id, req.Status, req.VerifierID, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification status updated successfully"})
}

// DeleteVerification deletes a verification
func (h *VerificationHandler) DeleteVerification(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification ID"})
		return
	}

	err = h.verificationService.DeleteVerification(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification deleted successfully"})
}

// GetVerificationHistory retrieves the history of a verification
func (h *VerificationHandler) GetVerificationHistory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification ID"})
		return
	}

	history, err := h.verificationService.GetVerificationHistory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// GetVerificationStats retrieves verification statistics
func (h *VerificationHandler) GetVerificationStats(c *gin.Context) {
	stats, err := h.verificationService.GetVerificationStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetVerificationSummary retrieves a summary of a verification
func (h *VerificationHandler) GetVerificationSummary(c *gin.Context) {
	documentID, err := uuid.Parse(c.Param("document_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document ID"})
		return
	}

	summary, err := h.verificationService.GetVerificationSummary(documentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetExpiredVerifications retrieves all expired verifications
func (h *VerificationHandler) GetExpiredVerifications(c *gin.Context) {
	verifications, err := h.verificationService.GetExpiredVerifications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, verifications)
}

// GetPendingVerifications retrieves all pending verifications
func (h *VerificationHandler) GetPendingVerifications(c *gin.Context) {
	verifications, err := h.verificationService.GetPendingVerifications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, verifications)
}

// GetFailedVerifications retrieves all failed verifications
func (h *VerificationHandler) GetFailedVerifications(c *gin.Context) {
	verifications, err := h.verificationService.GetFailedVerifications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, verifications)
}

// GetVerificationsByVerifier retrieves verifications by verifier ID
func (h *VerificationHandler) GetVerificationsByVerifier(c *gin.Context) {
	verifierID, err := uuid.Parse(c.Param("verifier_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verifier ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	verifications, total, err := h.verificationService.GetVerificationsByVerifier(verifierID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"verifications": verifications,
		"total":         total,
		"page":          page,
		"page_size":     pageSize,
	})
}

// GetVerificationsByMethod retrieves verifications by method
func (h *VerificationHandler) GetVerificationsByMethod(c *gin.Context) {
	method := model.VerificationMethod(c.Param("method"))
	if !method.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification method"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	verifications, total, err := h.verificationService.GetVerificationsByMethod(method, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"verifications": verifications,
		"total":         total,
		"page":          page,
		"page_size":     pageSize,
	})
}

// GetVerificationsByDateRange retrieves verifications within a date range
func (h *VerificationHandler) GetVerificationsByDateRange(c *gin.Context) {
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

	verifications, total, err := h.verificationService.GetVerificationsByDateRange(startDate, endDate, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"verifications": verifications,
		"total":         total,
		"page":          page,
		"page_size":     pageSize,
	})
}

// UpdateConfidenceScore updates the confidence score of a verification
func (h *VerificationHandler) UpdateConfidenceScore(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification ID"})
		return
	}

	var req UpdateConfidenceScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.verificationService.UpdateConfidenceScore(id, req.Score)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "confidence score updated successfully"})
}

// UpdateVerificationMetadata updates the metadata of a verification
func (h *VerificationHandler) UpdateVerificationMetadata(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification ID"})
		return
	}

	var req UpdateMetadataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.verificationService.UpdateVerificationMetadata(id, req.Metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "metadata updated successfully"})
}
