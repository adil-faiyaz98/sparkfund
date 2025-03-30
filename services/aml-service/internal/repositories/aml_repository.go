package repositories

import (
	"context"
	"errors"

	"aml-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AMLRepository interface {
	CreateTransaction(ctx context.Context, tx *models.Transaction) error
	GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	UpdateTransaction(ctx context.Context, tx *models.Transaction) error
	ListTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Transaction, error)
	CreateRiskFactor(ctx context.Context, rf *models.RiskFactor) error
	CreateAlert(ctx context.Context, alert *models.Alert) error
	CreateScreeningResult(ctx context.Context, sr *models.ScreeningResult) error
	GetRiskProfile(ctx context.Context, userID uuid.UUID) (*models.RiskProfile, error)
	UpdateRiskProfile(ctx context.Context, rp *models.RiskProfile) error
}

type amlRepository struct {
	db *gorm.DB
}

func NewAMLRepository(db *gorm.DB) AMLRepository {
	return &amlRepository{db: db}
}

func (r *amlRepository) CreateTransaction(ctx context.Context, tx *models.Transaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *amlRepository) GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.db.WithContext(ctx).
		Preload("RiskFactors").
		Preload("Alerts").
		Preload("ScreeningResult").
		First(&tx, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tx, nil
}

func (r *amlRepository) UpdateTransaction(ctx context.Context, tx *models.Transaction) error {
	return r.db.WithContext(ctx).Save(tx).Error
}

func (r *amlRepository) ListTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *amlRepository) CreateRiskFactor(ctx context.Context, rf *models.RiskFactor) error {
	return r.db.WithContext(ctx).Create(rf).Error
}

func (r *amlRepository) CreateAlert(ctx context.Context, alert *models.Alert) error {
	return r.db.WithContext(ctx).Create(alert).Error
}

func (r *amlRepository) CreateScreeningResult(ctx context.Context, sr *models.ScreeningResult) error {
	return r.db.WithContext(ctx).Create(sr).Error
}

func (r *amlRepository) GetRiskProfile(ctx context.Context, userID uuid.UUID) (*models.RiskProfile, error) {
	var rp models.RiskProfile
	err := r.db.WithContext(ctx).First(&rp, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rp, nil
}

func (r *amlRepository) UpdateRiskProfile(ctx context.Context, rp *models.RiskProfile) error {
	return r.db.WithContext(ctx).Save(rp).Error
}
