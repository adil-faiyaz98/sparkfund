package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"github.com/sparkfund/investment-service/internal/config"
	"go.uber.org/zap"
)

// NewDB creates a new database connection.
func NewDB(cfg config.DatabaseConfig, logger *zap.Logger) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		logger.Error("Error connecting to database", zap.Error(err))
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		logger.Error("Error pinging database", zap.Error(err))
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	logger.Info("Connected to database")
	return db, nil
}

// CloseDB closes the database connection
func CloseDB(db *sqlx.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}