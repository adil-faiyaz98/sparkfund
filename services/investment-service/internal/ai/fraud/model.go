package fraud

import (
	"math"
	"time"
)

// Transaction represents a financial transaction to be analyzed for fraud
type Transaction struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	TransactionType string    `json:"transaction_type"` // DEPOSIT, WITHDRAWAL, INVESTMENT, SALE
	AssetID         string    `json:"asset_id,omitempty"`
	Timestamp       time.Time `json:"timestamp"`
	IPAddress       string    `json:"ip_address"`
	DeviceID        string    `json:"device_id"`
	Location        Location  `json:"location"`
	UserAgent       string    `json:"user_agent"`
}

// Location represents geographical coordinates
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
	City      string  `json:"city"`
}

// UserProfile contains user information relevant for fraud detection
type UserProfile struct {
	UserID                string    `json:"user_id"`
	CreatedAt             time.Time `json:"created_at"`
	LastLogin             time.Time `json:"last_login"`
	UsualIPAddresses      []string  `json:"usual_ip_addresses"`
	UsualDeviceIDs        []string  `json:"usual_device_ids"`
	UsualLocations        []Location `json:"usual_locations"`
	AverageTransactionAmount float64 `json:"average_transaction_amount"`
	TransactionFrequency  float64   `json:"transaction_frequency"` // Transactions per day
	RiskScore             float64   `json:"risk_score"`           // 0.0 (low risk) to 1.0 (high risk)
}

// FraudDetectionResult represents the result of fraud detection analysis
type FraudDetectionResult struct {
	TransactionID  string    `json:"transaction_id"`
	UserID         string    `json:"user_id"`
	Timestamp      time.Time `json:"timestamp"`
	FraudScore     float64   `json:"fraud_score"`      // 0.0 (legitimate) to 1.0 (fraudulent)
	FraudLevel     string    `json:"fraud_level"`      // LOW, MEDIUM, HIGH
	FraudIndicators []string  `json:"fraud_indicators"` // Reasons for the fraud score
	Action         string    `json:"action"`           // APPROVE, REVIEW, REJECT
}

// FraudDetectionModel implements fraud detection for investment transactions
type FraudDetectionModel struct {
	// Thresholds for different fraud levels
	LowRiskThreshold    float64
	MediumRiskThreshold float64
	HighRiskThreshold   float64
}

// NewFraudDetectionModel creates a new fraud detection model with default thresholds
func NewFraudDetectionModel() *FraudDetectionModel {
	return &FraudDetectionModel{
		LowRiskThreshold:    0.3,
		MediumRiskThreshold: 0.6,
		HighRiskThreshold:   0.8,
	}
}

// DetectFraud analyzes a transaction for potential fraud
func (m *FraudDetectionModel) DetectFraud(transaction Transaction, userProfile UserProfile, userTransactionHistory []Transaction) FraudDetectionResult {
	// Initialize fraud indicators
	var fraudIndicators []string
	
	// Calculate base fraud score
	fraudScore := 0.0
	
	// Check for unusual amount
	amountScore := m.calculateAmountScore(transaction, userProfile, userTransactionHistory)
	if amountScore > 0.5 {
		fraudIndicators = append(fraudIndicators, "UNUSUAL_AMOUNT")
	}
	fraudScore += amountScore * 0.3 // 30% weight
	
	// Check for unusual location
	locationScore := m.calculateLocationScore(transaction, userProfile)
	if locationScore > 0.5 {
		fraudIndicators = append(fraudIndicators, "UNUSUAL_LOCATION")
	}
	fraudScore += locationScore * 0.25 // 25% weight
	
	// Check for unusual device or IP
	deviceScore := m.calculateDeviceScore(transaction, userProfile)
	if deviceScore > 0.5 {
		fraudIndicators = append(fraudIndicators, "UNUSUAL_DEVICE")
	}
	fraudScore += deviceScore * 0.2 // 20% weight
	
	// Check for unusual frequency
	frequencyScore := m.calculateFrequencyScore(transaction, userTransactionHistory)
	if frequencyScore > 0.5 {
		fraudIndicators = append(fraudIndicators, "UNUSUAL_FREQUENCY")
	}
	fraudScore += frequencyScore * 0.15 // 15% weight
	
	// Check for velocity (multiple transactions in short time)
	velocityScore := m.calculateVelocityScore(transaction, userTransactionHistory)
	if velocityScore > 0.5 {
		fraudIndicators = append(fraudIndicators, "VELOCITY_ALERT")
	}
	fraudScore += velocityScore * 0.1 // 10% weight
	
	// Determine fraud level and action
	var fraudLevel, action string
	
	if fraudScore < m.LowRiskThreshold {
		fraudLevel = "LOW"
		action = "APPROVE"
	} else if fraudScore < m.MediumRiskThreshold {
		fraudLevel = "MEDIUM"
		action = "REVIEW"
	} else {
		fraudLevel = "HIGH"
		if fraudScore >= m.HighRiskThreshold {
			action = "REJECT"
		} else {
			action = "REVIEW"
		}
	}
	
	// Create result
	result := FraudDetectionResult{
		TransactionID:   transaction.ID,
		UserID:          transaction.UserID,
		Timestamp:       time.Now(),
		FraudScore:      fraudScore,
		FraudLevel:      fraudLevel,
		FraudIndicators: fraudIndicators,
		Action:          action,
	}
	
	return result
}

