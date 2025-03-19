package service

import (
	"context"
	"fmt"

	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/domain"
	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/repository"
	"github.com/google/uuid"
)

type AccountService interface {
	CreateAccount(ctx context.Context, account *domain.Account) error
	GetAccount(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error)
	UpdateAccount(ctx context.Context, account *domain.Account) error
	DeleteAccount(ctx context.Context, id uuid.UUID) error
	GetAccountByNumber(ctx context.Context, accountNumber string) (*domain.Account, error)
}

type accountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

func (s *accountService) CreateAccount(ctx context.Context, account *domain.Account) error {
	if account.AccountNumber == "" {
		account.AccountNumber = generateAccountNumber(account.Type)
	}
	return s.repo.Create(ctx, account)
}

func (s *accountService) GetAccount(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *accountService) GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *accountService) UpdateAccount(ctx context.Context, account *domain.Account) error {
	existing, err := s.repo.GetByID(ctx, account.ID)
	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	if existing.Type != account.Type && !domain.IsValidAccountType(account.Type) {
		return fmt.Errorf("invalid account type: %s", account.Type)
	}

	return s.repo.Update(ctx, account)
}

func (s *accountService) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *accountService) GetAccountByNumber(ctx context.Context, accountNumber string) (*domain.Account, error) {
	return s.repo.GetByAccountNumber(ctx, accountNumber)
}

func generateAccountNumber(accountType domain.AccountType) string {
	prefix := "ACC"
	switch accountType {
	case domain.AccountTypeSavings:
		prefix = "SAV"
	case domain.AccountTypeChecking:
		prefix = "CHK"
	case domain.AccountTypeInvestment:
		prefix = "INV"
	case domain.AccountTypeCredit:
		prefix = "CRD"
	}
	return fmt.Sprintf("%s%s", prefix, uuid.New().String()[:8])
}
