package recommendation

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
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

// GetUserProfile retrieves a user profile by ID
func (r *PostgresRepository) GetUserProfile(ctx context.Context, userID string) (*UserProfile, error) {
	query := `
		SELECT * FROM user_profiles
		WHERE user_id = $1
	`

	var profile UserProfile
	err := r.db.GetContext(ctx, &profile, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return &profile, nil
}

// SaveUserProfile saves a user profile
func (r *PostgresRepository) SaveUserProfile(ctx context.Context, profile *UserProfile) error {
	query := `
		INSERT INTO user_profiles (
			user_id, risk_tolerance, investment_goals, time_horizon, age, income,
			liquid_net_worth, investment_amount, preferred_sectors, excluded_sectors, created_at, updated_at
		) VALUES (
			:user_id, :risk_tolerance, :investment_goals, :time_horizon, :age, :income,
			:liquid_net_worth, :investment_amount, :preferred_sectors, :excluded_sectors, :created_at, :updated_at
		)
		ON CONFLICT (user_id) DO UPDATE SET
			risk_tolerance = :risk_tolerance,
			investment_goals = :investment_goals,
			time_horizon = :time_horizon,
			age = :age,
			income = :income,
			liquid_net_worth = :liquid_net_worth,
			investment_amount = :investment_amount,
			preferred_sectors = :preferred_sectors,
			excluded_sectors = :excluded_sectors,
			updated_at = :updated_at
	`

	_, err := r.db.NamedExecContext(ctx, query, profile)
	if err != nil {
		return fmt.Errorf("failed to save user profile: %w", err)
	}

	return nil
}

// GetInvestmentAsset retrieves an investment asset by ID
func (r *PostgresRepository) GetInvestmentAsset(ctx context.Context, assetID string) (*InvestmentAsset, error) {
	query := `
		SELECT * FROM investment_assets
		WHERE id = $1
	`

	var asset InvestmentAsset
	err := r.db.GetContext(ctx, &asset, query, assetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInsufficientData
		}
		return nil, fmt.Errorf("failed to get investment asset: %w", err)
	}

	return &asset, nil
}

