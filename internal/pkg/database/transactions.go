package database

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/adilm/money-pulse/internal/pkg/models"
)

// TransactionRepository handles database operations for transactions
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// GetTransactions retrieves transactions for a user
func (r *TransactionRepository) GetTransactions(ctx context.Context, userID uint, filter *TransactionFilter) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("date DESC")

	if filter != nil {
		if filter.CategoryID != 0 {
			query = query.Where("category_id = ?", filter.CategoryID)
		}

		if filter.AccountID != 0 {
			query = query.Where("account_id = ?", filter.AccountID)
		}

		if !filter.StartDate.IsZero() {
			query = query.Where("date >= ?", filter.StartDate)
		}

		if !filter.EndDate.IsZero() {
			query = query.Where("date <= ?", filter.EndDate)
		}

		if filter.MinAmount != 0 {
			query = query.Where("amount >= ?", filter.MinAmount)
		}

		if filter.MaxAmount != 0 {
			query = query.Where("amount <= ?", filter.MaxAmount)
		}

		if filter.Type != "" {
			query = query.Where("type = ?", filter.Type)
		}
	}

	if err := query.Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

// TransactionFilter defines filtering options for transactions
type TransactionFilter struct {
	CategoryID uint
	AccountID  uint
	StartDate  time.Time
	EndDate    time.Time
	MinAmount  float64
	MaxAmount  float64
	Type       string // "income" or "expense"
}

// GetTransactionByID retrieves a transaction by ID
func (r *TransactionRepository) GetTransactionByID(ctx context.Context, id uint, userID uint) (*models.Transaction, error) {
	var transaction models.Transaction
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&transaction)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &transaction, nil
}

// CreateTransaction creates a new transaction
func (r *TransactionRepository) CreateTransaction(ctx context.Context, transaction *models.Transaction) error {
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()
	result := r.db.WithContext(ctx).Create(transaction)
	return result.Error
}

// UpdateTransaction updates an existing transaction
func (r *TransactionRepository) UpdateTransaction(ctx context.Context, transaction *models.Transaction) error {
	transaction.UpdatedAt = time.Now()
	result := r.db.WithContext(ctx).Save(transaction)
	return result.Error
}

// DeleteTransaction deletes a transaction
func (r *TransactionRepository) DeleteTransaction(ctx context.Context, id uint, userID uint) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.Transaction{})
	return result.Error
}
