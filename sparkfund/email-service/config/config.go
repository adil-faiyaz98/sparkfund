package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

// Config holds the service configuration
type Config struct {
	Port            int            `validate:"required,gt=0,lte=65535"`
	ShutdownTimeout time.Duration  `validate:"required"`
	Database        DatabaseConfig `validate:"required"`
	Jaeger          JaegerConfig
}

type JaegerConfig struct {
	Endpoint    string `validate:"required"`
	ServiceName string `validate:"required"`
}

// DatabaseConfig holds the database configuration
type DatabaseConfig struct {
	Host     string `validate:"required"`
	Port     int    `validate:"required,gt=0,lte=65535"`
	User     string `validate:"required"`
	Password string `validate:"required"` // Ideally, retrieve from secrets manager
	Name     string `validate:"required"`
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port:            getEnvOrDefaultInt("PORT", 8080),
		ShutdownTimeout: getEnvOrDefaultDuration("SHUTDOWN_TIMEOUT", 5*time.Second),
		Database: DatabaseConfig{
			Host: getEnvOrDefault("DB_HOST", "localhost"),
			Port: getEnvOrDefaultInt("DB_PORT", 5432),
			User: getEnvOrDefault("DB_USER", "postgres"),
			// Password: getEnvOrDefault("DB_PASS", "default_password"), // DO NOT DO THIS IN PRODUCTION
			Name: getEnvOrDefault("DB_NAME", "postgres"),
		},
		Jaeger: JaegerConfig{
			Endpoint:    getEnvOrDefault("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
			ServiceName: getEnvOrDefault("JAEGER_SERVICE_NAME", "email-service"),
		},
	}

	// Retrieve password from secrets manager (replace with your actual implementation)
	password, err := getSecret("DB_PASS")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve database password from secrets manager: %w", err)
	}
	cfg.Database.Password = password

	// Validate configuration
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		// Log the error or return a default value
		fmt.Printf("Error converting %s to int: %v, using default value %d\n", key, err, defaultValue)
		return defaultValue
	}
	return intValue
}

func getEnvOrDefaultDuration(key string, defaultValue time.Duration) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	durationValue, err := time.ParseDuration(value)
	if err != nil {
		// Log the error or return a default value
		fmt.Printf("Error converting %s to duration: %v, using default value %s\n", key, err, defaultValue)
		return defaultValue
	}
	return durationValue
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
