package timeseries

import (
	"context"
	"math"
	"sync"
	"time"
)

// MarketIndicator represents a market indicator
type MarketIndicator struct {
	Name  string `json:"name"`
	Value float64 `json:"value"`
	Type  string `json:"type"` // "TECHNICAL", "FUNDAMENTAL", "SENTIMENT"
}

// MarketPrediction represents a market prediction
type MarketPrediction struct {
	Symbol           string    `json:"symbol"`
	Timestamp        time.Time `json:"timestamp"`
	Direction        string    `json:"direction"`        // "UP", "DOWN", "SIDEWAYS"
	Probability      float64   `json:"probability"`      // 0.0 to 1.0
	ExpectedReturn   float64   `json:"expected_return"`  // Expected percentage return
	Volatility       float64   `json:"volatility"`       // Expected volatility
	TimePeriod       string    `json:"time_period"`      // "SHORT_TERM", "MEDIUM_TERM", "LONG_TERM"
	SupportingFactors []string  `json:"supporting_factors"`
}

// MarketPredictor predicts market movements using multiple indicators
type MarketPredictor struct {
	priceForecaster *PriceForecaster
	
	// Cache of predictions
	predictionCache map[string]MarketPrediction
	cacheMutex      sync.RWMutex
	cacheTime       time.Duration
}

// NewMarketPredictor creates a new market predictor
func NewMarketPredictor(priceForecaster *PriceForecaster) *MarketPredictor {
	return &MarketPredictor{
		priceForecaster: priceForecaster,
		predictionCache: make(map[string]MarketPrediction),
		cacheTime:       1 * time.Hour,
	}
}

// PredictMarket predicts market movements for a symbol
func (p *MarketPredictor) PredictMarket(ctx context.Context, symbol string, historicalPrices []PricePoint, indicators []MarketIndicator) (*MarketPrediction, error) {
	// Check cache first
	cacheKey := symbol
	p.cacheMutex.RLock()
	if prediction, ok := p.predictionCache[cacheKey]; ok {
		if time.Since(prediction.Timestamp) < p.cacheTime {
			p.cacheMutex.RUnlock()
			return &prediction, nil
		}
	}
	p.cacheMutex.RUnlock()
	
	// Get price forecast
	forecast, err := p.priceForecaster.ForecastPrice(ctx, symbol, historicalPrices, 30)
	if err != nil {
		return nil, err
	}
	
	// Calculate technical indicators
	technicalScore := calculateTechnicalScore(historicalPrices)
	
	// Calculate fundamental score
	fundamentalScore := calculateFundamentalScore(indicators)
	
	// Calculate sentiment score
	sentimentScore := calculateSentimentScore(indicators)
	
	// Combine scores
	combinedScore := (technicalScore*0.4 + fundamentalScore*0.3 + sentimentScore*0.3)
	
	// Determine direction
	direction := "SIDEWAYS"
	if combinedScore > 0.2 {
		direction = "UP"
	} else if combinedScore < -0.2 {
		direction = "DOWN"
	}
	
	// Calculate probability
	probability := math.Abs(combinedScore)
	if probability > 1.0 {
		probability = 1.0
	}
	
	// Calculate expected return
	// Use the 30-day forecast
	currentPrice := forecast.CurrentPrice
	forecastPrice := forecast.ForecastPoints[29].Price
	expectedReturn := (forecastPrice - currentPrice) / currentPrice
	
	// Calculate volatility
	volatility := calculateVolatility(historicalPrices)
	
	// Determine time period
	timePeriod := "MEDIUM_TERM" // Default
	
	// Identify supporting factors
	supportingFactors := identifySupportingFactors(direction, indicators, historicalPrices)
	
	// Create prediction
	prediction := MarketPrediction{
		Symbol:           symbol,
		Timestamp:        time.Now(),
		Direction:        direction,
		Probability:      probability,
		ExpectedReturn:   expectedReturn,
		Volatility:       volatility,
		TimePeriod:       timePeriod,
		SupportingFactors: supportingFactors,
	}
	
	// Cache prediction
	p.cacheMutex.Lock()
	p.predictionCache[cacheKey] = prediction
	p.cacheMutex.Unlock()
	
	return &prediction, nil
}

