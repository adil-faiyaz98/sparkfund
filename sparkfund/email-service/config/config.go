package config

import (
	"os"
)

// Config holds the service configuration
type Config struct {
	Port     string
	Database DatabaseConfig
}

// DatabaseConfig holds the database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (*Config, error) {
	return &Config{
		Port: getEnvOrDefault("PORT", "8080"),
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASS", "postgres"),
			Name:     getEnvOrDefault("DB_NAME", "postgres"), // Default to postgres
		},
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
