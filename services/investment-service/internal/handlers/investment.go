package handlers

import (
	"net/http"
	"time"

	"investment-service/internal/database"
	"investment-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateInvestment godoc
// @Summary      Create a new investment
// @Description  Create a new investment with the provided details
// @Tags         investments
// @Accept       json
// @Produce      json
// @Param        investment  body      models.Investment  true  "Investment object"
// @Success      201         {object}  models.Investment
// @Failure      400         {object}  map[string]string
// @Failure      401         {object}  map[string]string
// @Failure      403         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /investments [post]
// @Example      {object}  models.Investment
// @Example      {
// @Example        "user_id": 1,
// @Example        "portfolio_id": 1,
// @Example        "amount": 1000.00,
// @Example        "type": "STOCK",
// @Example        "status": "ACTIVE",
// @Example        "purchase_date": "2025-03-28T12:00:00Z",
// @Example        "purchase_price": 150.50,
// @Example        "symbol": "AAPL",
// @Example        "quantity": 10,
// @Example        "notes": "Initial purchase of Apple stock"
// @Example      }
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

// GetInvestment godoc
// @Summary      Get an investment by ID
// @Description  Get investment details by ID
// @Tags         investments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Investment ID"
// @Success      200  {object}  models.Investment
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /investments/{id} [get]
// @Example      {object}  models.Investment
// @Example      {
// @Example        "id": 1,
// @Example        "user_id": 1,
// @Example        "portfolio_id": 1,
// @Example        "amount": 1000.00,
// @Example        "type": "STOCK",
// @Example        "status": "ACTIVE",
// @Example        "purchase_date": "2025-03-28T12:00:00Z",
// @Example        "purchase_price": 150.50,
// @Example        "symbol": "AAPL",
// @Example        "quantity": 10,
// @Example        "notes": "Initial purchase of Apple stock"
// @Example      }
func GetInvestment(c *gin.Context) {
	id := c.Param("id")
	var investment models.Investment

	if err := database.DB.First(&investment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Investment not found"})
		return
	}

	c.JSON(http.StatusOK, investment)
}

// ListInvestments godoc
// @Summary      List all investments
// @Description  Get a list of all investments
// @Tags         investments
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Investment
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /investments [get]
// @Example      {array}  models.Investment
// @Example      [
// @Example        {
// @Example          "id": 1,
// @Example          "user_id": 1,
// @Example          "portfolio_id": 1,
// @Example          "amount": 1000.00,
// @Example          "type": "STOCK",
// @Example          "status": "ACTIVE",
// @Example          "purchase_date": "2025-03-28T12:00:00Z",
// @Example          "purchase_price": 150.50,
// @Example          "symbol": "AAPL",
// @Example          "quantity": 10,
// @Example          "notes": "Initial purchase of Apple stock"
// @Example        },
// @Example        {
// @Example          "id": 2,
// @Example          "user_id": 1,
// @Example          "portfolio_id": 1,
// @Example          "amount": 2000.00,
// @Example          "type": "STOCK",
// @Example          "status": "ACTIVE",
// @Example          "purchase_date": "2025-03-28T12:00:00Z",
// @Example          "purchase_price": 280.75,
// @Example          "symbol": "GOOGL",
// @Example          "quantity": 5,
// @Example          "notes": "Initial purchase of Google stock"
// @Example        }
// @Example      ]
func ListInvestments(c *gin.Context) {
	userID := c.GetUint("user_id")
	var investments []models.Investment

	if err := database.DB.Where("user_id = ?", userID).Find(&investments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch investments"})
		return
	}

	c.JSON(http.StatusOK, investments)
}

// UpdateInvestment godoc
// @Summary      Update an investment
// @Description  Update investment details by ID
// @Tags         investments
// @Accept       json
// @Produce      json
// @Param        id           path      int             true  "Investment ID"
// @Param        investment   body      models.Investment  true  "Updated investment object"
// @Success      200          {object}  models.Investment
// @Failure      400          {object}  map[string]string
// @Failure      401          {object}  map[string]string
// @Failure      403          {object}  map[string]string
// @Failure      404          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Router       /investments/{id} [put]
// @Example      {object}  models.Investment
// @Example      {
// @Example        "id": 1,
// @Example        "user_id": 1,
// @Example        "portfolio_id": 1,
// @Example        "amount": 1000.00,
// @Example        "type": "STOCK",
// @Example        "status": "ACTIVE",
// @Example        "purchase_date": "2025-03-28T12:00:00Z",
// @Example        "purchase_price": 150.50,
// @Example        "symbol": "AAPL",
// @Example        "quantity": 10,
// @Example        "notes": "Updated notes for Apple stock"
// @Example      }
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

// DeleteInvestment godoc
// @Summary      Delete an investment
// @Description  Delete an investment by ID
// @Tags         investments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Investment ID"
// @Success      204  "No Content"
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /investments/{id} [delete]
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

// CreateTransaction godoc
// @Summary      Create a new transaction
// @Description  Create a new transaction for an investment
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        transaction  body      models.Transaction  true  "Transaction object"
// @Success      201          {object}  models.Transaction
// @Failure      400          {object}  map[string]string
// @Failure      401          {object}  map[string]string
// @Failure      403          {object}  map[string]string
// @Failure      404          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Router       /transactions [post]
// @Example      {object}  models.Transaction
// @Example      {
// @Example        "user_id": 1,
// @Example        "investment_id": 1,
// @Example        "type": "BUY",
// @Example        "amount": 1000.00,
// @Example        "price": 150.50,
// @Example        "quantity": 10,
// @Example        "timestamp": "2025-03-28T12:00:00Z",
// @Example        "status": "PENDING"
// @Example      }
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

// GetPortfolio godoc
// @Summary      Get a portfolio by ID
// @Description  Get portfolio details by ID including its investments
// @Tags         portfolios
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Portfolio ID"
// @Success      200  {object}  models.Portfolio
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /portfolios/{id} [get]
func GetPortfolio(c *gin.Context) {
	id := c.Param("id")
	var portfolio models.Portfolio

	if err := database.DB.Preload("Investments").First(&portfolio, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

// CreatePortfolio godoc
// @Summary      Create a new portfolio
// @Description  Create a new portfolio with the provided details
// @Tags         portfolios
// @Accept       json
// @Produce      json
// @Param        portfolio  body      models.Portfolio  true  "Portfolio object"
// @Success      201        {object}  models.Portfolio
// @Failure      400        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Router       /portfolios [post]
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

// UpdatePortfolio godoc
// @Summary      Update a portfolio
// @Description  Update portfolio details by ID
// @Tags         portfolios
// @Accept       json
// @Produce      json
// @Param        id         path      string         true  "Portfolio ID"
// @Param        portfolio  body      models.Portfolio  true  "Updated portfolio object"
// @Success      200        {object}  models.Portfolio
// @Failure      400        {object}  map[string]string
// @Failure      404        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Router       /portfolios/{id} [put]
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

// DeletePortfolio godoc
// @Summary      Delete a portfolio
// @Description  Delete a portfolio by ID
// @Tags         portfolios
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Portfolio ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /portfolios/{id} [delete]
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
