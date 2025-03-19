package database

import (
	"fmt"
	"time"

	"github.com/adil-faiyaz98/sparkfund/email-service/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewDB creates a new database connection with retries
func NewDB(cfg *config.DatabaseConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	var db *sqlx.DB
	var err error

	// Retry connection with exponential backoff
	for i := 0; i < 5; i++ {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			break
		}

		// Wait before retrying
		time.Sleep(time.Second * time.Duration(1<<uint(i)))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after retries: %v", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetMaxIdleConns(cfg.MaxIdle)
	db.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

// CloseDB closes the database connection
func CloseDB(db *sqlx.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
