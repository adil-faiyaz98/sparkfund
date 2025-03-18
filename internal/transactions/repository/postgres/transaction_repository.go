package postgres

import (
	"fmt"

	"github.com/adil-faiyaz98/structgen/internal/transactions"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) transactions.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(transaction *transactions.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) GetByID(id uuid.UUID) (*transactions.Transaction, error) {
	var transaction transactions.Transaction
	if err := r.db.First(&transaction, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}
	return &transaction, nil
}

func (r *transactionRepository) GetByUserID(userID uuid.UUID) ([]*transactions.Transaction, error) {
	var transactionList []*transactions.Transaction
	if err := r.db.Where("user_id = ?", userID).Find(&transactionList).Error; err != nil {
		return nil, fmt.Errorf("failed to get user transactions: %w", err)
	}
	return transactionList, nil
}

func (r *transactionRepository) GetByAccountID(accountID uuid.UUID) ([]*transactions.Transaction, error) {
	var transactionList []*transactions.Transaction
	if err := r.db.Where("account_id = ?", accountID).Find(&transactionList).Error; err != nil {
		return nil, fmt.Errorf("failed to get account transactions: %w", err)
	}
	return transactionList, nil
}

func (r *transactionRepository) Update(transaction *transactions.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *transactionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&transactions.Transaction{}, "id = ?", id).Error
}
