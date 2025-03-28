package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sparkfund/aml-service/internal/model"
	"github.com/sparkfund/aml-service/internal/repository"
)

type AMLService interface {
	ProcessTransaction(ctx context.Context, userID uuid.UUID, req *model.TransactionRequest) (*model.TransactionResponse, error)
	AssessRisk(ctx context.Context, tx *model.Transaction) (*model.RiskAssessment, error)
	FlagTransaction(ctx context.Context, txID uuid.UUID, reason string, flaggedBy string) error
	ReviewTransaction(ctx context.Context, txID uuid.UUID, status model.TransactionStatus, notes string, reviewedBy string) error
	ListTransactions(ctx context.Context, filter *model.TransactionFilter) ([]*model.TransactionResponse, error)
	GetTransaction(ctx context.Context, txID uuid.UUID) (*model.TransactionResponse, error)
}

type amlService struct {
	txRepo repository.TransactionRepository
}

func NewAMLService(txRepo repository.TransactionRepository) AMLService {
	return &amlService{
		txRepo: txRepo,
	}
}

func (s *amlService) ProcessTransaction(ctx context.Context, userID uuid.UUID, req *model.TransactionRequest) (*model.TransactionResponse, error) {
	// Create transaction
	tx := &model.Transaction{
		UserID:             userID,
		Type:               req.Type,
		Amount:             req.Amount,
		Currency:           req.Currency,
		Status:             model.TransactionStatusPending,
		SourceAccount:      req.SourceAccount,
		DestinationAccount: req.DestinationAccount,
		Description:        req.Description,
		IPAddress:          req.IPAddress,
		DeviceID:           req.DeviceID,
		Location:           req.Location,
	}

	// Convert metadata to JSON
	if req.Metadata != nil {
		metadataJSON, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, err
		}
		tx.Metadata = string(metadataJSON)
	}

	// Assess risk
	assessment, err := s.AssessRisk(ctx, tx)
	if err != nil {
		return nil, err
	}

	tx.RiskLevel = assessment.RiskLevel

	// Save transaction
	if err := s.txRepo.Create(tx); err != nil {
		return nil, err
	}

	return s.toTransactionResponse(tx), nil
}

func (s *amlService) AssessRisk(ctx context.Context, tx *model.Transaction) (*model.RiskAssessment, error) {
	assessment := &model.RiskAssessment{
		TransactionID: tx.ID,
		RiskScore:     0,
		Factors:       make([]string, 0),
	}

	// Check transaction amount
	if tx.Amount > 10000 {
		assessment.RiskScore += 0.3
		assessment.Factors = append(assessment.Factors, "High transaction amount")
	}

	// Check for multiple transactions in short time
	recentTxs, err := s.txRepo.GetRecentTransactions(ctx, tx.UserID, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	if len(recentTxs) > 5 {
		assessment.RiskScore += 0.2
		assessment.Factors = append(assessment.Factors, "Multiple transactions in short time")
	}

	// Check for unusual locations
	if tx.Location != "US" {
		assessment.RiskScore += 0.1
		assessment.Factors = append(assessment.Factors, "Non-US location")
	}

	// Determine risk level
	switch {
	case assessment.RiskScore >= 0.6:
		assessment.RiskLevel = model.RiskLevelCritical
		assessment.Recommendation = "Manual review required"
	case assessment.RiskScore >= 0.4:
		assessment.RiskLevel = model.RiskLevelHigh
		assessment.Recommendation = "Enhanced due diligence required"
	case assessment.RiskScore >= 0.2:
		assessment.RiskLevel = model.RiskLevelMedium
		assessment.Recommendation = "Standard review"
	default:
		assessment.RiskLevel = model.RiskLevelLow
		assessment.Recommendation = "No special review required"
	}

	return assessment, nil
}

func (s *amlService) FlagTransaction(ctx context.Context, txID uuid.UUID, reason string, flaggedBy string) error {
	tx, err := s.txRepo.GetByID(txID)
	if err != nil {
		return err
	}

	tx.Status = model.TransactionStatusFlagged
	tx.FlaggedBy = &flaggedBy
	tx.FlagReason = &reason

	return s.txRepo.Update(tx)
}

func (s *amlService) ReviewTransaction(ctx context.Context, txID uuid.UUID, status model.TransactionStatus, notes string, reviewedBy string) error {
	tx, err := s.txRepo.GetByID(txID)
	if err != nil {
		return err
	}

	if tx.Status != model.TransactionStatusFlagged && tx.Status != model.TransactionStatusPending {
		return errors.New("transaction cannot be reviewed in current status")
	}

	tx.Status = status
	tx.ReviewedBy = &reviewedBy
	tx.ReviewNotes = &notes

	return s.txRepo.Update(tx)
}

func (s *amlService) ListTransactions(ctx context.Context, filter *model.TransactionFilter) ([]*model.TransactionResponse, error) {
	txs, err := s.txRepo.List(filter)
	if err != nil {
		return nil, err
	}

	responses := make([]*model.TransactionResponse, len(txs))
	for i, tx := range txs {
		responses[i] = s.toTransactionResponse(tx)
	}

	return responses, nil
}

func (s *amlService) GetTransaction(ctx context.Context, txID uuid.UUID) (*model.TransactionResponse, error) {
	tx, err := s.txRepo.GetByID(txID)
	if err != nil {
		return nil, err
	}

	return s.toTransactionResponse(tx), nil
}

func (s *amlService) toTransactionResponse(tx *model.Transaction) *model.TransactionResponse {
	response := &model.TransactionResponse{
		ID:                 tx.ID,
		UserID:             tx.UserID,
		Type:               tx.Type,
		Amount:             tx.Amount,
		Currency:           tx.Currency,
		Status:             tx.Status,
		RiskLevel:          tx.RiskLevel,
		SourceAccount:      tx.SourceAccount,
		DestinationAccount: tx.DestinationAccount,
		Description:        tx.Description,
		IPAddress:          tx.IPAddress,
		DeviceID:           tx.DeviceID,
		Location:           tx.Location,
		FlaggedBy:          tx.FlaggedBy,
		FlagReason:         tx.FlagReason,
		ReviewedBy:         tx.ReviewedBy,
		ReviewNotes:        tx.ReviewNotes,
		CreatedAt:          tx.CreatedAt,
		UpdatedAt:          tx.UpdatedAt,
	}

	// Parse metadata
	if tx.Metadata != "" {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(tx.Metadata), &metadata); err == nil {
			response.Metadata = metadata
		}
	}

	return response
}
