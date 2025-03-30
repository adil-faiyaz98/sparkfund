package handlers

import (
	"net/http"
	"os/exec"
	"encoding/json"
	"regexp"
	"time"

	"github.com/sparkfund/services/investment-service/internal/database"
	"github.com/sparkfund/services/investment-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Create investment
// @Summary      Create a new investment
// @Description  Create a new investment
// @Tags         investments

type InvestmentHandler struct {
}

func NewInvestmentHandler() *InvestmentHandler {
	return &InvestmentHandler{}
}

func (h *InvestmentHandler) RegisterRoutes(r *gin.Engine) {
	investments := r.Group("/investments")
	{
		investments.POST("", CreateInvestment)
		investments.GET("/:id", GetInvestment)
		investments.GET("", ListInvestments)
		investments.PUT("/:id", UpdateInvestment)
		investments.DELETE("/:id", DeleteInvestment)
	}

	transactions := r.Group("/transactions")
	{
		transactions.POST("", CreateTransaction)
	}

	portfolios := r.Group("/portfolios")
	{
		portfolios.POST("", CreatePortfolio)
		portfolios.GET("/:id", GetPortfolio)
		portfolios.PUT("/:id", UpdatePortfolio)
		portfolios.DELETE("/:id", DeletePortfolio)
	}
	r.POST("/stock-recommendation", h.GetStockRecommendation)
}
// @Accept       json
// @Produce      json
// @Param        investment  body      models.Investment  true  "Investment data"
// @Success      201         {object}  models.Investment
// @Failure      400         {object}  models.ErrorResponse
// @Failure      500         {object}  models.ErrorResponse
// @Router       /investments [post]
func CreateInvestment(c *gin.Context) {
	var investment models.Investment
	if err := c.ShouldBindJSON(&investment); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Validate required fields
	if investment.UserID == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "user_id is required"})
		return
	}

	if investment.PortfolioID == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "portfolio_id is required"})
		return
	}

	if investment.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "quantity must be greater than 0"})
		return
	}

	if investment.Type == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "type is required"})
		return
	}

	if investment.Symbol == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "symbol is required"})
		return
	}

	// Validate type enum
	validTypes := map[string]bool{
		"STOCK": true, "CRYPTO": true, "REAL_ESTATE": true,
		"ETF": true, "BOND": true, "MUTUAL_FUND": true,
	}
	if !validTypes[investment.Type] {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid investment type"})
		return
	}

	// Set default values
	now := time.Now()
	investment.CreatedAt = now
	investment.UpdatedAt = now
	investment.PurchaseDate = now
	investment.Status = "ACTIVE"

	// Create investment
	if err := database.DB.Create(&investment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create investment"})
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
// @Failure      401  {object}  models.ErrorResponse  "Unauthorized"
// @Failure      403  {object}  models.ErrorResponse  "Forbidden"
// @Failure      404  {object}  models.ErrorResponse  "Not found"
// @Failure      500  {object}  models.ErrorResponse  "Internal server error"
// @Router       /investments/{id} [get]
// @Example      response
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
// @Example        "notes": "Initial purchase of Apple stock",
// @Example        "created_at": "2025-03-28T12:00:00Z",
// @Example        "updated_at": "2025-03-28T12:00:00Z",
// @Example        "sell_date": null,
// @Example        "sell_price": null
// @Example      }
func GetInvestment(c *gin.Context) {
	id := c.Param("id")
	var investment models.Investment

	if err := database.DB.First(&investment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Investment not found"})
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
// @Failure      401  {object}  models.ErrorResponse  "Unauthorized"
// @Failure      403  {object}  models.ErrorResponse  "Forbidden"
// @Failure      500  {object}  models.ErrorResponse  "Internal server error"
// @Router       /investments [get]
// @Example      response
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
// @Example          "notes": "Initial purchase of Apple stock",
// @Example          "created_at": "2025-03-28T12:00:00Z",
// @Example          "updated_at": "2025-03-28T12:00:00Z",
// @Example          "sell_date": null,
// @Example          "sell_price": null
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
// @Example          "notes": "Initial purchase of Google stock",
// @Example          "created_at": "2025-03-28T12:00:00Z",
// @Example          "updated_at": "2025-03-28T12:00:00Z",
// @Example          "sell_date": null,
// @Example          "sell_price": null
// @Example        }
// @Example      ]
func ListInvestments(c *gin.Context) {
	userID := c.GetUint("user_id")
	var investments []models.Investment

	if err := database.DB.Where("user_id = ?", userID).Find(&investments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch investments"})
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
// @Failure      400          {object}  models.ErrorResponse  "Bad request"
// @Failure      401          {object}  models.ErrorResponse  "Unauthorized"
// @Failure      403          {object}  models.ErrorResponse  "Forbidden"
// @Failure      404          {object}  models.ErrorResponse  "Not found"
// @Failure      500          {object}  models.ErrorResponse  "Internal server error"
// @Router       /investments/{id} [put]
// @Example      request
// @Example      {
// @Example        "user_id": 1,
// @Example        "portfolio_id": 1,
// @Example        "amount": 1000.00,
// @Example        "type": "STOCK",
// @Example        "status": "ACTIVE",
// @Example        "purchase_price": 150.50,
// @Example        "symbol": "AAPL",
// @Example        "quantity": 10,
// @Example        "notes": "Updated notes for Apple stock"
// @Example      }
// @Example      response
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
// @Example        "notes": "Updated notes for Apple stock",
// @Example        "created_at": "2025-03-28T12:00:00Z",
// @Example        "updated_at": "2025-03-29T10:30:00Z",
// @Example        "sell_date": null,
// @Example        "sell_price": null
// @Example      }
func UpdateInvestment(c *gin.Context) {
	id := c.Param("id")
	var investment models.Investment

	if err := database.DB.First(&investment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Investment not found"})
		return
	}

	if err := c.ShouldBindJSON(&investment); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	if err := database.DB.Save(&investment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update investment"})
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
// @Success      200  {object}  SuccessResponse  "Successfully deleted"
// @Failure      401  {object}  models.ErrorResponse    "Unauthorized"
// @Failure      403  {object}  models.ErrorResponse    "Forbidden"
// @Failure      404  {object}  models.ErrorResponse    "Not found"
// @Failure      500  {object}  models.ErrorResponse    "Internal server error"
// @Router       /investments/{id} [delete]
// @Example      response
// @Example      {
// @Example        "message": "Investment deleted successfully"
// @Example      }
func DeleteInvestment(c *gin.Context) {
	id := c.Param("id")
	var investment models.Investment

	if err := database.DB.First(&investment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Investment not found"})
		return
	}

	if err := database.DB.Delete(&investment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete investment"})
		return
	}

	c.JSON(200, models.SuccessResponse{Message: "Successfully deleted"})
}

func (h *InvestmentHandler) GetStockRecommendation(c *gin.Context) {
	// Execute the Python script
	cmd := exec.Command("python", "scripts/stock_picking_model.py")

	// Capture the standard output
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to execute python script"})
		return
	}

	// Parse the output using regular expressions
	re := regexp.MustCompile(`\('(\d{4}-\d{2}-\d{2})',\s*([\d.-]+)\)`)
	matches := re.FindAllStringSubmatch(string(output), -1)

	// Define a struct for the recommendation
	type Recommendation struct {
		Date       string  `json:"date"`
		Prediction float64 `json:"prediction"`
	}

	// Create a slice to store the recommendations
	var recommendations []Recommendation

	// Iterate through the matches and create the recommendation objects
	for _, match := range matches {
		if len(match) == 3 {
			date := match[1]
			prediction, err := strconv.ParseFloat(match[2], 64)
			if err != nil {
				// Handle error parsing prediction
				continue
			}

			recommendations = append(recommendations, Recommendation{
				Date:       date,
				Prediction: prediction,
			})
		}
	}

	//Return and empty array if not predictions found
	if len(matches) == 0 {
		
		c.JSON(http.StatusOK, []Recommendation{})
		return

	}
	// Marshal the recommendations to JSON
	response, err := json.Marshal(recommendations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to marshal recommendations to JSON"})
		return
	}

	// Return the structured JSON response
	c.Data(http.StatusOK, "application/json", response)
}





// CreateTransaction godoc
// @Summary      Create a new transaction
// @Description  Create a new transaction for an investment
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        transaction  body      models.Transaction  true  "Transaction object"
// @Success      201          {object}  models.Transaction
// @Failure      400          {object}  models.ErrorResponse  "Bad request"
// @Failure      401          {object}  models.ErrorResponse  "Unauthorized"
// @Failure      403          {object}  models.ErrorResponse  "Forbidden"
// @Failure      404          {object}  models.ErrorResponse  "Not found"
// @Failure      500          {object}  models.ErrorResponse  "Internal server error"
// @Router       /transactions [post]
// @Example      request
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
// @Example      response
// @Example      {
// @Example        "id": 1,
// @Example        "user_id": 1,
// @Example        "investment_id": 1,
// @Example        "transaction_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
// @Example        "type": "BUY",
// @Example        "amount": 1000.00,
// @Example        "price": 150.50,
// @Example        "quantity": 10,
// @Example        "timestamp": "2025-03-28T12:00:00Z",
// @Example        "status": "PENDING",
// @Example        "created_at": "2025-03-28T12:00:00Z",
// @Example        "updated_at": "2025-03-28T12:00:00Z"
// @Example      }
func CreateTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Validate required fields
	if transaction.UserID == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "user_id is required"})
		return
	}

	if transaction.InvestmentID == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "investment_id is required"})
		return
	}

	if transaction.Type == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "type is required"})
		return
	}

	if transaction.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "quantity must be greater than 0"})
		return
	}

	// Validate type enum
	if transaction.Type != "BUY" && transaction.Type != "SELL" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "transaction type must be either BUY or SELL"})
		return
	}

	// Set default values
	now := time.Now()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now
	transaction.Timestamp = now
	transaction.TransactionID = uuid.New().String()
	transaction.Status = "PENDING"

	// Start transaction
	tx := database.DB.Begin()

	// Create transaction record
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create transaction"})
		return
	}

	// Update investment based on transaction type
	var investment models.Investment
	if err := tx.First(&investment, transaction.InvestmentID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Investment not found"})
		return
	}

	if transaction.Type == "SELL" {
		investment.Status = "SOLD"
		investment.SellDate = &transaction.Timestamp
		investment.SellPrice = &transaction.Price
	}

	if err := tx.Save(&investment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update investment"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to commit transaction"})
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
// @Failure      404  {object}  models.ErrorResponse  "Not found"
// @Failure      500  {object}  models.ErrorResponse  "Internal server error"
// @Router       /portfolios/{id} [get]
// @Example      response
// @Example      {
// @Example        "id": 1,
// @Example        "user_id": 1,
// @Example        "name": "Retirement Portfolio",
// @Example        "description": "Long term investments for retirement",
// @Example        "created_at": "2025-03-28T12:00:00Z",
// @Example        "updated_at": "2025-03-28T12:00:00Z",
// @Example        "last_updated": "2025-03-28T12:00:00Z",
// @Example        "investments": [
// @Example          {
// @Example            "id": 1,
// @Example            "user_id": 1,
// @Example            "portfolio_id": 1,
// @Example            "amount": 1000.00,
// @Example            "type": "STOCK",
// @Example            "status": "ACTIVE",
// @Example            "purchase_date": "2025-03-28T12:00:00Z",
// @Example            "purchase_price": 150.50,
// @Example            "symbol": "AAPL",
// @Example            "quantity": 10,
// @Example            "notes": "Initial purchase of Apple stock"
// @Example          }
// @Example        ]
// @Example      }
func GetPortfolio(c *gin.Context) {
	id := c.Param("id")
	var portfolio models.Portfolio

	if err := database.DB.Preload("Investments").First(&portfolio, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Portfolio not found"})
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
// @Failure      400        {object}  models.ErrorResponse  "Bad request"
// @Failure      500        {object}  models.ErrorResponse  "Internal server error"
// @Router       /portfolios [post]
// @Example      request
// @Example      {
// @Example        "user_id": 1,
// @Example        "name": "Tech Stocks",
// @Example        "description": "Portfolio focusing on technology sector"
// @Example      }
// @Example      response
// @Example      {
// @Example        "id": 1,
// @Example        "user_id": 1,
// @Example        "name": "Tech Stocks",
// @Example        "description": "Portfolio focusing on technology sector",
// @Example        "created_at": "2025-03-28T12:00:00Z",
// @Example        "updated_at": "2025-03-28T12:00:00Z",
// @Example        "last_updated": "2025-03-28T12:00:00Z"
// @Example      }
func CreatePortfolio(c *gin.Context) {
	var portfolio models.Portfolio
	if err := c.ShouldBindJSON(&portfolio); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Validate required fields
	if portfolio.UserID == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "user_id is required"})
		return
	}

	if portfolio.Name == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "name is required"})
		return
	}

	// Set default values
	now := time.Now()
	portfolio.CreatedAt = now
	portfolio.UpdatedAt = now
	portfolio.LastUpdated = now

	if err := database.DB.Create(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create portfolio"})
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
// @Failure      400        {object}  models.ErrorResponse  "Bad request"
// @Failure      404        {object}  models.ErrorResponse  "Not found"
// @Failure      500        {object}  models.ErrorResponse  "Internal server error"
// @Router       /portfolios/{id} [put]
// @Example      request
// @Example      {
// @Example        "name": "Updated Tech Portfolio",
// @Example        "description": "Updated portfolio description"
// @Example      }
// @Example      response
// @Example      {
// @Example        "id": 1,
// @Example        "user_id": 1,
// @Example        "name": "Updated Tech Portfolio",
// @Example        "description": "Updated portfolio description",
// @Example        "created_at": "2025-03-28T12:00:00Z",
// @Example        "updated_at": "2025-03-29T10:30:00Z",
// @Example        "last_updated": "2025-03-29T10:30:00Z"
// @Example      }
func UpdatePortfolio(c *gin.Context) {
	id := c.Param("id")
	var portfolio models.Portfolio

	if err := database.DB.First(&portfolio, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Portfolio not found"})
		return
	}

	if err := c.ShouldBindJSON(&portfolio); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	portfolio.LastUpdated = time.Now()

	if err := database.DB.Save(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update portfolio"})
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
// @Success      200  {object}  SuccessResponse  "Successfully deleted"
// @Failure      404  {object}  models.ErrorResponse    "Not found"
// @Failure      500  {object}  models.ErrorResponse    "Internal server error"
// @Router       /portfolios/{id} [delete]
// @Example      response
// @Example      {
// @Example        "message": "Portfolio deleted successfully"
// @Example      }
func DeletePortfolio(c *gin.Context) {
	id := c.Param("id")
	var portfolio models.Portfolio

	if err := database.DB.First(&portfolio, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Portfolio not found"})
		return
	}

	if err := database.DB.Delete(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete portfolio"})
		return
	}

	c.JSON(200, models.SuccessResponse{Message: "Successfully deleted"})
}
