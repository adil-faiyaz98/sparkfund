package security

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sparkfund/services/investment-service/internal/ai/anomaly"
	"github.com/sparkfund/services/investment-service/internal/ai/fraud"
)

// SecurityService provides fraud and anomaly detection services
type SecurityService struct {
	fraudModel   *fraud.FraudDetectionModel
	anomalyModel *anomaly.AnomalyDetectionModel
	repository   Repository
}

// Repository defines the interface for accessing security-related data
type Repository interface {
	// User profile methods
	GetUserSecurityProfile(ctx context.Context, userID string) (*fraud.UserProfile, error)
	SaveUserSecurityProfile(ctx context.Context, profile *fraud.UserProfile) error
	
	// Transaction methods
	GetUserTransactions(ctx context.Context, userID string, limit int) ([]fraud.Transaction, error)
	SaveTransaction(ctx context.Context, transaction *fraud.Transaction) error
	
	// Fraud detection methods
	SaveFraudDetectionResult(ctx context.Context, result *fraud.FraudDetectionResult) error
	GetRecentFraudDetectionResults(ctx context.Context, userID string, limit int) ([]fraud.FraudDetectionResult, error)
	
	// Anomaly detection methods
	GetUserBehaviorData(ctx context.Context, userID string) (*anomaly.UserBehaviorData, error)
	SaveUserBehaviorData(ctx context.Context, data *anomaly.UserBehaviorData) error
	SaveAnomalyDetectionResult(ctx context.Context, result *anomaly.AnomalyDetectionResult) error
	GetRecentAnomalyDetectionResults(ctx context.Context, userID string, limit int) ([]anomaly.AnomalyDetectionResult, error)
	
	// Market data methods
	GetLatestMarketData(ctx context.Context) (*anomaly.MarketData, error)
}

// NewSecurityService creates a new security service
func NewSecurityService(repository Repository) *SecurityService {
	return &SecurityService{
		fraudModel:   fraud.NewFraudDetectionModel(),
		anomalyModel: anomaly.NewAnomalyDetectionModel(),
		repository:   repository,
	}
}

// AnalyzeTransaction performs fraud and anomaly detection on a transaction
type AnalysisResult struct {
	TransactionID      string                        `json:"transaction_id"`
	UserID             string                        `json:"user_id"`
	Timestamp          time.Time                     `json:"timestamp"`
	FraudDetection     fraud.FraudDetectionResult    `json:"fraud_detection"`
	AnomalyDetection   anomaly.AnomalyDetectionResult `json:"anomaly_detection"`
	OverallRiskScore   float64                       `json:"overall_risk_score"`
	OverallRiskLevel   string                        `json:"overall_risk_level"`
	RecommendedAction  string                        `json:"recommended_action"`
}

