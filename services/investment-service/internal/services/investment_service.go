package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"investment-service/internal/metrics"
	"investment-service/internal/models"
	"investment-service/internal/repositories"
	"investment-service/internal/validation"

	"github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Business errors
var (
	ErrInvestmentNotFound    = errors.New("investment not found")
	ErrInsufficientQuantity  = errors.New("insufficient quantity for transaction")
	ErrInvalidInvestmentType = errors.New("invalid investment type")
)

type InvestmentService interface {
	CreateInvestment(ctx context.Context, investment *models.Investment) error
	GetInvestment(ctx context.Context, id uuid.UUID) (*models.Investment, error)
	UpdateInvestment(ctx context.Context, investment *models.Investment) error
	ListInvestments(ctx context.Context, userID uuid.UUID, portfolioID *uuid.UUID, limit, offset int) ([]models.Investment, error)
	CreatePortfolio(ctx context.Context, portfolio *models.Portfolio) error
	GetPortfolio(ctx context.Context, id uuid.UUID) (*models.Portfolio, error)
	UpdatePortfolio(ctx context.Context, portfolio *models.Portfolio) error
	ListPortfolios(ctx context.Context, userID uuid.UUID) ([]models.Portfolio, error)
	CreateAsset(ctx context.Context, asset *models.Asset) error
	GetAsset(ctx context.Context, id uuid.UUID) (*models.Asset, error)
	UpdateAsset(ctx context.Context, asset *models.Asset) error
	ListAssets(ctx context.Context, assetType string, riskLevel string, limit, offset int) ([]models.Asset, error)
	ProcessTransaction(ctx context.Context, transaction *models.Transaction) error
	GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	ListTransactions(ctx context.Context, investmentID uuid.UUID, limit, offset int) ([]models.Transaction, error)
	GetUserPreference(ctx context.Context, userID uuid.UUID) (*models.UserPreference, error)
	UpdateUserPreference(ctx context.Context, preference *models.UserPreference) error
}

type investmentService struct {
	repo   repositories.InvestmentRepository
	logger *zap.Logger
}

func NewInvestmentService(repo repositories.InvestmentRepository, logger *zap.Logger) InvestmentService {
	return &investmentService{
		repo:   repo,
		logger: logger,
	}
}

func (s *investmentService) CreateInvestment(ctx context.Context, investment *models.Investment) error {
	// Validate investment
	if err := validation.ValidateInvestment(investment); err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id": investment.UserID,
			"symbol":  investment.Symbol,
			"error":   err.Error(),
		}).Warn("Investment validation failed")
		metrics.RecordBusinessError("validation")
		return err
	}

	// Calculate amount if not provided
	if investment.Amount == 0 && investment.Quantity > 0 && investment.PurchasePrice > 0 {
		investment.Amount = investment.Quantity * investment.PurchasePrice
	}

	// Set default status if not provided
	if investment.Status == "" {
		investment.Status = "ACTIVE"
	}

	// Set purchase date if not provided
	if investment.PurchaseDate.IsZero() {
		investment.PurchaseDate = time.Now()
	}

	// Create investment
	err := s.repo.Create(ctx, investment)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id": investment.UserID,
			"symbol":  investment.Symbol,
			"error":   err.Error(),
		}).Error("Failed to create investment")
		return fmt.Errorf("failed to create investment: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"id":       investment.ID,
		"user_id":  investment.UserID,
		"symbol":   investment.Symbol,
		"quantity": investment.Quantity,
		"amount":   investment.Amount,
	}).Info("Investment created successfully")

	// Create transaction record
	transaction := &models.Transaction{
		InvestmentID:   investment.ID,
		Type:           "buy",
		Amount:         investment.Amount,
		Currency:       investment.Currency,
		Status:         "pending",
		TransactionFee: calculateTransactionFee(investment.Amount),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.CreateTransaction(ctx, transaction); err != nil {
		s.logger.Error("Failed to create transaction", zap.Error(err))
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

func (s *investmentService) GetInvestment(ctx context.Context, id uuid.UUID) (*models.Investment, error) {
	investment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"id":    id,
			"error": err.Error(),
		}).Error("Failed to retrieve investment")
		return nil, fmt.Errorf("failed to retrieve investment: %w", err)
	}

	if investment == nil {
		metrics.RecordBusinessError("not_found")
		return nil, ErrInvestmentNotFound
	}

	return investment, nil
}

func (s *investmentService) UpdateInvestment(ctx context.Context, investment *models.Investment) error {
	// Validate investment
	if err := validation.ValidateInvestment(investment); err != nil {
		metrics.RecordBusinessError("validation")
		return err
	}

	// Check if investment exists
	existing, err := s.repo.GetByID(ctx, investment.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve investment: %w", err)
	}
	if existing == nil {
		metrics.RecordBusinessError("not_found")
		return ErrInvestmentNotFound
	}

	// Update investment
	err = s.repo.Update(ctx, investment)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"id":     investment.ID,
			"symbol": investment.Symbol,
			"error":  err.Error(),
		}).Error("Failed to update investment")
		return fmt.Errorf("failed to update investment: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"id":       investment.ID,
		"user_id":  investment.UserID,
		"symbol":   investment.Symbol,
		"quantity": investment.Quantity,
		"amount":   investment.Amount,
	}).Info("Investment updated successfully")

	return nil
}

