package recommendation

import (
	"time"
)

// UserProfile represents a user's investment profile
type UserProfile struct {
	UserID           string    `json:"user_id"`
	RiskTolerance    float64   `json:"risk_tolerance"`    // 0.0 (very conservative) to 1.0 (very aggressive)
	InvestmentGoals  []string  `json:"investment_goals"`  // e.g., "RETIREMENT", "EDUCATION", "HOUSE", "WEALTH_GROWTH"
	TimeHorizon      int       `json:"time_horizon"`      // Investment horizon in years
	Age              int       `json:"age"`               // User's age
	Income           float64   `json:"income"`            // Annual income
	LiquidNetWorth   float64   `json:"liquid_net_worth"`  // Liquid net worth
	InvestmentAmount float64   `json:"investment_amount"` // Amount available to invest
	PreferredSectors []string  `json:"preferred_sectors"` // Preferred sectors
	ExcludedSectors  []string  `json:"excluded_sectors"`  // Sectors to exclude
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// InvestmentAsset represents an investment asset that can be recommended
type InvestmentAsset struct {
	ID                string    `json:"id"`
	Symbol            string    `json:"symbol"`
	Name              string    `json:"name"`
	AssetType         string    `json:"asset_type"`         // e.g., "STOCK", "BOND", "ETF", "CRYPTO", "REAL_ESTATE"
	Sector            string    `json:"sector"`             // e.g., "TECHNOLOGY", "HEALTHCARE", "FINANCE"
	RiskLevel         float64   `json:"risk_level"`         // 0.0 (very low risk) to 1.0 (very high risk)
	HistoricalReturns []float64 `json:"historical_returns"` // Historical annual returns
	Volatility        float64   `json:"volatility"`         // Standard deviation of returns
	CurrentPrice      float64   `json:"current_price"`
	OneYearTarget     float64   `json:"one_year_target"`
	FiveYearTarget    float64   `json:"five_year_target"`
	DividendYield     float64   `json:"dividend_yield"`
	MarketCap         float64   `json:"market_cap"`
	ESGScore          float64   `json:"esg_score"`          // Environmental, Social, and Governance score
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// MarketData represents current market conditions
type MarketData struct {
	Timestamp          time.Time           `json:"timestamp"`
	MarketTrends       map[string]float64  `json:"market_trends"`       // Trends by sector
	EconomicIndicators map[string]float64  `json:"economic_indicators"` // e.g., "INFLATION", "GDP_GROWTH", "UNEMPLOYMENT"
	InterestRates      map[string]float64  `json:"interest_rates"`      // e.g., "FED_RATE", "10Y_TREASURY"
	MarketSentiment    float64             `json:"market_sentiment"`    // -1.0 (very bearish) to 1.0 (very bullish)
	SectorPerformance  map[string]float64  `json:"sector_performance"`  // Performance by sector
	AssetCorrelations  map[string][]string `json:"asset_correlations"`  // Correlations between assets
}

// UserInteraction represents a user's interaction with investments
type UserInteraction struct {
	UserID       string    `json:"user_id"`
	AssetID      string    `json:"asset_id"`
	InteractType string    `json:"interact_type"` // e.g., "VIEW", "SAVE", "PURCHASE", "SELL"
	Timestamp    time.Time `json:"timestamp"`
	Amount       float64   `json:"amount"`       // Amount involved in the interaction
	Duration     int       `json:"duration"`     // Duration of interaction in seconds (for views)
	Feedback     int       `json:"feedback"`     // -1 (negative), 0 (neutral), 1 (positive)
}

// RecommendationRequest represents a request for investment recommendations
type RecommendationRequest struct {
	UserID           string   `json:"user_id"`
	Amount           float64  `json:"amount"`            // Amount to invest
	TimeHorizon      int      `json:"time_horizon"`      // Investment horizon in years
	RiskTolerance    float64  `json:"risk_tolerance"`    // 0.0 (very conservative) to 1.0 (very aggressive)
	PreferredSectors []string `json:"preferred_sectors"` // Preferred sectors
	ExcludedSectors  []string `json:"excluded_sectors"`  // Sectors to exclude
	Goals            []string `json:"goals"`             // Investment goals
	Constraints      []string `json:"constraints"`       // Investment constraints
}

// RecommendedAsset represents a recommended investment asset
type RecommendedAsset struct {
	Asset              InvestmentAsset `json:"asset"`
	AllocationPercent  float64         `json:"allocation_percent"`  // Recommended allocation percentage
	ExpectedReturn     float64         `json:"expected_return"`     // Expected annual return
	RiskContribution   float64         `json:"risk_contribution"`   // Contribution to portfolio risk
	ConfidenceScore    float64         `json:"confidence_score"`    // Confidence in the recommendation (0.0 to 1.0)
	RecommendationTags []string        `json:"recommendation_tags"` // e.g., "GROWTH", "INCOME", "DIVERSIFICATION"
	Reasoning          string          `json:"reasoning"`           // Explanation for the recommendation
}

// PortfolioRecommendation represents a recommended portfolio
type PortfolioRecommendation struct {
	UserID                string            `json:"user_id"`
	RecommendationID      string            `json:"recommendation_id"`
	RecommendedAssets     []RecommendedAsset `json:"recommended_assets"`
	TotalExpectedReturn   float64           `json:"total_expected_return"`
	PortfolioRiskLevel    float64           `json:"portfolio_risk_level"`
	DiversificationScore  float64           `json:"diversification_score"`
	RebalancingFrequency  string            `json:"rebalancing_frequency"` // e.g., "MONTHLY", "QUARTERLY", "ANNUALLY"
	TimeHorizon           int               `json:"time_horizon"`          // Investment horizon in years
	CreatedAt             time.Time         `json:"created_at"`
}

// ModelPerformanceMetrics represents the performance metrics of the recommendation model
type ModelPerformanceMetrics struct {
	ModelVersion       string    `json:"model_version"`
	Accuracy           float64   `json:"accuracy"`
	Precision          float64   `json:"precision"`
	Recall             float64   `json:"recall"`
	F1Score            float64   `json:"f1_score"`
	MeanAbsoluteError  float64   `json:"mean_absolute_error"`
	RootMeanSquareError float64  `json:"root_mean_square_error"`
	UserSatisfaction   float64   `json:"user_satisfaction"` // Average user satisfaction score
	LastEvaluatedAt    time.Time `json:"last_evaluated_at"`
}
