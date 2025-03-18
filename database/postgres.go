package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConfig creates a new database configuration from environment variables
func NewConfig() *Config {
	return &Config{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvAsIntOrDefault("DB_PORT", 5432),
		User:     getEnvOrDefault("DB_USER", "moneypulse"),
		Password: getEnvOrDefault("DB_PASSWORD", "moneypulse"),
		DBName:   getEnvOrDefault("DB_NAME", "moneypulse"),
		SSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
	}
}

// Connect establishes a connection to PostgreSQL
func Connect(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Hour)

	// Enable logging in development
	if getEnvOrDefault("ENV", "development") == "development" {
		db.LogMode(true)
	}

	return db, nil
}

// ConnectFromURL connects to PostgreSQL using a connection URL
func ConnectFromURL(url string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Hour)

	// Enable logging in development
	if getEnvOrDefault("ENV", "development") == "development" {
		db.LogMode(true)
	}

	return db, nil
}

// Close closes the database connection
func Close(db *gorm.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// Helper functions for environment variables
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