func (s *investmentService) ListInvestments(ctx context.Context, userID uuid.UUID, portfolioID *uuid.UUID, limit, offset int) ([]models.Investment, error) {
	investments, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to retrieve investments for user")
		return nil, fmt.Errorf("failed to retrieve investments: %w", err)
	}

	return investments, nil
}

func (s *investmentService) CreatePortfolio(ctx context.Context, portfolio *models.Portfolio) error {
	// Set default values
	portfolio.TotalValue = 0
	portfolio.CreatedAt = time.Now()
	portfolio.UpdatedAt = time.Now()

	if err := s.repo.CreatePortfolio(ctx, portfolio); err != nil {
		s.logger.Error("Failed to create portfolio", zap.Error(err))
		return fmt.Errorf("failed to create portfolio: %w", err)
	}
	return nil
}

func (s *investmentService) GetPortfolio(ctx context.Context, id uuid.UUID) (*models.Portfolio, error) {
	portfolio, err := s.repo.GetPortfolio(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get portfolio", zap.Error(err))
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}
	return portfolio, nil
}

func (s *investmentService) UpdatePortfolio(ctx context.Context, portfolio *models.Portfolio) error {
	portfolio.UpdatedAt = time.Now()
	if err := s.repo.UpdatePortfolio(ctx, portfolio); err != nil {
		s.logger.Error("Failed to update portfolio", zap.Error(err))
		return fmt.Errorf("failed to update portfolio: %w", err)
	}
	return nil
}

func (s *investmentService) ListPortfolios(ctx context.Context, userID uuid.UUID) ([]models.Portfolio, error) {
	portfolios, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("Failed to retrieve portfolios for user")
		return nil, fmt.Errorf("failed to retrieve portfolios: %w", err)
	}

	return portfolios, nil
}

func (s *investmentService) CreateAsset(ctx context.Context, asset *models.Asset) error {
	// Set default values
	asset.CreatedAt = time.Now()
	asset.UpdatedAt = time.Now()

	if err := s.repo.CreateAsset(ctx, asset); err != nil {
		s.logger.Error("Failed to create asset", zap.Error(err))
		return fmt.Errorf("failed to create asset: %w", err)
	}
	return nil
}

func (s *investmentService) GetAsset(ctx context.Context, id uuid.UUID) (*models.Asset, error) {
	asset, err := s.repo.GetAsset(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get asset", zap.Error(err))
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}
	return asset, nil
}

func (s *investmentService) UpdateAsset(ctx context.Context, asset *models.Asset) error {
	asset.UpdatedAt = time.Now()
	if err := s.repo.UpdateAsset(ctx, asset); err != nil {
		s.logger.Error("Failed to update asset", zap.Error(err))
		return fmt.Errorf("failed to update asset: %w", err)
	}
	return nil
}

func (s *investmentService) ListAssets(ctx context.Context, assetType string, riskLevel string, limit, offset int) ([]models.Asset, error) {
	assets, err := s.repo.ListAssets(ctx, assetType, riskLevel, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list assets", zap.Error(err))
		return nil, fmt.Errorf("failed to list assets: %w", err)
	}
	return assets, nil
}

func (s *investmentService) ProcessTransaction(ctx context.Context, transaction *models.Transaction) error {
	// Get investment details
	investment, err := s.repo.GetInvestment(ctx, transaction.InvestmentID)
	if err != nil {
		return fmt.Errorf("failed to get investment: %w", err)
	}

	// Update investment status based on transaction
	switch transaction.Status {
	case "completed":
		investment.Status = "active"
	case "failed":
		investment.Status = "failed"
	case "cancelled":
		investment.Status = "cancelled"
	}

	investment.UpdatedAt = time.Now()
	if err := s.repo.UpdateInvestment(ctx, investment); err != nil {
		return fmt.Errorf("failed to update investment: %w", err)
	}

	// Update transaction
	transaction.UpdatedAt = time.Now()
	if err := s.repo.CreateTransaction(ctx, transaction); err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

func (s *investmentService) GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	transaction, err := s.repo.GetTransaction(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	return transaction, nil
}

func (s *investmentService) ListTransactions(ctx context.Context, investmentID uuid.UUID, limit, offset int) ([]models.Transaction, error) {
	transactions, err := s.repo.ListTransactions(ctx, investmentID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list transactions", zap.Error(err))
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	return transactions, nil
}

func (s *investmentService) GetUserPreference(ctx context.Context, userID uuid.UUID) (*models.UserPreference, error) {
	preference, err := s.repo.GetUserPreference(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user preference", zap.Error(err))
		return nil, fmt.Errorf("failed to get user preference: %w", err)
	}
	return preference, nil
}

func (s *investmentService) UpdateUserPreference(ctx context.Context, preference *models.UserPreference) error {
	preference.UpdatedAt = time.Now()
	if err := s.repo.UpdateUserPreference(ctx, preference); err != nil {
		s.logger.Error("Failed to update user preference", zap.Error(err))
		return fmt.Errorf("failed to update user preference: %w", err)
	}
	return nil
}

func calculateTransactionFee(amount float64) float64 {
	// Example fee calculation: 0.1% for amounts over $1000, 0.2% otherwise
	if amount > 1000 {
		return amount * 0.001
	}
	return amount * 0.002
}
