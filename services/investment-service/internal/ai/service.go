package ai

import (
	"context"
	"log"
	"time"

	"github.com/sparkfund/services/investment-service/internal/ai/nlp"
	"github.com/sparkfund/services/investment-service/internal/ai/recommendation"
	"github.com/sparkfund/services/investment-service/internal/ai/rl"
	"github.com/sparkfund/services/investment-service/internal/ai/security"
	"github.com/sparkfund/services/investment-service/internal/ai/timeseries"
)

// AIService provides a unified interface to all AI capabilities
type AIService struct {
	recommendationService *recommendation.Service
	securityService       *security.SecurityService
	sentimentAnalyzer     *nlp.SentimentAnalyzer
	newsAnalyzer          *nlp.NewsAnalyzer
	priceForecaster       *timeseries.PriceForecaster
	marketPredictor       *timeseries.MarketPredictor
	portfolioOptimizer    *rl.PortfolioOptimizer
	rlAgent               *rl.RLAgent
}

// NewAIService creates a new AI service
func NewAIService(
	recommendationService *recommendation.Service,
	securityService *security.SecurityService,
) *AIService {
	// Create NLP components
	sentimentAnalyzer := nlp.NewSentimentAnalyzer()
	newsAnalyzer := nlp.NewNewsAnalyzer(sentimentAnalyzer)
	
	// Create time series components
	priceForecaster := timeseries.NewPriceForecaster()
	marketPredictor := timeseries.NewMarketPredictor(priceForecaster)
	
	// Create RL components
	portfolioOptimizer := rl.NewPortfolioOptimizer()
	rlAgent := rl.NewRLAgent()
	
	return &AIService{
		recommendationService: recommendationService,
		securityService:       securityService,
		sentimentAnalyzer:     sentimentAnalyzer,
		newsAnalyzer:          newsAnalyzer,
		priceForecaster:       priceForecaster,
		marketPredictor:       marketPredictor,
		portfolioOptimizer:    portfolioOptimizer,
		rlAgent:               rlAgent,
	}
}

// AnalyzeNews analyzes a news article and generates investment signals
func (s *AIService) AnalyzeNews(ctx context.Context, article nlp.NewsArticle) (*nlp.NewsAnalysisResult, error) {
	log.Printf("Analyzing news article: %s", article.Title)
	
	// Analyze the article
	result, err := s.newsAnalyzer.AnalyzeArticle(ctx, article)
	if err != nil {
		log.Printf("Error analyzing news article: %v", err)
		return nil, err
	}
	
	log.Printf("News analysis complete. Found %d investment signals", len(result.InvestmentSignals))
	
	return result, nil
}

// ForecastPrice forecasts future prices for a symbol
func (s *AIService) ForecastPrice(ctx context.Context, symbol string, historicalPrices []timeseries.PricePoint, days int) (*timeseries.Forecast, error) {
	log.Printf("Forecasting prices for %s for %d days", symbol, days)
	
	// Forecast prices
	forecast, err := s.priceForecaster.ForecastPrice(ctx, symbol, historicalPrices, days)
	if err != nil {
		log.Printf("Error forecasting prices: %v", err)
		return nil, err
	}
	
	log.Printf("Price forecast complete. Model accuracy: %.2f%%", forecast.ModelAccuracy*100)
	
	return forecast, nil
}

// PredictMarket predicts market movements for a symbol
func (s *AIService) PredictMarket(ctx context.Context, symbol string, historicalPrices []timeseries.PricePoint, indicators []timeseries.MarketIndicator) (*timeseries.MarketPrediction, error) {
	log.Printf("Predicting market for %s", symbol)
	
	// Predict market
	prediction, err := s.marketPredictor.PredictMarket(ctx, symbol, historicalPrices, indicators)
	if err != nil {
		log.Printf("Error predicting market: %v", err)
		return nil, err
	}
	
	log.Printf("Market prediction complete. Direction: %s, Probability: %.2f%%", prediction.Direction, prediction.Probability*100)
	
	return prediction, nil
}

