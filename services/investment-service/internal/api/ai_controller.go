package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sparkfund/services/investment-service/internal/ai/fraud"
	"github.com/sparkfund/services/investment-service/internal/ai/recommendation"
	"github.com/sparkfund/services/investment-service/internal/ai/security"
)

// AIController handles AI-related API endpoints
type AIController struct {
	recommendationService *recommendation.Service
	securityService       *security.SecurityService
}

// NewAIController creates a new AI controller
func NewAIController(recommendationService *recommendation.Service, securityService *security.SecurityService) *AIController {
	return &AIController{
		recommendationService: recommendationService,
		securityService:       securityService,
	}
}

// RegisterRoutes registers the AI controller routes
func (c *AIController) RegisterRoutes(router *gin.Engine) {
	aiGroup := router.Group("/api/v1/ai")
	{
		// Recommendation endpoints
		aiGroup.POST("/recommendations", c.GetRecommendation)
		aiGroup.GET("/recommendations/personalized/:userId", c.GetPersonalizedAssets)
		aiGroup.GET("/recommendations/similar/:assetId", c.GetSimilarAssets)
		aiGroup.GET("/recommendations/diversification/:userId", c.GetDiversificationSuggestions)
		aiGroup.GET("/recommendations/rebalancing/:userId", c.GetRebalancingSuggestions)
		
		// Security endpoints
		aiGroup.POST("/security/analyze", c.AnalyzeTransaction)
	}
}

// RecommendationRequest represents a request for investment recommendations
type RecommendationRequest struct {
	UserID           string   `json:"user_id" binding:"required"`
	Amount           float64  `json:"amount" binding:"required,gt=0"`
	TimeHorizon      int      `json:"time_horizon" binding:"required,gt=0"`
	RiskTolerance    float64  `json:"risk_tolerance" binding:"required,min=0,max=1"`
	PreferredSectors []string `json:"preferred_sectors"`
	ExcludedSectors  []string `json:"excluded_sectors"`
	Goals            []string `json:"goals"`
	Constraints      []string `json:"constraints"`
}

// GetRecommendation handles the recommendation request
func (c *AIController) GetRecommendation(ctx *gin.Context) {
	var request RecommendationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Convert to service request
	serviceRequest := recommendation.RecommendationRequest{
		UserID:           request.UserID,
		Amount:           request.Amount,
		TimeHorizon:      request.TimeHorizon,
		RiskTolerance:    request.RiskTolerance,
		PreferredSectors: request.PreferredSectors,
		ExcludedSectors:  request.ExcludedSectors,
		Goals:            request.Goals,
		Constraints:      request.Constraints,
	}
	
	// Get recommendation
	recommendation, err := c.recommendationService.GetRecommendation(ctx, serviceRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, recommendation)
}

// GetPersonalizedAssets handles the personalized assets request
func (c *AIController) GetPersonalizedAssets(ctx *gin.Context) {
	userID := ctx.Param("userId")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}
	
	count := 5 // Default count
	if countParam := ctx.Query("count"); countParam != "" {
		if _, err := ctx.GetQuery("count"); err {
			count = ctx.GetInt("count")
		}
	}
	
	// Get personalized assets
	assets, err := c.recommendationService.GetPersonalizedAssets(ctx, userID, count)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, assets)
}

// GetSimilarAssets handles the similar assets request
func (c *AIController) GetSimilarAssets(ctx *gin.Context) {
	assetID := ctx.Param("assetId")
	if assetID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "asset ID is required"})
		return
	}
	
	count := 5 // Default count
	if countParam := ctx.Query("count"); countParam != "" {
		if _, err := ctx.GetQuery("count"); err {
			count = ctx.GetInt("count")
		}
	}
	
	// Get similar assets
	assets, err := c.recommendationService.GetSimilarAssets(ctx, assetID, count)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, assets)
}

// GetDiversificationSuggestions handles the diversification suggestions request
func (c *AIController) GetDiversificationSuggestions(ctx *gin.Context) {
	userID := ctx.Param("userId")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}
	
	// Get diversification suggestions
	suggestions, err := c.recommendationService.GetDiversificationSuggestions(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, suggestions)
}

// GetRebalancingSuggestions handles the rebalancing suggestions request
func (c *AIController) GetRebalancingSuggestions(ctx *gin.Context) {
	userID := ctx.Param("userId")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}
	
	// Get rebalancing suggestions
	suggestions, err := c.recommendationService.GetRebalancingSuggestions(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, suggestions)
}

// TransactionAnalysisRequest represents a request to analyze a transaction
type TransactionAnalysisRequest struct {
	UserID          string    `json:"user_id" binding:"required"`
	Amount          float64   `json:"amount" binding:"required,gt=0"`
	Currency        string    `json:"currency" binding:"required"`
	TransactionType string    `json:"transaction_type" binding:"required"`
	AssetID         string    `json:"asset_id"`
	IPAddress       string    `json:"ip_address"`
	DeviceID        string    `json:"device_id"`
	Location        Location  `json:"location"`
	UserAgent       string    `json:"user_agent"`
}

// Location represents geographical coordinates
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
	City      string  `json:"city"`
}

// AnalyzeTransaction handles the transaction analysis request
func (c *AIController) AnalyzeTransaction(ctx *gin.Context) {
	var request TransactionAnalysisRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Convert to service request
	transaction := fraud.Transaction{
		ID:              uuid.New().String(),
		UserID:          request.UserID,
		Amount:          request.Amount,
		Currency:        request.Currency,
		TransactionType: request.TransactionType,
		AssetID:         request.AssetID,
		Timestamp:       time.Now(),
		IPAddress:       request.IPAddress,
		DeviceID:        request.DeviceID,
		Location: fraud.Location{
			Latitude:  request.Location.Latitude,
			Longitude: request.Location.Longitude,
			Country:   request.Location.Country,
			City:      request.Location.City,
		},
		UserAgent: request.UserAgent,
	}
	
	// Analyze transaction
	result, err := c.securityService.AnalyzeTransaction(ctx, transaction)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, result)
}
