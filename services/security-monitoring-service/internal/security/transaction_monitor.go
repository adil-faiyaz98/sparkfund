package security

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TransactionMonitor handles transaction monitoring and control
type TransactionMonitor struct {
	ubaManager      *UBAManager
	notificationSvc *NotificationService
	config          *TransactionMonitorConfig
	mu              sync.RWMutex
}

// TransactionMonitorConfig defines transaction monitoring configuration
type TransactionMonitorConfig struct {
	// Industry standard thresholds (5% less than standard)
	IndustryThresholds struct {
		MaxAmountPerDay          float64
		MaxAmountPerHour         float64
		MaxAmountPerMinute       float64
		MaxTransactionsPerDay    int
		MaxTransactionsPerHour   int
		MaxTransactionsPerMinute int
	}

	// Risk thresholds
	RiskThresholds struct {
		HighRiskThreshold   float64
		MediumRiskThreshold float64
		LowRiskThreshold    float64
	}

	// Monitoring settings
	Monitoring struct {
		CheckInterval    time.Duration
		HistoryRetention time.Duration
		MaxPendingTime   time.Duration
		MaxReviewTime    time.Duration
	}

	// Notification settings
	Notifications struct {
		EnableEmailAlerts bool
		EnableSMSAlerts   bool
		EnablePhoneAlerts bool
		AlertThreshold    float64
	}
}

// NewTransactionMonitor creates a new transaction monitor
func NewTransactionMonitor(config TransactionMonitorConfig, ubaManager *UBAManager, notificationSvc *NotificationService) *TransactionMonitor {
	return &TransactionMonitor{
		ubaManager:      ubaManager,
		notificationSvc: notificationSvc,
		config:          &config,
	}
}

// MonitorTransaction monitors a transaction for potential risks
func (m *TransactionMonitor) MonitorTransaction(ctx context.Context, tx *Transaction) (*TransactionDecision, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Analyze transaction using UBA
	analysis, err := m.ubaManager.AnalyzeTransaction(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze transaction: %v", err)
	}

	// Create decision
	decision := &TransactionDecision{
		TransactionID: tx.ID,
		Timestamp:     time.Now(),
		RiskScore:     analysis.RiskScore,
		Flags:         analysis.Flags,
		Status:        "pending",
	}

	// Check transaction limits against industry standard (5% less)
	if m.exceedsIndustryLimits(tx) {
		decision.Status = "rejected"
		decision.Reason = "industry_limits_exceeded"
		m.notifyLimitExceeded(ctx, tx, analysis)
		return decision, nil
	}

	// Apply risk-based decision
	switch {
	case analysis.RiskScore >= m.config.RiskThresholds.HighRiskThreshold:
		decision.Status = "blocked"
		decision.Reason = "high_risk"
		m.notifyHighRiskTransaction(ctx, tx, analysis)
	case analysis.RiskScore >= m.config.RiskThresholds.MediumRiskThreshold:
		decision.Status = "pending_review"
		decision.Reason = "medium_risk"
		m.notifyMediumRiskTransaction(ctx, tx, analysis)
	case analysis.RiskScore >= m.config.RiskThresholds.LowRiskThreshold:
		decision.Status = "approved"
		decision.Reason = "low_risk"
	default:
		decision.Status = "approved"
		decision.Reason = "normal"
	}

	// Check for unusual patterns
	if m.hasUnusualPatterns(tx) {
		decision.Status = "pending_review"
		decision.Reason = "unusual_patterns"
		m.notifyUnusualPatterns(ctx, tx, analysis)
	}

	// Check for suspicious amounts
	if m.isSuspiciousAmount(tx) {
		decision.Status = "pending_review"
		decision.Reason = "suspicious_amount"
		m.notifySuspiciousAmount(ctx, tx, analysis)
	}

	// Check for unauthorized recipients
	if m.isUnauthorizedRecipient(tx) {
		decision.Status = "blocked"
		decision.Reason = "unauthorized_recipient"
		m.notifyUnauthorizedRecipient(ctx, tx, analysis)
	}

	return decision, nil
}