// OptimizePortfolio optimizes a portfolio using reinforcement learning
func (s *AIService) OptimizePortfolio(ctx context.Context, portfolio rl.Portfolio, availableAssets []rl.Asset) (*rl.OptimizationResult, error) {
	log.Printf("Optimizing portfolio for user %s", portfolio.UserID)
	
	// Optimize portfolio
	result, err := s.portfolioOptimizer.OptimizePortfolio(ctx, portfolio, availableAssets)
	if err != nil {
		log.Printf("Error optimizing portfolio: %v", err)
		return nil, err
	}
	
	log.Printf("Portfolio optimization complete. Improvement: %.2f%%", result.Improvement*100)
	
	return result, nil
}

// GetRLAction gets the best action for a given state using reinforcement learning
func (s *AIService) GetRLAction(ctx context.Context, state rl.State) rl.Action {
	log.Printf("Getting RL action for user %s", state.Portfolio.UserID)
	
	// Get action
	action := s.rlAgent.GetAction(ctx, state)
	
	log.Printf("RL action: %s", action.Type)
	
	return action
}

// TrainRLAgent trains the reinforcement learning agent
func (s *AIService) TrainRLAgent(ctx context.Context) error {
	log.Printf("Training RL agent")
	
	// Train agent
	err := s.rlAgent.Train(ctx)
	if err != nil {
		log.Printf("Error training RL agent: %v", err)
		return err
	}
	
	log.Printf("RL agent training complete")
	
	return nil
}

// AddRLExperience adds an experience to the RL agent's replay buffer
func (s *AIService) AddRLExperience(experience rl.Experience) {
	log.Printf("Adding RL experience: %s action, reward: %.2f", experience.Action.Type, experience.Reward)
	
	// Add experience
	s.rlAgent.AddExperience(experience)
}

// AutomatedInvestment represents an automated investment decision
type AutomatedInvestment struct {
	UserID          string    `json:"user_id"`
	Symbol          string    `json:"symbol"`
	Action          string    `json:"action"`          // "BUY", "SELL"
	Amount          float64   `json:"amount"`          // Amount to invest/sell
	Confidence      float64   `json:"confidence"`      // 0.0 to 1.0
	Reasoning       string    `json:"reasoning"`
	DataSources     []string  `json:"data_sources"`    // e.g., "NEWS", "PRICE_FORECAST", "MARKET_PREDICTION"
	Timestamp       time.Time `json:"timestamp"`
}

