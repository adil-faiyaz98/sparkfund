package handlers

import (
	"net/http"
	"strconv"

	"your-project/internal/investments"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InvestmentHandler struct {
	service investments.InvestmentService
}

func NewInvestmentHandler(service investments.InvestmentService) *InvestmentHandler {
	return &InvestmentHandler{service: service}
}

// CreateInvestment handles the creation of a new investment
func (h *InvestmentHandler) CreateInvestment(c *gin.Context) {
	var investment investments.Investment
	if err := c.ShouldBindJSON(&investment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateInvestment(c.Request.Context(), &investment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, investment)
}

// GetInvestment retrieves an investment by ID
func (h *InvestmentHandler) GetInvestment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid investment ID"})
		return
	}

	investment, err := h.service.GetInvestment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "investment not found"})
		return
	}

	c.JSON(http.StatusOK, investment)
}

// GetUserInvestments retrieves all investments for a user
func (h *InvestmentHandler) GetUserInvestments(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	investments, err := h.service.GetUserInvestments(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, investments)
}

// GetAccountInvestments retrieves all investments for an account
func (h *InvestmentHandler) GetAccountInvestments(c *gin.Context) {
	accountID, err := uuid.Parse(c.Param("accountId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	investments, err := h.service.GetAccountInvestments(c.Request.Context(), accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, investments)
}

// UpdateInvestment updates an existing investment
func (h *InvestmentHandler) UpdateInvestment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid investment ID"})
		return
	}

	var investment investments.Investment
	if err := c.ShouldBindJSON(&investment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	investment.ID = id
	if err := h.service.UpdateInvestment(c.Request.Context(), &investment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, investment)
}

// DeleteInvestment deletes an investment
func (h *InvestmentHandler) DeleteInvestment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid investment ID"})
		return
	}

	if err := h.service.DeleteInvestment(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetInvestmentsBySymbol retrieves investments by symbol
func (h *InvestmentHandler) GetInvestmentsBySymbol(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symbol is required"})
		return
	}

	investments, err := h.service.GetInvestmentsBySymbol(c.Request.Context(), symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, investments)
}

// UpdateInvestmentPrice updates the current price of an investment
func (h *InvestmentHandler) UpdateInvestmentPrice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid investment ID"})
		return
	}

	priceStr := c.Param("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
		return
	}

	if err := h.service.UpdateInvestmentPrice(c.Request.Context(), id, price); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
