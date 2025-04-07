package tests

import (
	"context"
	"time"

	"github.com/sparkfund/services/investment-service/internal/ai/recommendation"
)

// MockRepository is a mock implementation of the recommendation.Repository interface for testing
type MockRepository struct {
	UserProfiles            map[string]recommendation.UserProfile
	InvestmentAssets        map[string]recommendation.InvestmentAsset
	MarketData              []recommendation.MarketData
	UserInteractions        []recommendation.UserInteraction
	PortfolioRecommendations map[string][]recommendation.PortfolioRecommendation
	ModelMetrics            recommendation.ModelPerformanceMetrics
	SimilarUsers            map[string][]string
	PopularAssets           []recommendation.InvestmentAsset
	UserPortfolios          map[string][]recommendation.InvestmentAsset
	
	// Error simulation
	SimulateErrors          bool
	ErrorOnGetUserProfile   bool
	ErrorOnGetAssets        bool
	ErrorOnGetMarketData    bool
	ErrorOnSaveRecommendation bool
}

// NewMockRepository creates a new MockRepository
func NewMockRepository() *MockRepository {
	return &MockRepository{
		UserProfiles:            make(map[string]recommendation.UserProfile),
		InvestmentAssets:        make(map[string]recommendation.InvestmentAsset),
		MarketData:              []recommendation.MarketData{},
		UserInteractions:        []recommendation.UserInteraction{},
		PortfolioRecommendations: make(map[string][]recommendation.PortfolioRecommendation),
		SimilarUsers:            make(map[string][]string),
		UserPortfolios:          make(map[string][]recommendation.InvestmentAsset),
		SimulateErrors:          false,
	}
}

// GetUserProfile retrieves a user profile by ID
func (r *MockRepository) GetUserProfile(ctx context.Context, userID string) (*recommendation.UserProfile, error) {
	if r.SimulateErrors && r.ErrorOnGetUserProfile {
		return nil, recommendation.ErrUserNotFound
	}
	
	profile, exists := r.UserProfiles[userID]
	if !exists {
		return nil, recommendation.ErrUserNotFound
	}
	
	return &profile, nil
}

// SaveUserProfile saves a user profile
func (r *MockRepository) SaveUserProfile(ctx context.Context, profile *recommendation.UserProfile) error {
	r.UserProfiles[profile.UserID] = *profile
	return nil
}

// GetInvestmentAsset retrieves an investment asset by ID
func (r *MockRepository) GetInvestmentAsset(ctx context.Context, assetID string) (*recommendation.InvestmentAsset, error) {
	asset, exists := r.InvestmentAssets[assetID]
	if !exists {
		return nil, recommendation.ErrInsufficientData
	}
	
	return &asset, nil
}

// GetInvestmentAssets retrieves investment assets based on a filter
func (r *MockRepository) GetInvestmentAssets(ctx context.Context, filter map[string]interface{}, limit int) ([]recommendation.InvestmentAsset, error) {
	if r.SimulateErrors && r.ErrorOnGetAssets {
		return nil, recommendation.ErrInsufficientData
	}
	
	// Simple implementation that ignores the filter and just returns all assets up to the limit
	assets := make([]recommendation.InvestmentAsset, 0, len(r.InvestmentAssets))
	for _, asset := range r.InvestmentAssets {
		assets = append(assets, asset)
		if len(assets) >= limit {
			break
		}
	}
	
	return assets, nil
}

// SaveInvestmentAsset saves an investment asset
func (r *MockRepository) SaveInvestmentAsset(ctx context.Context, asset *recommendation.InvestmentAsset) error {
	r.InvestmentAssets[asset.ID] = *asset
	return nil
}

// GetLatestMarketData retrieves the latest market data
func (r *MockRepository) GetLatestMarketData(ctx context.Context) (*recommendation.MarketData, error) {
	if r.SimulateErrors && r.ErrorOnGetMarketData {
		return nil, recommendation.ErrInsufficientData
	}
	
	if len(r.MarketData) == 0 {
		return nil, recommendation.ErrInsufficientData
	}
	
	// Return the most recent market data
	latestData := r.MarketData[len(r.MarketData)-1]
	return &latestData, nil
}

