package repository

import (
	"context"
	"fmt"
	"time"

	"your-project/internal/accounts"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new PostgreSQL repository instance
func NewPostgresRepository(db *gorm.DB) accounts.AccountRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(account *accounts.Account) error {
	account.ID = uuid.New()
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	return r.db.Create(account).Error
}

func (r *postgresRepository) GetByID(id uuid.UUID) (*accounts.Account, error) {
	var account accounts.Account
	err := r.db.Where("id = ?", id).First(&account).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get account by ID: %w", err)
	}
	return &account, nil
}

func (r *postgresRepository) GetByUserID(userID uuid.UUID) ([]*accounts.Account, error) {
	var accounts []*accounts.Account
	err := r.db.Where("user_id = ?", userID).Find(&accounts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts by user ID: %w", err)
	}
	return accounts, nil
}

func (r *postgresRepository) Update(account *accounts.Account) error {
	account.UpdatedAt = time.Now()
	return r.db.Save(account).Error
}

func (r *postgresRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&accounts.Account{}).Error
}

func (r *postgresRepository) GetByAccountNumber(accountNumber string) (*accounts.Account, error) {
	var account accounts.Account
	err := r.db.Where("account_number = ?", accountNumber).First(&account).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get account by account number: %w", err)
	}
	return &account, nil
}

// Migrate performs database migrations
func (r *postgresRepository) Migrate(ctx context.Context) error {
	return r.db.AutoMigrate(&accounts.Account{}).Error
}
