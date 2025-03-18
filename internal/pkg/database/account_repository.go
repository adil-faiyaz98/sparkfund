package database

import (
	"context"

	"gorm.io/gorm"

	"github.com/adilm/money-pulse/internal/pkg/models"
)

// AccountRepository handles database operations for accounts
type AccountRepository struct {
	db *gorm.DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

// CreateAccount creates a new account in the database
func (r *AccountRepository) CreateAccount(ctx context.Context, account *models.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

// GetAccounts retrieves accounts for a user
func (r *AccountRepository) GetAccounts(ctx context.Context, userID uint) ([]models.Account, error) {
	var accounts []models.Account

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}

// GetAccountByID retrieves a single account by ID
func (r *AccountRepository) GetAccountByID(ctx context.Context, id uint) (*models.Account, error) {
	var account models.Account

	if err := r.db.WithContext(ctx).First(&account, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &account, nil
}

// UpdateAccount updates an existing account
func (r *AccountRepository) UpdateAccount(ctx context.Context, account *models.Account) error {
	return r.db.WithContext(ctx).Save(account).Error
}

// DeleteAccount deletes an account by ID
func (r *AccountRepository) DeleteAccount(ctx context.Context, id uint, userID uint) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.Account{}).Error
}

// UpdateAccountBalance updates the current balance of an account
func (r *AccountRepository) UpdateAccountBalance(ctx context.Context, accountID uint, amount float64) error {
	return r.db.WithContext(ctx).Model(&models.Account{}).
		Where("id = ?", accountID).
		UpdateColumn("current_balance", gorm.Expr("current_balance + ?", amount)).
		Error
}
