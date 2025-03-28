package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sparkfund/credit-scoring-service/internal/model"
	"github.com/sparkfund/credit-scoring-service/internal/service"
)

type CreditHandler struct {
	creditService service.CreditService
}

func NewCreditHandler(creditService service.CreditService) *CreditHandler {
	return &CreditHandler{
		creditService: creditService,
	}
}

func (h *CreditHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1/credit")
	{
		api.POST("/check", h.ProcessCreditCheck)
		api.GET("/score/:user_id", h.GetCreditScore)
		api.GET("/history/:user_id", h.GetCreditHistory)
		api.PUT("/history/:id/status", h.UpdateCreditHistoryStatus)
	}
}

func (h *CreditHandler) ProcessCreditCheck(c *gin.Context) {
	var req model.CreditCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if req.UserID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	// Process credit check
	response, err := h.creditService.ProcessCreditCheck(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *CreditHandler) GetCreditScore(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	score, err := h.creditService.GetCreditScore(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, score)
}

func (h *CreditHandler) GetCreditHistory(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	history, err := h.creditService.GetCreditHistory(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

func (h *CreditHandler) UpdateCreditHistoryStatus(c *gin.Context) {
	historyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid history_id"})
		return
	}

	var req struct {
		Status model.CreditHistoryStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.creditService.UpdateCreditHistory(c.Request.Context(), historyID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "credit history status updated successfully"})
}
