package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/sparkfund/credit-scoring-service/internal/model"
	"github.com/sparkfund/credit-scoring-service/internal/repository"
)

type CreditService interface {
	ProcessCreditCheck(ctx context.Context, req *model.CreditCheckRequest) (*model.CreditHistoryResponse, error)
	CalculateCreditScore(ctx context.Context, userID uuid.UUID) (*model.CreditScoreResponse, error)
	GetCreditHistory(ctx context.Context, userID uuid.UUID) ([]*model.CreditHistoryResponse, error)
	GetCreditScore(ctx context.Context, userID uuid.UUID) (*model.CreditScoreResponse, error)
	UpdateCreditHistory(ctx context.Context, historyID uuid.UUID, status model.CreditHistoryStatus) error
}

type creditService struct {
	creditRepo repository.CreditRepository
}

func NewCreditService(creditRepo repository.CreditRepository) CreditService {
	return &creditService{
		creditRepo: creditRepo,
	}
}

func (s *creditService) ProcessCreditCheck(ctx context.Context, req *model.CreditCheckRequest) (*model.CreditHistoryResponse, error) {
	// Create credit history record
	history := &model.CreditHistory{
		UserID:            req.UserID,
		AccountType:       req.AccountType,
		Institution:       req.Institution,
		AccountNumber:     req.AccountNumber,
		Status:            model.CreditHistoryStatusActive,
		CreditLimit:       req.CreditLimit,
		CurrentBalance:    req.CurrentBalance,
		OpenDate:          req.OpenDate,
		CloseDate:         req.CloseDate,
		LastPaymentDate:   req.LastPaymentDate,
		LastPaymentAmount: req.LastPaymentAmount,
	}

	// Initialize payment history
	paymentHistory := []model.PaymentRecord{}
	if req.LastPaymentDate != nil && req.LastPaymentAmount != nil {
		paymentHistory = append(paymentHistory, model.PaymentRecord{
			Date:   *req.LastPaymentDate,
			Amount: *req.LastPaymentAmount,
			Status: "on_time",
		})
	}

	// Convert payment history to JSON
	paymentHistoryJSON, err := json.Marshal(paymentHistory)
	if err != nil {
		return nil, err
	}
	history.PaymentHistory = string(paymentHistoryJSON)

	// Save credit history
	if err := s.creditRepo.CreateCreditHistory(history); err != nil {
		return nil, err
	}

	// Calculate new credit score
	score, err := s.CalculateCreditScore(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return s.toCreditHistoryResponse(history), nil
}

func (s *creditService) CalculateCreditScore(ctx context.Context, userID uuid.UUID) (*model.CreditScoreResponse, error) {
	// Get all credit history for user
	histories, err := s.creditRepo.GetCreditHistories(userID)
	if err != nil {
		return nil, err
	}

	// Initialize scoring factors
	factors := &model.CreditScoreFactors{
		PaymentHistoryWeight:    0.35,
		CreditUtilizationWeight: 0.30,
		CreditHistoryWeight:     0.15,
		AccountMixWeight:        0.10,
		NewCreditWeight:         0.10,
	}

	// Calculate individual scores
	paymentHistoryScore := s.calculatePaymentHistoryScore(histories)
	creditUtilizationScore := s.calculateCreditUtilizationScore(histories)
	creditHistoryScore := s.calculateCreditHistoryScore(histories)
	accountMixScore := s.calculateAccountMixScore(histories)
	newCreditScore := s.calculateNewCreditScore(histories)

	// Calculate weighted average
	totalScore := int(
		paymentHistoryScore*factors.PaymentHistoryWeight +
			creditUtilizationScore*factors.CreditUtilizationWeight +
			creditHistoryScore*factors.CreditHistoryWeight +
			accountMixScore*factors.AccountMixWeight +
			newCreditScore*factors.NewCreditWeight,
	)

	// Ensure score is within range
	if totalScore < 300 {
		totalScore = 300
	} else if totalScore > 850 {
		totalScore = 850
	}

	// Determine score range
	var scoreRange model.CreditScoreRange
	switch {
	case totalScore >= 800:
		scoreRange = model.CreditScoreRangeExcellent
	case totalScore >= 740:
		scoreRange = model.CreditScoreRangeVeryGood
	case totalScore >= 670:
		scoreRange = model.CreditScoreRangeGood
	case totalScore >= 580:
		scoreRange = model.CreditScoreRangeFair
	default:
		scoreRange = model.CreditScoreRangePoor
	}

	// Generate recommendations
	recommendations := s.generateRecommendations(histories, totalScore)

	// Create or update credit score
	score := &model.CreditScore{
		UserID:          userID,
		Score:           totalScore,
		ScoreRange:      scoreRange,
		LastUpdated:     time.Now(),
		Factors:         s.serializeFactors(factors),
		Recommendations: s.serializeRecommendations(recommendations),
	}

	if err := s.creditRepo.UpsertCreditScore(score); err != nil {
		return nil, err
	}

	return s.toCreditScoreResponse(score), nil
}

func (s *creditService) GetCreditHistory(ctx context.Context, userID uuid.UUID) ([]*model.CreditHistoryResponse, error) {
	histories, err := s.creditRepo.GetCreditHistories(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*model.CreditHistoryResponse, len(histories))
	for i, history := range histories {
		responses[i] = s.toCreditHistoryResponse(history)
	}

	return responses, nil
}

func (s *creditService) GetCreditScore(ctx context.Context, userID uuid.UUID) (*model.CreditScoreResponse, error) {
	score, err := s.creditRepo.GetCreditScore(userID)
	if err != nil {
		return nil, err
	}

	return s.toCreditScoreResponse(score), nil
}

func (s *creditService) UpdateCreditHistory(ctx context.Context, historyID uuid.UUID, status model.CreditHistoryStatus) error {
	history, err := s.creditRepo.GetCreditHistory(historyID)
	if err != nil {
		return err
	}

	history.Status = status
	if status == model.CreditHistoryStatusClosed {
		now := time.Now()
		history.CloseDate = &now
	}

	return s.creditRepo.UpdateCreditHistory(history)
}

// Helper functions for credit score calculation
func (s *creditService) calculatePaymentHistoryScore(histories []*model.CreditHistory) float64 {
	if len(histories) == 0 {
		return 0.0
	}

	var totalScore float64
	var totalWeight float64

	for _, history := range histories {
		var paymentHistory []model.PaymentRecord
		json.Unmarshal([]byte(history.PaymentHistory), &paymentHistory)

		if len(paymentHistory) == 0 {
			continue
		}

		// Calculate payment history score based on last 12 months
		var onTimePayments, latePayments, missedPayments int
		var totalPayments float64
		var paymentScore float64

		for _, payment := range paymentHistory {
			totalPayments++
			switch payment.Status {
			case "on_time":
				onTimePayments++
			case "late":
				latePayments++
			case "missed":
				missedPayments++
			}
		}

		// Calculate payment score (0-100)
		if totalPayments > 0 {
			paymentScore = float64(onTimePayments) / totalPayments * 100
			paymentScore -= float64(latePayments) * 10
			paymentScore -= float64(missedPayments) * 30
		}

		// Weight based on account age
		weight := calculateAccountAgeWeight(history.OpenDate)
		totalScore += paymentScore * weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		return 0.0
	}

	return totalScore / totalWeight
}

func (s *creditService) calculateCreditUtilizationScore(histories []*model.CreditHistory) float64 {
	if len(histories) == 0 {
		return 0.0
	}

	var totalUtilization float64
	var totalWeight float64

	for _, history := range histories {
		if history.CreditLimit <= 0 {
			continue
		}

		utilization := history.CurrentBalance / history.CreditLimit
		weight := calculateAccountAgeWeight(history.OpenDate)

		// Score based on utilization ratio (0-100)
		var utilizationScore float64
		switch {
		case utilization <= 0.1:
			utilizationScore = 100
		case utilization <= 0.3:
			utilizationScore = 80
		case utilization <= 0.5:
			utilizationScore = 60
		case utilization <= 0.7:
			utilizationScore = 40
		case utilization <= 0.9:
			utilizationScore = 20
		default:
			utilizationScore = 10
		}

		totalUtilization += utilizationScore * weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		return 0.0
	}

	return totalUtilization / totalWeight
}

func (s *creditService) calculateCreditHistoryScore(histories []*model.CreditHistory) float64 {
	if len(histories) == 0 {
		return 0.0
	}

	var oldestAccount time.Time
	var totalAccounts int
	var totalCreditLimit float64

	for _, history := range histories {
		if history.OpenDate.Before(oldestAccount) || oldestAccount.IsZero() {
			oldestAccount = history.OpenDate
		}
		totalAccounts++
		totalCreditLimit += history.CreditLimit
	}

	// Calculate average account age in years
	accountAge := time.Since(oldestAccount).Hours() / 8760 // 8760 hours in a year

	// Score based on account age (0-100)
	var ageScore float64
	switch {
	case accountAge >= 10:
		ageScore = 100
	case accountAge >= 7:
		ageScore = 80
	case accountAge >= 5:
		ageScore = 60
	case accountAge >= 3:
		ageScore = 40
	case accountAge >= 1:
		ageScore = 20
	default:
		ageScore = 10
	}

	// Score based on number of accounts (0-100)
	var accountScore float64
	switch {
	case totalAccounts >= 5:
		accountScore = 100
	case totalAccounts >= 3:
		accountScore = 80
	case totalAccounts >= 2:
		accountScore = 60
	case totalAccounts == 1:
		accountScore = 40
	default:
		accountScore = 20
	}

	// Combine scores with weights
	return (ageScore*0.6 + accountScore*0.4)
}

func (s *creditService) calculateAccountMixScore(histories []*model.CreditHistory) float64 {
	if len(histories) == 0 {
		return 0.0
	}

	accountTypes := make(map[string]int)
	for _, history := range histories {
		accountTypes[history.AccountType]++
	}

	// Score based on account mix (0-100)
	var mixScore float64
	totalTypes := len(accountTypes)
	switch {
	case totalTypes >= 4:
		mixScore = 100
	case totalTypes >= 3:
		mixScore = 80
	case totalTypes >= 2:
		mixScore = 60
	case totalTypes == 1:
		mixScore = 40
	default:
		mixScore = 20
	}

	return mixScore
}

func (s *creditService) calculateNewCreditScore(histories []*model.CreditHistory) float64 {
	if len(histories) == 0 {
		return 0.0
	}

	var recentAccounts int
	var totalAccounts int
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)

	for _, history := range histories {
		totalAccounts++
		if history.OpenDate.After(sixMonthsAgo) {
			recentAccounts++
		}
	}

	// Score based on new credit (0-100)
	var newCreditScore float64
	if totalAccounts > 0 {
		newCreditRatio := float64(recentAccounts) / float64(totalAccounts)
		switch {
		case newCreditRatio <= 0.1:
			newCreditScore = 100
		case newCreditRatio <= 0.2:
			newCreditScore = 80
		case newCreditRatio <= 0.3:
			newCreditScore = 60
		case newCreditRatio <= 0.4:
			newCreditScore = 40
		case newCreditRatio <= 0.5:
			newCreditScore = 20
		default:
			newCreditScore = 10
		}
	}

	return newCreditScore
}

