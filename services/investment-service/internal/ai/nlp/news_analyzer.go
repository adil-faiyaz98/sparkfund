package nlp

import (
	"context"
	"log"
	"sync"
	"time"
)

// NewsSource represents a source of market news
type NewsSource struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Reliability float64 `json:"reliability"` // 0.0 to 1.0
	Categories  []string `json:"categories"` // e.g., "FINANCE", "TECHNOLOGY"
}

// NewsArticle represents a market news article
type NewsArticle struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Summary     string    `json:"summary"`
	URL         string    `json:"url"`
	SourceID    string    `json:"source_id"`
	PublishedAt time.Time `json:"published_at"`
	Categories  []string  `json:"categories"`
	Entities    []string  `json:"entities"` // Companies, sectors mentioned
}

// NewsAnalysisResult represents the result of news analysis
type NewsAnalysisResult struct {
	ArticleID       string         `json:"article_id"`
	Sentiment       SentimentResult `json:"sentiment"`
	MarketImpact    MarketImpact    `json:"market_impact"`
	InvestmentSignals []InvestmentSignal `json:"investment_signals"`
	Timestamp       time.Time      `json:"timestamp"`
}

// InvestmentSignal represents a signal for automated investment
type InvestmentSignal struct {
	Symbol        string  `json:"symbol"`
	Action        string  `json:"action"`        // "BUY", "SELL", "HOLD"
	Strength      float64 `json:"strength"`      // 0.0 to 1.0
	Confidence    float64 `json:"confidence"`    // 0.0 to 1.0
	Reasoning     string  `json:"reasoning"`
	RecommendedAmount float64 `json:"recommended_amount,omitempty"`
}

// NewsAnalyzer analyzes market news for investment signals
type NewsAnalyzer struct {
	sentimentAnalyzer *SentimentAnalyzer
	
	// Configuration
	signalThreshold float64 // Minimum strength for a signal
	
	// News sources
	sources map[string]NewsSource
	
	// Cache
	analysisCache      map[string]NewsAnalysisResult
	cacheMutex         sync.RWMutex
	cacheTime          time.Duration
}

// NewNewsAnalyzer creates a new news analyzer
func NewNewsAnalyzer(sentimentAnalyzer *SentimentAnalyzer) *NewsAnalyzer {
	return &NewsAnalyzer{
		sentimentAnalyzer: sentimentAnalyzer,
		signalThreshold:   0.6,
		sources:           loadNewsSources(),
		analysisCache:     make(map[string]NewsAnalysisResult),
		cacheTime:         1 * time.Hour,
	}
}

// AnalyzeArticle analyzes a news article
func (a *NewsAnalyzer) AnalyzeArticle(ctx context.Context, article NewsArticle) (*NewsAnalysisResult, error) {
	// Check cache first
	a.cacheMutex.RLock()
	if result, ok := a.analysisCache[article.ID]; ok {
		if time.Since(result.Timestamp) < a.cacheTime {
			a.cacheMutex.RUnlock()
			return &result, nil
		}
	}
	a.cacheMutex.RUnlock()
	
	// Analyze sentiment
	sentiment := a.sentimentAnalyzer.AnalyzeText(article.Title + " " + article.Content)
	
	// Predict market impact
	marketImpact := a.sentimentAnalyzer.PredictMarketImpact(sentiment)
	
	// Generate investment signals
	signals := a.generateInvestmentSignals(article, sentiment, marketImpact)
	
	// Create result
	result := NewsAnalysisResult{
		ArticleID:        article.ID,
		Sentiment:        sentiment,
		MarketImpact:     marketImpact,
		InvestmentSignals: signals,
		Timestamp:        time.Now(),
	}
	
	// Cache result
	a.cacheMutex.Lock()
	a.analysisCache[article.ID] = result
	a.cacheMutex.Unlock()
	
	return &result, nil
}

