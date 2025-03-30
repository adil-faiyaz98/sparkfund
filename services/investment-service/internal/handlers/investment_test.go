package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"investment-service/internal/database"
	"investment-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type InvestmentHandlerTestSuite struct {
	suite.Suite
	router *gin.Engine
	db     *gorm.DB
}

func (suite *InvestmentHandlerTestSuite) SetupSuite() {
	// Use in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		suite.T().Fatal(err)
	}

	// Migrate models
	err = db.AutoMigrate(&models.Portfolio{}, &models.Investment{}, &models.Transaction{})
	if err != nil {
		suite.T().Fatal(err)
	}

	// Set DB for tests
	database.DB = db
	suite.db = db

	// Setup test router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// Register routes for testing
	r.POST("/investments", CreateInvestment)
	r.GET("/investments/:id", GetInvestment)
	r.GET("/investments", ListInvestments)
	r.PUT("/investments/:id", UpdateInvestment)
	r.DELETE("/investments/:id", DeleteInvestment)

	suite.router = r
}

func (suite *InvestmentHandlerTestSuite) TearDownSuite() {
	// Close DB connection
	sqlDB, err := suite.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func (suite *InvestmentHandlerTestSuite) SetupTest() {
	// Clean up tables between tests
	suite.db.Where("1 = 1").Delete(&models.Investment{})
	suite.db.Where("1 = 1").Delete(&models.Portfolio{})
}

func (suite *InvestmentHandlerTestSuite) TestCreateInvestment() {
	// Create a portfolio first
	portfolio := models.Portfolio{
		UserID:      1,
		Name:        "Test Portfolio",
		Description: "Test Description",
		TotalValue:  0,
		LastUpdated: time.Now(),
	}
	result := suite.db.Create(&portfolio)
	assert.NoError(suite.T(), result.Error)

	// Test investment payload
	investment := models.Investment{
		UserID:        1,
		PortfolioID:   portfolio.ID,
		Amount:        1000.0,
		Type:          "STOCK",
		Status:        "ACTIVE",
		PurchaseDate:  time.Now(),
		PurchasePrice: 150.50,
		Symbol:        "AAPL",
		Quantity:      6.64,
		Notes:         "Test investment",
	}

	// Convert to JSON
	jsonValue, _ := json.Marshal(investment)

	// Create request
	req := httptest.NewRequest("POST", "/investments", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	suite.router.ServeHTTP(w, req)

	// Parse response
	var response models.Investment
	err := json.Unmarshal(w.Body.Bytes(), &response)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	assert.NotZero(suite.T(), response.ID)
	assert.Equal(suite.T(), investment.UserID, response.UserID)
	assert.Equal(suite.T(), investment.Amount, response.Amount)
	assert.Equal(suite.T(), investment.Type, response.Type)
	assert.Equal(suite.T(), investment.Symbol, response.Symbol)
}

func (suite *InvestmentHandlerTestSuite) TestGetInvestment() {
	// Create a portfolio first
	portfolio := models.Portfolio{
		UserID:      1,
		Name:        "Test Portfolio",
		Description: "Test Description",
		TotalValue:  0,
		LastUpdated: time.Now(),
	}
	suite.db.Create(&portfolio)

	// Create test investment
	investment := models.Investment{
		UserID:        1,
		PortfolioID:   portfolio.ID,
		Amount:        1000.0,
		Type:          "STOCK",
		Status:        "ACTIVE",
		PurchaseDate:  time.Now(),
		PurchasePrice: 150.50,
		Symbol:        "AAPL",
		Quantity:      6.64,
		Notes:         "Test investment",
	}
	result := suite.db.Create(&investment)
	assert.NoError(suite.T(), result.Error)

	// Create request
	req := httptest.NewRequest("GET", "/investments/"+string(rune(investment.ID)), nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	suite.router.ServeHTTP(w, req)

	// Parse response
	var response models.Investment
	err := json.Unmarshal(w.Body.Bytes(), &response)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), investment.ID, response.ID)
	assert.Equal(suite.T(), investment.Symbol, response.Symbol)
	assert.Equal(suite.T(), investment.Amount, response.Amount)
}

// Additional test methods for other endpoints...

func TestInvestmentHandlerSuite(t *testing.T) {
	suite.Run(t, new(InvestmentHandlerTestSuite))
}