// GenerateAutomatedInvestment generates an automated investment decision based on multiple data sources
func (s *AIService) GenerateAutomatedInvestment(
	ctx context.Context,
	userID string,
	newsAnalysis *nlp.NewsAnalysisResult,
	priceForecast *timeseries.Forecast,
	marketPrediction *timeseries.MarketPrediction,
) (*AutomatedInvestment, error) {
	log.Printf("Generating automated investment for user %s", userID)
	
	// Initialize data sources
	var dataSources []string
	
	// Initialize confidence and action
	var confidence float64
	var action string
	var symbol string
	var reasoning string
	
	// Combine signals from news analysis
	if newsAnalysis != nil {
		dataSources = append(dataSources, "NEWS")
		
		// Find the strongest signal
		var strongestSignal *nlp.InvestmentSignal
		
		for i, signal := range newsAnalysis.InvestmentSignals {
			if strongestSignal == nil || signal.Strength > strongestSignal.Strength {
				strongestSignal = &newsAnalysis.InvestmentSignals[i]
			}
		}
		
		if strongestSignal != nil {
			symbol = strongestSignal.Symbol
			action = strongestSignal.Action
			confidence = strongestSignal.Confidence * 0.4 // News contributes 40% to confidence
			reasoning = "Based on news analysis: " + strongestSignal.Reasoning
		}
	}
	
	// Combine signals from price forecast
	if priceForecast != nil && priceForecast.Symbol == symbol {
		dataSources = append(dataSources, "PRICE_FORECAST")
		
		// Calculate expected return over forecast period
		lastForecastPoint := priceForecast.ForecastPoints[len(priceForecast.ForecastPoints)-1]
		expectedReturn := (lastForecastPoint.Price - priceForecast.CurrentPrice) / priceForecast.CurrentPrice
		
		// Determine action based on expected return
		forecastAction := "HOLD"
		if expectedReturn > 0.05 { // 5% threshold for buy
			forecastAction = "BUY"
		} else if expectedReturn < -0.05 { // -5% threshold for sell
			forecastAction = "SELL"
		}
		
		// Add to confidence if actions align
		if action == "" {
			action = forecastAction
			confidence = priceForecast.ModelAccuracy * 0.3 // Price forecast contributes 30% to confidence
			reasoning = "Based on price forecast: " + formatReturn(expectedReturn) + " expected return over forecast period."
		} else if action == forecastAction {
			confidence += priceForecast.ModelAccuracy * 0.3
			reasoning += " Price forecast confirms with " + formatReturn(expectedReturn) + " expected return."
		} else {
			confidence -= priceForecast.ModelAccuracy * 0.1 // Conflicting signals reduce confidence
			reasoning += " However, price forecast suggests " + forecastAction + " with " + formatReturn(expectedReturn) + " expected return."
		}
	}
	
	// Combine signals from market prediction
	if marketPrediction != nil && marketPrediction.Symbol == symbol {
		dataSources = append(dataSources, "MARKET_PREDICTION")
		
		// Determine action based on market prediction
		marketAction := "HOLD"
		if marketPrediction.Direction == "UP" {
			marketAction = "BUY"
		} else if marketPrediction.Direction == "DOWN" {
			marketAction = "SELL"
		}
		
		// Add to confidence if actions align
		if action == "" {
			action = marketAction
			confidence = marketPrediction.Probability * 0.3 // Market prediction contributes 30% to confidence
			reasoning = "Based on market prediction: " + marketPrediction.Direction + " direction with " + 
				formatReturn(marketPrediction.ExpectedReturn) + " expected return."
		} else if action == marketAction {
			confidence += marketPrediction.Probability * 0.3
			reasoning += " Market prediction confirms with " + marketPrediction.Direction + " direction and " + 
				formatReturn(marketPrediction.ExpectedReturn) + " expected return."
		} else {
			confidence -= marketPrediction.Probability * 0.1 // Conflicting signals reduce confidence
			reasoning += " However, market prediction suggests " + marketAction + " with " + 
				formatReturn(marketPrediction.ExpectedReturn) + " expected return."
		}
	}
	
	// If no action determined or confidence too low, default to HOLD
	if action == "" || action == "HOLD" || confidence < 0.6 {
		log.Printf("No strong investment signal found. Action: %s, Confidence: %.2f", action, confidence)
		return nil, nil
	}
	
	// Calculate investment amount based on confidence
	// Higher confidence = higher percentage of available funds
	amount := 1000.0 + confidence * 9000.0 // $1,000 to $10,000 based on confidence
	
	// Create automated investment
	investment := &AutomatedInvestment{
		UserID:      userID,
		Symbol:      symbol,
		Action:      action,
		Amount:      amount,
		Confidence:  confidence,
		Reasoning:   reasoning,
		DataSources: dataSources,
		Timestamp:   time.Now(),
	}
	
	log.Printf("Generated automated investment: %s %s $%.2f with %.2f%% confidence", 
		investment.Action, investment.Symbol, investment.Amount, investment.Confidence*100)
	
	return investment, nil
}

// formatReturn formats a return value as a percentage string
func formatReturn(value float64) string {
	if value >= 0 {
		return "+" + fmt.Sprintf("%.2f%%", value*100)
	}
	return fmt.Sprintf("%.2f%%", value*100)
}
