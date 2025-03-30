package handler

import (
	"net/http"

	"aml-service/internal/model"
	"aml-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AMLHandler struct {
	amlService service.AMLService
}

func NewAMLHandler(amlService service.AMLService) *AMLHandler {
	return &AMLHandler{
		amlService: amlService,
	}
}

func (h *AMLHandler) ProcessTransaction(c *gin.Context) {
	var req model.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Process transaction
	tx, err := h.amlService.ProcessTransaction(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process transaction"})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

func (h *AMLHandler) FlagTransaction(c *gin.Context) {
	txID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var req struct {
		Reason    string `json:"reason" binding:"required"`
		FlaggedBy string `json:"flaggedBy" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.amlService.FlagTransaction(c.Request.Context(), txID, req.Reason, req.FlaggedBy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to flag transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction flagged successfully"})
}

func (h *AMLHandler) ReviewTransaction(c *gin.Context) {
	txID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var req struct {
		Status     model.TransactionStatus `json:"status" binding:"required,oneof=approved rejected"`
		Notes      string                  `json:"notes" binding:"required"`
		ReviewedBy string                  `json:"reviewedBy" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.amlService.ReviewTransaction(c.Request.Context(), txID, req.Status, req.Notes, req.ReviewedBy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to review transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction reviewed successfully"})
}

func (h *AMLHandler) ListTransactions(c *gin.Context) {
	var filter model.TransactionFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Set user ID in filter
	uid := userID.(uuid.UUID)
	filter.UserID = &uid

	txs, err := h.amlService.ListTransactions(c.Request.Context(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list transactions"})
		return
	}

	c.JSON(http.StatusOK, txs)
}

func (h *AMLHandler) GetTransaction(c *gin.Context) {
	txID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	tx, err := h.amlService.GetTransaction(c.Request.Context(), txID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, tx)
}

func (h *AMLHandler) RegisterRoutes(router *gin.Engine) {
	aml := router.Group("/api/v1/aml")
	{
		aml.POST("/transactions", h.ProcessTransaction)
		aml.GET("/transactions", h.ListTransactions)
		aml.GET("/transactions/:id", h.GetTransaction)
		aml.POST("/transactions/:id/flag", h.FlagTransaction)
		aml.POST("/transactions/:id/review", h.ReviewTransaction)
	}
}
