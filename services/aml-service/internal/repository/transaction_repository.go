package repository

import (
	"context"
	"time"

	"aml-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(tx *model.Transaction) error
	GetByID(id uuid.UUID) (*model.Transaction, error)
	GetRecentTransactions(ctx context.Context, userID uuid.UUID, duration time.Duration) ([]*model.Transaction, error)
	List(filter *model.TransactionFilter) ([]*model.Transaction, error)
	Update(tx *model.Transaction) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) Create(tx *model.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepository) GetByID(id uuid.UUID) (*model.Transaction, error) {
	var tx model.Transaction
	err := r.db.Where("id = ?", id).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) GetRecentTransactions(ctx context.Context, userID uuid.UUID, duration time.Duration) ([]*model.Transaction, error) {
	var txs []*model.Transaction
	cutoff := time.Now().Add(-duration)

	err := r.db.Where("user_id = ? AND created_at >= ?", userID, cutoff).Find(&txs).Error
	if err != nil {
		return nil, err
	}

	return txs, nil
}

func (r *transactionRepository) List(filter *model.TransactionFilter) ([]*model.Transaction, error) {
	var txs []*model.Transaction
	query := r.db.Model(&model.Transaction{})

	if filter != nil {
		if filter.UserID != nil {
			query = query.Where("user_id = ?", filter.UserID)
		}
		if filter.Type != nil {
			query = query.Where("type = ?", filter.Type)
		}
		if filter.Status != nil {
			query = query.Where("status = ?", filter.Status)
		}
		if filter.RiskLevel != nil {
			query = query.Where("risk_level = ?", filter.RiskLevel)
		}
		if filter.StartDate != nil {
			query = query.Where("created_at >= ?", filter.StartDate)
		}
		if filter.EndDate != nil {
			query = query.Where("created_at <= ?", filter.EndDate)
		}
		if filter.MinAmount != nil {
			query = query.Where("amount >= ?", filter.MinAmount)
		}
		if filter.MaxAmount != nil {
			query = query.Where("amount <= ?", filter.MaxAmount)
		}
		if filter.Currency != nil {
			query = query.Where("currency = ?", filter.Currency)
		}
		if filter.FlaggedOnly != nil && *filter.FlaggedOnly {
			query = query.Where("status = ?", model.TransactionStatusFlagged)
		}
	}

	err := query.Find(&txs).Error
	if err != nil {
		return nil, err
	}

	return txs, nil
}

func (r *transactionRepository) Update(tx *model.Transaction) error {
	return r.db.Save(tx).Error
}
