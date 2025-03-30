package repositories

import (
	"context"
	"errors"
	"time"

	"investment-service/internal/database"
	"investment-service/internal/metrics"
	"investment-service/internal/models"

	"gorm.io/gorm"
)

// InvestmentRepository handles database operations for investments
type InvestmentRepository interface {
	Create(ctx context.Context, investment *models.Investment) error
	GetByID(ctx context.Context, id uint) (*models.Investment, error)
	GetByUserID(ctx context.Context, userID uint) ([]models.Investment, error)
	GetByPortfolioID(ctx context.Context, portfolioID uint) ([]models.Investment, error)
	Update(ctx context.Context, investment *models.Investment) error
	Delete(ctx context.Context, id uint) error
	GetAll(ctx context.Context, page, pageSize int) ([]models.Investment, int64, error)
}

// GormInvestmentRepository implements InvestmentRepository using GORM
type GormInvestmentRepository struct {
	db *gorm.DB
}

// NewInvestmentRepository creates a new investment repository
func NewInvestmentRepository() InvestmentRepository {
	return &GormInvestmentRepository{
		db: database.DB,
	}
}

// Create inserts a new investment
func (r *GormInvestmentRepository) Create(ctx context.Context, investment *models.Investment) error {
	defer metrics.TrackDBQuery("investment_create")()

	// Set timestamps
	investment.CreatedAt = time.Now()
	investment.UpdatedAt = time.Now()

	return r.db.WithContext(ctx).Create(investment).Error
}

// GetByID retrieves an investment by ID
func (r *GormInvestmentRepository) GetByID(ctx context.Context, id uint) (*models.Investment, error) {
	defer metrics.TrackDBQuery("investment_get_by_id")()

	var investment models.Investment
	if err := r.db.WithContext(ctx).First(&investment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &investment, nil
}

// GetByUserID retrieves all investments for a user
func (r *GormInvestmentRepository) GetByUserID(ctx context.Context, userID uint) ([]models.Investment, error) {
	defer metrics.TrackDBQuery("investment_get_by_user_id")()

	var investments []models.Investment
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&investments).Error; err != nil {
		return nil, err
	}
	return investments, nil
}

// GetByPortfolioID retrieves all investments in a portfolio
func (r *GormInvestmentRepository) GetByPortfolioID(ctx context.Context, portfolioID uint) ([]models.Investment, error) {
	defer metrics.TrackDBQuery("investment_get_by_portfolio_id")()

	var investments []models.Investment
	if err := r.db.WithContext(ctx).Where("portfolio_id = ?", portfolioID).Find(&investments).Error; err != nil {
		return nil, err
	}
	return investments, nil
}

// Update updates an investment
func (r *GormInvestmentRepository) Update(ctx context.Context, investment *models.Investment) error {
	defer metrics.TrackDBQuery("investment_update")()

	// Set updated timestamp
	investment.UpdatedAt = time.Now()

	return r.db.WithContext(ctx).Save(investment).Error
}

// Delete removes an investment
func (r *GormInvestmentRepository) Delete(ctx context.Context, id uint) error {
	defer metrics.TrackDBQuery("investment_delete")()

	return r.db.WithContext(ctx).Delete(&models.Investment{}, id).Error
}

// GetAll retrieves all investments with pagination
func (r *GormInvestmentRepository) GetAll(ctx context.Context, page, pageSize int) ([]models.Investment, int64, error) {
	defer metrics.TrackDBQuery("investment_get_all")()

	var investments []models.Investment
	var total int64

	// Calculate offset
	offset := (page - 1) * pageSize

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Investment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&investments).Error; err != nil {
		return nil, 0, err
	}

	return investments, total, nil
}
