package services

import (
	"context"
	"fmt"
	"time"

	"aml-service/internal/models"
	"aml-service/internal/repositories"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AMLService interface {
	ProcessTransaction(ctx context.Context, tx *models.Transaction) error
	GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	ListTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Transaction, error)
	GetRiskProfile(ctx context.Context, userID uuid.UUID) (*models.RiskProfile, error)
}

type amlService struct {
	repo   repositories.AMLRepository
	logger *zap.Logger
}

func NewAMLService(repo repositories.AMLRepository, logger *zap.Logger) AMLService {
	return &amlService{
		repo:   repo,
		logger: logger,
	}
}

func (s *amlService) ProcessTransaction(ctx context.Context, tx *models.Transaction) error {
	// Calculate risk score and level
	riskScore, riskLevel := s.calculateRisk(tx)
	tx.RiskScore = riskScore
	tx.RiskLevel = riskLevel

	// Create transaction record
	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		s.logger.Error("Failed to create transaction", zap.Error(err))
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// Perform screening
	screeningResult := s.performScreening(tx)
	screeningResult.TransactionID = tx.ID
	if err := s.repo.CreateScreeningResult(ctx, screeningResult); err != nil {
		s.logger.Error("Failed to create screening result", zap.Error(err))
		return fmt.Errorf("failed to create screening result: %w", err)
	}

	// Create risk factors
	riskFactors := s.identifyRiskFactors(tx, screeningResult)
	for _, rf := range riskFactors {
		if err := s.repo.CreateRiskFactor(ctx, rf); err != nil {
			s.logger.Error("Failed to create risk factor", zap.Error(err))
			return fmt.Errorf("failed to create risk factor: %w", err)
		}
	}

	// Generate alerts if necessary
	alerts := s.generateAlerts(tx, screeningResult, riskFactors)
	for _, alert := range alerts {
		if err := s.repo.CreateAlert(ctx, alert); err != nil {
			s.logger.Error("Failed to create alert", zap.Error(err))
			return fmt.Errorf("failed to create alert: %w", err)
		}
	}

	// Update user's risk profile
	if err := s.updateRiskProfile(ctx, tx); err != nil {
		s.logger.Error("Failed to update risk profile", zap.Error(err))
		return fmt.Errorf("failed to update risk profile: %w", err)
	}

	return nil
}

func (s *amlService) GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	tx, err := s.repo.GetTransaction(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	return tx, nil
}

func (s *amlService) ListTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Transaction, error) {
	transactions, err := s.repo.ListTransactions(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list transactions", zap.Error(err))
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	return transactions, nil
}

func (s *amlService) GetRiskProfile(ctx context.Context, userID uuid.UUID) (*models.RiskProfile, error) {
	rp, err := s.repo.GetRiskProfile(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get risk profile", zap.Error(err))
		return nil, fmt.Errorf("failed to get risk profile: %w", err)
	}
	return rp, nil
}

func (s *amlService) calculateRisk(tx *models.Transaction) (float64, string) {
	// Implement risk calculation logic
	var riskScore float64

	// Example risk factors:
	// 1. Transaction amount
	if tx.Amount > 10000 {
		riskScore += 0.3
	} else if tx.Amount > 5000 {
		riskScore += 0.2
	}

	// 2. Transaction type
	if tx.Type == "international" {
		riskScore += 0.2
	}

	// Determine risk level based on score
	var riskLevel string
	switch {
	case riskScore >= 0.7:
		riskLevel = "high"
	case riskScore >= 0.4:
		riskLevel = "medium"
	default:
		riskLevel = "low"
	}

	return riskScore, riskLevel
}

func (s *amlService) performScreening(tx *models.Transaction) *models.ScreeningResult {
	// Implement screening logic
	return &models.ScreeningResult{
		ID:           uuid.New(),
		SanctionList: false,
		PEPList:      false,
		WatchList:    false,
		Details:      "No matches found",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (s *amlService) identifyRiskFactors(tx *models.Transaction, sr *models.ScreeningResult) []*models.RiskFactor {
	var riskFactors []*models.RiskFactor

	// Example risk factor identification
	if tx.Amount > 10000 {
		riskFactors = append(riskFactors, &models.RiskFactor{
			ID:            uuid.New(),
			TransactionID: tx.ID,
			Type:          "large_amount",
			Description:   "Transaction amount exceeds high-risk threshold",
			Weight:        0.3,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
	}

	if sr.SanctionList {
		riskFactors = append(riskFactors, &models.RiskFactor{
			ID:            uuid.New(),
			TransactionID: tx.ID,
			Type:          "sanctions_match",
			Description:   "Match found in sanctions list",
			Weight:        0.5,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
	}

	return riskFactors
}

func (s *amlService) generateAlerts(tx *models.Transaction, sr *models.ScreeningResult, rfs []*models.RiskFactor) []*models.Alert {
	var alerts []*models.Alert

	// Generate alerts based on risk factors and screening results
	if tx.RiskLevel == "high" {
		alerts = append(alerts, &models.Alert{
			ID:            uuid.New(),
			TransactionID: tx.ID,
			Type:          "high_risk",
			Severity:      "high",
			Description:   "High-risk transaction detected",
			Status:        "open",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
	}

	if sr.SanctionList {
		alerts = append(alerts, &models.Alert{
			ID:            uuid.New(),
			TransactionID: tx.ID,
			Type:          "sanctions",
			Severity:      "critical",
			Description:   "Sanctions list match detected",
			Status:        "open",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
	}

	return alerts
}

func (s *amlService) updateRiskProfile(ctx context.Context, tx *models.Transaction) error {
	rp, err := s.repo.GetRiskProfile(ctx, tx.UserID)
	if err != nil {
		return err
	}

	if rp == nil {
		// Create new risk profile
		rp = &models.RiskProfile{
			ID:          uuid.New(),
			UserID:      tx.UserID,
			RiskScore:   tx.RiskScore,
			RiskLevel:   tx.RiskLevel,
			LastUpdated: time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
	} else {
		// Update existing risk profile
		// Simple averaging of risk scores
		rp.RiskScore = (rp.RiskScore + tx.RiskScore) / 2
		rp.RiskLevel = tx.RiskLevel // Use latest transaction's risk level
		rp.LastUpdated = time.Now()
		rp.UpdatedAt = time.Now()
	}

	return s.repo.UpdateRiskProfile(ctx, rp)
}
