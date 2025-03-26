package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"sparkfund/aml-service/internal/models"
	"sparkfund/aml-service/internal/services"
	"sparkfund/aml-service/internal/errors"
)

type Handler struct {
	logger  *zap.Logger
	service services.Service
}

func NewHandler(logger *zap.Logger, service services.Service) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

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

// CheckCompliance godoc
// @Summary      Check AML compliance
// @Description  Check if a transaction or client complies with AML regulations
// @ID           checkCompliance 
// @Tags         compliance
// @Accept       json
// @Produce      json
// @Param        request body models.ComplianceRequest true "Compliance check request"
// @Success      200 {object} models.ComplianceResponse
// @Failure      400 {object} models.Error
// @Failure      500 {object} models.Error
// @Router       /compliance/check [post]
// @Security     BearerAuth
func (h *Handler) CheckCompliance(c *gin.Context) {
	var req models.ComplianceRequest
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

	response, err := h.service.CheckCompliance(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}