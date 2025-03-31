package ai

import (
	"context"
	"time"

	"github.com/sparkfund/services/investment-service/internal/models"
	"github.com/sparkfund/services/investment-service/internal/repositories"
)

// FeatureService handles feature extraction and engineering for ML models
type FeatureService struct {
	userRepo       repositories.UserRepository
	portfolioRepo  repositories.PortfolioRepository
	investmentRepo repositories.InvestmentRepository
	marketDataSvc  MarketDataService
	logger         *logrus.Logger
}

// NewFeatureService creates a new feature service
func NewFeatureService(
	userRepo repositories.UserRepository,
	portfolioRepo repositories.PortfolioRepository,
	investmentRepo repositories.InvestmentRepository,
	marketDataSvc MarketDataService,
	logger *logrus.Logger,
) *FeatureService {
	return &FeatureService{
		userRepo:       userRepo,
		portfolioRepo:  portfolioRepo,
		investmentRepo: investmentRepo,
		marketDataSvc:  marketDataSvc,
		logger:         logger,
	}
}

// GetUserFeatures extracts features for user-based predictions
func (s *FeatureService) GetUserFeatures(ctx context.Context, userID uint) (map[string]interface{}, error) {
	features := make(map[string]interface{})

	// Get user profile data (anonymized/numerical features only)
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Basic user features
	features["account_age_days"] = time.Since(user.CreatedAt).Hours() / 24
	features["risk_profile"] = user.RiskTolerance

	// Get user's historical investment behavior
	investments, err := s.investmentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Calculate behavioral features
	var totalInvested float64
	investmentTypes := make(map[string]float64)
	sectorAllocation := make(map[string]float64)

	for _, inv := range investments {
		totalInvested += inv.Amount
		investmentTypes[inv.Type] += inv.Amount

		// Get sector for this investment
		sector, err := s.marketDataSvc.GetEquitySector(ctx, inv.Symbol)
		if err == nil {
			sectorAllocation[sector] += inv.Amount
		}
	}

	// Add derived features
	features["total_invested"] = totalInvested
	features["avg_investment_size"] = totalInvested / float64(max(1, len(investments)))
	features["investment_count"] = len(investments)

	// Convert investment types to features
	for typ, amount := range investmentTypes {
		features["investment_type_"+typ] = amount / totalInvested
	}

	// Convert sector allocation to features
	for sector, amount := range sectorAllocation {
		features["sector_"+sector] = amount / totalInvested
	}

	// Get market context features
	marketFeatures, err := s.marketDataSvc.GetMarketFeatures(ctx)
	if err == nil {
		for k, v := range marketFeatures {
			features["market_"+k] = v
		}
	}

	return features, nil
}

// GetPortfolioFeatures extracts features for portfolio analysis
func (s *FeatureService) GetPortfolioFeatures(ctx context.Context, portfolio models.Portfolio) (map[string]interface{}, error) {
	features := make(map[string]interface{})

	// Basic portfolio features
	features["portfolio_id"] = portfolio.ID
	features["portfolio_age_days"] = time.Since(portfolio.CreatedAt).Hours() / 24

	// Get investments in this portfolio
	investments, err := s.investmentRepo.GetByPortfolioID(ctx, portfolio.ID)
	if err != nil {
		return nil, err
	}

	// Calculate portfolio composition
	var totalValue float64
	assetTypes := make(map[string]float64)
	symbols := make(map[string]float64)
	sectors := make(map[string]float64)

	for _, inv := range investments {
		// Get current value of investment
		currentPrice, err := s.marketDataSvc.GetCurrentPrice(ctx, inv.Symbol)
		if err != nil {
			s.logger.WithError(err).WithField("symbol", inv.Symbol).Warn("Failed to get current price")
			currentPrice = inv.PurchasePrice // Fallback to purchase price
		}

		value := currentPrice * inv.Quantity
		totalValue += value

		assetTypes[inv.Type] += value
		symbols[inv.Symbol] += value

		// Get sector data
		sector, err := s.marketDataSvc.GetEquitySector(ctx, inv.Symbol)
		if err == nil {
			sectors[sector] += value
		}

		// Add individual holding features
		features["holding_"+inv.Symbol+"_weight"] = value
		features["holding_"+inv.Symbol+"_price_change"] = (currentPrice - inv.PurchasePrice) / inv.PurchasePrice
	}

	// Add composition features
	for typ, value := range assetTypes {
		features["asset_type_"+typ] = value / totalValue
	}

	// Add concentration features
	for sector, value := range sectors {
		features["sector_"+sector] = value / totalValue
	}

	// Calculate concentration metrics
	features["top_holding_concentration"] = calculateTopConcentration(symbols, totalValue, 3)
	features["sector_concentration"] = calculateHerfindahlIndex(sectors, totalValue)

	// Get correlation features between holdings
	correlationMatrix, err := s.marketDataSvc.GetCorrelationMatrix(ctx, maps.Keys(symbols))
	if err == nil {
		features["avg_correlation"] = calculateAverageCorrelation(correlationMatrix)
	}

	// Market environment features
	marketFeatures, err := s.marketDataSvc.GetMarketFeatures(ctx)
	if err == nil {
		for k, v := range marketFeatures {
			features["market_"+k] = v
		}
	}

	return features, nil
}
