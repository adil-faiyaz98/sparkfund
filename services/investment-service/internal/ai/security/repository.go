package security

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sparkfund/services/investment-service/internal/ai/anomaly"
	"github.com/sparkfund/services/investment-service/internal/ai/fraud"
)

// PostgresRepository implements the Repository interface using PostgreSQL
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// GetUserSecurityProfile retrieves a user's security profile
func (r *PostgresRepository) GetUserSecurityProfile(ctx context.Context, userID string) (*fraud.UserProfile, error) {
	query := `
		SELECT * FROM user_security_profiles
		WHERE user_id = $1
	`

	var profile fraud.UserProfile
	err := r.db.GetContext(ctx, &profile, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user security profile not found")
		}
		return nil, fmt.Errorf("failed to get user security profile: %w", err)
	}

	return &profile, nil
}

// SaveUserSecurityProfile saves a user's security profile
func (r *PostgresRepository) SaveUserSecurityProfile(ctx context.Context, profile *fraud.UserProfile) error {
	// Serialize arrays to JSON
	usualIPAddressesJSON, err := json.Marshal(profile.UsualIPAddresses)
	if err != nil {
		return fmt.Errorf("failed to marshal usual IP addresses: %w", err)
	}

	usualDeviceIDsJSON, err := json.Marshal(profile.UsualDeviceIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal usual device IDs: %w", err)
	}

	usualLocationsJSON, err := json.Marshal(profile.UsualLocations)
	if err != nil {
		return fmt.Errorf("failed to marshal usual locations: %w", err)
	}

	// Create a map for the query
	params := map[string]interface{}{
		"user_id":                profile.UserID,
		"created_at":             profile.CreatedAt,
		"last_login":             profile.LastLogin,
		"usual_ip_addresses":     usualIPAddressesJSON,
		"usual_device_ids":       usualDeviceIDsJSON,
		"usual_locations":        usualLocationsJSON,
		"average_transaction_amount": profile.AverageTransactionAmount,
		"transaction_frequency":  profile.TransactionFrequency,
		"risk_score":             profile.RiskScore,
	}

	query := `
		INSERT INTO user_security_profiles (
			user_id, created_at, last_login, usual_ip_addresses, usual_device_ids,
			usual_locations, average_transaction_amount, transaction_frequency, risk_score
		) VALUES (
			:user_id, :created_at, :last_login, :usual_ip_addresses, :usual_device_ids,
			:usual_locations, :average_transaction_amount, :transaction_frequency, :risk_score
		)
		ON CONFLICT (user_id) DO UPDATE SET
			last_login = :last_login,
			usual_ip_addresses = :usual_ip_addresses,
			usual_device_ids = :usual_device_ids,
			usual_locations = :usual_locations,
			average_transaction_amount = :average_transaction_amount,
			transaction_frequency = :transaction_frequency,
			risk_score = :risk_score
	`

	_, err = r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to save user security profile: %w", err)
	}

	return nil
}