func (s *creditService) generateRecommendations(histories []*model.CreditHistory, score int) *model.CreditScoreRecommendations {
	recommendations := &model.CreditScoreRecommendations{
		PaymentHistory:    []string{},
		CreditUtilization: []string{},
		CreditHistory:     []string{},
		AccountMix:        []string{},
		NewCredit:         []string{},
	}

	// Payment History Recommendations
	var latePayments, missedPayments int
	for _, history := range histories {
		var paymentHistory []model.PaymentRecord
		json.Unmarshal([]byte(history.PaymentHistory), &paymentHistory)
		for _, payment := range paymentHistory {
			switch payment.Status {
			case "late":
				latePayments++
			case "missed":
				missedPayments++
			}
		}
	}

	if latePayments > 0 {
		recommendations.PaymentHistory = append(recommendations.PaymentHistory,
			"Make payments on time to improve your credit score")
	}
	if missedPayments > 0 {
		recommendations.PaymentHistory = append(recommendations.PaymentHistory,
			"Avoid missing payments as they significantly impact your credit score")
	}

	// Credit Utilization Recommendations
	for _, history := range histories {
		if history.CreditLimit > 0 {
			utilization := history.CurrentBalance / history.CreditLimit
			if utilization > 0.7 {
				recommendations.CreditUtilization = append(recommendations.CreditUtilization,
					"Reduce your credit utilization to below 70%")
			}
		}
	}

	// Credit History Recommendations
	var oldestAccount time.Time
	for _, history := range histories {
		if history.OpenDate.Before(oldestAccount) || oldestAccount.IsZero() {
			oldestAccount = history.OpenDate
		}
	}

	accountAge := time.Since(oldestAccount).Hours() / 8760
	if accountAge < 2 {
		recommendations.CreditHistory = append(recommendations.CreditHistory,
			"Maintain your existing credit accounts to build a longer credit history")
	}

	// Account Mix Recommendations
	accountTypes := make(map[string]int)
	for _, history := range histories {
		accountTypes[history.AccountType]++
	}

	if len(accountTypes) < 2 {
		recommendations.AccountMix = append(recommendations.AccountMix,
			"Consider diversifying your credit mix with different types of accounts")
	}

	// New Credit Recommendations
	var recentAccounts int
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)
	for _, history := range histories {
		if history.OpenDate.After(sixMonthsAgo) {
			recentAccounts++
		}
	}

	if recentAccounts > 2 {
		recommendations.NewCredit = append(recommendations.NewCredit,
			"Avoid opening new credit accounts in the next 6 months")
	}

	return recommendations
}