// GetInvestmentAssets retrieves investment assets based on a filter
func (r *PostgresRepository) GetInvestmentAssets(ctx context.Context, filter map[string]interface{}, limit int) ([]InvestmentAsset, error) {
	query := `SELECT * FROM investment_assets`
	
	// Build WHERE clause if filter is provided
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter != nil {
		for key, value := range filter {
			conditions = append(conditions, fmt.Sprintf("%s = $%d", key, argIndex))
			args = append(args, value)
			argIndex++
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	// Add limit
	query += fmt.Sprintf(" LIMIT %d", limit)

	var assets []InvestmentAsset
	err := r.db.SelectContext(ctx, &assets, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get investment assets: %w", err)
	}

	return assets, nil
}

// SaveInvestmentAsset saves an investment asset
func (r *PostgresRepository) SaveInvestmentAsset(ctx context.Context, asset *InvestmentAsset) error {
	query := `
		INSERT INTO investment_assets (
			id, symbol, name, asset_type, sector, risk_level, historical_returns,
			volatility, current_price, one_year_target, five_year_target,
			dividend_yield, market_cap, esg_score, features, created_at, updated_at
		) VALUES (
			:id, :symbol, :name, :asset_type, :sector, :risk_level, :historical_returns,
			:volatility, :current_price, :one_year_target, :five_year_target,
			:dividend_yield, :market_cap, :esg_score, :features, :created_at, :updated_at
		)
		ON CONFLICT (id) DO UPDATE SET
			symbol = :symbol,
			name = :name,
			asset_type = :asset_type,
			sector = :sector,
			risk_level = :risk_level,
			historical_returns = :historical_returns,
			volatility = :volatility,
			current_price = :current_price,
			one_year_target = :one_year_target,
			five_year_target = :five_year_target,
			dividend_yield = :dividend_yield,
			market_cap = :market_cap,
			esg_score = :esg_score,
			features = :features,
			updated_at = :updated_at
	`

	_, err := r.db.NamedExecContext(ctx, query, asset)
	if err != nil {
		return fmt.Errorf("failed to save investment asset: %w", err)
	}

	return nil
}

// GetLatestMarketData retrieves the latest market data
func (r *PostgresRepository) GetLatestMarketData(ctx context.Context) (*MarketData, error) {
	query := `
		SELECT * FROM market_data
		ORDER BY date DESC
		LIMIT 1
	`

	var marketData MarketData
	err := r.db.GetContext(ctx, &marketData, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInsufficientData
		}
		return nil, fmt.Errorf("failed to get latest market data: %w", err)
	}

	return &marketData, nil
}

// GetHistoricalMarketData retrieves historical market data within a time range
func (r *PostgresRepository) GetHistoricalMarketData(ctx context.Context, startTime, endTime time.Time) ([]MarketData, error) {
	query := `
		SELECT * FROM market_data
		WHERE date >= $1 AND date <= $2
		ORDER BY date
	`

	var marketData []MarketData
	err := r.db.SelectContext(ctx, &marketData, query, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical market data: %w", err)
	}

	return marketData, nil
}

// SaveMarketData saves market data
func (r *PostgresRepository) SaveMarketData(ctx context.Context, data *MarketData) error {
	query := `
		INSERT INTO market_data (
			date, market_trends, economic_indicators, sector_performance
		) VALUES (
			:date, :market_trends, :economic_indicators, :sector_performance
		)
		ON CONFLICT (date) DO UPDATE SET
			market_trends = :market_trends,
			economic_indicators = :economic_indicators,
			sector_performance = :sector_performance
	`

	_, err := r.db.NamedExecContext(ctx, query, data)
	if err != nil {
		return fmt.Errorf("failed to save market data: %w", err)
	}

	return nil
}

// GetUserInteractions retrieves user interactions for a user
func (r *PostgresRepository) GetUserInteractions(ctx context.Context, userID string, limit int) ([]UserInteraction, error) {
	query := `
		SELECT * FROM user_interactions
		WHERE user_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	var interactions []UserInteraction
	err := r.db.SelectContext(ctx, &interactions, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user interactions: %w", err)
	}

	return interactions, nil
}

// SaveUserInteraction saves a user interaction
func (r *PostgresRepository) SaveUserInteraction(ctx context.Context, interaction *UserInteraction) error {
	query := `
		INSERT INTO user_interactions (
			id, user_id, asset_id, interact_type, rating, timestamp, amount, duration, feedback
		) VALUES (
			:id, :user_id, :asset_id, :interact_type, :rating, :timestamp, :amount, :duration, :feedback
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, interaction)
	if err != nil {
		return fmt.Errorf("failed to save user interaction: %w", err)
	}

	return nil
}

// SavePortfolioRecommendation saves a portfolio recommendation
func (r *PostgresRepository) SavePortfolioRecommendation(ctx context.Context, recommendation *PortfolioRecommendation) error {
	// Serialize recommended assets to JSON
	recommendedAssetsJSON, err := json.Marshal(recommendation.RecommendedAssets)
	if err != nil {
		return fmt.Errorf("failed to marshal recommended assets: %w", err)
	}

	// Create a map for the query
	params := map[string]interface{}{
		"id":                    recommendation.ID,
		"user_id":               recommendation.UserID,
		"recommended_assets":    recommendedAssetsJSON,
		"total_expected_return": recommendation.TotalExpectedReturn,
		"portfolio_risk_level":  recommendation.PortfolioRiskLevel,
		"diversification_score": recommendation.DiversificationScore,
		"rebalancing_frequency": recommendation.RebalancingFrequency,
		"time_horizon":          recommendation.TimeHorizon,
		"created_at":            recommendation.CreatedAt,
	}

	query := `
		INSERT INTO portfolio_recommendations (
			id, user_id, recommended_assets, total_expected_return,
			portfolio_risk_level, diversification_score, rebalancing_frequency,
			time_horizon, created_at
		) VALUES (
			:id, :user_id, :recommended_assets, :total_expected_return,
			:portfolio_risk_level, :diversification_score, :rebalancing_frequency,
			:time_horizon, :created_at
		)
	`

	_, err = r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to save portfolio recommendation: %w", err)
	}

	return nil
}