// exceedsIndustryLimits checks if transaction exceeds industry standard limits (5% less)
func (m *TransactionMonitor) exceedsIndustryLimits(tx *Transaction) bool {
	// Check daily limits
	if m.getDailyAmount(tx.UserID) > m.config.IndustryThresholds.MaxAmountPerDay {
		return true
	}

	// Check hourly limits
	if m.getHourlyAmount(tx.UserID) > m.config.IndustryThresholds.MaxAmountPerHour {
		return true
	}

	// Check minute limits
	if m.getMinuteAmount(tx.UserID) > m.config.IndustryThresholds.MaxAmountPerMinute {
		return true
	}

	// Check transaction count limits
	if m.getDailyTransactionCount(tx.UserID) > m.config.IndustryThresholds.MaxTransactionsPerDay {
		return true
	}

	if m.getHourlyTransactionCount(tx.UserID) > m.config.IndustryThresholds.MaxTransactionsPerHour {
		return true
	}

	if m.getMinuteTransactionCount(tx.UserID) > m.config.IndustryThresholds.MaxTransactionsPerMinute {
		return true
	}

	return false
}

// hasUnusualPatterns checks for unusual transaction patterns
func (m *TransactionMonitor) hasUnusualPatterns(tx *Transaction) bool {
	// Check for unusual time patterns
	if m.isUnusualTime(tx) {
		return true
	}

	// Check for unusual location patterns
	if m.isUnusualLocation(tx) {
		return true
	}

	// Check for unusual device patterns
	if m.isUnusualDevice(tx) {
		return true
	}

	// Check for unusual recipient patterns
	if m.isUnusualRecipient(tx) {
		return true
	}

	return false
}

// isSuspiciousAmount checks for suspicious transaction amounts
func (m *TransactionMonitor) isSuspiciousAmount(tx *Transaction) bool {
	// Check for round numbers
	if m.isRoundNumber(tx.Amount) {
		return true
	}

	// Check for common fraud amounts
	if m.isCommonFraudAmount(tx.Amount) {
		return true
	}

	// Check for unusual decimal places
	if m.hasUnusualDecimals(tx.Amount) {
		return true
	}

	// Check for significant deviation from average
	if m.hasSignificantDeviation(tx) {
		return true
	}

	return false
}

// isUnauthorizedRecipient checks if recipient is authorized
func (m *TransactionMonitor) isUnauthorizedRecipient(tx *Transaction) bool {
	// Check against authorized recipients list
	if !m.isAuthorizedRecipient(tx.Recipient) {
		return true
	}

	// Check for high-risk recipients
	if m.isHighRiskRecipient(tx.Recipient) {
		return true
	}

	// Check for suspicious recipient patterns
	if m.hasSuspiciousRecipientPatterns(tx) {
		return true
	}

	return false
}

// Notification methods
func (m *TransactionMonitor) notifyHighRiskTransaction(ctx context.Context, tx *Transaction, analysis *TransactionAnalysis) {
	if analysis.RiskScore >= m.config.Notifications.AlertThreshold {
		alert := &TransactionAlert{
			Type:          "high_risk",
			UserID:        tx.UserID,
			TransactionID: tx.ID,
			Amount:        tx.Amount,
			Currency:      tx.Currency,
			Timestamp:     tx.Timestamp,
			Reason:        "High risk transaction detected",
			Details:       make(map[string]interface{}),
		}
		m.notificationSvc.NotifyTransactionAlert(ctx, alert)
	}
}

func (m *TransactionMonitor) notifyMediumRiskTransaction(ctx context.Context, tx *Transaction, analysis *TransactionAnalysis) {
	if analysis.RiskScore >= m.config.Notifications.AlertThreshold {
		alert := &TransactionAlert{
			Type:          "unusual",
			UserID:        tx.UserID,
			TransactionID: tx.ID,
			Amount:        tx.Amount,
			Currency:      tx.Currency,
			Timestamp:     tx.Timestamp,
			Reason:        "Medium risk transaction detected",
			Details:       make(map[string]interface{}),
		}
		m.notificationSvc.NotifyTransactionAlert(ctx, alert)
	}
}

func (m *TransactionMonitor) notifyUnusualPatterns(ctx context.Context, tx *Transaction, analysis *TransactionAnalysis) {
	alert := &TransactionAlert{
		Type:          "unusual",
		UserID:        tx.UserID,
		TransactionID: tx.ID,
		Amount:        tx.Amount,
		Currency:      tx.Currency,
		Timestamp:     tx.Timestamp,
		Reason:        "Unusual transaction patterns detected",
		Details:       make(map[string]interface{}),
	}
	m.notificationSvc.NotifyTransactionAlert(ctx, alert)
}

