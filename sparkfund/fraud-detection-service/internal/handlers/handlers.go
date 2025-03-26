package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"sparkfund/fraud-detection-service/internal/models"
	"sparkfund/fraud-detection-service/internal/services"
	"sparkfund/fraud-detection-service/internal/errors"
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

// DetectFraud godoc
// @Summary      Detect potential fraud
// @Description  Analyze a transaction or activity for potential fraud
// @ID           detectFraud
// @Tags         fraud-detection
// @Accept       json
// @Produce      json
// @Param        request body models.FraudDetectionRequest true "Fraud detection request"
// @Success      200 {object} models.FraudDetectionResponse
// @Failure      400 {object} models.Error
// @Failure      500 {object} models.Error
// @Router       /fraud/detect [post]
// @Security     BearerAuth
func (h *Handler) DetectFraud(c *gin.Context) {
	var req models.FraudDetectionRequest
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

	response, err := h.service.DetectFraud(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}