// generateInvestmentSignals generates investment signals from news analysis
func (a *NewsAnalyzer) generateInvestmentSignals(article NewsArticle, sentiment SentimentResult, impact MarketImpact) []InvestmentSignal {
	var signals []InvestmentSignal
	
	// Get source reliability
	sourceReliability := 0.5 // Default
	if source, ok := a.sources[article.SourceID]; ok {
		sourceReliability = source.Reliability
	}
	
	// Generate signals for companies
	for company, companyImpact := range impact.CompanyImpacts {
		// Calculate signal strength based on price impact and sentiment magnitude
		strength := math.Abs(companyImpact.PriceImpact) * sentiment.Magnitude
		
		// Calculate confidence based on source reliability and sentiment magnitude
		confidence := sourceReliability * sentiment.Magnitude
		
		// Determine action based on price impact
		action := "HOLD"
		if companyImpact.PriceImpact > 0.01 { // 1% threshold for buy
			action = "BUY"
		} else if companyImpact.PriceImpact < -0.01 { // -1% threshold for sell
			action = "SELL"
		}
		
		// Only include signals above threshold
		if strength >= a.signalThreshold {
			// Generate reasoning
			reasoning := generateReasoning(action, company, sentiment, companyImpact)
			
			// Create signal
			signal := InvestmentSignal{
				Symbol:     company,
				Action:     action,
				Strength:   strength,
				Confidence: confidence,
				Reasoning:  reasoning,
			}
			
			signals = append(signals, signal)
		}
	}
	
	return signals
}

// generateReasoning generates reasoning for an investment signal
func generateReasoning(action, symbol string, sentiment SentimentResult, impact Impact) string {
	var reasoning string
	
	switch action {
	case "BUY":
		reasoning = "Positive sentiment detected for " + symbol + " with expected price increase of " + 
			fmt.Sprintf("%.2f%%", impact.PriceImpact*100) + ". "
		
		if sentiment.Magnitude > 0.7 {
			reasoning += "Strong sentiment magnitude indicates high confidence. "
		}
		
		if impact.TimePeriod == "SHORT_TERM" {
			reasoning += "Impact expected in the short term."
		} else {
			reasoning += "Impact expected over the " + strings.ToLower(impact.TimePeriod) + "."
		}
		
	case "SELL":
		reasoning = "Negative sentiment detected for " + symbol + " with expected price decrease of " + 
			fmt.Sprintf("%.2f%%", -impact.PriceImpact*100) + ". "
		
		if sentiment.Magnitude > 0.7 {
			reasoning += "Strong sentiment magnitude indicates high confidence. "
		}
		
		if impact.TimePeriod == "SHORT_TERM" {
			reasoning += "Impact expected in the short term."
		} else {
			reasoning += "Impact expected over the " + strings.ToLower(impact.TimePeriod) + "."
		}
		
	case "HOLD":
		reasoning = "Mixed or neutral sentiment detected for " + symbol + ". "
		reasoning += "No significant price movement expected."
	}
	
	return reasoning
}

// loadNewsSources loads news sources
func loadNewsSources() map[string]NewsSource {
	// In a real implementation, this would load from a file or database
	// This is a simplified version with a few examples
	sources := []NewsSource{
		{
			ID:          "bloomberg",
			Name:        "Bloomberg",
			URL:         "https://www.bloomberg.com",
			Reliability: 0.9,
			Categories:  []string{"FINANCE", "MARKETS", "ECONOMY"},
		},
		{
			ID:          "reuters",
			Name:        "Reuters",
			URL:         "https://www.reuters.com",
			Reliability: 0.9,
			Categories:  []string{"FINANCE", "MARKETS", "ECONOMY", "WORLD"},
		},
		{
			ID:          "wsj",
			Name:        "Wall Street Journal",
			URL:         "https://www.wsj.com",
			Reliability: 0.85,
			Categories:  []string{"FINANCE", "MARKETS", "ECONOMY", "BUSINESS"},
		},
		{
			ID:          "ft",
			Name:        "Financial Times",
			URL:         "https://www.ft.com",
			Reliability: 0.85,
			Categories:  []string{"FINANCE", "MARKETS", "ECONOMY", "BUSINESS"},
		},
		{
			ID:          "cnbc",
			Name:        "CNBC",
			URL:         "https://www.cnbc.com",
			Reliability: 0.75,
			Categories:  []string{"FINANCE", "MARKETS", "BUSINESS"},
		},
	}
	
	sourceMap := make(map[string]NewsSource)
	for _, source := range sources {
		sourceMap[source.ID] = source
	}
	
	return sourceMap
}
