package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/domain"
)

type AccountEntity struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID        uuid.UUID `gorm:"type:uuid;not null"`
	Name          string    `gorm:"not null"`
	AccountNumber string    `gorm:"unique;not null"`
	Type          string    `gorm:"not null"`
	Balance       float64   `gorm:"not null;default:0"`
	Currency      string    `gorm:"not null"`
	CreatedAt     int64     `gorm:"not null"`
	UpdatedAt     int64     `gorm:"not null"`
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(dsn string) (*Adapter, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&AccountEntity{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &Adapter{db: db}, nil
}

func (a *Adapter) Create(ctx context.Context, account *domain.Account) error {
	entity := toEntity(account)
	if err := a.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}
	return nil
}

func (a *Adapter) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	var entity AccountEntity
	if err := a.db.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return toDomain(&entity), nil
}

func (a *Adapter) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error) {
	var entities []AccountEntity
	if err := a.db.WithContext(ctx).Where("user_id = ?", userID).Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("failed to get user accounts: %w", err)
	}

	accounts := make([]*domain.Account, len(entities))
	for i, entity := range entities {
		accounts[i] = toDomain(&entity)
	}
	return accounts, nil
}

func (a *Adapter) Update(ctx context.Context, account *domain.Account) error {
	entity := toEntity(account)
	if err := a.db.WithContext(ctx).Save(entity).Error; err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}
	return nil
}

func (a *Adapter) Delete(ctx context.Context, id uuid.UUID) error {
	if err := a.db.WithContext(ctx).Delete(&AccountEntity{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}
	return nil
}

func (a *Adapter) GetByAccountNumber(ctx context.Context, accountNumber string) (*domain.Account, error) {
	var entity AccountEntity
	if err := a.db.WithContext(ctx).First(&entity, "account_number = ?", accountNumber).Error; err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return toDomain(&entity), nil
}

func toEntity(account *domain.Account) *AccountEntity {
	return &AccountEntity{
		ID:            account.ID,
		UserID:        account.UserID,
		Name:          account.Name,
		AccountNumber: account.AccountNumber,
		Type:          string(account.Type),
		Balance:       account.Balance,
		Currency:      account.Currency,
		CreatedAt:     account.CreatedAt.Unix(),
		UpdatedAt:     account.UpdatedAt.Unix(),
	}
}

func toDomain(entity *AccountEntity) *domain.Account {
	return &domain.Account{
		ID:            entity.ID,
		UserID:        entity.UserID,
		Name:          entity.Name,
		AccountNumber: entity.AccountNumber,
		Type:          domain.AccountType(entity.Type),
		Balance:       entity.Balance,
		Currency:      entity.Currency,
		CreatedAt:     time.Unix(entity.CreatedAt, 0),
		UpdatedAt:     time.Unix(entity.UpdatedAt, 0),
	}
}
