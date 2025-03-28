package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sparkfund/investment-service/internal/database"
	"github.com/sparkfund/investment-service/internal/models"
)

// CreateInvestment creates a new investment
func CreateInvestment(c *gin.Context) {
	var investment models.Investment
	if err := c.ShouldBindJSON(&investment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default values
	investment.PurchaseDate = time.Now()
	investment.Status = "ACTIVE"

	// Create investment
	if err := database.DB.Create(&investment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create investment"})
		return
	}

	c.JSON(http.StatusCreated, investment)
}

// GetInvestment retrieves an investment by ID
func GetInvestment(c *gin.Context) {
	id := c.Param("id")
	var investment models.Investment

	if err := database.DB.First(&investment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	c.JSON(http.StatusOK, investment)
}

// ListInvestments retrieves all investments for a user
func ListInvestments(c *gin.Context) {
	userID := c.GetUint("user_id")
	var investments []models.Investment

	if err := database.DB.Where("user_id = ?", userID).Find(&investments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch investments"})
		return
	}

	c.JSON(http.StatusOK, investments)
}

// UpdateInvestment updates an existing investment
func UpdateInvestment(c *gin.Context) {
	id := c.Param("id")
	var investment models.Investment

	if err := database.DB.First(&investment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	if err := c.ShouldBindJSON(&investment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&investment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update investment"})
		return
	}

	c.JSON(http.StatusOK, investment)
}

// DeleteInvestment deletes an investment
func DeleteInvestment(c *gin.Context) {
	id := c.Param("id")
	var investment models.Investment

	if err := database.DB.First(&investment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	if err := database.DB.Delete(&investment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete investment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Investment deleted successfully"})
}

// CreateTransaction creates a new transaction
func CreateTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default values
	transaction.Timestamp = time.Now()
	transaction.TransactionID = uuid.New().String()
	transaction.Status = "PENDING"

	// Start transaction
	tx := database.DB.Begin()

	// Create transaction record
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	// Update investment based on transaction type
	var investment models.Investment
	if err := tx.First(&investment, transaction.InvestmentID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	if transaction.Type == "SELL" {
		investment.Status = "SOLD"
		investment.SellDate = &transaction.Timestamp
		investment.SellPrice = &transaction.Price
	}

	if err := tx.Save(&investment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update investment"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// GetPortfolio retrieves a portfolio by ID
func GetPortfolio(c *gin.Context) {
	id := c.Param("id")
	var portfolio models.Portfolio

	if err := database.DB.Preload("Investments").First(&portfolio, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

// CreatePortfolio creates a new portfolio
func CreatePortfolio(c *gin.Context) {
	var portfolio models.Portfolio
	if err := c.ShouldBindJSON(&portfolio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default values
	portfolio.LastUpdated = time.Now()

	if err := database.DB.Create(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create portfolio"})
		return
	}

	c.JSON(http.StatusCreated, portfolio)
}

// UpdatePortfolio updates an existing portfolio
func UpdatePortfolio(c *gin.Context) {
	id := c.Param("id")
	var portfolio models.Portfolio

	if err := database.DB.First(&portfolio, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
		return
	}

	if err := c.ShouldBindJSON(&portfolio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	portfolio.LastUpdated = time.Now()

	if err := database.DB.Save(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update portfolio"})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

// DeletePortfolio deletes a portfolio
func DeletePortfolio(c *gin.Context) {
	id := c.Param("id")
	var portfolio models.Portfolio

	if err := database.DB.First(&portfolio, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
		return
	}

	if err := database.DB.Delete(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete portfolio"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Portfolio deleted successfully"})
}
