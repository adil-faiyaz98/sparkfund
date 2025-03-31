package security

import (
	"context"
	"time"
)

// RiskAssessment handles comprehensive risk evaluation
type RiskAssessment struct {
	config *RiskAssessmentConfig
	store  RiskAssessmentStore
}

// RiskAssessmentConfig defines configuration for risk assessment
type RiskAssessmentConfig struct {
	BehaviorThreshold    float64
	TransactionThreshold float64
	LocationThreshold    float64
	DeviceThreshold      float64
	TimeThreshold        float64
	HistoryWindow        time.Duration
	MaxFailedAttempts    int
	RiskWeights          map[string]float64
}

// RiskAssessmentStore defines the interface for risk assessment storage
type RiskAssessmentStore interface {
	GetUserBehavior(ctx context.Context, userID string) ([]BehaviorEvent, error)
	GetTransactionHistory(ctx context.Context, userID string) ([]TransactionEvent, error)
	GetLocationHistory(ctx context.Context, userID string) ([]LocationEvent, error)
	GetDeviceHistory(ctx context.Context, userID string) ([]DeviceEvent, error)
	GetFailedAttempts(ctx context.Context, userID string) ([]FailedAttempt, error)
}

// BehaviorEvent represents a user behavior event
type BehaviorEvent struct {
	Timestamp  time.Time
	EventType  string
	Details    map[string]interface{}
	RiskScore  float64
	Confidence float64
}

// TransactionEvent represents a transaction event
type TransactionEvent struct {
	Timestamp time.Time
	Amount    float64
	Type      string
	Status    string
	RiskScore float64
	Location  Location
	Device    DeviceInfo
}

// LocationEvent represents a location event
type LocationEvent struct {
	Timestamp  time.Time
	Location   Location
	RiskScore  float64
	Confidence float64
}

// DeviceEvent represents a device event
type DeviceEvent struct {
	Timestamp  time.Time
	Device     DeviceInfo
	RiskScore  float64
	Confidence float64
}

// FailedAttempt represents a failed authentication attempt
type FailedAttempt struct {
	Timestamp time.Time
	IP        string
	Device    DeviceInfo
	Location  Location
	Reason    string
}

// NewRiskAssessment creates a new risk assessment instance
func NewRiskAssessment(config RiskAssessmentConfig, store RiskAssessmentStore) *RiskAssessment {
	return &RiskAssessment{
		config: &config,
		store:  store,
	}
}