func (m *TransactionMonitor) notifySuspiciousAmount(ctx context.Context, tx *Transaction, analysis *TransactionAnalysis) {
	alert := &TransactionAlert{
		Type:          "unusual",
		UserID:        tx.UserID,
		TransactionID: tx.ID,
		Amount:        tx.Amount,
		Currency:      tx.Currency,
		Timestamp:     tx.Timestamp,
		Reason:        "Suspicious transaction amount detected",
		Details:       make(map[string]interface{}),
	}
	m.notificationSvc.NotifyTransactionAlert(ctx, alert)
}

func (m *TransactionMonitor) notifyUnauthorizedRecipient(ctx context.Context, tx *Transaction, analysis *TransactionAnalysis) {
	alert := &TransactionAlert{
		Type:          "blocked",
		UserID:        tx.UserID,
		TransactionID: tx.ID,
		Amount:        tx.Amount,
		Currency:      tx.Currency,
		Timestamp:     tx.Timestamp,
		Reason:        "Unauthorized recipient detected",
		Details:       make(map[string]interface{}),
	}
	m.notificationSvc.NotifyTransactionAlert(ctx, alert)
}

func (m *TransactionMonitor) notifyLimitExceeded(ctx context.Context, tx *Transaction, analysis *TransactionAnalysis) {
	alert := &TransactionAlert{
		Type:          "limit_exceeded",
		UserID:        tx.UserID,
		TransactionID: tx.ID,
		Amount:        tx.Amount,
		Currency:      tx.Currency,
		Timestamp:     tx.Timestamp,
		Reason:        "Transaction limit exceeded",
		Details:       make(map[string]interface{}),
	}
	m.notificationSvc.NotifyTransactionAlert(ctx, alert)
}

// Helper functions
func (m *TransactionMonitor) getDailyAmount(userID string) float64 {
	// Implement daily amount calculation
	return 0.0
}

func (m *TransactionMonitor) getHourlyAmount(userID string) float64 {
	// Implement hourly amount calculation
	return 0.0
}

func (m *TransactionMonitor) getMinuteAmount(userID string) float64 {
	// Implement minute amount calculation
	return 0.0
}

func (m *TransactionMonitor) getDailyTransactionCount(userID string) int {
	// Implement daily transaction count calculation
	return 0
}

func (m *TransactionMonitor) getHourlyTransactionCount(userID string) int {
	// Implement hourly transaction count calculation
	return 0
}

func (m *TransactionMonitor) getMinuteTransactionCount(userID string) int {
	// Implement minute transaction count calculation
	return 0
}

func (m *TransactionMonitor) isUnusualTime(tx *Transaction) bool {
	// Implement unusual time check
	return false
}

func (m *TransactionMonitor) isUnusualLocation(tx *Transaction) bool {
	// Implement unusual location check
	return false
}

func (m *TransactionMonitor) isUnusualDevice(tx *Transaction) bool {
	// Implement unusual device check
	return false
}

func (m *TransactionMonitor) isUnusualRecipient(tx *Transaction) bool {
	// Implement unusual recipient check
	return false
}

func (m *TransactionMonitor) isRoundNumber(amount float64) bool {
	// Implement round number check
	return false
}

func (m *TransactionMonitor) isCommonFraudAmount(amount float64) bool {
	// Implement common fraud amount check
	return false
}

func (m *TransactionMonitor) hasUnusualDecimals(amount float64) bool {
	// Implement unusual decimals check
	return false
}

func (m *TransactionMonitor) hasSignificantDeviation(tx *Transaction) bool {
	// Implement significant deviation check
	return false
}

func (m *TransactionMonitor) isAuthorizedRecipient(recipient string) bool {
	// Implement authorized recipient check
	return false
}

func (m *TransactionMonitor) isHighRiskRecipient(recipient string) bool {
	// Implement high-risk recipient check
	return false
}

func (m *TransactionMonitor) hasSuspiciousRecipientPatterns(tx *Transaction) bool {
	// Implement suspicious recipient patterns check
	return false
}

// TransactionDecision represents the decision made for a transaction
type TransactionDecision struct {
	TransactionID string
	Timestamp     time.Time
	RiskScore     float64
	Flags         []string
	Status        string
	Reason        string
}
