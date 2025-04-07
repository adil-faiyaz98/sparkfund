package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"sparkfund/services/kyc-service/internal/api/dto"
	"sparkfund/services/kyc-service/internal/domain"
	"sparkfund/services/kyc-service/internal/model"
	"sparkfund/services/kyc-service/internal/service"
)

// KYCHandler handles KYC-related HTTP requests
type KYCHandler struct {
	kycService *service.KYCService
}

// NewKYCHandler creates a new KYC handler
func NewKYCHandler(kycService *service.KYCService) *KYCHandler {
	return &KYCHandler{
		kycService: kycService,
	}
}

// RegisterRoutes registers the KYC routes
func (h *KYCHandler) RegisterRoutes(router *gin.RouterGroup) {
	kyc := router.Group("/kyc")
	{
		kyc.POST("", h.CreateKYC)
		kyc.GET("/:id", h.GetKYC)
		kyc.GET("/user/:user_id", h.GetKYCByUserID)
		kyc.PUT("/:id/status", h.UpdateKYCStatus)
		kyc.PUT("/:id/risk", h.UpdateKYCRiskLevel)
		kyc.GET("", h.ListKYCs)
		kyc.GET("/by-status/:status", h.GetKYCsByStatus)
		kyc.GET("/by-risk-level/:risk_level", h.GetKYCsByRiskLevel)
	}
}

// CreateKYC handles KYC creation
// @Summary Create a KYC verification
// @Description Create a new KYC verification
// @Tags kyc
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Param request body dto.KYCRequest true "KYC request"
// @Success 201 {object} dto.KYCResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /kyc [post]
func (h *KYCHandler) CreateKYC(c *gin.Context) {
	// Parse user ID
	userIDStr := c.Query("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	// Parse request
	var req dto.KYCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Convert to model request
	modelReq := &model.KYCRequest{
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		DateOfBirth:       req.DateOfBirth,
		Nationality:       req.Nationality,
		Email:             req.Email,
		PhoneNumber:       req.PhoneNumber,
		Address:           req.Address,
		City:              req.City,
		State:             req.State,
		Country:           req.Country,
		PostalCode:        req.PostalCode,
		DocumentType:      req.DocumentType,
		DocumentNumber:    req.DocumentNumber,
		DocumentFront:     req.DocumentFront,
		DocumentBack:      req.DocumentBack,
		SelfieImage:       req.SelfieImage,
		DocumentExpiry:    req.DocumentExpiry,
		TransactionAmount: req.TransactionAmount,
	}

	// Create KYC
	kyc, err := h.kycService.CreateKYC(c.Request.Context(), userID, modelReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to create KYC verification",
		})
		return
	}

	// Return response
	c.JSON(http.StatusCreated, dto.FromDomainKYC(kyc))
}

// GetKYC handles KYC retrieval
// @Summary Get a KYC verification
// @Description Get a KYC verification by ID
// @Tags kyc
// @Produce json
// @Param id path string true "KYC ID"
// @Success 200 {object} dto.KYCDetailResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /kyc/{id} [get]
func (h *KYCHandler) GetKYC(c *gin.Context) {
	// Parse KYC ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid KYC ID",
		})
		return
	}

	// Get KYC
	kyc, err := h.kycService.GetKYC(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "KYC verification not found",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.FromDomainKYCDetail(kyc))
}

// GetKYCByUserID handles KYC retrieval by user ID
// @Summary Get a KYC verification by user ID
// @Description Get a KYC verification by user ID
// @Tags kyc
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} dto.KYCDetailResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /kyc/user/{user_id} [get]
func (h *KYCHandler) GetKYCByUserID(c *gin.Context) {
	// Parse user ID
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	// Get KYC
	kyc, err := h.kycService.GetKYCByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "KYC verification not found",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.FromDomainKYCDetail(kyc))
}