// AssessRisk performs comprehensive risk assessment
func (r *RiskAssessment) AssessRisk(ctx context.Context, userID string) (*RiskAssessmentResult, error) {
	// Get user behavior
	behavior, err := r.store.GetUserBehavior(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get transaction history
	transactions, err := r.store.GetTransactionHistory(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get location history
	locations, err := r.store.GetLocationHistory(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get device history
	devices, err := r.store.GetDeviceHistory(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get failed attempts
	failedAttempts, err := r.store.GetFailedAttempts(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Calculate risk scores
	behaviorScore := r.calculateBehaviorRisk(behavior)
	transactionScore := r.calculateTransactionRisk(transactions)
	locationScore := r.calculateLocationRisk(locations)
	deviceScore := r.calculateDeviceRisk(devices)
	failedAttemptScore := r.calculateFailedAttemptRisk(failedAttempts)

	// Calculate overall risk score
	overallScore := r.calculateOverallRisk(
		behaviorScore,
		transactionScore,
		locationScore,
		deviceScore,
		failedAttemptScore,
	)

	// Determine risk level
	riskLevel := r.determineRiskLevel(overallScore)

	return &RiskAssessmentResult{
		OverallScore:       overallScore,
		BehaviorScore:      behaviorScore,
		TransactionScore:   transactionScore,
		LocationScore:      locationScore,
		DeviceScore:        deviceScore,
		FailedAttemptScore: failedAttemptScore,
		RiskLevel:          riskLevel,
		Factors:            r.identifyRiskFactors(behavior, transactions, locations, devices, failedAttempts),
		Recommendations:    r.generateRecommendations(riskLevel),
	}, nil
}

// RiskAssessmentResult represents the result of a risk assessment
type RiskAssessmentResult struct {
	OverallScore       float64
	BehaviorScore      float64
	TransactionScore   float64
	LocationScore      float64
	DeviceScore        float64
	FailedAttemptScore float64
	RiskLevel          RiskLevel
	Factors            []RiskFactor
	Recommendations    []string
}

// RiskFactor represents a specific risk factor
type RiskFactor struct {
	Type        string
	Description string
	Score       float64
	Confidence  float64
	Details     map[string]interface{}
}

// Helper functions

func (r *RiskAssessment) calculateBehaviorRisk(behavior []BehaviorEvent) float64 {
	if len(behavior) == 0 {
		return 0.0
	}

	var totalScore float64
	var totalConfidence float64

	for _, event := range behavior {
		totalScore += event.RiskScore * event.Confidence
		totalConfidence += event.Confidence
	}

	if totalConfidence == 0 {
		return 0.0
	}

	return totalScore / totalConfidence
}

func (r *RiskAssessment) calculateTransactionRisk(transactions []TransactionEvent) float64 {
	if len(transactions) == 0 {
		return 0.0
	}

	var totalScore float64
	for _, tx := range transactions {
		totalScore += tx.RiskScore
	}

	return totalScore / float64(len(transactions))
}

func (r *RiskAssessment) calculateLocationRisk(locations []LocationEvent) float64 {
	if len(locations) == 0 {
		return 0.0
	}

	var totalScore float64
	var totalConfidence float64

	for _, loc := range locations {
		totalScore += loc.RiskScore * loc.Confidence
		totalConfidence += loc.Confidence
	}

	if totalConfidence == 0 {
		return 0.0
	}

	return totalScore / totalConfidence
}

func (r *RiskAssessment) calculateDeviceRisk(devices []DeviceEvent) float64 {
	if len(devices) == 0 {
		return 0.0
	}

	var totalScore float64
	var totalConfidence float64

	for _, device := range devices {
		totalScore += device.RiskScore * device.Confidence
		totalConfidence += device.Confidence
	}

	if totalConfidence == 0 {
		return 0.0
	}

	return totalScore / totalConfidence
}

func (r *RiskAssessment) calculateFailedAttemptRisk(attempts []FailedAttempt) float64 {
	if len(attempts) == 0 {
		return 0.0
	}

	// Calculate risk based on recent failed attempts
	recentAttempts := r.filterRecentAttempts(attempts)
	if len(recentAttempts) == 0 {
		return 0.0
	}

	// Higher risk for more recent and frequent attempts
	riskScore := float64(len(recentAttempts)) / float64(r.config.MaxFailedAttempts)
	return riskScore
}

func (r *RiskAssessment) calculateOverallRisk(
	behaviorScore, transactionScore, locationScore, deviceScore, failedAttemptScore float64,
) float64 {
	weights := r.config.RiskWeights
	if weights == nil {
		weights = map[string]float64{
			"behavior":       0.3,
			"transaction":    0.25,
			"location":       0.2,
			"device":         0.15,
			"failed_attempt": 0.1,
		}
	}

	return behaviorScore*weights["behavior"] +
		transactionScore*weights["transaction"] +
		locationScore*weights["location"] +
		deviceScore*weights["device"] +
		failedAttemptScore*weights["failed_attempt"]
}

func (r *RiskAssessment) determineRiskLevel(score float64) RiskLevel {
	switch {
	case score >= 0.8:
		return RiskLevel{
			Level:  "high",
			Score:  score,
			Reason: "High risk score based on multiple factors",
		}
	case score >= 0.5:
		return RiskLevel{
			Level:  "medium",
			Score:  score,
			Reason: "Medium risk score based on multiple factors",
		}
	default:
		return RiskLevel{
			Level:  "low",
			Score:  score,
			Reason: "Low risk score based on multiple factors",
		}
	}
}

func (r *RiskAssessment) filterRecentAttempts(attempts []FailedAttempt) []FailedAttempt {
	cutoff := time.Now().Add(-r.config.HistoryWindow)
	var recent []FailedAttempt
	for _, attempt := range attempts {
		if attempt.Timestamp.After(cutoff) {
			recent = append(recent, attempt)
		}
	}
	return recent
}

func (r *RiskAssessment) identifyRiskFactors(
	behavior []BehaviorEvent,
	transactions []TransactionEvent,
	locations []LocationEvent,
	devices []DeviceEvent,
	attempts []FailedAttempt,
) []RiskFactor {
	var factors []RiskFactor

	// Add behavior risk factors
	if len(behavior) > 0 {
		factors = append(factors, RiskFactor{
			Type:        "behavior",
			Description: "Unusual behavior patterns detected",
			Score:       r.calculateBehaviorRisk(behavior),
			Confidence:  r.calculateConfidence(behavior),
			Details:     map[string]interface{}{"event_count": len(behavior)},
		})
	}

	// Add transaction risk factors
	if len(transactions) > 0 {
		factors = append(factors, RiskFactor{
			Type:        "transaction",
			Description: "Suspicious transaction patterns detected",
			Score:       r.calculateTransactionRisk(transactions),
			Confidence:  r.calculateConfidence(transactions),
			Details:     map[string]interface{}{"transaction_count": len(transactions)},
		})
	}

	// Add location risk factors
	if len(locations) > 0 {
		factors = append(factors, RiskFactor{
			Type:        "location",
			Description: "Unusual location patterns detected",
			Score:       r.calculateLocationRisk(locations),
			Confidence:  r.calculateConfidence(locations),
			Details:     map[string]interface{}{"location_count": len(locations)},
		})
	}

	// Add device risk factors
	if len(devices) > 0 {
		factors = append(factors, RiskFactor{
			Type:        "device",
			Description: "Suspicious device patterns detected",
			Score:       r.calculateDeviceRisk(devices),
			Confidence:  r.calculateConfidence(devices),
			Details:     map[string]interface{}{"device_count": len(devices)},
		})
	}

	// Add failed attempt risk factors
	if len(attempts) > 0 {
		factors = append(factors, RiskFactor{
			Type:        "failed_attempt",
			Description: "Multiple failed authentication attempts detected",
			Score:       r.calculateFailedAttemptRisk(attempts),
			Confidence:  1.0,
			Details:     map[string]interface{}{"attempt_count": len(attempts)},
		})
	}

	return factors
}

func (r *RiskAssessment) generateRecommendations(riskLevel RiskLevel) []string {
	var recommendations []string

	switch riskLevel.Level {
	case "high":
		recommendations = append(recommendations,
			"Require additional verification for all transactions",
			"Implement stricter access controls",
			"Monitor all activities closely",
			"Consider temporary account restrictions",
		)
	case "medium":
		recommendations = append(recommendations,
			"Implement enhanced monitoring",
			"Require additional verification for high-value transactions",
			"Review recent activities",
			"Consider implementing additional security measures",
		)
	default:
		recommendations = append(recommendations,
			"Continue regular monitoring",
			"Maintain standard security measures",
			"Review security policies periodically",
		)
	}

	return recommendations
}

func (r *RiskAssessment) calculateConfidence(events interface{}) float64 {
	// Implement confidence calculation based on event type and count
	return 0.8
}
