package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sparkfund/kyc-service/internal/model"
	"github.com/sparkfund/kyc-service/internal/repository"
)

type KYCService struct {
	repo *repository.KYCRepository
}

func NewKYCService(repo *repository.KYCRepository) *KYCService {
	return &KYCService{repo: repo}
}

func (s *KYCService) SubmitKYC(userID uuid.UUID, req *model.KYCRequest) (*model.KYCResponse, error) {
	// Check if user already has a KYC submission
	existingKYC, err := s.repo.GetByUserID(userID)
	if err == nil && existingKYC != nil {
		return nil, errors.New("user already has a KYC submission")
	}

	kyc := &model.KYC{
		UserID:         userID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		DateOfBirth:    req.DateOfBirth,
		Address:        req.Address,
		City:           req.City,
		Country:        req.Country,
		PostalCode:     req.PostalCode,
		DocumentType:   req.DocumentType,
		DocumentNumber: req.DocumentNumber,
		DocumentFront:  req.DocumentFront,
		DocumentBack:   req.DocumentBack,
		SelfieImage:    req.SelfieImage,
		Status:         model.KYCStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.Create(kyc); err != nil {
		return nil, err
	}

	return &model.KYCResponse{
		ID:         kyc.ID,
		UserID:     kyc.UserID,
		Status:     kyc.Status,
		VerifiedAt: kyc.VerifiedAt,
		CreatedAt:  kyc.CreatedAt,
		UpdatedAt:  kyc.UpdatedAt,
	}, nil
}

func (s *KYCService) GetKYCStatus(id uuid.UUID) (*model.KYCResponse, error) {
	kyc, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &model.KYCResponse{
		ID:              kyc.ID,
		UserID:          kyc.UserID,
		Status:          kyc.Status,
		RejectionReason: kyc.RejectionReason,
		VerifiedAt:      kyc.VerifiedAt,
		CreatedAt:       kyc.CreatedAt,
		UpdatedAt:       kyc.UpdatedAt,
	}, nil
}

func (s *KYCService) VerifyKYC(id uuid.UUID, verifiedBy uuid.UUID) error {
	kyc, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if kyc.Status != model.KYCStatusPending {
		return errors.New("KYC is not in pending status")
	}

	return s.repo.Verify(id, verifiedBy)
}

func (s *KYCService) RejectKYC(id uuid.UUID, reason string) error {
	kyc, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if kyc.Status != model.KYCStatusPending {
		return errors.New("KYC is not in pending status")
	}

	return s.repo.UpdateStatus(id, model.KYCStatusRejected, reason)
}

func (s *KYCService) ListPendingKYC() ([]model.KYC, error) {
	return s.repo.ListPending()
}