// UpdateKYCStatus handles KYC status update
// @Summary Update KYC status
// @Description Update the status of a KYC verification
// @Tags kyc
// @Accept json
// @Produce json
// @Param id path string true "KYC ID"
// @Param request body dto.KYCStatusUpdateRequest true "Status update request"
// @Success 200 {object} dto.KYCResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /kyc/{id}/status [put]
func (h *KYCHandler) UpdateKYCStatus(c *gin.Context) {
	// Parse KYC ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid KYC ID",
		})
		return
	}

	// Parse request
	var req dto.KYCStatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Update KYC status
	err = h.kycService.UpdateKYCStatus(
		c.Request.Context(),
		id,
		domain.KYCStatus(req.Status),
		req.Notes,
		req.ReviewerID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to update KYC status",
		})
		return
	}

	// Get updated KYC
	kyc, err := h.kycService.GetKYC(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "KYC verification not found",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.FromDomainKYC(kyc))
}

// UpdateKYCRiskLevel handles KYC risk level update
// @Summary Update KYC risk level
// @Description Update the risk level of a KYC verification
// @Tags kyc
// @Accept json
// @Produce json
// @Param id path string true "KYC ID"
// @Param request body dto.KYCRiskUpdateRequest true "Risk update request"
// @Success 200 {object} dto.KYCResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /kyc/{id}/risk [put]
func (h *KYCHandler) UpdateKYCRiskLevel(c *gin.Context) {
	// Parse KYC ID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid KYC ID",
		})
		return
	}

	// Parse request
	var req dto.KYCRiskUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Update KYC risk level
	err = h.kycService.UpdateKYCRiskLevel(
		c.Request.Context(),
		id,
		domain.RiskLevel(req.RiskLevel),
		req.RiskScore,
		req.Notes,
		req.ReviewerID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to update KYC risk level",
		})
		return
	}

	// Get updated KYC
	kyc, err := h.kycService.GetKYC(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "KYC verification not found",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.FromDomainKYC(kyc))
}

// ListKYCs handles KYC listing
// @Summary List KYC verifications
// @Description List KYC verifications with pagination
// @Tags kyc
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} dto.KYCListResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /kyc [get]
func (h *KYCHandler) ListKYCs(c *gin.Context) {
	// Parse pagination parameters
	page, pageSize := getPaginationParams(c)

	// Get KYCs
	kycs, total, err := h.kycService.ListKYCs(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to list KYC verifications",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.KYCListResponse{
		KYCs:     dto.FromDomainKYCs(kycs),
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetKYCsByStatus handles KYC retrieval by status
// @Summary Get KYC verifications by status
// @Description Get KYC verifications by status with pagination
// @Tags kyc
// @Produce json
// @Param status path string true "KYC status"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} dto.KYCListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /kyc/by-status/{status} [get]
func (h *KYCHandler) GetKYCsByStatus(c *gin.Context) {
	// Parse status
	status := domain.KYCStatus(c.Param("status"))

	// Parse pagination parameters
	page, pageSize := getPaginationParams(c)

	// Get KYCs
	kycs, total, err := h.kycService.GetKYCsByStatus(c.Request.Context(), status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get KYC verifications by status",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.KYCListResponse{
		KYCs:     dto.FromDomainKYCs(kycs),
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetKYCsByRiskLevel handles KYC retrieval by risk level
// @Summary Get KYC verifications by risk level
// @Description Get KYC verifications by risk level with pagination
// @Tags kyc
// @Produce json
// @Param risk_level path string true "Risk level"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} dto.KYCListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /kyc/by-risk-level/{risk_level} [get]
func (h *KYCHandler) GetKYCsByRiskLevel(c *gin.Context) {
	// Parse risk level
	riskLevel := domain.RiskLevel(c.Param("risk_level"))

	// Parse pagination parameters
	page, pageSize := getPaginationParams(c)

	// Get KYCs
	kycs, total, err := h.kycService.GetKYCsByRiskLevel(c.Request.Context(), riskLevel, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to get KYC verifications by risk level",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, dto.KYCListResponse{
		KYCs:     dto.FromDomainKYCs(kycs),
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
