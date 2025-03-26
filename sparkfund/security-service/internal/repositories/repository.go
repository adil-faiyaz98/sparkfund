package repositories

import (
	"context"
	"database/sql"

	"go.uber.org/zap"
)

// Repository defines the interface for database operations
type Repository interface {
	SaveToken(ctx context.Context, userID string, token string, refreshToken string, expiresAt int64) error
	GetTokenByRefresh(ctx context.Context, refreshToken string) (string, string, int64, error)
	RevokeToken(ctx context.Context, token string) error
}

// postgresRepository implements the Repository interface with PostgreSQL
type postgresRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewRepository creates a new repository with PostgreSQL
func NewRepository(db *sql.DB, logger *zap.Logger) Repository {
	return &postgresRepository{
		db:     db,
		logger: logger,
	}
}

// SaveToken stores token information in the database
func (r *postgresRepository) SaveToken(ctx context.Context, userID string, token string, refreshToken string, expiresAt int64) error {
	query := `
		INSERT INTO tokens (user_id, token, refresh_token, expires_at, created_at)
		VALUES ($1, $2, $3, to_timestamp($4), NOW())
	`

	_, err := r.db.ExecContext(ctx, query, userID, token, refreshToken, expiresAt)
	if err != nil {
		r.logger.Error("Failed to save token",
			zap.String("user_id", userID),
			zap.Error(err))
		return err
	}

	return nil
}

// GetTokenByRefresh retrieves token information by refresh token
func (r *postgresRepository) GetTokenByRefresh(ctx context.Context, refreshToken string) (string, string, int64, error) {
	query := `
		SELECT user_id, token, extract(epoch from expires_at)
		FROM tokens
		WHERE refresh_token = $1 AND revoked = false AND expires_at > NOW()
	`

	var userID, token string
	var expiresAt int64

	err := r.db.QueryRowContext(ctx, query, refreshToken).Scan(&userID, &token, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn("Refresh token not found or expired",
				zap.String("refresh_token", refreshToken))
		} else {
			r.logger.Error("Failed to get token by refresh",
				zap.String("refresh_token", refreshToken),
				zap.Error(err))
		}
		return "", "", 0, err
	}

	return userID, token, expiresAt, nil
}

// RevokeToken marks a token as revoked in the database
func (r *postgresRepository) RevokeToken(ctx context.Context, token string) error {
	query := `
		UPDATE tokens
		SET revoked = true
		WHERE token = $1
	`

	result, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		r.logger.Error("Failed to revoke token", zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		r.logger.Warn("Token not found for revocation")
	}

	return nil
}