// GetUserRecommendations retrieves portfolio recommendations for a user
func (r *PostgresRepository) GetUserRecommendations(ctx context.Context, userID string, limit int) ([]PortfolioRecommendation, error) {
	query := `
		SELECT * FROM portfolio_recommendations
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryxContext(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user recommendations: %w", err)
	}
	defer rows.Close()

	var recommendations []PortfolioRecommendation
	for rows.Next() {
		var rec struct {
			ID                   string          `db:"id"`
			UserID               string          `db:"user_id"`
			RecommendedAssetsJSON json.RawMessage `db:"recommended_assets"`
			TotalExpectedReturn  float64         `db:"total_expected_return"`
			PortfolioRiskLevel   float64         `db:"portfolio_risk_level"`
			DiversificationScore float64         `db:"diversification_score"`
			RebalancingFrequency string          `db:"rebalancing_frequency"`
			TimeHorizon          int             `db:"time_horizon"`
			CreatedAt            time.Time       `db:"created_at"`
		}

		err := rows.StructScan(&rec)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recommendation: %w", err)
		}

		var recommendedAssets []RecommendedAsset
		err = json.Unmarshal(rec.RecommendedAssetsJSON, &recommendedAssets)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal recommended assets: %w", err)
		}

		recommendations = append(recommendations, PortfolioRecommendation{
			ID:                   rec.ID,
			UserID:               rec.UserID,
			RecommendedAssets:    recommendedAssets,
			TotalExpectedReturn:  rec.TotalExpectedReturn,
			PortfolioRiskLevel:   rec.PortfolioRiskLevel,
			DiversificationScore: rec.DiversificationScore,
			RebalancingFrequency: rec.RebalancingFrequency,
			TimeHorizon:          rec.TimeHorizon,
			CreatedAt:            rec.CreatedAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating recommendations: %w", err)
	}

	return recommendations, nil
}

// SaveModelMetrics saves model performance metrics
func (r *PostgresRepository) SaveModelMetrics(ctx context.Context, metrics *ModelPerformanceMetrics) error {
	query := `
		INSERT INTO model_performance_metrics (
			model_version, accuracy, precision, recall, f1_score,
			mean_absolute_error, root_mean_square_error, user_satisfaction,
			last_evaluated_at
		) VALUES (
			:model_version, :accuracy, :precision, :recall, :f1_score,
			:mean_absolute_error, :root_mean_square_error, :user_satisfaction,
			:last_evaluated_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, metrics)
	if err != nil {
		return fmt.Errorf("failed to save model metrics: %w", err)
	}

	return nil
}

// GetLatestModelMetrics retrieves the latest model performance metrics
func (r *PostgresRepository) GetLatestModelMetrics(ctx context.Context) (*ModelPerformanceMetrics, error) {
	query := `
		SELECT * FROM model_performance_metrics
		ORDER BY last_evaluated_at DESC
		LIMIT 1
	`

	var metrics ModelPerformanceMetrics
	err := r.db.GetContext(ctx, &metrics, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrModelNotReady
		}
		return nil, fmt.Errorf("failed to get latest model metrics: %w", err)
	}

	return &metrics, nil
}

// GetSimilarUsers retrieves users similar to the given user
func (r *PostgresRepository) GetSimilarUsers(ctx context.Context, userID string, limit int) ([]string, error) {
	// This is a simplified implementation that finds users with similar risk tolerance
	query := `
		WITH user_profile AS (
			SELECT risk_tolerance, time_horizon FROM user_profiles WHERE user_id = $1
		)
		SELECT up.user_id
		FROM user_profiles up, user_profile ref
		WHERE up.user_id != $1
		ORDER BY ABS(up.risk_tolerance - ref.risk_tolerance) + ABS(up.time_horizon - ref.time_horizon) / 10.0
		LIMIT $2
	`

	var userIDs []string
	err := r.db.SelectContext(ctx, &userIDs, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get similar users: %w", err)
	}

	return userIDs, nil
}

// GetPopularAssets retrieves popular investment assets
func (r *PostgresRepository) GetPopularAssets(ctx context.Context, limit int) ([]InvestmentAsset, error) {
	query := `
		SELECT a.*
		FROM investment_assets a
		JOIN (
			SELECT asset_id, COUNT(*) as interaction_count
			FROM user_interactions
			WHERE interact_type IN ('PURCHASE', 'SAVE')
			GROUP BY asset_id
			ORDER BY interaction_count DESC
			LIMIT $1
		) popular ON a.id = popular.asset_id
	`

	var assets []InvestmentAsset
	err := r.db.SelectContext(ctx, &assets, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular assets: %w", err)
	}

	return assets, nil
}

// GetUserPortfolio retrieves a user's current portfolio
func (r *PostgresRepository) GetUserPortfolio(ctx context.Context, userID string) ([]InvestmentAsset, error) {
	query := `
		SELECT a.*
		FROM investment_assets a
		JOIN user_portfolio p ON a.id = p.asset_id
		WHERE p.user_id = $1
	`

	var assets []InvestmentAsset
	err := r.db.SelectContext(ctx, &assets, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user portfolio: %w", err)
	}

	return assets, nil
}

// GetAllUserProfiles retrieves all user profiles
func (r *PostgresRepository) GetAllUserProfiles(ctx context.Context) ([]UserProfile, error) {
	query := `SELECT * FROM user_profiles`

	var profiles []UserProfile
	err := r.db.SelectContext(ctx, &profiles, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all user profiles: %w", err)
	}

	return profiles, nil
}

// GetRecentUserInteractions retrieves recent user interactions
func (r *PostgresRepository) GetRecentUserInteractions(ctx context.Context, limit int) ([]UserInteraction, error) {
	query := `
		SELECT * FROM user_interactions
		ORDER BY timestamp DESC
		LIMIT $1
	`

	var interactions []UserInteraction
	err := r.db.SelectContext(ctx, &interactions, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent user interactions: %w", err)
	}

	return interactions, nil
}
