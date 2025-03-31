package security

import (
	"context"
	"fmt"
	"sort"
	"time"
)

// TransactionHistory provides transaction history and trend analysis
type TransactionHistory struct {
	store  TransactionStore
	config *TransactionHistoryConfig
}

// TransactionHistoryConfig defines configuration for transaction history
type TransactionHistoryConfig struct {
	MaxHistoryMonths int
	TrendAnalysis    struct {
		EnableAmountTrends      bool
		EnableLocationTrends    bool
		EnableTimeTrends        bool
		EnableRecipientTrends   bool
		MinTransactionsForTrend int
	}
}

// MonthlyTransactionSummary represents a summary of transactions for a month
type MonthlyTransactionSummary struct {
	Month            time.Time
	TotalAmount      float64
	TransactionCount int
	AverageAmount    float64
	MaxAmount        float64
	MinAmount        float64
	Locations        []string
	TopRecipients    []RecipientSummary
	TimeDistribution map[string]int // Distribution by hour
	RiskDistribution map[string]int // Distribution by risk level
}

// RecipientSummary represents summary of transactions with a recipient
type RecipientSummary struct {
	RecipientID     string
	TotalAmount     float64
	Count           int
	LastTransaction time.Time
}

// TransactionTrend represents a trend analysis
type TransactionTrend struct {
	Type           string
	Period         string
	StartDate      time.Time
	EndDate        time.Time
	Values         []float64
	Labels         []string
	TrendDirection string // "increasing", "decreasing", "stable"
	ChangePercent  float64
}

// NewTransactionHistory creates a new transaction history service
func NewTransactionHistory(store TransactionStore, config TransactionHistoryConfig) *TransactionHistory {
	return &TransactionHistory{
		store:  store,
		config: &config,
	}
}

// GetMonthlyHistory retrieves transaction history for a specific month
func (h *TransactionHistory) GetMonthlyHistory(ctx context.Context, userID string, month time.Time) (*MonthlyTransactionSummary, error) {
	// Get start and end of month
	startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	// Get transactions for the month
	transactions, err := h.store.GetTransactionsByTimeRange(ctx, userID, startOfMonth, endOfMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %v", err)
	}

	// Calculate summary
	summary := &MonthlyTransactionSummary{
		Month:            startOfMonth,
		TimeDistribution: make(map[string]int),
		RiskDistribution: make(map[string]int),
	}

	// Process transactions
	for _, tx := range transactions {
		summary.TotalAmount += tx.Amount
		summary.TransactionCount++

		// Update max/min amounts
		if tx.Amount > summary.MaxAmount {
			summary.MaxAmount = tx.Amount
		}
		if summary.MinAmount == 0 || tx.Amount < summary.MinAmount {
			summary.MinAmount = tx.Amount
		}

		// Add location if not already present
		locationKey := fmt.Sprintf("%s,%s", tx.Location.Country, tx.Location.City)
		if !containsString(summary.Locations, locationKey) {
			summary.Locations = append(summary.Locations, locationKey)
		}

		// Update time distribution
		hour := tx.Timestamp.Format("15")
		summary.TimeDistribution[hour]++

		// Update risk distribution
		summary.RiskDistribution[tx.RiskLevel]++

		// Update recipient summary
		h.updateRecipientSummary(summary, tx)
	}

	// Calculate average amount
	if summary.TransactionCount > 0 {
		summary.AverageAmount = summary.TotalAmount / float64(summary.TransactionCount)
	}

	// Sort recipients by total amount
	sort.Slice(summary.TopRecipients, func(i, j int) bool {
		return summary.TopRecipients[i].TotalAmount > summary.TopRecipients[j].TotalAmount
	})

	// Keep only top 10 recipients
	if len(summary.TopRecipients) > 10 {
		summary.TopRecipients = summary.TopRecipients[:10]
	}

	return summary, nil
}

// GetTransactionTrends analyzes transaction trends over time
func (h *TransactionHistory) GetTransactionTrends(ctx context.Context, userID string, months int) ([]TransactionTrend, error) {
	if months > h.config.MaxHistoryMonths {
		months = h.config.MaxHistoryMonths
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, -months, 0)

	// Get transactions for the period
	transactions, err := h.store.GetTransactionsByTimeRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %v", err)
	}

	var trends []TransactionTrend

	// Amount trends
	if h.config.TrendAnalysis.EnableAmountTrends {
		amountTrend := h.analyzeAmountTrend(transactions)
		trends = append(trends, amountTrend)
	}

	// Location trends
	if h.config.TrendAnalysis.EnableLocationTrends {
		locationTrend := h.analyzeLocationTrend(transactions)
		trends = append(trends, locationTrend)
	}

	// Time distribution trends
	if h.config.TrendAnalysis.EnableTimeTrends {
		timeTrend := h.analyzeTimeTrend(transactions)
		trends = append(trends, timeTrend)
	}

	// Recipient trends
	if h.config.TrendAnalysis.EnableRecipientTrends {
		recipientTrend := h.analyzeRecipientTrend(transactions)
		trends = append(trends, recipientTrend)
	}

	return trends, nil
}

