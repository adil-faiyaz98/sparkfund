package ports

import (
	"context"

	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/domain"
	"github.com/google/uuid"
)

// APIPort defines the interface for the application core
type APIPort interface {
	CreateAccount(ctx context.Context, userID uuid.UUID, name string, accountType domain.AccountType, currency string) (*domain.Account, error)
	GetAccount(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error)
	UpdateAccount(ctx context.Context, id uuid.UUID, name string, accountType domain.AccountType) error
	DeleteAccount(ctx context.Context, id uuid.UUID) error
	GetAccountByNumber(ctx context.Context, accountNumber string) (*domain.Account, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, amount float64) error
}

// DBPort defines the interface for database operations
type DBPort interface {
	Create(ctx context.Context, account *domain.Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error)
	Update(ctx context.Context, account *domain.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByAccountNumber(ctx context.Context, accountNumber string) (*domain.Account, error)
}

// GRPCPort defines the interface for gRPC operations
type GRPCPort interface {
	Run() error
}