// GetUserTransactions retrieves a user's transaction history
func (r *PostgresRepository) GetUserTransactions(ctx context.Context, userID string, limit int) ([]fraud.Transaction, error) {
	query := `
		SELECT * FROM transactions
		WHERE user_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	var transactions []fraud.Transaction
	err := r.db.SelectContext(ctx, &transactions, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user transactions: %w", err)
	}

	return transactions, nil
}

// SaveTransaction saves a transaction
func (r *PostgresRepository) SaveTransaction(ctx context.Context, transaction *fraud.Transaction) error {
	// Serialize location to JSON
	locationJSON, err := json.Marshal(transaction.Location)
	if err != nil {
		return fmt.Errorf("failed to marshal location: %w", err)
	}

	// Create a map for the query
	params := map[string]interface{}{
		"id":               transaction.ID,
		"user_id":          transaction.UserID,
		"amount":           transaction.Amount,
		"currency":         transaction.Currency,
		"transaction_type": transaction.TransactionType,
		"asset_id":         transaction.AssetID,
		"timestamp":        transaction.Timestamp,
		"ip_address":       transaction.IPAddress,
		"device_id":        transaction.DeviceID,
		"location":         locationJSON,
		"user_agent":       transaction.UserAgent,
	}

	query := `
		INSERT INTO transactions (
			id, user_id, amount, currency, transaction_type, asset_id,
			timestamp, ip_address, device_id, location, user_agent
		) VALUES (
			:id, :user_id, :amount, :currency, :transaction_type, :asset_id,
			:timestamp, :ip_address, :device_id, :location, :user_agent
		)
	`

	_, err = r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to save transaction: %w", err)
	}

	return nil
}

// SaveFraudDetectionResult saves a fraud detection result
func (r *PostgresRepository) SaveFraudDetectionResult(ctx context.Context, result *fraud.FraudDetectionResult) error {
	// Serialize fraud indicators to JSON
	fraudIndicatorsJSON, err := json.Marshal(result.FraudIndicators)
	if err != nil {
		return fmt.Errorf("failed to marshal fraud indicators: %w", err)
	}

	// Create a map for the query
	params := map[string]interface{}{
		"transaction_id":   result.TransactionID,
		"user_id":          result.UserID,
		"timestamp":        result.Timestamp,
		"fraud_score":      result.FraudScore,
		"fraud_level":      result.FraudLevel,
		"fraud_indicators": fraudIndicatorsJSON,
		"action":           result.Action,
	}

	query := `
		INSERT INTO fraud_detection_results (
			transaction_id, user_id, timestamp, fraud_score,
			fraud_level, fraud_indicators, action
		) VALUES (
			:transaction_id, :user_id, :timestamp, :fraud_score,
			:fraud_level, :fraud_indicators, :action
		)
	`

	_, err = r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to save fraud detection result: %w", err)
	}

	return nil
}

// SaveAnomalyDetectionResult saves an anomaly detection result
func (r *PostgresRepository) SaveAnomalyDetectionResult(ctx context.Context, result *anomaly.AnomalyDetectionResult) error {
	// Serialize anomaly indicators to JSON
	anomalyIndicatorsJSON, err := json.Marshal(result.AnomalyIndicators)
	if err != nil {
		return fmt.Errorf("failed to marshal anomaly indicators: %w", err)
	}

	// Create a map for the query
	params := map[string]interface{}{
		"transaction_id":     result.TransactionID,
		"user_id":            result.UserID,
		"timestamp":          result.Timestamp,
		"anomaly_score":      result.AnomalyScore,
		"anomaly_level":      result.AnomalyLevel,
		"anomaly_indicators": anomalyIndicatorsJSON,
		"action":             result.Action,
	}

	query := `
		INSERT INTO anomaly_detection_results (
			transaction_id, user_id, timestamp, anomaly_score,
			anomaly_level, anomaly_indicators, action
		) VALUES (
			:transaction_id, :user_id, :timestamp, :anomaly_score,
			:anomaly_level, :anomaly_indicators, :action
		)
	`

	_, err = r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to save anomaly detection result: %w", err)
	}

	return nil
}

// GetLatestMarketData retrieves the latest market data
func (r *PostgresRepository) GetLatestMarketData(ctx context.Context) (*anomaly.MarketData, error) {
	query := `
		SELECT * FROM market_data
		ORDER BY timestamp DESC
		LIMIT 1
	`

	var marketDataDB struct {
		Timestamp          time.Time       `db:"timestamp"`
		MarketTrendsJSON   json.RawMessage `db:"market_trends"`
		EconomicIndicatorsJSON json.RawMessage `db:"economic_indicators"`
		SectorPerformanceJSON  json.RawMessage `db:"sector_performance"`
	}

	err := r.db.GetContext(ctx, &marketDataDB, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("market data not found")
		}
		return nil, fmt.Errorf("failed to get market data: %w", err)
	}

	// Deserialize JSON fields
	var marketTrends map[string]float64
	err = json.Unmarshal(marketDataDB.MarketTrendsJSON, &marketTrends)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal market trends: %w", err)
	}

	var economicIndicators map[string]float64
	err = json.Unmarshal(marketDataDB.EconomicIndicatorsJSON, &economicIndicators)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal economic indicators: %w", err)
	}

	var sectorPerformance map[string]float64
	err = json.Unmarshal(marketDataDB.SectorPerformanceJSON, &sectorPerformance)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal sector performance: %w", err)
	}

	// Create market data
	marketData := &anomaly.MarketData{
		Timestamp:          marketDataDB.Timestamp,
		MarketTrends:       marketTrends,
		EconomicIndicators: economicIndicators,
		SectorPerformance:  sectorPerformance,
	}

	return marketData, nil
}
