package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds application configuration
type Config struct {
	// Server settings
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// Database settings
	DatabaseURL string

	// Security settings
	JWTSecret         string
	PasswordHashCost  int
	SessionExpiry     time.Duration
	MaxLoginAttempts  int
	LockoutDuration   time.Duration
	RequireMFA        bool
	AllowedCountries  []string
	MinPasswordLength int

	// Email settings
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		// Server settings
		Port:         getEnvOrDefault("PORT", "8080"),
		ReadTimeout:  getDurationEnvOrDefault("READ_TIMEOUT", 5*time.Second),
		WriteTimeout: getDurationEnvOrDefault("WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:  getDurationEnvOrDefault("IDLE_TIMEOUT", 120*time.Second),

		// Database settings
		DatabaseURL: getEnvOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/sparkfund?sslmode=disable"),

		// Security settings
		JWTSecret:         getEnvOrDefault("JWT_SECRET", "your-secret-key"),
		PasswordHashCost:  getIntEnvOrDefault("PASSWORD_HASH_COST", 10),
		SessionExpiry:     getDurationEnvOrDefault("SESSION_EXPIRY", 24*time.Hour),
		MaxLoginAttempts:  getIntEnvOrDefault("MAX_LOGIN_ATTEMPTS", 5),
		LockoutDuration:   getDurationEnvOrDefault("LOCKOUT_DURATION", 30*time.Minute),
		RequireMFA:        getBoolEnvOrDefault("REQUIRE_MFA", false),
		AllowedCountries:  getStringSliceEnvOrDefault("ALLOWED_COUNTRIES", []string{"US", "GB", "CA"}),
		MinPasswordLength: getIntEnvOrDefault("MIN_PASSWORD_LENGTH", 8),

		// Email settings
		SMTPHost:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getIntEnvOrDefault("SMTP_PORT", 587),
		SMTPUsername: getEnvOrDefault("SMTP_USERNAME", ""),
		SMTPPassword: getEnvOrDefault("SMTP_PASSWORD", ""),
		FromEmail:    getEnvOrDefault("FROM_EMAIL", "noreply@sparkfund.com"),
	}
}

// Helper functions for environment variables

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnvOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnvOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getBoolEnvOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getStringSliceEnvOrDefault(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