// GetHistoricalMarketData retrieves historical market data within a time range
func (r *MockRepository) GetHistoricalMarketData(ctx context.Context, startTime, endTime time.Time) ([]recommendation.MarketData, error) {
	if r.SimulateErrors && r.ErrorOnGetMarketData {
		return nil, recommendation.ErrInsufficientData
	}
	
	// Filter market data by time range
	var filteredData []recommendation.MarketData
	for _, data := range r.MarketData {
		if (data.Timestamp.Equal(startTime) || data.Timestamp.After(startTime)) &&
		   (data.Timestamp.Equal(endTime) || data.Timestamp.Before(endTime)) {
			filteredData = append(filteredData, data)
		}
	}
	
	return filteredData, nil
}

// SaveMarketData saves market data
func (r *MockRepository) SaveMarketData(ctx context.Context, data *recommendation.MarketData) error {
	r.MarketData = append(r.MarketData, *data)
	return nil
}

// GetUserInteractions retrieves user interactions for a user
func (r *MockRepository) GetUserInteractions(ctx context.Context, userID string, limit int) ([]recommendation.UserInteraction, error) {
	var userInteractions []recommendation.UserInteraction
	count := 0
	
	for _, interaction := range r.UserInteractions {
		if interaction.UserID == userID {
			userInteractions = append(userInteractions, interaction)
			count++
			if count >= limit {
				break
			}
		}
	}
	
	return userInteractions, nil
}

// SaveUserInteraction saves a user interaction
func (r *MockRepository) SaveUserInteraction(ctx context.Context, interaction *recommendation.UserInteraction) error {
	r.UserInteractions = append(r.UserInteractions, *interaction)
	return nil
}

// SavePortfolioRecommendation saves a portfolio recommendation
func (r *MockRepository) SavePortfolioRecommendation(ctx context.Context, recommendation *recommendation.PortfolioRecommendation) error {
	if r.SimulateErrors && r.ErrorOnSaveRecommendation {
		return recommendation.ErrRecommendationFailed
	}
	
	userID := recommendation.UserID
	if _, exists := r.PortfolioRecommendations[userID]; !exists {
		r.PortfolioRecommendations[userID] = []recommendation.PortfolioRecommendation{}
	}
	
	r.PortfolioRecommendations[userID] = append(r.PortfolioRecommendations[userID], *recommendation)
	return nil
}

// GetUserRecommendations retrieves portfolio recommendations for a user
func (r *MockRepository) GetUserRecommendations(ctx context.Context, userID string, limit int) ([]recommendation.PortfolioRecommendation, error) {
	recommendations, exists := r.PortfolioRecommendations[userID]
	if !exists {
		return []recommendation.PortfolioRecommendation{}, nil
	}
	
	if len(recommendations) <= limit {
		return recommendations, nil
	}
	
	return recommendations[len(recommendations)-limit:], nil
}

// SaveModelMetrics saves model performance metrics
func (r *MockRepository) SaveModelMetrics(ctx context.Context, metrics *recommendation.ModelPerformanceMetrics) error {
	r.ModelMetrics = *metrics
	return nil
}

// GetLatestModelMetrics retrieves the latest model performance metrics
func (r *MockRepository) GetLatestModelMetrics(ctx context.Context) (*recommendation.ModelPerformanceMetrics, error) {
	return &r.ModelMetrics, nil
}

// GetSimilarUsers retrieves users similar to the given user
func (r *MockRepository) GetSimilarUsers(ctx context.Context, userID string, limit int) ([]string, error) {
	similarUsers, exists := r.SimilarUsers[userID]
	if !exists {
		return []string{}, nil
	}
	
	if len(similarUsers) <= limit {
		return similarUsers, nil
	}
	
	return similarUsers[:limit], nil
}

// GetPopularAssets retrieves popular investment assets
func (r *MockRepository) GetPopularAssets(ctx context.Context, limit int) ([]recommendation.InvestmentAsset, error) {
	if len(r.PopularAssets) <= limit {
		return r.PopularAssets, nil
	}
	
	return r.PopularAssets[:limit], nil
}

// GetUserPortfolio retrieves a user's current portfolio
func (r *MockRepository) GetUserPortfolio(ctx context.Context, userID string) ([]recommendation.InvestmentAsset, error) {
	portfolio, exists := r.UserPortfolios[userID]
	if !exists {
		return []recommendation.InvestmentAsset{}, nil
	}
	
	return portfolio, nil
}
