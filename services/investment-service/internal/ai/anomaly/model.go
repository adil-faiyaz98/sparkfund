package anomaly

import (
	"math"
	"sort"
	"time"
)

// TransactionData represents a financial transaction for anomaly detection
type TransactionData struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	TransactionType string    `json:"transaction_type"` // DEPOSIT, WITHDRAWAL, INVESTMENT, SALE
	AssetID         string    `json:"asset_id,omitempty"`
	AssetType       string    `json:"asset_type,omitempty"`
	AssetSector     string    `json:"asset_sector,omitempty"`
	Timestamp       time.Time `json:"timestamp"`
}

// MarketData represents market conditions at a point in time
type MarketData struct {
	Timestamp          time.Time           `json:"timestamp"`
	MarketTrends       map[string]float64  `json:"market_trends"`       // Trends by sector
	EconomicIndicators map[string]float64  `json:"economic_indicators"` // e.g., "INFLATION", "GDP_GROWTH"
	SectorPerformance  map[string]float64  `json:"sector_performance"`  // Performance by sector
}

// UserBehaviorData represents a user's historical behavior
type UserBehaviorData struct {
	UserID                  string                `json:"user_id"`
	TransactionHistory      []TransactionData     `json:"transaction_history"`
	AssetPreferences        map[string]float64    `json:"asset_preferences"`        // Asset type to preference score
	SectorPreferences       map[string]float64    `json:"sector_preferences"`       // Sector to preference score
	TransactionPatterns     map[string]Pattern    `json:"transaction_patterns"`     // Type to pattern
	AverageTransactionSizes map[string]float64    `json:"average_transaction_sizes"` // Type to average size
}