// calculateAmountScore calculates a score based on transaction amount
func (m *FraudDetectionModel) calculateAmountScore(transaction Transaction, userProfile UserProfile, history []Transaction) float64 {
	// If no history or average amount is zero, use moderate score
	if len(history) == 0 || userProfile.AverageTransactionAmount == 0 {
		// For new users, large amounts are more suspicious
		if transaction.Amount > 10000 {
			return 0.7
		}
		return 0.3
	}
	
	// Calculate z-score (standard deviations from mean)
	var sum, sumSquared float64
	for _, tx := range history {
		sum += tx.Amount
		sumSquared += tx.Amount * tx.Amount
	}
	
	mean := sum / float64(len(history))
	variance := (sumSquared / float64(len(history))) - (mean * mean)
	stdDev := math.Sqrt(variance)
	
	// Avoid division by zero
	if stdDev == 0 {
		stdDev = 1.0
	}
	
	zScore := math.Abs(transaction.Amount - mean) / stdDev
	
	// Convert z-score to a 0-1 score
	// A z-score of 3 (3 standard deviations) or more is considered unusual
	score := math.Min(1.0, zScore/3.0)
	
	return score
}

// calculateLocationScore calculates a score based on transaction location
func (m *FraudDetectionModel) calculateLocationScore(transaction Transaction, userProfile UserProfile) float64 {
	// If no usual locations, use moderate score
	if len(userProfile.UsualLocations) == 0 {
		return 0.5
	}
	
	// Find minimum distance to any usual location
	minDistance := math.MaxFloat64
	for _, usualLocation := range userProfile.UsualLocations {
		distance := calculateDistance(
			transaction.Location.Latitude, 
			transaction.Location.Longitude,
			usualLocation.Latitude,
			usualLocation.Longitude,
		)
		if distance < minDistance {
			minDistance = distance
		}
	}
	
	// Convert distance to a 0-1 score
	// Distances over 1000km are considered unusual
	score := math.Min(1.0, minDistance/1000.0)
	
	return score
}

// calculateDeviceScore calculates a score based on device and IP
func (m *FraudDetectionModel) calculateDeviceScore(transaction Transaction, userProfile UserProfile) float64 {
	// Check if device ID is in usual devices
	deviceKnown := false
	for _, usualDeviceID := range userProfile.UsualDeviceIDs {
		if transaction.DeviceID == usualDeviceID {
			deviceKnown = true
			break
		}
	}
	
	// Check if IP address is in usual IPs
	ipKnown := false
	for _, usualIP := range userProfile.UsualIPAddresses {
		if transaction.IPAddress == usualIP {
			ipKnown = true
			break
		}
	}
	
	// Calculate score based on device and IP familiarity
	if deviceKnown && ipKnown {
		return 0.0 // Both known, low risk
	} else if deviceKnown || ipKnown {
		return 0.4 // One known, moderate risk
	} else {
		return 0.8 // Neither known, high risk
	}
}

// calculateFrequencyScore calculates a score based on transaction frequency
func (m *FraudDetectionModel) calculateFrequencyScore(transaction Transaction, history []Transaction) float64 {
	// If no history, use moderate score
	if len(history) < 2 {
		return 0.5
	}
	
	// Calculate average time between transactions
	var totalDuration time.Duration
	for i := 1; i < len(history); i++ {
		duration := history[i-1].Timestamp.Sub(history[i].Timestamp)
		totalDuration += duration
	}
	
	averageDuration := totalDuration / time.Duration(len(history)-1)
	
	// Calculate time since last transaction
	var timeSinceLast time.Duration
	if len(history) > 0 {
		timeSinceLast = transaction.Timestamp.Sub(history[0].Timestamp)
	}
	
	// If time since last is much shorter than average, it's suspicious
	if averageDuration > 0 && timeSinceLast < averageDuration/10 {
		return 0.9 // Very suspicious
	} else if timeSinceLast < averageDuration/2 {
		return 0.5 // Somewhat suspicious
	}
	
	return 0.1 // Not suspicious
}

// calculateVelocityScore calculates a score based on transaction velocity
func (m *FraudDetectionModel) calculateVelocityScore(transaction Transaction, history []Transaction) float64 {
	// Count transactions in the last hour
	var recentCount int
	oneHourAgo := transaction.Timestamp.Add(-1 * time.Hour)
	
	for _, tx := range history {
		if tx.Timestamp.After(oneHourAgo) {
			recentCount++
		}
	}
	
	// More than 5 transactions in an hour is suspicious
	if recentCount > 10 {
		return 1.0 // Very suspicious
	} else if recentCount > 5 {
		return 0.7 // Suspicious
	} else if recentCount > 3 {
		return 0.4 // Somewhat suspicious
	}
	
	return 0.0 // Not suspicious
}

// calculateDistance calculates the distance between two points in kilometers
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0 // Earth radius in kilometers
	
	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180.0
	lon1Rad := lon1 * math.Pi / 180.0
	lat2Rad := lat2 * math.Pi / 180.0
	lon2Rad := lon2 * math.Pi / 180.0
	
	// Haversine formula
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c
	
	return distance
}
