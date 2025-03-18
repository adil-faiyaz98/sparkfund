package service

import (
	"context"
	"fmt"
	"time"

	"your-project/internal/investments"

	"github.com/google/uuid"
)

type investmentService struct {
	repo investments.InvestmentRepository
}

func NewInvestmentService(repo investments.InvestmentRepository) investments.InvestmentService {
	return &investmentService{repo: repo}
}

func (s *investmentService) CreateInvestment(ctx context.Context, investment *investments.Investment) error {
	// Validate investment type
	if !isValidInvestmentType(investment.Type) {
		return fmt.Errorf("invalid investment type: %s", investment.Type)
	}

	// Set initial values
	investment.PurchaseDate = time.Now()
	investment.CurrentPrice = investment.PurchasePrice

	return s.repo.Create(investment)
}

func (s *investmentService) GetInvestment(ctx context.Context, id uuid.UUID) (*investments.Investment, error) {
	return s.repo.GetByID(id)
}

func (s *investmentService) GetUserInvestments(ctx context.Context, userID uuid.UUID) ([]*investments.Investment, error) {
	return s.repo.GetByUserID(userID)
}

func (s *investmentService) GetAccountInvestments(ctx context.Context, accountID uuid.UUID) ([]*investments.Investment, error) {
	return s.repo.GetByAccountID(accountID)
}

func (s *investmentService) UpdateInvestment(ctx context.Context, investment *investments.Investment) error {
	// Validate investment exists
	existing, err := s.repo.GetByID(investment.ID)
	if err != nil {
		return fmt.Errorf("investment not found: %w", err)
	}

	// Validate investment type if changed
	if existing.Type != investment.Type && !isValidInvestmentType(investment.Type) {
		return fmt.Errorf("invalid investment type: %s", investment.Type)
	}

	return s.repo.Update(investment)
}

func (s *investmentService) DeleteInvestment(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *investmentService) GetInvestmentsBySymbol(ctx context.Context, symbol string) ([]*investments.Investment, error) {
	return s.repo.GetBySymbol(symbol)
}

func (s *investmentService) UpdateInvestmentPrice(ctx context.Context, id uuid.UUID, newPrice float64) error {
	investment, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("investment not found: %w", err)
	}

	investment.CurrentPrice = newPrice
	investment.LastUpdated = time.Now()

	return s.repo.Update(investment)
}

func isValidInvestmentType(investmentType investments.InvestmentType) bool {
	switch investmentType {
	case investments.InvestmentTypeStock,
		investments.InvestmentTypeBond,
		investments.InvestmentTypeMutualFund,
		investments.InvestmentTypeETF,
		investments.InvestmentTypeCrypto:
		return true
	default:
		return false
	}
}
