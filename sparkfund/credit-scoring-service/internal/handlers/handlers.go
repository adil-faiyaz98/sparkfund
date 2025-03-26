package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"sparkfund/credit-scoring-service/internal/models"
	"sparkfund/credit-scoring-service/internal/services"
	"sparkfund/credit-scoring-service/internal/errors"
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

// CalculateScore godoc
// @Summary      Calculate credit score
// @Description  Calculate credit score for a client based on their financial data
// @ID           calculateCreditScore
// @Tags         credit-scoring
// @Accept       json
// @Produce      json
// @Param        request body models.CreditScoreRequest true "Credit score calculation request"
// @Success      200 {object} models.CreditScoreResponse
// @Failure      400 {object} models.Error
// @Failure      500 {object} models.Error
// @Router       /credit-score/calculate [post]
// @Security     BearerAuth
func (h *Handler) CalculateScore(c *gin.Context) {
	var req models.CreditScoreRequest
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

	response, err := h.service.CalculateScore(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetCreditHistory godoc
// @Summary      Get credit history
// @Description  Retrieve credit history for a client
// @ID           getCreditHistory
// @Tags         credit-scoring
// @Accept       json
// @Produce      json
// @Param        clientId path string true "Client ID"
// @Success      200 {object} models.CreditHistory
// @Failure      400 {object} models.Error
// @Failure      404 {object} models.Error
// @Failure      500 {object} models.Error
// @Router       /credit-score/{clientId}/history [get]
// @Security     BearerAuth
func (h *Handler) GetCreditHistory(c *gin.Context) {
	clientId := c.Param("clientId")
	
	history, err := h.service.GetCreditHistory(c.Request.Context(), clientId)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, history)
}