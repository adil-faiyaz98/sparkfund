package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	Environment string
	Port        string
	Database    DatabaseConfig
	Service     string
}

// DatabaseConfig holds the database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig(service string) (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
		Port:        getEnvOrDefault("PORT", "8080"),
		Service:     service,
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
			DBName:   getEnvOrDefault("DB_NAME", service),
			SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		},
	}

	return config, nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
