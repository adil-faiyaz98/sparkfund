package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"sparkfund/services/kyc-service/internal/api/dto"
	"sparkfund/services/kyc-service/internal/domain"
	"sparkfund/services/kyc-service/internal/service"
)

// VerificationHandler handles verification-related HTTP requests
type VerificationHandler struct {
	verificationService *service.VerificationService
}

// NewVerificationHandler creates a new verification handler
func NewVerificationHandler(verificationService *service.VerificationService) *VerificationHandler {
	return &VerificationHandler{
		verificationService: verificationService,
	}
}

// RegisterRoutes registers the verification routes
func (h *VerificationHandler) RegisterRoutes(router *gin.RouterGroup) {
	verifications := router.Group("/verifications")
	{
		verifications.POST("", h.CreateVerification)
		verifications.GET("/:id", h.GetVerification)
		verifications.GET("", h.ListVerifications)
		verifications.PUT("/:id/status", h.UpdateVerificationStatus)
		verifications.POST("/:id/result", h.CreateVerificationResult)
		verifications.GET("/document/:document_id", h.GetVerificationsByDocument)
		verifications.GET("/kyc/:kyc_id", h.GetVerificationsByKYC)
	}
}

// CreateVerification handles verification creation
// @Summary Create a verification
// @Description Create a new verification for a document
// @Tags verifications
// @Accept json
// @Produce json
// @Param request body dto.VerificationRequest true "Verification request"
// @Success 201 {object} dto.VerificationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /verifications [post]
func (h *VerificationHandler) CreateVerification(c *gin.Context) {
	// Parse request
	var req dto.VerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Create verification
	verification, err := h.verificationService.CreateVerification(
		c.Request.Context(),
		req.DocumentID,
		domain.VerificationMethod(req.Method),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to create verification",
		})
		return
	}

	// Return response
	c.JSON(http.StatusCreated, dto.FromDomainVerification(verification))
}

// GetVerification handles verification retrieval
// @Summary Get a verification
// @Description Get a verification by ID
// @Tags verifications
// @Produce json
// @Param id path string true "Verification ID"
// @Success 200 {object} dto.VerificationResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /verifications/{id} [get]
func (h *VerificationHandler) GetVerification(c *gin.Context) {
	// Parse verification ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid verification ID",
		})
		return
	}

	// Get verification
	verification, err := h.verificationService.GetVerification(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "Verification not found",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.FromDomainVerification(verification))
}

// ListVerifications handles verification listing
// @Summary List verifications
// @Description List verifications with pagination
// @Tags verifications
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} dto.VerificationListResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /verifications [get]
func (h *VerificationHandler) ListVerifications(c *gin.Context) {
	// Parse pagination parameters
	page, pageSize := getPaginationParams(c)

	// Get verifications
	verifications, total, err := h.verificationService.ListVerifications(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to list verifications",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.VerificationListResponse{
		Verifications: dto.FromDomainVerifications(verifications),
		Total:         total,
		Page:          page,
		PageSize:      pageSize,
	})
}

// UpdateVerificationStatus handles verification status update
// @Summary Update verification status
// @Description Update the status of a verification
// @Tags verifications
// @Accept json
// @Produce json
// @Param id path string true "Verification ID"
// @Param request body dto.VerificationStatusUpdateRequest true "Status update request"
// @Success 200 {object} dto.VerificationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /verifications/{id}/status [put]
func (h *VerificationHandler) UpdateVerificationStatus(c *gin.Context) {
	// Parse verification ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid verification ID",
		})
		return
	}

	// Parse request
	var req dto.VerificationStatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Update verification status
	err = h.verificationService.UpdateVerificationStatus(
		c.Request.Context(),
		id,
		domain.VerificationStatus(req.Status),
		req.ConfidenceScore,
		req.Notes,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to update verification status",
		})
		return
	}

	// Get updated verification
	verification, err := h.verificationService.GetVerification(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "Verification not found",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.FromDomainVerification(verification))
}

// CreateVerificationResult handles verification result creation
// @Summary Create a verification result
// @Description Create a result for a verification
// @Tags verifications
// @Accept json
// @Produce json
// @Param id path string true "Verification ID"
// @Param request body dto.VerificationResultRequest true "Result request"
// @Success 201 {object} dto.VerificationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /verifications/{id}/result [post]
func (h *VerificationHandler) CreateVerificationResult(c *gin.Context) {
	// Parse verification ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid verification ID",
		})
		return
	}

	// Parse request
	var req dto.VerificationResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Create verification result
	err = h.verificationService.CreateVerificationResult(
		c.Request.Context(),
		id,
		dto.ToDomainVerificationResult(&req),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to create verification result",
		})
		return
	}

	// Get updated verification
	verification, err := h.verificationService.GetVerification(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "Verification not found",
		})
		return
	}

	// Return response
	c.JSON(http.StatusCreated, dto.FromDomainVerification(verification))
}

// GetVerificationsByDocument handles verification retrieval by document
// @Summary Get verifications by document
// @Description Get verifications for a document
// @Tags verifications
// @Produce json
// @Param document_id path string true "Document ID"
// @Success 200 {array} dto.VerificationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /verifications/document/{document_id} [get]
func (h *VerificationHandler) GetVerificationsByDocument(c *gin.Context) {
	// Parse document ID
	documentID, err := uuid.Parse(c.Param("document_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid document ID",
		})
		return
	}

	// Get verifications
	verifications, err := h.verificationService.GetVerificationsByDocument(c.Request.Context(), documentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get verifications by document",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.FromDomainVerifications(verifications))
}

// GetVerificationsByKYC handles verification retrieval by KYC
// @Summary Get verifications by KYC
// @Description Get verifications for a KYC verification
// @Tags verifications
// @Produce json
// @Param kyc_id path string true "KYC ID"
// @Success 200 {array} dto.VerificationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /verifications/kyc/{kyc_id} [get]
func (h *VerificationHandler) GetVerificationsByKYC(c *gin.Context) {
	// Parse KYC ID
	kycID, err := uuid.Parse(c.Param("kyc_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid KYC ID",
		})
		return
	}

	// Get verifications
	verifications, err := h.verificationService.GetVerificationsByKYC(c.Request.Context(), kycID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get verifications by KYC",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.FromDomainVerifications(verifications))
}
