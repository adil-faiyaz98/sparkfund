package repository

import (
	"context"

	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/domain/model"
)

// AccountRepository defines the interface for account data operations
type AccountRepository interface {
	Create(ctx context.Context, account *model.Account) (*model.Account, error)
	GetByID(ctx context.Context, id string) (*model.Account, error)
	Update(ctx context.Context, account *model.Account) (*model.Account, error)
	Delete(ctx context.Context, id string) error
	ListByUserID(ctx context.Context, userID string, page, pageSize int) ([]*model.Account, int, error)
	Close() error
}
