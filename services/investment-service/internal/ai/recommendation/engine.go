package recommendation

import (
	"context"
	"errors"
)

// Errors
var (
	ErrUserNotFound          = errors.New("user profile not found")
	ErrInsufficientData      = errors.New("insufficient data for recommendation")
	ErrModelNotReady         = errors.New("recommendation model not ready")
	ErrInvalidRequest        = errors.New("invalid recommendation request")
	ErrRecommendationFailed  = errors.New("recommendation generation failed")
)

// Engine defines the interface for the recommendation engine
type Engine interface {
	// GetRecommendation generates investment recommendations for a user
	GetRecommendation(ctx context.Context, request RecommendationRequest) (*PortfolioRecommendation, error)
	
	// GetPersonalizedAssets returns personalized investment assets for a user
	GetPersonalizedAssets(ctx context.Context, userID string, count int) ([]RecommendedAsset, error)
	
	// GetSimilarAssets returns assets similar to the given asset
	GetSimilarAssets(ctx context.Context, assetID string, count int) ([]RecommendedAsset, error)
	
	// GetDiversificationSuggestions returns suggestions to diversify a portfolio
	GetDiversificationSuggestions(ctx context.Context, userID string) ([]RecommendedAsset, error)
	
	// GetRebalancingSuggestions returns suggestions to rebalance a portfolio
	GetRebalancingSuggestions(ctx context.Context, userID string) ([]RecommendedAsset, error)
	
	// TrainModel trains the recommendation model with new data
	TrainModel(ctx context.Context) error
	
	// GetModelMetrics returns the performance metrics of the recommendation model
	GetModelMetrics(ctx context.Context) (*ModelPerformanceMetrics, error)
}
