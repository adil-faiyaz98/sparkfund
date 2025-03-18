package repository

import (
	"context"
	"fmt"
	"time"

	"your-project/internal/investments"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new PostgreSQL repository instance
func NewPostgresRepository(db *gorm.DB) investments.InvestmentRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(investment *investments.Investment) error {
	investment.ID = uuid.New()
	investment.CreatedAt = time.Now()
	investment.UpdatedAt = time.Now()
	investment.LastUpdated = time.Now()

	return r.db.Create(investment).Error
}

func (r *postgresRepository) GetByID(id uuid.UUID) (*investments.Investment, error) {
	var investment investments.Investment
	err := r.db.Where("id = ?", id).First(&investment).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get investment by ID: %w", err)
	}
	return &investment, nil
}

func (r *postgresRepository) GetByUserID(userID uuid.UUID) ([]*investments.Investment, error) {
	var investments []*investments.Investment
	err := r.db.Where("user_id = ?", userID).Find(&investments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get investments by user ID: %w", err)
	}
	return investments, nil
}

func (r *postgresRepository) GetByAccountID(accountID uuid.UUID) ([]*investments.Investment, error) {
	var investments []*investments.Investment
	err := r.db.Where("account_id = ?", accountID).Find(&investments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get investments by account ID: %w", err)
	}
	return investments, nil
}

func (r *postgresRepository) Update(investment *investments.Investment) error {
	investment.UpdatedAt = time.Now()
	return r.db.Save(investment).Error
}

func (r *postgresRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&investments.Investment{}).Error
}

func (r *postgresRepository) GetBySymbol(symbol string) ([]*investments.Investment, error) {
	var investments []*investments.Investment
	err := r.db.Where("symbol = ?", symbol).Find(&investments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get investments by symbol: %w", err)
	}
	return investments, nil
}

// Migrate performs database migrations
func (r *postgresRepository) Migrate(ctx context.Context) error {
	return r.db.AutoMigrate(&investments.Investment{}).Error
}
