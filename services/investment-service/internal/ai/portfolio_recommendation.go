package ai

import (
    "context"
    "fmt"
    "time"

    "github.com/sirupsen/logrus"
    "github.com/sparkfund/services/investment-service/internal/models"
)

// PortfolioRecommendationService uses AI for generating portfolio recommendations
type PortfolioRecommendationService struct {
    portfolioClient PortfolioAIClient
    marketDataAPI   MarketDataAPIClient
    userProfileAPI  UserProfileAPIClient
    logger          *logrus.Logger
    config          *config.AIConfig
}

// NewPortfolioRecommendationService creates a new portfolio recommendation service
func NewPortfolioRecommendationService(
    portfolioClient PortfolioAIClient,
    marketDataAPI MarketDataAPIClient,
    userProfileAPI UserProfileAPIClient,
    logger *logrus.Logger,
    config *config.AIConfig,
) *PortfolioRecommendationService {
    return &PortfolioRecommendationService{
        portfolioClient: portfolioClient,
        marketDataAPI:   marketDataAPI,
        userProfileAPI:  userProfileAPI,
        logger:          logger,
        config:          config,
    }
}

// GenerateRecommendations generates personalized investment recommendations
func (s *PortfolioRecommendationService) GenerateRecommendations(
    ctx context.Context, 
    userID uint, 
    amount float64,
    riskTolerance string,
    investmentHorizon int, // in months
) (*models.PortfolioRecommendation, error) {
    span, ctx := opentracing.StartSpanFromContext(ctx, "ai.GenerateRecommendations")
    defer span.Finish()
    
    // Get user profile data
    userProfile, err := s.userProfileAPI.GetUserProfile(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user profile: %w", err)
    }
    
    // Get market data
    marketData, err := s.marketDataAPI.GetMarketData(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get market data: %w", err)
    }
    
    // Get user's existing investments
    existingInvestments, err := s.getUserExistingInvestments(ctx, userID)
    if err != nil {
        s.logger.WithError(err).Warn("Failed to get existing investments, proceeding with limited data")
    }
    
    // Prepare recommendation request
    request := &PortfolioRecommendationRequest{
        UserID:            userID,
        InvestmentAmount:  amount,
        RiskTolerance:     riskTolerance,
        InvestmentHorizon: investmentHorizon,
        Age:               userProfile.Age,
        Income:            userProfile.Income,
        Location:          userProfile.Location,
        InvestmentGoals:   userProfile.InvestmentGoals,
        ExistingHoldings:  existingInvestments,
        MarketConditions:  marketData.MarketConditions,
        SectorOutlooks:    marketData.SectorOutlooks,
        EconomicIndicators: marketData.EconomicIndicators,
    }
    
    // Get AI recommendation
    recommendation, err := s.portfolioClient.RecommendPortfolio(ctx, request)
    if err != nil {
        return nil, fmt.Errorf("AI recommendation failed: %w", err)
    }
    
    // Transform to response model
    result := &models.PortfolioRecommendation{
        UserID:           userID,
        RecommendationID: generateUUID(),
        Amount:           amount,
        RiskLevel:        riskTolerance,
        ExpectedReturn:   recommendation.ExpectedReturn,
        Volatility:       recommendation.Volatility,
        TimeHorizon:      investmentHorizon,
        GeneratedAt:      time.Now().UTC(),
        ModelVersion:     s.portfolioClient.GetModelVersion(),
        Allocations:      make([]models.AssetAllocation, 0, len(recommendation.Allocations)),
        Rationale:        recommendation.Rationale,
    }
    
    // Convert allocations
    for _, alloc := range recommendation.Allocations {
        result.Allocations = append(result.Allocations, models.AssetAllocation{
            AssetType:    alloc.AssetType,
            Symbol:       alloc.Symbol,
            Name:         alloc.Name,
            Percentage:   alloc.Percentage,
            Amount:       amount * alloc.Percentage / 100,
            Reason:       alloc.Reason,
        })
    }
    
    // Log the recommendation
    s.logger.WithFields// filepath: c:\Users\adilm\repositories\Go\sparkfund\services\investment-service\internal\ai\portfolio_recommendation.go
package ai

import (
    "context"
    "fmt"
    "time"

    "github.com/sirupsen/logrus"
    "github.com/sparkfund/services/investment-service/internal/models"
)

// PortfolioRecommendationService uses AI for generating portfolio recommendations
type PortfolioRecommendationService struct {
    portfolioClient PortfolioAIClient
    marketDataAPI   MarketDataAPIClient
    userProfileAPI  UserProfileAPIClient
    logger          *logrus.Logger
    config          *config.AIConfig
}

// NewPortfolioRecommendationService creates a new portfolio recommendation service
func NewPortfolioRecommendationService(
    portfolioClient PortfolioAIClient,
    marketDataAPI MarketDataAPIClient,
    userProfileAPI UserProfileAPIClient,
    logger *logrus.Logger,
    config *config.AIConfig,
) *PortfolioRecommendationService {
    return &PortfolioRecommendationService{
        portfolioClient: portfolioClient,
        marketDataAPI:   marketDataAPI,
        userProfileAPI:  userProfileAPI,
        logger:          logger,
        config:          config,
    }
}

// GenerateRecommendations generates personalized investment recommendations
func (s *PortfolioRecommendationService) GenerateRecommendations(
    ctx context.Context, 
    userID uint, 
    amount float64,
    riskTolerance string,
    investmentHorizon int, // in months
) (*models.PortfolioRecommendation, error) {
    span, ctx := opentracing.StartSpanFromContext(ctx, "ai.GenerateRecommendations")
    defer span.Finish()
    
    // Get user profile data
    userProfile, err := s.userProfileAPI.GetUserProfile(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user profile: %w", err)
    }
    
    // Get market data
    marketData, err := s.marketDataAPI.GetMarketData(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get market data: %w", err)
    }
    
    // Get user's existing investments
    existingInvestments, err := s.getUserExistingInvestments(ctx, userID)
    if err != nil {
        s.logger.WithError(err).Warn("Failed to get existing investments, proceeding with limited data")
    }
    
    // Prepare recommendation request
    request := &PortfolioRecommendationRequest{
        UserID:            userID,
        InvestmentAmount:  amount,
        RiskTolerance:     riskTolerance,
        InvestmentHorizon: investmentHorizon,
        Age:               userProfile.Age,
        Income:            userProfile.Income,
        Location:          userProfile.Location,
        InvestmentGoals:   userProfile.InvestmentGoals,
        ExistingHoldings:  existingInvestments,
        MarketConditions:  marketData.MarketConditions,
        SectorOutlooks:    marketData.SectorOutlooks,
        EconomicIndicators: marketData.EconomicIndicators,
    }
    
    // Get AI recommendation
    recommendation, err := s.portfolioClient.RecommendPortfolio(ctx, request)
    if err != nil {
        return nil, fmt.Errorf("AI recommendation failed: %w", err)
    }
    
    // Transform to response model
    result := &models.PortfolioRecommendation{
        UserID:           userID,
        RecommendationID: generateUUID(),
        Amount:           amount,
        RiskLevel:        riskTolerance,
        ExpectedReturn:   recommendation.ExpectedReturn,
        Volatility:       recommendation.Volatility,
        TimeHorizon:      investmentHorizon,
        GeneratedAt:      time.Now().UTC(),
        ModelVersion:     s.portfolioClient.GetModelVersion(),
        Allocations:      make([]models.AssetAllocation, 0, len(recommendation.Allocations)),
        Rationale:        recommendation.Rationale,
    }
    
    // Convert allocations
    for _, alloc := range recommendation.Allocations {
        result.Allocations = append(result.Allocations, models.AssetAllocation{
            AssetType:    alloc.AssetType,
            Symbol:       alloc.Symbol,
            Name:         alloc.Name,
            Percentage:   alloc.Percentage,
            Amount:       amount * alloc.Percentage / 100,
            Reason:       alloc.Reason,
        })
    }
    
    // Log the recommendation
    s.logger.WithFields