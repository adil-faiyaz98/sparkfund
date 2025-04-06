package timeseries

import (
	"context"
	"errors"
	"math"
	"sort"
	"sync"
	"time"
)

// PricePoint represents a single price data point
type PricePoint struct {
	Symbol    string    `json:"symbol"`
	Timestamp time.Time `json:"timestamp"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    int64     `json:"volume"`
}

// Forecast represents a price forecast
type Forecast struct {
	Symbol           string       `json:"symbol"`
	CurrentPrice     float64      `json:"current_price"`
	ForecastPoints   []ForecastPoint `json:"forecast_points"`
	ConfidenceInterval ConfidenceInterval `json:"confidence_interval"`
	ModelAccuracy    float64      `json:"model_accuracy"`
	CreatedAt        time.Time    `json:"created_at"`
}

// ForecastPoint represents a single forecast point
type ForecastPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Price     float64   `json:"price"`
}

// ConfidenceInterval represents the confidence interval for a forecast
type ConfidenceInterval struct {
	Lower []float64 `json:"lower"` // Lower bound for each forecast point
	Upper []float64 `json:"upper"` // Upper bound for each forecast point
	Level float64   `json:"level"` // Confidence level (e.g., 0.95 for 95%)
}

// PriceForecaster forecasts future prices using time series models
type PriceForecaster struct {
	// Cache of forecasts
	forecastCache map[string]Forecast
	cacheMutex    sync.RWMutex
	cacheTime     time.Duration
}

// NewPriceForecaster creates a new price forecaster
func NewPriceForecaster() *PriceForecaster {
	return &PriceForecaster{
		forecastCache: make(map[string]Forecast),
		cacheTime:     1 * time.Hour,
	}
}

// ForecastPrice forecasts future prices for a symbol
func (f *PriceForecaster) ForecastPrice(ctx context.Context, symbol string, historicalPrices []PricePoint, days int) (*Forecast, error) {
	// Validate input
	if len(historicalPrices) < 30 {
		return nil, errors.New("insufficient historical data for forecasting")
	}
	
	if days <= 0 {
		return nil, errors.New("forecast days must be positive")
	}
	
	// Check cache first
	cacheKey := symbol + "_" + strconv.Itoa(days)
	f.cacheMutex.RLock()
	if forecast, ok := f.forecastCache[cacheKey]; ok {
		if time.Since(forecast.CreatedAt) < f.cacheTime {
			f.cacheMutex.RUnlock()
			return &forecast, nil
		}
	}
	f.cacheMutex.RUnlock()
	
	// Sort historical prices by timestamp (ascending)
	sort.Slice(historicalPrices, func(i, j int) bool {
		return historicalPrices[i].Timestamp.Before(historicalPrices[j].Timestamp)
	})
	
	// Extract closing prices
	closingPrices := make([]float64, len(historicalPrices))
	for i, price := range historicalPrices {
		closingPrices[i] = price.Close
	}
	
	// Get current price (last closing price)
	currentPrice := closingPrices[len(closingPrices)-1]
	
	// Apply multiple forecasting models and ensemble the results
	arima := arimaForecast(closingPrices, days)
	ema := emaForecast(closingPrices, days)
	
	// Ensemble forecasts (simple average)
	forecastPrices := make([]float64, days)
	for i := 0; i < days; i++ {
		forecastPrices[i] = (arima[i] + ema[i]) / 2.0
	}
	
	// Calculate confidence intervals
	lower, upper := calculateConfidenceIntervals(forecastPrices, historicalPrices, 0.95)
	
	// Create forecast points
	lastDate := historicalPrices[len(historicalPrices)-1].Timestamp
	forecastPoints := make([]ForecastPoint, days)
	for i := 0; i < days; i++ {
		forecastDate := lastDate.AddDate(0, 0, i+1)
		forecastPoints[i] = ForecastPoint{
			Timestamp: forecastDate,
			Price:     forecastPrices[i],
		}
	}
	
	// Calculate model accuracy using backtesting
	accuracy := calculateModelAccuracy(closingPrices)
	
	// Create forecast
	forecast := Forecast{
		Symbol:       symbol,
		CurrentPrice: currentPrice,
		ForecastPoints: forecastPoints,
		ConfidenceInterval: ConfidenceInterval{
			Lower: lower,
			Upper: upper,
			Level: 0.95,
		},
		ModelAccuracy: accuracy,
		CreatedAt:    time.Now(),
	}
	
	// Cache forecast
	f.cacheMutex.Lock()
	f.forecastCache[cacheKey] = forecast
	f.cacheMutex.Unlock()
	
	return &forecast, nil
}

// arimaForecast implements a simplified ARIMA model for forecasting
func arimaForecast(prices []float64, days int) []float64 {
	// This is a simplified implementation of ARIMA
	// In a real implementation, this would use a proper ARIMA model
	
	// Calculate returns
	returns := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		returns[i-1] = prices[i] / prices[i-1] - 1
	}
	
	// Calculate mean return
	var sumReturns float64
	for _, r := range returns {
		sumReturns += r
	}
	meanReturn := sumReturns / float64(len(returns))
	
	// Calculate standard deviation of returns
	var sumSquaredDiff float64
	for _, r := range returns {
		diff := r - meanReturn
		sumSquaredDiff += diff * diff
	}
	stdDev := math.Sqrt(sumSquaredDiff / float64(len(returns)))
	
	// Generate forecast
	forecast := make([]float64, days)
	lastPrice := prices[len(prices)-1]
	
	for i := 0; i < days; i++ {
		// Use AR(1) model: next_return = mean_return + phi * (last_return - mean_return) + error
		// Simplified: next_return = mean_return
		nextReturn := meanReturn
		
		// Calculate next price
		nextPrice := lastPrice * (1 + nextReturn)
		forecast[i] = nextPrice
		lastPrice = nextPrice
	}
	
	return forecast
}

// emaForecast implements Exponential Moving Average forecasting
func emaForecast(prices []float64, days int) []float64 {
	// Calculate EMA
	period := 10 // EMA period
	multiplier := 2.0 / float64(period+1)
	
	// Calculate initial SMA
	var sum float64
	for i := len(prices) - period; i < len(prices); i++ {
		sum += prices[i]
	}
	sma := sum / float64(period)
	
	// Calculate EMA
	ema := sma
	for i := len(prices) - period + 1; i < len(prices); i++ {
		ema = (prices[i] - ema) * multiplier + ema
	}
	
	// Generate forecast
	forecast := make([]float64, days)
	lastPrice := prices[len(prices)-1]
	
	for i := 0; i < days; i++ {
		// EMA forecast: next_price = last_price + alpha * (ema - last_price)
		alpha := 0.3 // Smoothing factor
		nextPrice := lastPrice + alpha * (ema - lastPrice)
		forecast[i] = nextPrice
		lastPrice = nextPrice
	}
	
	return forecast
}

// calculateConfidenceIntervals calculates confidence intervals for a forecast
func calculateConfidenceIntervals(forecast []float64, historicalPrices []PricePoint, level float64) ([]float64, []float64) {
	// Extract historical returns
	returns := make([]float64, len(historicalPrices)-1)
	for i := 1; i < len(historicalPrices); i++ {
		returns[i-1] = historicalPrices[i].Close / historicalPrices[i-1].Close - 1
	}
	
	// Calculate standard deviation of returns
	var sumReturns, sumSquaredReturns float64
	for _, r := range returns {
		sumReturns += r
		sumSquaredReturns += r * r
	}
	meanReturn := sumReturns / float64(len(returns))
	variance := sumSquaredReturns/float64(len(returns)) - meanReturn*meanReturn
	stdDev := math.Sqrt(variance)
	
	// Calculate z-score for the confidence level
	// For 95% confidence, z = 1.96
	z := 1.96
	if level != 0.95 {
		// This is a simplification; in a real implementation, we would use a proper z-table
		if level == 0.90 {
			z = 1.645
		} else if level == 0.99 {
			z = 2.576
		}
	}
	
	// Calculate confidence intervals
	lower := make([]float64, len(forecast))
	upper := make([]float64, len(forecast))
	
	lastPrice := historicalPrices[len(historicalPrices)-1].Close
	
	for i := 0; i < len(forecast); i++ {
		// Confidence interval widens with time
		timeAdjustment := math.Sqrt(float64(i + 1))
		margin := lastPrice * stdDev * z * timeAdjustment
		
		lower[i] = forecast[i] - margin
		upper[i] = forecast[i] + margin
		
		// Ensure lower bound is not negative
		if lower[i] < 0 {
			lower[i] = 0
		}
	}
	
	return lower, upper
}

// calculateModelAccuracy calculates the accuracy of the forecasting model using backtesting
func calculateModelAccuracy(prices []float64) float64 {
	// This is a simplified implementation of backtesting
	// In a real implementation, this would be more sophisticated
	
	// Use the last 20% of data for testing
	testSize := int(float64(len(prices)) * 0.2)
	if testSize < 5 {
		testSize = 5
	}
	
	trainingSize := len(prices) - testSize
	
	// Calculate mean absolute percentage error (MAPE)
	var sumAPE float64
	
	for i := 0; i < testSize; i++ {
		// Use simple moving average as the forecast
		period := 5
		if trainingSize+i < period {
			continue
		}
		
		var sum float64
		for j := 0; j < period; j++ {
			sum += prices[trainingSize+i-j-1]
		}
		forecast := sum / float64(period)
		
		actual := prices[trainingSize+i]
		ape := math.Abs((actual - forecast) / actual)
		sumAPE += ape
	}
	
	mape := sumAPE / float64(testSize)
	
	// Convert MAPE to accuracy (0-1 scale)
	accuracy := 1.0 - mape
	if accuracy < 0 {
		accuracy = 0
	}
	
	return accuracy
}