// calculateTechnicalScore calculates a score based on technical indicators
func calculateTechnicalScore(prices []PricePoint) float64 {
	if len(prices) < 50 {
		return 0
	}
	
	// Calculate moving averages
	sma20 := calculateSMA(prices, 20)
	sma50 := calculateSMA(prices, 50)
	
	// Calculate RSI
	rsi := calculateRSI(prices, 14)
	
	// Calculate MACD
	macd := calculateMACD(prices)
	
	// Calculate score
	var score float64
	
	// Moving average crossover
	if sma20 > sma50 {
		score += 0.3 // Bullish
	} else {
		score -= 0.3 // Bearish
	}
	
	// RSI
	if rsi < 30 {
		score += 0.2 // Oversold, bullish
	} else if rsi > 70 {
		score -= 0.2 // Overbought, bearish
	}
	
	// MACD
	if macd > 0 {
		score += 0.2 // Bullish
	} else {
		score -= 0.2 // Bearish
	}
	
	return score
}

// calculateFundamentalScore calculates a score based on fundamental indicators
func calculateFundamentalScore(indicators []MarketIndicator) float64 {
	var score float64
	var count int
	
	for _, indicator := range indicators {
		if indicator.Type != "FUNDAMENTAL" {
			continue
		}
		
		switch indicator.Name {
		case "PE_RATIO":
			// Lower P/E is better (value investing)
			if indicator.Value < 15 {
				score += 0.3
			} else if indicator.Value > 30 {
				score -= 0.2
			}
			count++
			
		case "PB_RATIO":
			// Lower P/B is better (value investing)
			if indicator.Value < 1.5 {
				score += 0.2
			} else if indicator.Value > 3 {
				score -= 0.1
			}
			count++
			
		case "DIVIDEND_YIELD":
			// Higher dividend yield is better
			if indicator.Value > 0.04 { // 4%
				score += 0.2
			}
			count++
			
		case "REVENUE_GROWTH":
			// Higher revenue growth is better
			if indicator.Value > 0.2 { // 20%
				score += 0.3
			} else if indicator.Value < 0 {
				score -= 0.3
			}
			count++
			
		case "PROFIT_MARGIN":
			// Higher profit margin is better
			if indicator.Value > 0.15 { // 15%
				score += 0.2
			} else if indicator.Value < 0 {
				score -= 0.3
			}
			count++
		}
	}
	
	if count == 0 {
		return 0
	}
	
	return score / float64(count)
}

// calculateSentimentScore calculates a score based on sentiment indicators
func calculateSentimentScore(indicators []MarketIndicator) float64 {
	var score float64
	var count int
	
	for _, indicator := range indicators {
		if indicator.Type != "SENTIMENT" {
			continue
		}
		
		switch indicator.Name {
		case "NEWS_SENTIMENT":
			// Range: -1.0 to 1.0
			score += indicator.Value
			count++
			
		case "SOCIAL_MEDIA_SENTIMENT":
			// Range: -1.0 to 1.0
			score += indicator.Value * 0.7 // Less weight than news
			count++
			
		case "ANALYST_RECOMMENDATIONS":
			// Range: 1.0 (Strong Sell) to 5.0 (Strong Buy)
			// Convert to -1.0 to 1.0
			score += (indicator.Value - 3.0) / 2.0
			count++
			
		case "INSIDER_TRADING":
			// Positive values indicate buying, negative indicate selling
			score += indicator.Value * 0.5
			count++
			
		case "FEAR_GREED_INDEX":
			// Range: 0 (Extreme Fear) to 100 (Extreme Greed)
			// Convert to -1.0 to 1.0
			score += (indicator.Value - 50) / 50.0
			count++
		}
	}
	
	if count == 0 {
		return 0
	}
	
	return score / float64(count)
}

// calculateSMA calculates Simple Moving Average
func calculateSMA(prices []PricePoint, period int) float64 {
	if len(prices) < period {
		return 0
	}
	
	var sum float64
	for i := len(prices) - period; i < len(prices); i++ {
		sum += prices[i].Close
	}
	
	return sum / float64(period)
}

// calculateRSI calculates Relative Strength Index
func calculateRSI(prices []PricePoint, period int) float64 {
	if len(prices) < period + 1 {
		return 50 // Neutral
	}
	
	var gains, losses float64
	
	for i := len(prices) - period; i < len(prices); i++ {
		change := prices[i].Close - prices[i-1].Close
		if change > 0 {
			gains += change
		} else {
			losses -= change
		}
	}
	
	if losses == 0 {
		return 100 // All gains
	}
	
	rs := gains / losses
	rsi := 100 - (100 / (1 + rs))
	
	return rsi
}

