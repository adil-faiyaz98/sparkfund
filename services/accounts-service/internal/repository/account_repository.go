package repository

import (
	"context"

	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/domain"
	"github.com/google/uuid"
)

type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error)
	GetByAccountNumber(ctx context.Context, accountNumber string) (*domain.Account, error)
	Update(ctx context.Context, account *domain.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
}
