package recommendation

import (
	"context"
	"time"
)

// Repository defines the interface for accessing data needed by the recommendation engine
type Repository interface {
	// User profile methods
	GetUserProfile(ctx context.Context, userID string) (*UserProfile, error)
	SaveUserProfile(ctx context.Context, profile *UserProfile) error
	
	// Investment asset methods
	GetInvestmentAsset(ctx context.Context, assetID string) (*InvestmentAsset, error)
	GetInvestmentAssets(ctx context.Context, filter map[string]interface{}, limit int) ([]InvestmentAsset, error)
	SaveInvestmentAsset(ctx context.Context, asset *InvestmentAsset) error
	
	// Market data methods
	GetLatestMarketData(ctx context.Context) (*MarketData, error)
	GetHistoricalMarketData(ctx context.Context, startTime, endTime time.Time) ([]MarketData, error)
	SaveMarketData(ctx context.Context, data *MarketData) error
	
	// User interaction methods
	GetUserInteractions(ctx context.Context, userID string, limit int) ([]UserInteraction, error)
	SaveUserInteraction(ctx context.Context, interaction *UserInteraction) error
	
	// Recommendation methods
	SavePortfolioRecommendation(ctx context.Context, recommendation *PortfolioRecommendation) error
	GetUserRecommendations(ctx context.Context, userID string, limit int) ([]PortfolioRecommendation, error)
	
	// Model performance methods
	SaveModelMetrics(ctx context.Context, metrics *ModelPerformanceMetrics) error
	GetLatestModelMetrics(ctx context.Context) (*ModelPerformanceMetrics, error)
	
	// Collaborative filtering methods
	GetSimilarUsers(ctx context.Context, userID string, limit int) ([]string, error)
	GetPopularAssets(ctx context.Context, limit int) ([]InvestmentAsset, error)
	GetUserPortfolio(ctx context.Context, userID string) ([]InvestmentAsset, error)
}