// AnalyzeTransaction analyzes a transaction for fraud and anomalies
func (s *SecurityService) AnalyzeTransaction(ctx context.Context, transaction fraud.Transaction) (*AnalysisResult, error) {
	// Convert to anomaly transaction format
	anomalyTx := anomaly.TransactionData{
		ID:              transaction.ID,
		UserID:          transaction.UserID,
		Amount:          transaction.Amount,
		Currency:        transaction.Currency,
		TransactionType: transaction.TransactionType,
		AssetID:         transaction.AssetID,
		Timestamp:       transaction.Timestamp,
	}
	
	// Get user security profile
	userProfile, err := s.repository.GetUserSecurityProfile(ctx, transaction.UserID)
	if err != nil {
		log.Printf("Warning: failed to get user security profile: %v", err)
		// Create a default profile
		userProfile = &fraud.UserProfile{
			UserID:    transaction.UserID,
			CreatedAt: time.Now(),
			LastLogin: time.Now(),
			RiskScore: 0.5, // Default moderate risk
		}
	}
	
	// Get user transaction history
	transactionHistory, err := s.repository.GetUserTransactions(ctx, transaction.UserID, 100)
	if err != nil {
		log.Printf("Warning: failed to get user transaction history: %v", err)
		transactionHistory = []fraud.Transaction{}
	}
	
	// Get user behavior data
	userBehaviorData, err := s.repository.GetUserBehaviorData(ctx, transaction.UserID)
	if err != nil {
		log.Printf("Warning: failed to get user behavior data: %v", err)
		// Convert transaction history to anomaly format
		var anomalyTxHistory []anomaly.TransactionData
		for _, tx := range transactionHistory {
			anomalyTxHistory = append(anomalyTxHistory, anomaly.TransactionData{
				ID:              tx.ID,
				UserID:          tx.UserID,
				Amount:          tx.Amount,
				Currency:        tx.Currency,
				TransactionType: tx.TransactionType,
				AssetID:         tx.AssetID,
				Timestamp:       tx.Timestamp,
			})
		}
		// Build user behavior data
		userData := s.anomalyModel.BuildUserBehaviorData(transaction.UserID, anomalyTxHistory)
		userBehaviorData = &userData
	}
	
	// Get market data
	marketData, err := s.repository.GetLatestMarketData(ctx)
	if err != nil {
		log.Printf("Warning: failed to get market data: %v", err)
		marketData = &anomaly.MarketData{
			Timestamp:          time.Now(),
			MarketTrends:       make(map[string]float64),
			EconomicIndicators: make(map[string]float64),
			SectorPerformance:  make(map[string]float64),
		}
	}
	
	// Perform fraud detection
	fraudResult := s.fraudModel.DetectFraud(transaction, *userProfile, transactionHistory)
	
	// Perform anomaly detection
	anomalyResult := s.anomalyModel.DetectAnomaly(anomalyTx, *userBehaviorData, *marketData)
	
	// Calculate overall risk
	overallRiskScore := (fraudResult.FraudScore * 0.6) + (anomalyResult.AnomalyScore * 0.4)
	
	// Determine overall risk level
	var overallRiskLevel, recommendedAction string
	
	if overallRiskScore < 0.3 {
		overallRiskLevel = "LOW"
		recommendedAction = "APPROVE"
	} else if overallRiskScore < 0.6 {
		overallRiskLevel = "MEDIUM"
		recommendedAction = "REVIEW"
	} else {
		overallRiskLevel = "HIGH"
		if overallRiskScore >= 0.8 {
			recommendedAction = "REJECT"
		} else {
			recommendedAction = "REVIEW"
		}
	}
	
	// Create analysis result
	result := &AnalysisResult{
		TransactionID:     transaction.ID,
		UserID:            transaction.UserID,
		Timestamp:         time.Now(),
		FraudDetection:    fraudResult,
		AnomalyDetection:  anomalyResult,
		OverallRiskScore:  overallRiskScore,
		OverallRiskLevel:  overallRiskLevel,
		RecommendedAction: recommendedAction,
	}
	
	// Save results
	err = s.repository.SaveFraudDetectionResult(ctx, &fraudResult)
	if err != nil {
		log.Printf("Warning: failed to save fraud detection result: %v", err)
	}
	
	err = s.repository.SaveAnomalyDetectionResult(ctx, &anomalyResult)
	if err != nil {
		log.Printf("Warning: failed to save anomaly detection result: %v", err)
	}
	
	// Update user behavior data
	userBehaviorData.TransactionHistory = append(userBehaviorData.TransactionHistory, anomalyTx)
	err = s.repository.SaveUserBehaviorData(ctx, userBehaviorData)
	if err != nil {
		log.Printf("Warning: failed to save user behavior data: %v", err)
	}
	
	// Update user security profile
	// Add current IP and device to usual lists if not already present
	addToStringSlice := func(slice []string, item string) []string {
		for _, s := range slice {
			if s == item {
				return slice // Already present
			}
		}
		return append(slice, item)
	}
	
	userProfile.UsualIPAddresses = addToStringSlice(userProfile.UsualIPAddresses, transaction.IPAddress)
	userProfile.UsualDeviceIDs = addToStringSlice(userProfile.UsualDeviceIDs, transaction.DeviceID)
	
	// Add location if not already present
	locationPresent := false
	for _, loc := range userProfile.UsualLocations {
		if loc.Latitude == transaction.Location.Latitude && loc.Longitude == transaction.Location.Longitude {
			locationPresent = true
			break
		}
	}
	if !locationPresent {
		userProfile.UsualLocations = append(userProfile.UsualLocations, transaction.Location)
	}
	
	// Update average transaction amount
	if userProfile.AverageTransactionAmount == 0 {
		userProfile.AverageTransactionAmount = transaction.Amount
	} else {
		// Simple moving average
		userProfile.AverageTransactionAmount = (userProfile.AverageTransactionAmount*0.9 + transaction.Amount*0.1)
	}
	
	// Update risk score based on fraud detection
	userProfile.RiskScore = (userProfile.RiskScore*0.9 + fraudResult.FraudScore*0.1)
	
	// Save updated profile
	err = s.repository.SaveUserSecurityProfile(ctx, userProfile)
	if err != nil {
		log.Printf("Warning: failed to save user security profile: %v", err)
	}
	
	return result, nil
}

