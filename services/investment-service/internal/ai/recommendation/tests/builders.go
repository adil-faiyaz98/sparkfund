package tests

import (
	"time"

	"github.com/google/uuid"
	"github.com/sparkfund/services/investment-service/internal/ai/recommendation"
)

// UserProfileBuilder helps build UserProfile objects for testing
type UserProfileBuilder struct {
	profile recommendation.UserProfile
}

// NewUserProfileBuilder creates a new UserProfileBuilder with default values
func NewUserProfileBuilder() *UserProfileBuilder {
	return &UserProfileBuilder{
		profile: recommendation.UserProfile{
			UserID:           uuid.New().String(),
			RiskTolerance:    0.5,
			InvestmentGoals:  []string{"RETIREMENT", "WEALTH_GROWTH"},
			TimeHorizon:      10,
			Age:              35,
			Income:           75000,
			LiquidNetWorth:   100000,
			InvestmentAmount: 10000,
			PreferredSectors: []string{"TECHNOLOGY", "HEALTHCARE"},
			ExcludedSectors:  []string{"TOBACCO", "GAMBLING"},
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
	}
}

// WithUserID sets the UserID
func (b *UserProfileBuilder) WithUserID(userID string) *UserProfileBuilder {
	b.profile.UserID = userID
	return b
}

// WithRiskTolerance sets the RiskTolerance
func (b *UserProfileBuilder) WithRiskTolerance(riskTolerance float64) *UserProfileBuilder {
	b.profile.RiskTolerance = riskTolerance
	return b
}

// WithInvestmentGoals sets the InvestmentGoals
func (b *UserProfileBuilder) WithInvestmentGoals(goals []string) *UserProfileBuilder {
	b.profile.InvestmentGoals = goals
	return b
}

// WithTimeHorizon sets the TimeHorizon
func (b *UserProfileBuilder) WithTimeHorizon(timeHorizon int) *UserProfileBuilder {
	b.profile.TimeHorizon = timeHorizon
	return b
}

// WithAge sets the Age
func (b *UserProfileBuilder) WithAge(age int) *UserProfileBuilder {
	b.profile.Age = age
	return b
}

// WithIncome sets the Income
func (b *UserProfileBuilder) WithIncome(income float64) *UserProfileBuilder {
	b.profile.Income = income
	return b
}

// WithLiquidNetWorth sets the LiquidNetWorth
func (b *UserProfileBuilder) WithLiquidNetWorth(netWorth float64) *UserProfileBuilder {
	b.profile.LiquidNetWorth = netWorth
	return b
}

// WithInvestmentAmount sets the InvestmentAmount
func (b *UserProfileBuilder) WithInvestmentAmount(amount float64) *UserProfileBuilder {
	b.profile.InvestmentAmount = amount
	return b
}

// WithPreferredSectors sets the PreferredSectors
func (b *UserProfileBuilder) WithPreferredSectors(sectors []string) *UserProfileBuilder {
	b.profile.PreferredSectors = sectors
	return b
}

// WithExcludedSectors sets the ExcludedSectors
func (b *UserProfileBuilder) WithExcludedSectors(sectors []string) *UserProfileBuilder {
	b.profile.ExcludedSectors = sectors
	return b
}

// Build returns the built UserProfile
func (b *UserProfileBuilder) Build() recommendation.UserProfile {
	return b.profile
}

// InvestmentAssetBuilder helps build InvestmentAsset objects for testing
type InvestmentAssetBuilder struct {
	asset recommendation.InvestmentAsset
}

// NewInvestmentAssetBuilder creates a new InvestmentAssetBuilder with default values
func NewInvestmentAssetBuilder() *InvestmentAssetBuilder {
	return &InvestmentAssetBuilder{
		asset: recommendation.InvestmentAsset{
			ID:                uuid.New().String(),
			Symbol:            "AAPL",
			Name:              "Apple Inc.",
			AssetType:         "STOCK",
			Sector:            "TECHNOLOGY",
			RiskLevel:         0.6,
			HistoricalReturns: []float64{0.12, 0.15, 0.08, 0.10, 0.14},
			Volatility:        0.2,
			CurrentPrice:      150.0,
			OneYearTarget:     165.0,
			FiveYearTarget:    200.0,
			DividendYield:     0.01,
			MarketCap:         2500000000000,
			ESGScore:          0.8,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
	}
}

// WithID sets the ID
func (b *InvestmentAssetBuilder) WithID(id string) *InvestmentAssetBuilder {
	b.asset.ID = id
	return b
}

// WithSymbol sets the Symbol
func (b *InvestmentAssetBuilder) WithSymbol(symbol string) *InvestmentAssetBuilder {
	b.asset.Symbol = symbol
	return b
}

// WithName sets the Name
func (b *InvestmentAssetBuilder) WithName(name string) *InvestmentAssetBuilder {
	b.asset.Name = name
	return b
}

// WithAssetType sets the AssetType
func (b *InvestmentAssetBuilder) WithAssetType(assetType string) *InvestmentAssetBuilder {
	b.asset.AssetType = assetType
	return b
}

// WithSector sets the Sector
func (b *InvestmentAssetBuilder) WithSector(sector string) *InvestmentAssetBuilder {
	b.asset.Sector = sector
	return b
}

// WithRiskLevel sets the RiskLevel
func (b *InvestmentAssetBuilder) WithRiskLevel(riskLevel float64) *InvestmentAssetBuilder {
	b.asset.RiskLevel = riskLevel
	return b
}

// WithHistoricalReturns sets the HistoricalReturns
func (b *InvestmentAssetBuilder) WithHistoricalReturns(returns []float64) *InvestmentAssetBuilder {
	b.asset.HistoricalReturns = returns
	return b
}

// WithVolatility sets the Volatility
func (b *InvestmentAssetBuilder) WithVolatility(volatility float64) *InvestmentAssetBuilder {
	b.asset.Volatility = volatility
	return b
}

// WithCurrentPrice sets the CurrentPrice
func (b *InvestmentAssetBuilder) WithCurrentPrice(price float64) *InvestmentAssetBuilder {
	b.asset.CurrentPrice = price
	return b
}

// WithOneYearTarget sets the OneYearTarget
func (b *InvestmentAssetBuilder) WithOneYearTarget(target float64) *InvestmentAssetBuilder {
	b.asset.OneYearTarget = target
	return b
}

// WithFiveYearTarget sets the FiveYearTarget
func (b *InvestmentAssetBuilder) WithFiveYearTarget(target float64) *InvestmentAssetBuilder {
	b.asset.FiveYearTarget = target
	return b
}

// WithDividendYield sets the DividendYield
func (b *InvestmentAssetBuilder) WithDividendYield(yield float64) *InvestmentAssetBuilder {
	b.asset.DividendYield = yield
	return b
}

// WithMarketCap sets the MarketCap
func (b *InvestmentAssetBuilder) WithMarketCap(marketCap float64) *InvestmentAssetBuilder {
	b.asset.MarketCap = marketCap
	return b
}

// WithESGScore sets the ESGScore
func (b *InvestmentAssetBuilder) WithESGScore(score float64) *InvestmentAssetBuilder {
	b.asset.ESGScore = score
	return b
}

// Build returns the built InvestmentAsset
func (b *InvestmentAssetBuilder) Build() recommendation.InvestmentAsset {
	return b.asset
}

// UserInteractionBuilder helps build UserInteraction objects for testing
type UserInteractionBuilder struct {
	interaction recommendation.UserInteraction
}

// NewUserInteractionBuilder creates a new UserInteractionBuilder with default values
func NewUserInteractionBuilder() *UserInteractionBuilder {
	return &UserInteractionBuilder{
		interaction: recommendation.UserInteraction{
			UserID:       uuid.New().String(),
			AssetID:      uuid.New().String(),
			InteractType: "VIEW",
			Timestamp:    time.Now(),
			Amount:       0,
			Duration:     30,
			Feedback:     0,
		},
	}
}

// WithUserID sets the UserID
func (b *UserInteractionBuilder) WithUserID(userID string) *UserInteractionBuilder {
	b.interaction.UserID = userID
	return b
}

// WithAssetID sets the AssetID
func (b *UserInteractionBuilder) WithAssetID(assetID string) *UserInteractionBuilder {
	b.interaction.AssetID = assetID
	return b
}

// WithInteractType sets the InteractType
func (b *UserInteractionBuilder) WithInteractType(interactType string) *UserInteractionBuilder {
	b.interaction.InteractType = interactType
	return b
}

// WithTimestamp sets the Timestamp
func (b *UserInteractionBuilder) WithTimestamp(timestamp time.Time) *UserInteractionBuilder {
	b.interaction.Timestamp = timestamp
	return b
}

// WithAmount sets the Amount
func (b *UserInteractionBuilder) WithAmount(amount float64) *UserInteractionBuilder {
	b.interaction.Amount = amount
	return b
}

// WithDuration sets the Duration
func (b *UserInteractionBuilder) WithDuration(duration int) *UserInteractionBuilder {
	b.interaction.Duration = duration
	return b
}

// WithFeedback sets the Feedback
func (b *UserInteractionBuilder) WithFeedback(feedback int) *UserInteractionBuilder {
	b.interaction.Feedback = feedback
	return b
}

// Build returns the built UserInteraction
func (b *UserInteractionBuilder) Build() recommendation.UserInteraction {
	return b.interaction
}

// RecommendationRequestBuilder helps build RecommendationRequest objects for testing
type RecommendationRequestBuilder struct {
	request recommendation.RecommendationRequest
}

// NewRecommendationRequestBuilder creates a new RecommendationRequestBuilder with default values
func NewRecommendationRequestBuilder() *RecommendationRequestBuilder {
	return &RecommendationRequestBuilder{
		request: recommendation.RecommendationRequest{
			UserID:           uuid.New().String(),
			Amount:           10000,
			TimeHorizon:      10,
			RiskTolerance:    0.5,
			PreferredSectors: []string{"TECHNOLOGY", "HEALTHCARE"},
			ExcludedSectors:  []string{"TOBACCO", "GAMBLING"},
			Goals:            []string{"RETIREMENT", "WEALTH_GROWTH"},
			Constraints:      []string{"LIQUIDITY", "TAX_EFFICIENCY"},
		},
	}
}

// WithUserID sets the UserID
func (b *RecommendationRequestBuilder) WithUserID(userID string) *RecommendationRequestBuilder {
	b.request.UserID = userID
	return b
}

// WithAmount sets the Amount
func (b *RecommendationRequestBuilder) WithAmount(amount float64) *RecommendationRequestBuilder {
	b.request.Amount = amount
	return b
}

// WithTimeHorizon sets the TimeHorizon
func (b *RecommendationRequestBuilder) WithTimeHorizon(timeHorizon int) *RecommendationRequestBuilder {
	b.request.TimeHorizon = timeHorizon
	return b
}

// WithRiskTolerance sets the RiskTolerance
func (b *RecommendationRequestBuilder) WithRiskTolerance(riskTolerance float64) *RecommendationRequestBuilder {
	b.request.RiskTolerance = riskTolerance
	return b
}

// WithPreferredSectors sets the PreferredSectors
func (b *RecommendationRequestBuilder) WithPreferredSectors(sectors []string) *RecommendationRequestBuilder {
	b.request.PreferredSectors = sectors
	return b
}

// WithExcludedSectors sets the ExcludedSectors
func (b *RecommendationRequestBuilder) WithExcludedSectors(sectors []string) *RecommendationRequestBuilder {
	b.request.ExcludedSectors = sectors
	return b
}

// WithGoals sets the Goals
func (b *RecommendationRequestBuilder) WithGoals(goals []string) *RecommendationRequestBuilder {
	b.request.Goals = goals
	return b
}

// WithConstraints sets the Constraints
func (b *RecommendationRequestBuilder) WithConstraints(constraints []string) *RecommendationRequestBuilder {
	b.request.Constraints = constraints
	return b
}

// Build returns the built RecommendationRequest
func (b *RecommendationRequestBuilder) Build() recommendation.RecommendationRequest {
	return b.request
}

// MarketDataBuilder helps build MarketData objects for testing
type MarketDataBuilder struct {
	data recommendation.MarketData
}

// NewMarketDataBuilder creates a new MarketDataBuilder with default values
func NewMarketDataBuilder() *MarketDataBuilder {
	return &MarketDataBuilder{
		data: recommendation.MarketData{
			Timestamp: time.Now(),
			MarketTrends: map[string]float64{
				"TECHNOLOGY": 0.05,
				"HEALTHCARE": 0.03,
				"FINANCE":    0.02,
				"ENERGY":     -0.01,
				"CONSUMER":   0.01,
			},
			EconomicIndicators: map[string]float64{
				"INFLATION":    0.03,
				"GDP_GROWTH":   0.025,
				"UNEMPLOYMENT": 0.04,
			},
			InterestRates: map[string]float64{
				"FED_RATE":    0.0025,
				"10Y_TREASURY": 0.015,
			},
			MarketSentiment: 0.2,
			SectorPerformance: map[string]float64{
				"TECHNOLOGY": 0.15,
				"HEALTHCARE": 0.10,
				"FINANCE":    0.08,
				"ENERGY":     0.05,
				"CONSUMER":   0.07,
			},
			AssetCorrelations: map[string][]string{
				"AAPL": {"MSFT", "GOOGL"},
				"MSFT": {"AAPL", "AMZN"},
			},
		},
	}
}

// WithTimestamp sets the Timestamp
func (b *MarketDataBuilder) WithTimestamp(timestamp time.Time) *MarketDataBuilder {
	b.data.Timestamp = timestamp
	return b
}

// WithMarketTrends sets the MarketTrends
func (b *MarketDataBuilder) WithMarketTrends(trends map[string]float64) *MarketDataBuilder {
	b.data.MarketTrends = trends
	return b
}

// WithEconomicIndicators sets the EconomicIndicators
func (b *MarketDataBuilder) WithEconomicIndicators(indicators map[string]float64) *MarketDataBuilder {
	b.data.EconomicIndicators = indicators
	return b
}

// WithInterestRates sets the InterestRates
func (b *MarketDataBuilder) WithInterestRates(rates map[string]float64) *MarketDataBuilder {
	b.data.InterestRates = rates
	return b
}

// WithMarketSentiment sets the MarketSentiment
func (b *MarketDataBuilder) WithMarketSentiment(sentiment float64) *MarketDataBuilder {
	b.data.MarketSentiment = sentiment
	return b
}

// WithSectorPerformance sets the SectorPerformance
func (b *MarketDataBuilder) WithSectorPerformance(performance map[string]float64) *MarketDataBuilder {
	b.data.SectorPerformance = performance
	return b
}

// WithAssetCorrelations sets the AssetCorrelations
func (b *MarketDataBuilder) WithAssetCorrelations(correlations map[string][]string) *MarketDataBuilder {
	b.data.AssetCorrelations = correlations
	return b
}

// Build returns the built MarketData
func (b *MarketDataBuilder) Build() recommendation.MarketData {
	return b.data
}
