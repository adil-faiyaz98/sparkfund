package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"sparkfund/security-service/internal/errors"
	"sparkfund/security-service/internal/models"
	"sparkfund/security-service/internal/services"
)

// Handler is responsible for handling HTTP requests
type Handler struct {
	logger  *zap.Logger
	service services.Service
}

// NewHandler creates a new handler instance
func NewHandler(logger *zap.Logger, service services.Service) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

// handleError handles error responses in a consistent way
func (h *Handler) handleError(c *gin.Context, err error) {
	var appError *errors.Error

	if errors.As(err, &appError) {
		h.logger.Error("Request failed",
			zap.Int("status", appError.Status),
			zap.String("message", appError.Message),
			zap.Error(err))
		c.JSON(appError.Status, models.Error{
			Code:    appError.Status,
			Message: appError.Message,
		})
		return
	}

	h.logger.Error("Request failed",
		zap.Int("status", http.StatusInternalServerError),
		zap.Error(err))
	c.JSON(http.StatusInternalServerError, models.Error{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
	})
}

// ValidateToken godoc
// @Summary      Validate JWT token
// @Description  Validate and decode a JWT token
// @ID           validateToken
// @Tags         security
// @Accept       json
// @Produce      json
// @Param        request body models.TokenValidationRequest true "Token validation request"
// @Success      200 {object} models.TokenValidationResponse
// @Failure      400 {object} models.Error
// @Failure      401 {object} models.Error
// @Failure      500 {object} models.Error
// @Router       /security/validate-token [post]
// @Security     BearerAuth
func (h *Handler) ValidateToken(c *gin.Context) {
	var req models.TokenValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request format",
		})
		return
	}

	if err := models.ValidateStruct(req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	response, err := h.service.ValidateToken(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GenerateToken godoc
// @Summary      Generate a new JWT token
// @Description  Generate a new JWT token based on user credentials
// @ID           generateToken
// @Tags         security
// @Accept       json
// @Produce      json
// @Param        request body models.TokenGenerationRequest true "Token generation request"
// @Success      200 {object} models.TokenGenerationResponse
// @Failure      400 {object} models.Error
// @Failure      401 {object} models.Error
// @Failure      500 {object} models.Error
// @Router       /security/generate-token [post]
func (h *Handler) GenerateToken(c *gin.Context) {
	var req models.TokenGenerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request format",
		})
		return
	}

	if err := models.ValidateStruct(req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	response, err := h.service.GenerateToken(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken godoc
// @Summary      Refresh JWT token
// @Description  Generate a new JWT token using a refresh token
// @ID           refreshToken
// @Tags         security
// @Accept       json
// @Produce      json
// @Param        request body models.TokenRefreshRequest true "Token refresh request"
// @Success      200 {object} models.TokenRefreshResponse
// @Failure      400 {object} models.Error
// @Failure      401 {object} models.Error
// @Failure      500 {object} models.Error
// @Router       /security/refresh-token [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req models.TokenRefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request format",
		})
		return
	}

	if err := models.ValidateStruct(req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	response, err := h.service.RefreshToken(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