// Pattern represents a statistical pattern
type Pattern struct {
	Mean   float64 `json:"mean"`
	StdDev float64 `json:"std_dev"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
}

// AnomalyDetectionResult represents the result of anomaly detection
type AnomalyDetectionResult struct {
	TransactionID    string    `json:"transaction_id"`
	UserID           string    `json:"user_id"`
	Timestamp        time.Time `json:"timestamp"`
	AnomalyScore     float64   `json:"anomaly_score"`      // 0.0 (normal) to 1.0 (highly anomalous)
	AnomalyLevel     string    `json:"anomaly_level"`      // LOW, MEDIUM, HIGH
	AnomalyIndicators []string  `json:"anomaly_indicators"` // Reasons for the anomaly score
	Action           string    `json:"action"`             // MONITOR, ALERT, INVESTIGATE
}

// AnomalyDetectionModel implements anomaly detection for investment transactions
type AnomalyDetectionModel struct {
	// Thresholds for different anomaly levels
	LowAnomalyThreshold    float64
	MediumAnomalyThreshold float64
	HighAnomalyThreshold   float64
}

// NewAnomalyDetectionModel creates a new anomaly detection model with default thresholds
func NewAnomalyDetectionModel() *AnomalyDetectionModel {
	return &AnomalyDetectionModel{
		LowAnomalyThreshold:    0.3,
		MediumAnomalyThreshold: 0.6,
		HighAnomalyThreshold:   0.8,
	}
}

// DetectAnomaly analyzes a transaction for anomalies
func (m *AnomalyDetectionModel) DetectAnomaly(transaction TransactionData, userData UserBehaviorData, marketData MarketData) AnomalyDetectionResult {
	// Initialize anomaly indicators
	var anomalyIndicators []string
	
	// Calculate base anomaly score
	anomalyScore := 0.0
	
	// Check for unusual amount
	amountScore := m.calculateAmountAnomaly(transaction, userData)
	if amountScore > 0.5 {
		anomalyIndicators = append(anomalyIndicators, "UNUSUAL_AMOUNT")
	}
	anomalyScore += amountScore * 0.3 // 30% weight
	
	// Check for unusual asset type or sector
	assetScore := m.calculateAssetAnomaly(transaction, userData)
	if assetScore > 0.5 {
		anomalyIndicators = append(anomalyIndicators, "UNUSUAL_ASSET_CHOICE")
	}
	anomalyScore += assetScore * 0.25 // 25% weight
	
	// Check for unusual timing
	timingScore := m.calculateTimingAnomaly(transaction, userData)
	if timingScore > 0.5 {
		anomalyIndicators = append(anomalyIndicators, "UNUSUAL_TIMING")
	}
	anomalyScore += timingScore * 0.2 // 20% weight
	
	// Check for market-contrary behavior
	marketScore := m.calculateMarketContraryAnomaly(transaction, marketData)
	if marketScore > 0.5 {
		anomalyIndicators = append(anomalyIndicators, "MARKET_CONTRARY")
	}
	anomalyScore += marketScore * 0.15 // 15% weight
	
	// Check for pattern break
	patternScore := m.calculatePatternBreakAnomaly(transaction, userData)
	if patternScore > 0.5 {
		anomalyIndicators = append(anomalyIndicators, "PATTERN_BREAK")
	}
	anomalyScore += patternScore * 0.1 // 10% weight
	
	// Determine anomaly level and action
	var anomalyLevel, action string
	
	if anomalyScore < m.LowAnomalyThreshold {
		anomalyLevel = "LOW"
		action = "MONITOR"
	} else if anomalyScore < m.MediumAnomalyThreshold {
		anomalyLevel = "MEDIUM"
		action = "ALERT"
	} else {
		anomalyLevel = "HIGH"
		if anomalyScore >= m.HighAnomalyThreshold {
			action = "INVESTIGATE"
		} else {
			action = "ALERT"
		}
	}
	
	// Create result
	result := AnomalyDetectionResult{
		TransactionID:     transaction.ID,
		UserID:            transaction.UserID,
		Timestamp:         time.Now(),
		AnomalyScore:      anomalyScore,
		AnomalyLevel:      anomalyLevel,
		AnomalyIndicators: anomalyIndicators,
		Action:            action,
	}
	
	return result
}

// calculateAmountAnomaly calculates a score based on transaction amount
func (m *AnomalyDetectionModel) calculateAmountAnomaly(transaction TransactionData, userData UserBehaviorData) float64 {
	// Get pattern for this transaction type
	pattern, exists := userData.TransactionPatterns[transaction.TransactionType]
	if !exists || pattern.StdDev == 0 {
		// If no pattern exists, check against average transaction size
		avgSize, exists := userData.AverageTransactionSizes[transaction.TransactionType]
		if !exists || avgSize == 0 {
			return 0.5 // No baseline, moderate anomaly score
		}
		
		// Calculate ratio of current amount to average
		ratio := transaction.Amount / avgSize
		if ratio > 5.0 {
			return 0.9 // Very large compared to average
		} else if ratio > 2.0 {
			return 0.6 // Moderately large
		} else if ratio < 0.2 {
			return 0.7 // Very small compared to average
		}
		
		return 0.3 // Somewhat normal
	}
	
	// Calculate z-score (standard deviations from mean)
	zScore := math.Abs(transaction.Amount - pattern.Mean) / pattern.StdDev
	
	// Convert z-score to a 0-1 score
	// A z-score of 3 (3 standard deviations) or more is considered unusual
	score := math.Min(1.0, zScore/3.0)
	
	return score
}

// calculateAssetAnomaly calculates a score based on asset type and sector
func (m *AnomalyDetectionModel) calculateAssetAnomaly(transaction TransactionData, userData UserBehaviorData) float64 {
	// Skip if not an investment transaction or missing asset info
	if transaction.TransactionType != "INVESTMENT" || transaction.AssetID == "" {
		return 0.0
	}
	
	// Check asset type preference
	assetTypeScore := 1.0
	if preference, exists := userData.AssetPreferences[transaction.AssetType]; exists {
		assetTypeScore = 1.0 - preference // Invert preference to get anomaly score
	}
	
	// Check sector preference
	sectorScore := 1.0
	if preference, exists := userData.SectorPreferences[transaction.AssetSector]; exists {
		sectorScore = 1.0 - preference // Invert preference to get anomaly score
	}
	
	// Combine scores (weighted average)
	combinedScore := (assetTypeScore * 0.6) + (sectorScore * 0.4)
	
	return combinedScore
}

// calculateTimingAnomaly calculates a score based on transaction timing
func (m *AnomalyDetectionModel) calculateTimingAnomaly(transaction TransactionData, userData UserBehaviorData) float64 {
	// Need at least a few transactions for timing analysis
	if len(userData.TransactionHistory) < 3 {
		return 0.3 // Not enough history, moderate-low score
	}
	
	// Extract hour of day for all transactions
	var hours []int
	for _, tx := range userData.TransactionHistory {
		hours = append(hours, tx.Timestamp.Hour())
	}
	
	// Calculate hour frequency
	hourFreq := make(map[int]int)
	for _, hour := range hours {
		hourFreq[hour]++
	}
	
	// Check if current transaction hour is common
	currentHour := transaction.Timestamp.Hour()
	currentHourFreq := hourFreq[currentHour]
	totalTransactions := len(userData.TransactionHistory)
	
	// Calculate frequency ratio
	freqRatio := float64(currentHourFreq) / float64(totalTransactions)
	
	// Uncommon hours have low frequency ratios
	if freqRatio < 0.05 {
		return 0.8 // Very uncommon hour
	} else if freqRatio < 0.1 {
		return 0.5 // Somewhat uncommon hour
	}
	
	return 0.2 // Common hour
}

// calculateMarketContraryAnomaly calculates a score based on market conditions
func (m *AnomalyDetectionModel) calculateMarketContraryAnomaly(transaction TransactionData, marketData MarketData) float64 {
	// Skip if not an investment transaction or missing asset info
	if transaction.TransactionType != "INVESTMENT" && transaction.TransactionType != "SALE" {
		return 0.0
	}
	
	// Check sector performance
	sectorPerformance, exists := marketData.SectorPerformance[transaction.AssetSector]
	if !exists {
		return 0.3 // No sector data, moderate-low score
	}
	
	// For investments, buying in declining sectors is contrary
	// For sales, selling in rising sectors is contrary
	var contraryScore float64
	
	if transaction.TransactionType == "INVESTMENT" {
		// Buying in declining sector
		if sectorPerformance < -0.05 {
			contraryScore = 0.8 // Strongly declining sector
		} else if sectorPerformance < -0.02 {
			contraryScore = 0.5 // Moderately declining sector
		} else {
			contraryScore = 0.2 // Stable or rising sector
		}
	} else if transaction.TransactionType == "SALE" {
		// Selling in rising sector
		if sectorPerformance > 0.05 {
			contraryScore = 0.8 // Strongly rising sector
		} else if sectorPerformance > 0.02 {
			contraryScore = 0.5 // Moderately rising sector
		} else {
			contraryScore = 0.2 // Stable or declining sector
		}
	}
	
	return contraryScore
}

// calculatePatternBreakAnomaly calculates a score based on pattern breaks
func (m *AnomalyDetectionModel) calculatePatternBreakAnomaly(transaction TransactionData, userData UserBehaviorData) float64 {
	// Need sufficient history for pattern analysis
	if len(userData.TransactionHistory) < 5 {
		return 0.3 // Not enough history, moderate-low score
	}
	
	// Check for unusual sequence
	// For example, multiple withdrawals in short succession
	if transaction.TransactionType == "WITHDRAWAL" {
		// Count recent withdrawals
		recentWithdrawals := 0
		oneDayAgo := transaction.Timestamp.Add(-24 * time.Hour)
		
		for _, tx := range userData.TransactionHistory {
			if tx.TransactionType == "WITHDRAWAL" && tx.Timestamp.After(oneDayAgo) {
				recentWithdrawals++
			}
		}
		
		// Multiple withdrawals in a day is unusual
		if recentWithdrawals >= 3 {
			return 0.9 // Very unusual pattern
		} else if recentWithdrawals >= 2 {
			return 0.6 // Somewhat unusual pattern
		}
	}
	
	// Check for unusual investment pattern
	if transaction.TransactionType == "INVESTMENT" {
		// Get recent investments
		var recentInvestments []TransactionData
		oneWeekAgo := transaction.Timestamp.Add(-7 * 24 * time.Hour)
		
		for _, tx := range userData.TransactionHistory {
			if tx.TransactionType == "INVESTMENT" && tx.Timestamp.After(oneWeekAgo) {
				recentInvestments = append(recentInvestments, tx)
			}
		}
		
		// Calculate average investment amount
		var totalAmount float64
		for _, tx := range recentInvestments {
			totalAmount += tx.Amount
		}
		
		var avgAmount float64
		if len(recentInvestments) > 0 {
			avgAmount = totalAmount / float64(len(recentInvestments))
		}
		
		// Check if current investment is much larger than recent average
		if avgAmount > 0 && transaction.Amount > avgAmount*3 {
			return 0.8 // Much larger than recent average
		} else if avgAmount > 0 && transaction.Amount > avgAmount*2 {
			return 0.5 // Moderately larger than recent average
		}
	}
	
	return 0.2 // No pattern break detected
}

// BuildUserBehaviorData builds a user behavior profile from transaction history
func (m *AnomalyDetectionModel) BuildUserBehaviorData(userID string, transactions []TransactionData) UserBehaviorData {
	userData := UserBehaviorData{
		UserID:                  userID,
		TransactionHistory:      transactions,
		AssetPreferences:        make(map[string]float64),
		SectorPreferences:       make(map[string]float64),
		TransactionPatterns:     make(map[string]Pattern),
		AverageTransactionSizes: make(map[string]float64),
	}
	
	// Group transactions by type
	typeGroups := make(map[string][]TransactionData)
	for _, tx := range transactions {
		typeGroups[tx.TransactionType] = append(typeGroups[tx.TransactionType], tx)
	}
	
	// Calculate patterns for each transaction type
	for txType, txs := range typeGroups {
		// Calculate amount statistics
		var amounts []float64
		var sum, sumSquared float64
		
		for _, tx := range txs {
			amounts = append(amounts, tx.Amount)
			sum += tx.Amount
			sumSquared += tx.Amount * tx.Amount
		}
		
		// Calculate mean and standard deviation
		mean := sum / float64(len(txs))
		variance := (sumSquared / float64(len(txs))) - (mean * mean)
		stdDev := math.Sqrt(variance)
		
		// Find min and max
		sort.Float64s(amounts)
		min := amounts[0]
		max := amounts[len(amounts)-1]
		
		// Store pattern
		userData.TransactionPatterns[txType] = Pattern{
			Mean:   mean,
			StdDev: stdDev,
			Min:    min,
			Max:    max,
		}
		
		// Store average transaction size
		userData.AverageTransactionSizes[txType] = mean
	}
	
	// Calculate asset preferences
	assetTypeCounts := make(map[string]int)
	sectorCounts := make(map[string]int)
	
	var investmentCount int
	for _, tx := range transactions {
		if tx.TransactionType == "INVESTMENT" && tx.AssetType != "" {
			assetTypeCounts[tx.AssetType]++
			sectorCounts[tx.AssetSector]++
			investmentCount++
		}
	}
	
	// Convert counts to preferences (0.0 to 1.0)
	if investmentCount > 0 {
		for assetType, count := range assetTypeCounts {
			userData.AssetPreferences[assetType] = float64(count) / float64(investmentCount)
		}
		
		for sector, count := range sectorCounts {
			userData.SectorPreferences[sector] = float64(count) / float64(investmentCount)
		}
	}
	
	return userData
}
