package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/domain/model"
	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/domain/repository"
)

// Common service errors
var (
	ErrInvalidAccountType = errors.New("invalid account type")
	ErrInvalidCurrency    = errors.New("invalid currency")
	ErrInsufficientFunds  = errors.New("insufficient funds")
	ErrAccountNotActive   = errors.New("account is not active")
)

// AccountService provides account-related operations
type AccountService struct {
	repo repository.AccountRepository
}

// NewAccountService creates a new account service
func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

// CreateAccount creates a new account
func (s *AccountService) CreateAccount(ctx context.Context, userID, accountType, currency string, initialDeposit float64) (*model.Account, error) {
	// Validate account type
	if !isValidAccountType(accountType) {
		return nil, ErrInvalidAccountType
	}

	// Validate currency
	if !isValidCurrency(currency) {
		return nil, ErrInvalidCurrency
	}

	// Create account with initial values
	account := &model.Account{
		UserID:      userID,
		AccountType: accountType,
		Currency:    currency,
		Balance:     initialDeposit,
		IsActive:    true,
	}

	// Save to repository
	return s.repo.Create(ctx, account)
}

// GetAccount retrieves an account by ID
func (s *AccountService) GetAccount(ctx context.Context, id string) (*model.Account, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateAccountStatus updates the status of an account
func (s *AccountService) UpdateAccountStatus(ctx context.Context, id string, isActive bool) (*model.Account, error) {
	account, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	account.IsActive = isActive
	return s.repo.Update(ctx, account)
}

// DeleteAccount removes an account
func (s *AccountService) DeleteAccount(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// ListUserAccounts retrieves all accounts for a specific user
func (s *AccountService) ListUserAccounts(ctx context.Context, userID string, page, pageSize int) ([]*model.Account, int, error) {
	return s.repo.ListByUserID(ctx, userID, page, pageSize)
}

// DepositFunds adds funds to an account
func (s *AccountService) DepositFunds(ctx context.Context, accountID string, amount float64) (*model.Account, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("deposit amount must be positive")
	}

	account, err := s.repo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if !account.IsActive {
		return nil, ErrAccountNotActive
	}

	account.Balance += amount
	return s.repo.Update(ctx, account)
}

// WithdrawFunds removes funds from an account
func (s *AccountService) WithdrawFunds(ctx context.Context, accountID string, amount float64) (*model.Account, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("withdrawal amount must be positive")
	}

	account, err := s.repo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if !account.IsActive {
		return nil, ErrAccountNotActive
	}

	if account.Balance < amount {
		return nil, ErrInsufficientFunds
	}

	account.Balance -= amount
	return s.repo.Update(ctx, account)
}

// TransferFunds transfers money between accounts
func (s *AccountService) TransferFunds(ctx context.Context, sourceID, destinationID string, amount float64) (*model.Account, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("transfer amount must be positive")
	}

	sourceAccount, err := s.repo.GetByID(ctx, sourceID)
	if err != nil {
		return nil, fmt.Errorf("source account: %w", err)
	}

	if !sourceAccount.IsActive {
		return nil, ErrAccountNotActive
	}

	if sourceAccount.Balance < amount {
		return nil, ErrInsufficientFunds
	}

	destAccount, err := s.repo.GetByID(ctx, destinationID)
	if err != nil {
		return nil, fmt.Errorf("destination account: %w", err)
	}

	if !destAccount.IsActive {
		return nil, ErrAccountNotActive
	}

	// Ensure same currency for simplicity
	if sourceAccount.Currency != destAccount.Currency {
		return nil, fmt.Errorf("cannot transfer between accounts with different currencies")
	}

	// Update source account
	sourceAccount.Balance -= amount
	updatedSource, err := s.repo.Update(ctx, sourceAccount)
	if err != nil {
		return nil, err
	}

	// Update destination account
	destAccount.Balance += amount
	_, err = s.repo.Update(ctx, destAccount)
	if err != nil {
		// In a real-world scenario, we would need to implement a transaction to rollback
		// the source account update if the destination update fails
		return nil, err
	}

	return updatedSource, nil
}

// Helper functions
func isValidAccountType(accountType string) bool {
	validTypes := []string{
		model.AccountTypes.Checking,
		model.AccountTypes.Savings,
		model.AccountTypes.Credit,
		model.AccountTypes.Investment,
	}

	for _, t := range validTypes {
		if accountType == t {
			return true
		}
	}
	return false
}

func isValidCurrency(currency string) bool {
	validCurrencies := []string{
		model.Currencies.USD,
		model.Currencies.EUR,
		model.Currencies.GBP,
	}

	for _, c := range validCurrencies {
		if currency == c {
			return true
		}
	}
	return false
}
