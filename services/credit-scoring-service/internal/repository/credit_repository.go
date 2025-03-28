package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/sparkfund/credit-scoring-service/internal/model"
	"gorm.io/gorm"
)

type CreditRepository interface {
	CreateCreditHistory(history *model.CreditHistory) error
	GetCreditHistory(id uuid.UUID) (*model.CreditHistory, error)
	GetCreditHistories(userID uuid.UUID) ([]*model.CreditHistory, error)
	UpdateCreditHistory(history *model.CreditHistory) error
	GetCreditScore(userID uuid.UUID) (*model.CreditScore, error)
	UpsertCreditScore(score *model.CreditScore) error
}

type creditRepository struct {
	db *gorm.DB
}

func NewCreditRepository(db *gorm.DB) CreditRepository {
	return &creditRepository{
		db: db,
	}
}

func (r *creditRepository) CreateCreditHistory(history *model.CreditHistory) error {
	return r.db.Create(history).Error
}

func (r *creditRepository) GetCreditHistory(id uuid.UUID) (*model.CreditHistory, error) {
	var history model.CreditHistory
	err := r.db.First(&history, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

func (r *creditRepository) GetCreditHistories(userID uuid.UUID) ([]*model.CreditHistory, error) {
	var histories []*model.CreditHistory
	err := r.db.Where("user_id = ?", userID).Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

func (r *creditRepository) UpdateCreditHistory(history *model.CreditHistory) error {
	history.UpdatedAt = time.Now()
	return r.db.Save(history).Error
}

func (r *creditRepository) GetCreditScore(userID uuid.UUID) (*model.CreditScore, error) {
	var score model.CreditScore
	err := r.db.Where("user_id = ?", userID).First(&score).Error
	if err != nil {
		return nil, err
	}
	return &score, nil
}

func (r *creditRepository) UpsertCreditScore(score *model.CreditScore) error {
	score.UpdatedAt = time.Now()
	return r.db.Save(score).Error
}
