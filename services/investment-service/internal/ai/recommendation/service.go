package recommendation

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

// Service provides recommendation functionality to the application
type Service struct {
	engine    Engine
	repository Repository
}

// NewService creates a new recommendation service
func NewService(engine Engine, repository Repository) *Service {
	return &Service{
		engine:    engine,
		repository: repository,
	}
}

// GetRecommendation generates investment recommendations for a user
func (s *Service) GetRecommendation(ctx context.Context, request RecommendationRequest) (*PortfolioRecommendation, error) {
	// Validate request
	if err := s.validateRecommendationRequest(request); err != nil {
		return nil, err
	}

	// Log request
	log.Printf("Generating recommendation for user %s with risk tolerance %.2f and time horizon %d years",
		request.UserID, request.RiskTolerance, request.TimeHorizon)

	// Get recommendation from engine
	recommendation, err := s.engine.GetRecommendation(ctx, request)
	if err != nil {
		log.Printf("Error generating recommendation: %v", err)
		return nil, err
	}

	// Log success
	log.Printf("Successfully generated recommendation with %d assets for user %s",
		len(recommendation.RecommendedAssets), request.UserID)

	return recommendation, nil
}

// GetPersonalizedAssets returns personalized investment assets for a user
func (s *Service) GetPersonalizedAssets(ctx context.Context, userID string, count int) ([]RecommendedAsset, error) {
	// Validate parameters
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if count <= 0 {
		count = 5 // Default count
	}

	// Log request
	log.Printf("Getting %d personalized assets for user %s", count, userID)

	// Get personalized assets from engine
	assets, err := s.engine.GetPersonalizedAssets(ctx, userID, count)
	if err != nil {
		log.Printf("Error getting personalized assets: %v", err)
		return nil, err
	}

	// Log success
	log.Printf("Successfully retrieved %d personalized assets for user %s", len(assets), userID)

	return assets, nil
}

// GetSimilarAssets returns assets similar to the given asset
func (s *Service) GetSimilarAssets(ctx context.Context, assetID string, count int) ([]RecommendedAsset, error) {
	// Validate parameters
	if assetID == "" {
		return nil, errors.New("asset ID is required")
	}
	if count <= 0 {
		count = 5 // Default count
	}

	// Log request
	log.Printf("Getting %d similar assets for asset %s", count, assetID)

	// Get similar assets from engine
	assets, err := s.engine.GetSimilarAssets(ctx, assetID, count)
	if err != nil {
		log.Printf("Error getting similar assets: %v", err)
		return nil, err
	}

	// Log success
	log.Printf("Successfully retrieved %d similar assets for asset %s", len(assets), assetID)

	return assets, nil
}

// GetDiversificationSuggestions returns suggestions to diversify a portfolio
func (s *Service) GetDiversificationSuggestions(ctx context.Context, userID string) ([]RecommendedAsset, error) {
	// Validate parameters
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Log request
	log.Printf("Getting diversification suggestions for user %s", userID)

	// Get diversification suggestions from engine
	assets, err := s.engine.GetDiversificationSuggestions(ctx, userID)
	if err != nil {
		log.Printf("Error getting diversification suggestions: %v", err)
		return nil, err
	}

	// Log success
	log.Printf("Successfully retrieved %d diversification suggestions for user %s", len(assets), userID)

	return assets, nil
}

// GetRebalancingSuggestions returns suggestions to rebalance a portfolio
func (s *Service) GetRebalancingSuggestions(ctx context.Context, userID string) ([]RecommendedAsset, error) {
	// Validate parameters
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Log request
	log.Printf("Getting rebalancing suggestions for user %s", userID)

	// Get rebalancing suggestions from engine
	assets, err := s.engine.GetRebalancingSuggestions(ctx, userID)
	if err != nil {
		log.Printf("Error getting rebalancing suggestions: %v", err)
		return nil, err
	}

	// Log success
	log.Printf("Successfully retrieved %d rebalancing suggestions for user %s", len(assets), userID)

	return assets, nil
}

// TrackUserInteraction tracks a user's interaction with an investment asset
func (s *Service) TrackUserInteraction(ctx context.Context, userID, assetID, interactType string, amount float64, duration int, feedback int) error {
	// Validate parameters
	if userID == "" {
		return errors.New("user ID is required")
	}
	if assetID == "" {
		return errors.New("asset ID is required")
	}
	if interactType == "" {
		return errors.New("interaction type is required")
	}

	// Create interaction
	interaction := UserInteraction{
		ID:           uuid.New().String(),
		UserID:       userID,
		AssetID:      assetID,
		InteractType: interactType,
		Rating:       0, // Will be calculated by the engine
		Timestamp:    time.Now(),
		Amount:       amount,
		Duration:     duration,
		Feedback:     feedback,
	}

	// Log interaction
	log.Printf("Tracking user interaction: user %s, asset %s, type %s", userID, assetID, interactType)

	// Save interaction
	err := s.repository.SaveUserInteraction(ctx, &interaction)
	if err != nil {
		log.Printf("Error saving user interaction: %v", err)
		return err
	}

	return nil
}

// UpdateUserProfile updates a user's investment profile
func (s *Service) UpdateUserProfile(ctx context.Context, profile UserProfile) error {
	// Validate profile
	if profile.UserID == "" {
		return errors.New("user ID is required")
	}
	if profile.RiskTolerance < 0 || profile.RiskTolerance > 1 {
		return errors.New("risk tolerance must be between 0 and 1")
	}
	if profile.TimeHorizon <= 0 {
		return errors.New("time horizon must be positive")
	}

	// Set timestamps
	profile.UpdatedAt = time.Now()
	if profile.CreatedAt.IsZero() {
		profile.CreatedAt = profile.UpdatedAt
	}

	// Log update
	log.Printf("Updating user profile for user %s", profile.UserID)

	// Save profile
	err := s.repository.SaveUserProfile(ctx, &profile)
	if err != nil {
		log.Printf("Error saving user profile: %v", err)
		return err
	}

	return nil
}

// TrainModel trains the recommendation model
func (s *Service) TrainModel(ctx context.Context) error {
	// Log training start
	log.Printf("Starting model training")

	// Train model
	err := s.engine.TrainModel(ctx)
	if err != nil {
		log.Printf("Error training model: %v", err)
		return err
	}

	// Log training completion
	log.Printf("Model training completed successfully")

	return nil
}

// GetModelMetrics returns the performance metrics of the recommendation model
func (s *Service) GetModelMetrics(ctx context.Context) (*ModelPerformanceMetrics, error) {
	// Log request
	log.Printf("Getting model metrics")

	// Get metrics
	metrics, err := s.engine.GetModelMetrics(ctx)
	if err != nil {
		log.Printf("Error getting model metrics: %v", err)
		return nil, err
	}

	return metrics, nil
}

// Helper methods

// validateRecommendationRequest validates a recommendation request
func (s *Service) validateRecommendationRequest(request RecommendationRequest) error {
	if request.UserID == "" {
		return errors.New("user ID is required")
	}
	if request.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	if request.TimeHorizon <= 0 {
		return errors.New("time horizon must be positive")
	}
	if request.RiskTolerance < 0 || request.RiskTolerance > 1 {
		return errors.New("risk tolerance must be between 0 and 1")
	}
	return nil
}