// GetUserRiskProfile gets a user's risk profile
func (s *SecurityService) GetUserRiskProfile(ctx context.Context, userID string) (*UserRiskProfile, error) {
	// Get user security profile
	userProfile, err := s.repository.GetUserSecurityProfile(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	// Get recent fraud detection results
	fraudResults, err := s.repository.GetRecentFraudDetectionResults(ctx, userID, 10)
	if err != nil {
		log.Printf("Warning: failed to get recent fraud detection results: %v", err)
		fraudResults = []fraud.FraudDetectionResult{}
	}
	
	// Get recent anomaly detection results
	anomalyResults, err := s.repository.GetRecentAnomalyDetectionResults(ctx, userID, 10)
	if err != nil {
		log.Printf("Warning: failed to get recent anomaly detection results: %v", err)
		anomalyResults = []anomaly.AnomalyDetectionResult{}
	}
	
	// Calculate average fraud and anomaly scores
	var avgFraudScore, avgAnomalyScore float64
	
	if len(fraudResults) > 0 {
		var sum float64
		for _, result := range fraudResults {
			sum += result.FraudScore
		}
		avgFraudScore = sum / float64(len(fraudResults))
	}
	
	if len(anomalyResults) > 0 {
		var sum float64
		for _, result := range anomalyResults {
			sum += result.AnomalyScore
		}
		avgAnomalyScore = sum / float64(len(anomalyResults))
	}
	
	// Create risk profile
	riskProfile := &UserRiskProfile{
		UserID:                userID,
		RiskScore:             userProfile.RiskScore,
		AverageFraudScore:     avgFraudScore,
		AverageAnomalyScore:   avgAnomalyScore,
		RecentFraudIndicators: getRecentIndicators(fraudResults),
		RecentAnomalyIndicators: getRecentIndicators(anomalyResults),
		LastUpdated:           time.Now(),
	}
	
	return riskProfile, nil
}

// UserRiskProfile represents a user's risk profile
type UserRiskProfile struct {
	UserID                 string    `json:"user_id"`
	RiskScore              float64   `json:"risk_score"`
	AverageFraudScore      float64   `json:"average_fraud_score"`
	AverageAnomalyScore    float64   `json:"average_anomaly_score"`
	RecentFraudIndicators  []string  `json:"recent_fraud_indicators"`
	RecentAnomalyIndicators []string  `json:"recent_anomaly_indicators"`
	LastUpdated            time.Time `json:"last_updated"`
}

// Helper functions

// getRecentIndicators extracts unique indicators from recent results
func getRecentIndicators(fraudResults []fraud.FraudDetectionResult) []string {
	// Use a map to deduplicate indicators
	indicatorMap := make(map[string]bool)
	
	for _, result := range fraudResults {
		for _, indicator := range result.FraudIndicators {
			indicatorMap[indicator] = true
		}
	}
	
	// Convert map keys to slice
	var indicators []string
	for indicator := range indicatorMap {
		indicators = append(indicators, indicator)
	}
	
	return indicators
}

// getRecentIndicators extracts unique indicators from recent results
func getRecentIndicators(anomalyResults []anomaly.AnomalyDetectionResult) []string {
	// Use a map to deduplicate indicators
	indicatorMap := make(map[string]bool)
	
	for _, result := range anomalyResults {
		for _, indicator := range result.AnomalyIndicators {
			indicatorMap[indicator] = true
		}
	}
	
	// Convert map keys to slice
	var indicators []string
	for indicator := range indicatorMap {
		indicators = append(indicators, indicator)
	}
	
	return indicators
}