// Helper functions

func (h *TransactionHistory) updateRecipientSummary(summary *MonthlyTransactionSummary, tx *Transaction) {
	for i, recipient := range summary.TopRecipients {
		if recipient.RecipientID == tx.RecipientID {
			summary.TopRecipients[i].TotalAmount += tx.Amount
			summary.TopRecipients[i].Count++
			if tx.Timestamp.After(recipient.LastTransaction) {
				summary.TopRecipients[i].LastTransaction = tx.Timestamp
			}
			return
		}
	}

	summary.TopRecipients = append(summary.TopRecipients, RecipientSummary{
		RecipientID:     tx.RecipientID,
		TotalAmount:     tx.Amount,
		Count:           1,
		LastTransaction: tx.Timestamp,
	})
}

func (h *TransactionHistory) analyzeAmountTrend(transactions []*Transaction) TransactionTrend {
	// Group transactions by month
	monthlyAmounts := make(map[string]float64)
	monthlyLabels := make(map[string]string)

	for _, tx := range transactions {
		monthKey := tx.Timestamp.Format("2006-01")
		monthlyAmounts[monthKey] += tx.Amount
		monthlyLabels[monthKey] = tx.Timestamp.Format("Jan 2006")
	}

	// Convert to sorted slice
	var values []float64
	var labels []string
	for monthKey := range monthlyAmounts {
		values = append(values, monthlyAmounts[monthKey])
		labels = append(labels, monthlyLabels[monthKey])
	}

	// Calculate trend direction and change
	direction, change := h.calculateTrendDirection(values)

	return TransactionTrend{
		Type:           "amount",
		Period:         "monthly",
		Values:         values,
		Labels:         labels,
		TrendDirection: direction,
		ChangePercent:  change,
	}
}

func (h *TransactionHistory) analyzeLocationTrend(transactions []*Transaction) TransactionTrend {
	// Group transactions by location
	locationCounts := make(map[string]int)
	locationLabels := make(map[string]string)

	for _, tx := range transactions {
		locationKey := fmt.Sprintf("%s,%s", tx.Location.Country, tx.Location.City)
		locationCounts[locationKey]++
		locationLabels[locationKey] = fmt.Sprintf("%s, %s", tx.Location.City, tx.Location.Country)
	}

	// Convert to slices
	var values []float64
	var labels []string
	for locationKey := range locationCounts {
		values = append(values, float64(locationCounts[locationKey]))
		labels = append(labels, locationLabels[locationKey])
	}

	return TransactionTrend{
		Type:           "location",
		Period:         "total",
		Values:         values,
		Labels:         labels,
		TrendDirection: "distribution",
		ChangePercent:  0,
	}
}

func (h *TransactionHistory) analyzeTimeTrend(transactions []*Transaction) TransactionTrend {
	// Group transactions by hour
	hourlyCounts := make(map[string]int)
	hourlyLabels := make(map[string]string)

	for _, tx := range transactions {
		hourKey := tx.Timestamp.Format("15")
		hourlyCounts[hourKey]++
		hourlyLabels[hourKey] = fmt.Sprintf("%s:00", hourKey)
	}

	// Convert to slices
	var values []float64
	var labels []string
	for hourKey := range hourlyCounts {
		values = append(values, float64(hourlyCounts[hourKey]))
		labels = append(labels, hourlyLabels[hourKey])
	}

	return TransactionTrend{
		Type:           "time",
		Period:         "hourly",
		Values:         values,
		Labels:         labels,
		TrendDirection: "distribution",
		ChangePercent:  0,
	}
}

func (h *TransactionHistory) analyzeRecipientTrend(transactions []*Transaction) TransactionTrend {
	// Group transactions by recipient
	recipientAmounts := make(map[string]float64)
	recipientLabels := make(map[string]string)

	for _, tx := range transactions {
		recipientAmounts[tx.RecipientID] += tx.Amount
		recipientLabels[tx.RecipientID] = tx.RecipientName
	}

	// Convert to slices
	var values []float64
	var labels []string
	for recipientID := range recipientAmounts {
		values = append(values, recipientAmounts[recipientID])
		labels = append(labels, recipientLabels[recipientID])
	}

	// Sort by amount
	sort.Slice(values, func(i, j int) bool {
		return values[i] > values[j]
	})
	sort.Slice(labels, func(i, j int) bool {
		return recipientAmounts[labels[i]] > recipientAmounts[labels[j]]
	})

	// Keep only top 10
	if len(values) > 10 {
		values = values[:10]
		labels = labels[:10]
	}

	return TransactionTrend{
		Type:           "recipient",
		Period:         "total",
		Values:         values,
		Labels:         labels,
		TrendDirection: "distribution",
		ChangePercent:  0,
	}
}

func (h *TransactionHistory) calculateTrendDirection(values []float64) (string, float64) {
	if len(values) < 2 {
		return "stable", 0
	}

	first := values[0]
	last := values[len(values)-1]
	change := ((last - first) / first) * 100

	if change > 5 {
		return "increasing", change
	} else if change < -5 {
		return "decreasing", change
	}
	return "stable", change
}

func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
