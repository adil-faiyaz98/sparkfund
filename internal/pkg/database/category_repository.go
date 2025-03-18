package database

import (
	"context"

	"gorm.io/gorm"

	"github.com/adilm/money-pulse/internal/pkg/models"
)

// CategoryRepository handles database operations for categories
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

// CreateCategory creates a new category in the database
func (r *CategoryRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// GetCategories retrieves categories for a user
func (r *CategoryRepository) GetCategories(ctx context.Context, userID uint) ([]models.Category, error) {
	var categories []models.Category

	// Get system categories and user's custom categories
	if err := r.db.WithContext(ctx).Where("user_id = ? OR is_system = true", userID).Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// GetCategoryByID retrieves a single category by ID
func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id uint) (*models.Category, error) {
	var category models.Category

	if err := r.db.WithContext(ctx).First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

// UpdateCategory updates an existing category
func (r *CategoryRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	// Only allow updating non-system categories
	return r.db.WithContext(ctx).Where("is_system = false").Save(category).Error
}

// DeleteCategory deletes a category by ID if it's not a system category
func (r *CategoryRepository) DeleteCategory(ctx context.Context, id uint, userID uint) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ? AND is_system = false", id, userID).Delete(&models.Category{}).Error
}
