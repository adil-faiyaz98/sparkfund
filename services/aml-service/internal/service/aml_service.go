package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"aml-service/internal/ai"
	"aml-service/internal/model"
	"aml-service/internal/repository"

	"github.com/google/uuid"
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
	txRepo     repository.TransactionRepository
	fraudModel ai.FraudModel
}

func NewAMLService(txRepo repository.TransactionRepository, fraudModel ai.FraudModel) AMLService {
	return &amlService{
		txRepo:     txRepo,
		fraudModel: fraudModel,
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

	// Assess risk using AI model
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
	// Get recent transactions for the user
	recentTxs, err := s.txRepo.GetRecentTransactions(ctx, tx.UserID, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	// Calculate transaction statistics
	var totalAmount float64
	for _, t := range recentTxs {
		totalAmount += t.Amount
	}
	avgAmount := totalAmount / float64(len(recentTxs)+1)
	amountDeviation := calculateAmountDeviation(tx.Amount, avgAmount)

	// Prepare features for fraud detection
	features := ai.FraudFeatures{
		Amount:           tx.Amount,
		Currency:         tx.Currency,
		TransactionTime:  time.Now(),
		TransactionType:  string(tx.Type),
		TransactionCount: len(recentTxs) + 1,
		AverageAmount:    avgAmount,
		AmountDeviation:  amountDeviation,
		UserAge:          25,    // TODO: Get from user service
		AccountAge:       30,    // TODO: Get from user service
		PreviousFraud:    false, // TODO: Get from fraud history
		RiskProfile:      0.5,   // TODO: Get from risk profile
		CountryRisk:      0.3,   // TODO: Get from country risk service
		IPRisk:           0.2,   // TODO: Get from IP risk service
		LocationMismatch: false, // TODO: Check against user's usual location
		LoginAttempts:    1,     // TODO: Get from auth service
		FailedLogins:     0,     // TODO: Get from auth service
		DeviceChanges:    0,     // TODO: Get from device tracking service
		PEPStatus:        false, // TODO: Get from PEP service
		SanctionList:     false, // TODO: Get from sanctions service
		WatchList:        false, // TODO: Get from watchlist service
	}

	// Get fraud prediction
	prediction, err := s.fraudModel.Predict(features)
	if err != nil {
		return nil, fmt.Errorf("failed to get fraud prediction: %w", err)
	}

	// Convert AI risk level to model risk level
	var riskLevel model.RiskLevel
	switch prediction.RiskLevel {
	case ai.FraudRiskLow:
		riskLevel = model.RiskLevelLow
	case ai.FraudRiskMedium:
		riskLevel = model.RiskLevelMedium
	case ai.FraudRiskHigh:
		riskLevel = model.RiskLevelHigh
	default:
		riskLevel = model.RiskLevelLow
	}

	// Create risk assessment
	assessment := &model.RiskAssessment{
		TransactionID:  tx.ID,
		RiskScore:      float64(prediction.RiskScore),
		RiskLevel:      riskLevel,
		Factors:        prediction.RiskFactors,
		Recommendation: prediction.Explanation,
	}

	return assessment, nil
}

func calculateAmountDeviation(amount, avgAmount float64) float64 {
	if avgAmount == 0 {
		return 0
	}
	return (amount - avgAmount) / avgAmount
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
