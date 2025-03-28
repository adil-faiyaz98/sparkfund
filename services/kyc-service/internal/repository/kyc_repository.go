package repository

import (
	"github.com/google/uuid"
	"github.com/sparkfund/kyc-service/internal/model"
	"gorm.io/gorm"
)

type KYCRepository struct {
	db *gorm.DB
}

func NewKYCRepository(db *gorm.DB) *KYCRepository {
	return &KYCRepository{db: db}
}

func (r *KYCRepository) Create(kyc *model.KYC) error {
	return r.db.Create(kyc).Error
}

func (r *KYCRepository) GetByID(id uuid.UUID) (*model.KYC, error) {
	var kyc model.KYC
	err := r.db.First(&kyc, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &kyc, nil
}

func (r *KYCRepository) GetByUserID(userID uuid.UUID) (*model.KYC, error) {
	var kyc model.KYC
	err := r.db.Where("user_id = ?", userID).First(&kyc).Error
	if err != nil {
		return nil, err
	}
	return &kyc, nil
}

func (r *KYCRepository) Update(kyc *model.KYC) error {
	return r.db.Save(kyc).Error
}

func (r *KYCRepository) UpdateStatus(id uuid.UUID, status model.KYCStatus, rejectionReason string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if rejectionReason != "" {
		updates["rejection_reason"] = rejectionReason
	}
	return r.db.Model(&model.KYC{}).Where("id = ?", id).Updates(updates).Error
}

func (r *KYCRepository) Verify(id uuid.UUID, verifiedBy uuid.UUID) error {
	return r.db.Model(&model.KYC{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      model.KYCStatusVerified,
		"verified_by": verifiedBy,
		"verified_at": gorm.Now(),
	}).Error
}

func (r *KYCRepository) ListPending() ([]model.KYC, error) {
	var kycs []model.KYC
	err := r.db.Where("status = ?", model.KYCStatusPending).Find(&kycs).Error
	return kycs, err
}
