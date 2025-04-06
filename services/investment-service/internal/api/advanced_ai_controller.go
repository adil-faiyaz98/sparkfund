package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/services/investment-service/internal/ai"
	"github.com/sparkfund/services/investment-service/internal/ai/nlp"
	"github.com/sparkfund/services/investment-service/internal/ai/rl"
	"github.com/sparkfund/services/investment-service/internal/ai/timeseries"
)

// AdvancedAIController handles advanced AI-related API endpoints
type AdvancedAIController struct {
	aiService *ai.AIService
}

// NewAdvancedAIController creates a new advanced AI controller
func NewAdvancedAIController(aiService *ai.AIService) *AdvancedAIController {
	return &AdvancedAIController{
		aiService: aiService,
	}
}

// RegisterRoutes registers the advanced AI controller routes
func (c *AdvancedAIController) RegisterRoutes(router *gin.Engine) {
	aiGroup := router.Group("/api/v1/ai/advanced")
	{
		// NLP endpoints
		aiGroup.POST("/news/analyze", c.AnalyzeNews)
		
		// Time series endpoints
		aiGroup.POST("/price/forecast", c.ForecastPrice)
		aiGroup.POST("/market/predict", c.PredictMarket)
		
		// RL endpoints
		aiGroup.POST("/portfolio/optimize", c.OptimizePortfolio)
		aiGroup.POST("/portfolio/action", c.GetPortfolioAction)
		
		// Automated investment endpoints
		aiGroup.POST("/investment/automated", c.GenerateAutomatedInvestment)
	}
}

// NewsArticleRequest represents a request to analyze a news article
type NewsArticleRequest struct {
	Article nlp.NewsArticle `json:"article" binding:"required"`
}

// AnalyzeNews handles the news analysis request
func (c *AdvancedAIController) AnalyzeNews(ctx *gin.Context) {
	var request NewsArticleRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Analyze news
	result, err := c.aiService.AnalyzeNews(ctx, request.Article)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, result)
}

// PriceForecastRequest represents a request to forecast prices
type PriceForecastRequest struct {
	Symbol           string                 `json:"symbol" binding:"required"`
	HistoricalPrices []timeseries.PricePoint `json:"historical_prices" binding:"required"`
	Days             int                    `json:"days" binding:"required,gt=0"`
}

// ForecastPrice handles the price forecast request
func (c *AdvancedAIController) ForecastPrice(ctx *gin.Context) {
	var request PriceForecastRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Forecast prices
	forecast, err := c.aiService.ForecastPrice(ctx, request.Symbol, request.HistoricalPrices, request.Days)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, forecast)
}

// MarketPredictionRequest represents a request to predict market movements
type MarketPredictionRequest struct {
	Symbol           string                     `json:"symbol" binding:"required"`
	HistoricalPrices []timeseries.PricePoint     `json:"historical_prices" binding:"required"`
	Indicators       []timeseries.MarketIndicator `json:"indicators"`
}

// PredictMarket handles the market prediction request
func (c *AdvancedAIController) PredictMarket(ctx *gin.Context) {
	var request MarketPredictionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Predict market
	prediction, err := c.aiService.PredictMarket(ctx, request.Symbol, request.HistoricalPrices, request.Indicators)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, prediction)
}

// PortfolioOptimizationRequest represents a request to optimize a portfolio
type PortfolioOptimizationRequest struct {
	Portfolio       rl.Portfolio `json:"portfolio" binding:"required"`
	AvailableAssets []rl.Asset   `json:"available_assets" binding:"required"`
}

// OptimizePortfolio handles the portfolio optimization request
func (c *AdvancedAIController) OptimizePortfolio(ctx *gin.Context) {
	var request PortfolioOptimizationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Optimize portfolio
	result, err := c.aiService.OptimizePortfolio(ctx, request.Portfolio, request.AvailableAssets)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, result)
}

// PortfolioActionRequest represents a request to get a portfolio action
type PortfolioActionRequest struct {
	State rl.State `json:"state" binding:"required"`
}

// GetPortfolioAction handles the portfolio action request
func (c *AdvancedAIController) GetPortfolioAction(ctx *gin.Context) {
	var request PortfolioActionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get action
	action := c.aiService.GetRLAction(ctx, request.State)
	
	ctx.JSON(http.StatusOK, action)
}

// AutomatedInvestmentRequest represents a request to generate an automated investment
type AutomatedInvestmentRequest struct {
	UserID           string                      `json:"user_id" binding:"required"`
	NewsAnalysisID   string                      `json:"news_analysis_id,omitempty"`
	NewsArticle      *nlp.NewsArticle             `json:"news_article,omitempty"`
	Symbol           string                      `json:"symbol" binding:"required"`
	HistoricalPrices []timeseries.PricePoint      `json:"historical_prices,omitempty"`
	Indicators       []timeseries.MarketIndicator `json:"indicators,omitempty"`
}

// GenerateAutomatedInvestment handles the automated investment request
func (c *AdvancedAIController) GenerateAutomatedInvestment(ctx *gin.Context) {
	var request AutomatedInvestmentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Process news analysis
	var newsAnalysis *nlp.NewsAnalysisResult
	if request.NewsArticle != nil {
		var err error
		newsAnalysis, err = c.aiService.AnalyzeNews(ctx, *request.NewsArticle)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze news: " + err.Error()})
			return
		}
	}
	
	// Process price forecast
	var priceForecast *timeseries.Forecast
	if len(request.HistoricalPrices) > 0 {
		var err error
		priceForecast, err = c.aiService.ForecastPrice(ctx, request.Symbol, request.HistoricalPrices, 30)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forecast prices: " + err.Error()})
			return
		}
	}
	
	// Process market prediction
	var marketPrediction *timeseries.MarketPrediction
	if len(request.HistoricalPrices) > 0 {
		var err error
		marketPrediction, err = c.aiService.PredictMarket(ctx, request.Symbol, request.HistoricalPrices, request.Indicators)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to predict market: " + err.Error()})
			return
		}
	}
	
	// Generate automated investment
	investment, err := c.aiService.GenerateAutomatedInvestment(ctx, request.UserID, newsAnalysis, priceForecast, marketPrediction)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	if investment == nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "No strong investment signal found"})
		return
	}
	
	ctx.JSON(http.StatusOK, investment)
}
