package ai

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/sparkfund/services/investment-service/internal/config"
	"github.com/sparkfund/services/investment-service/internal/models"
)

// ModelService manages machine learning model predictions
type ModelService struct {
	modelClient    MLClient
	featureService FeatureService
	cache          *redis.Client
	logger         *logrus.Logger
	config         *config.AIConfig
}

// NewModelService creates a new model service
func NewModelService(
	modelClient MLClient,
	featureService FeatureService,
	cache *redis.Client,
	logger *logrus.Logger,
	config *config.AIConfig,
) *ModelService {
	return &ModelService{
		modelClient:    modelClient,
		featureService: featureService,
		cache:          cache,
		logger:         logger,
		config:         config,
	}
}

// GetInvestmentRecommendations returns personalized investment recommendations
func (s *ModelService) GetInvestmentRecommendations(ctx context.Context, userID uint, amount float64) ([]models.InvestmentRecommendation, error) {
	// Start tracing span for observability
	span, ctx := opentracing.StartSpanFromContext(ctx, "ai.GetInvestmentRecommendations")
	defer span.Finish()

	// Try to get from cache first
	cacheKey := fmt.Sprintf("investment_rec:%d:%f", userID, amount)
	cachedResult, err := s.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var recommendations []models.InvestmentRecommendation
		if err := json.Unmarshal([]byte(cachedResult), &recommendations); err == nil {
			s.logger.WithField("user_id", userID).Debug("Retrieved investment recommendations from cache")
			return recommendations, nil
		}
	}

	// Get user features for the model
	features, err := s.featureService.GetUserFeatures(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user features: %w", err)
	}

	// Add context for the recommendation
	features["investment_amount"] = amount
	features["timestamp"] = float64(time.Now().Unix())

	// Call the ML model
	prediction, metadata, err := s.modelClient.Predict(ctx, "investment_recommendations", features)
	if err != nil {
		return nil, fmt.Errorf("model prediction failed: %w", err)
	}

	// Log metrics for monitoring
	s.recordPredictionMetrics("investment_recommendations", metadata)

	// Convert prediction to recommendations
	recommendations, err := s.convertToRecommendations(prediction)
	if err != nil {
		return nil, fmt.Errorf("failed to parse model output: %w", err)
	}

	// Store in cache
	if len(recommendations) > 0 {
		jsonData, _ := json.Marshal(recommendations)
		s.cache.Set(ctx, cacheKey, jsonData, 1*time.Hour)
	}

	return recommendations, nil
}

// GetRiskAssessment calculates investment risk for a portfolio
func (s *ModelService) GetRiskAssessment(ctx context.Context, portfolio models.Portfolio) (*models.RiskAssessment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ai.GetRiskAssessment")
	defer span.Finish()

	// Extract features from portfolio
	features, err := s.featureService.GetPortfolioFeatures(ctx, portfolio)
	if err != nil {
		return nil, fmt.Errorf("failed to extract portfolio features: %w", err)
	}

	// Call risk scoring model
	prediction, metadata, err := s.modelClient.Predict(ctx, "portfolio_risk_assessment", features)
	if err != nil {
		return nil, fmt.Errorf("risk model prediction failed: %w", err)
	}

	// Record model performance metrics
	s.recordPredictionMetrics("portfolio_risk_assessment", metadata)

	// Parse model output
	riskAssessment := &models.RiskAssessment{
		OverallScore:   prediction["risk_score"].(float64),
		VolatilityRisk: prediction["volatility_risk"].(float64),
		MarketRisk:     prediction["market_risk"].(float64),
		SectorRisk:     prediction["sector_risk"].(float64),
		ModelVersion:   metadata.ModelVersion,
		CalculatedAt:   time.Now(),
		Factors:        make(map[string]float64),
	}

	// Parse risk factors
	for k, v := range prediction {
		if strings.HasPrefix(k, "factor_") {
			riskAssessment.Factors[strings.TrimPrefix(k, "factor_")] = v.(float64)
		}
	}

	// Get explanations for regulatory compliance
	explanations, err := s.getModelExplanations(ctx, "portfolio_risk_assessment", features)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to get model explanations")
	} else {
		riskAssessment.Explanations = explanations
	}

	return riskAssessment, nil
}

// recordPredictionMetrics records model performance metrics
func (s *ModelService) recordPredictionMetrics(modelName string, metadata ModelMetadata) {
	metrics.ModelLatency.WithLabelValues(modelName, metadata.ModelVersion).Observe(metadata.LatencyMs)
	metrics.ModelPredictionCounter.WithLabelValues(modelName, metadata.ModelVersion).Inc()
}