// calculateMACD calculates Moving Average Convergence Divergence
func calculateMACD(prices []PricePoint) float64 {
	if len(prices) < 26 {
		return 0
	}
	
	// Calculate 12-day EMA
	ema12 := calculateEMA(prices, 12)
	
	// Calculate 26-day EMA
	ema26 := calculateEMA(prices, 26)
	
	// MACD Line = 12-day EMA - 26-day EMA
	macd := ema12 - ema26
	
	return macd
}

// calculateEMA calculates Exponential Moving Average
func calculateEMA(prices []PricePoint, period int) float64 {
	if len(prices) < period {
		return 0
	}
	
	// Calculate initial SMA
	var sum float64
	for i := len(prices) - period; i < len(prices); i++ {
		sum += prices[i].Close
	}
	sma := sum / float64(period)
	
	// Calculate multiplier
	multiplier := 2.0 / float64(period+1)
	
	// Calculate EMA
	ema := sma
	for i := len(prices) - period + 1; i < len(prices); i++ {
		ema = (prices[i].Close - ema) * multiplier + ema
	}
	
	return ema
}

// calculateVolatility calculates historical volatility
func calculateVolatility(prices []PricePoint) float64 {
	if len(prices) < 20 {
		return 0
	}
	
	// Calculate daily returns
	returns := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		returns[i-1] = prices[i].Close / prices[i-1].Close - 1
	}
	
	// Calculate mean return
	var sumReturns float64
	for _, r := range returns {
		sumReturns += r
	}
	meanReturn := sumReturns / float64(len(returns))
	
	// Calculate variance
	var sumSquaredDiff float64
	for _, r := range returns {
		diff := r - meanReturn
		sumSquaredDiff += diff * diff
	}
	variance := sumSquaredDiff / float64(len(returns))
	
	// Calculate annualized volatility
	dailyVolatility := math.Sqrt(variance)
	annualizedVolatility := dailyVolatility * math.Sqrt(252) // 252 trading days in a year
	
	return annualizedVolatility
}

// identifySupportingFactors identifies factors supporting the prediction
func identifySupportingFactors(direction string, indicators []MarketIndicator, prices []PricePoint) []string {
	var factors []string
	
	// Technical factors
	if len(prices) >= 50 {
		sma20 := calculateSMA(prices, 20)
		sma50 := calculateSMA(prices, 50)
		
		if direction == "UP" && sma20 > sma50 {
			factors = append(factors, "20-day SMA above 50-day SMA (Golden Cross)")
		} else if direction == "DOWN" && sma20 < sma50 {
			factors = append(factors, "20-day SMA below 50-day SMA (Death Cross)")
		}
		
		rsi := calculateRSI(prices, 14)
		if direction == "UP" && rsi < 30 {
			factors = append(factors, "RSI indicates oversold conditions")
		} else if direction == "DOWN" && rsi > 70 {
			factors = append(factors, "RSI indicates overbought conditions")
		}
		
		macd := calculateMACD(prices)
		if direction == "UP" && macd > 0 {
			factors = append(factors, "MACD is positive")
		} else if direction == "DOWN" && macd < 0 {
			factors = append(factors, "MACD is negative")
		}
	}
	
	// Fundamental factors
	for _, indicator := range indicators {
		if indicator.Type != "FUNDAMENTAL" {
			continue
		}
		
		switch indicator.Name {
		case "PE_RATIO":
			if direction == "UP" && indicator.Value < 15 {
				factors = append(factors, "Low P/E ratio indicates undervaluation")
			} else if direction == "DOWN" && indicator.Value > 30 {
				factors = append(factors, "High P/E ratio indicates overvaluation")
			}
			
		case "REVENUE_GROWTH":
			if direction == "UP" && indicator.Value > 0.2 {
				factors = append(factors, "Strong revenue growth")
			} else if direction == "DOWN" && indicator.Value < 0 {
				factors = append(factors, "Declining revenue")
			}
		}
	}
	
	// Sentiment factors
	for _, indicator := range indicators {
		if indicator.Type != "SENTIMENT" {
			continue
		}
		
		switch indicator.Name {
		case "NEWS_SENTIMENT":
			if direction == "UP" && indicator.Value > 0.3 {
				factors = append(factors, "Positive news sentiment")
			} else if direction == "DOWN" && indicator.Value < -0.3 {
				factors = append(factors, "Negative news sentiment")
			}
			
		case "ANALYST_RECOMMENDATIONS":
			if direction == "UP" && indicator.Value > 4 {
				factors = append(factors, "Strong analyst buy recommendations")
			} else if direction == "DOWN" && indicator.Value < 2 {
				factors = append(factors, "Strong analyst sell recommendations")
			}
		}
	}
	
	return factors
}
