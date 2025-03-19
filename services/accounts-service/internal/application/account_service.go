package application

import (
	"context"
	"fmt"

	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/domain"
	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/ports"
	"github.com/google/uuid"
)

type AccountService struct {
	db ports.DBPort
}

func NewAccountService(db ports.DBPort) *AccountService {
	return &AccountService{
		db: db,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, userID uuid.UUID, name string, accountType domain.AccountType, currency string) (*domain.Account, error) {
	account := domain.NewAccount(userID, name, accountType, currency)
	if err := s.db.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}
	return account, nil
}

func (s *AccountService) GetAccount(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	account, err := s.db.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return account, nil
}

func (s *AccountService) GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error) {
	accounts, err := s.db.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user accounts: %w", err)
	}
	return accounts, nil
}

func (s *AccountService) UpdateAccount(ctx context.Context, id uuid.UUID, name string, accountType domain.AccountType) error {
	account, err := s.db.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}

	if err := account.Update(name, accountType); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	if err := s.db.Update(ctx, account); err != nil {
		return fmt.Errorf("failed to save account: %w", err)
	}
	return nil
}

func (s *AccountService) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	if err := s.db.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}
	return nil
}

func (s *AccountService) GetAccountByNumber(ctx context.Context, accountNumber string) (*domain.Account, error) {
	account, err := s.db.GetByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return account, nil
}

func (s *AccountService) UpdateBalance(ctx context.Context, id uuid.UUID, amount float64) error {
	account, err := s.db.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}

	account.Balance += amount
	if err := s.db.Update(ctx, account); err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}
	return nil
}
