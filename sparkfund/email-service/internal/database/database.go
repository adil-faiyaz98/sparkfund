package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"github.com/sparkfund/email-service/internal/config"
	"go.uber.org/zap" // Import zap
)

// NewDB creates a new database connection.
func NewDB(cfg config.DatabaseConfig, logger *zap.Logger) (*sqlx.DB, error) {
	// Retrieve password from secrets manager (replace with your actual implementation)
	password, err := getSecret("DB_PASSWORD")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve database password: %w", err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", // Enforce SSL
		cfg.Host, cfg.Port, cfg.User, password, cfg.DBName)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Configure connection pooling
	db.SetMaxOpenConns(10)           // Adjust as needed
	db.SetMaxIdleConns(5)            // Adjust as needed
	db.SetConnMaxLifetime(time.Hour) // Adjust as needed

	// Test the connection with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	logger.Info("Connected to database") // Use zap logger
	return db, nil
}

// CloseDB closes the database connection
func CloseDB(db *sqlx.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// Replace with your actual secrets manager implementation
func getSecret(key string) (string, error) {
	// Example using environment variable as a fallback (NOT RECOMMENDED FOR PRODUCTION)
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("secret %s not found", key)
	}
	return value, nil
}
