package service

import (
	"context"
	"fmt"

	"github.com/adil-faiyaz98/structgen/internal/accounts"

	"github.com/google/uuid"
)

type accountService struct {
	repo accounts.AccountRepository
}

func NewAccountService(repo accounts.AccountRepository) accounts.AccountService {
	return &accountService{repo: repo}
}

func (s *accountService) CreateAccount(ctx context.Context, account *accounts.Account) error {
	// Validate account type
	if !isValidAccountType(account.Type) {
		return fmt.Errorf("invalid account type: %s", account.Type)
	}

	// Generate account number if not provided
	if account.AccountNumber == "" {
		account.AccountNumber = generateAccountNumber()
	}

	return s.repo.Create(account)
}

func (s *accountService) GetAccount(ctx context.Context, id uuid.UUID) (*accounts.Account, error) {
	return s.repo.GetByID(id)
}

func (s *accountService) GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]*accounts.Account, error) {
	return s.repo.GetByUserID(userID)
}

func (s *accountService) UpdateAccount(ctx context.Context, account *accounts.Account) error {
	// Validate account exists
	existing, err := s.repo.GetByID(account.ID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Validate account type if changed
	if existing.Type != account.Type && !isValidAccountType(account.Type) {
		return fmt.Errorf("invalid account type: %s", account.Type)
	}

	return s.repo.Update(account)
}

func (s *accountService) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *accountService) GetAccountByNumber(ctx context.Context, accountNumber string) (*accounts.Account, error) {
	return s.repo.GetByAccountNumber(accountNumber)
}

func isValidAccountType(accountType accounts.AccountType) bool {
	switch accountType {
	case accounts.AccountTypeSavings,
		accounts.AccountTypeChecking,
		accounts.AccountTypeInvestment,
		accounts.AccountTypeCredit:
		return true
	default:
		return false
	}
}

func generateAccountNumber() string {
	// In a real application, this would be more sophisticated
	// and might include a prefix based on account type
	return fmt.Sprintf("ACC%s", uuid.New().String()[:8])
}