func calculateAccountAgeWeight(openDate time.Time) float64 {
	age := time.Since(openDate).Hours() / 8760 // Convert to years
	switch {
	case age >= 10:
		return 1.0
	case age >= 7:
		return 0.8
	case age >= 5:
		return 0.6
	case age >= 3:
		return 0.4
	case age >= 1:
		return 0.2
	default:
		return 0.1
	}
}

func (s *creditService) serializeFactors(factors *model.CreditScoreFactors) string {
	json, _ := json.Marshal(factors)
	return string(json)
}

func (s *creditService) serializeRecommendations(recommendations *model.CreditScoreRecommendations) string {
	json, _ := json.Marshal(recommendations)
	return string(json)
}

func (s *creditService) toCreditScoreResponse(score *model.CreditScore) *model.CreditScoreResponse {
	var factors model.CreditScoreFactors
	var recommendations model.CreditScoreRecommendations

	json.Unmarshal([]byte(score.Factors), &factors)
	json.Unmarshal([]byte(score.Recommendations), &recommendations)

	return &model.CreditScoreResponse{
		ID:              score.ID,
		UserID:          score.UserID,
		Score:           score.Score,
		ScoreRange:      score.ScoreRange,
		LastUpdated:     score.LastUpdated,
		Factors:         []string{}, // Convert factors to string array
		Recommendations: []string{}, // Convert recommendations to string array
		CreatedAt:       score.CreatedAt,
		UpdatedAt:       score.UpdatedAt,
	}
}

func (s *creditService) toCreditHistoryResponse(history *model.CreditHistory) *model.CreditHistoryResponse {
	var paymentHistory []model.PaymentRecord
	json.Unmarshal([]byte(history.PaymentHistory), &paymentHistory)

	return &model.CreditHistoryResponse{
		ID:                history.ID,
		UserID:            history.UserID,
		AccountType:       history.AccountType,
		Institution:       history.Institution,
		AccountNumber:     history.AccountNumber,
		Status:            history.Status,
		CreditLimit:       history.CreditLimit,
		CurrentBalance:    history.CurrentBalance,
		PaymentHistory:    paymentHistory,
		OpenDate:          history.OpenDate,
		CloseDate:         history.CloseDate,
		LastPaymentDate:   history.LastPaymentDate,
		LastPaymentAmount: history.LastPaymentAmount,
		CreatedAt:         history.CreatedAt,
		UpdatedAt:         history.UpdatedAt,
	}
}
