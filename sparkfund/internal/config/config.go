package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	Port            int           `json:"port"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`

	// Database configuration
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
	DBSSLMode  string `json:"db_ssl_mode"`

	// Kafka configuration
	KafkaBrokers []string `json:"kafka_brokers"`
	KafkaTopic   string   `json:"kafka_topic"`

	// SMTP configuration
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUsername string `json:"smtp_username"`
	SMTPPassword string `json:"smtp_password"`
	SMTPFrom     string `json:"smtp_from"`

	// Jaeger configuration
	JaegerEndpoint string `json:"jaeger_endpoint"`
	JaegerService  string `json:"jaeger_service"`

	// Rate limiting
	RateLimit      int `json:"rate_limit"`
	RateLimitBurst int `json:"rate_limit_burst"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{}

	// Server configuration
	config.Port = getEnvInt("PORT", 8080)
	config.ShutdownTimeout = getEnvDuration("SHUTDOWN_TIMEOUT", 30*time.Second)

	// Database configuration
	config.DBHost = getEnv("DB_HOST", "localhost")
	config.DBPort = getEnvInt("DB_PORT", 5432)
	config.DBUser = getEnv("DB_USER", "postgres")
	config.DBPassword = getEnv("DB_PASSWORD", "postgres")
	config.DBName = getEnv("DB_NAME", "email_service")
	config.DBSSLMode = getEnv("DB_SSL_MODE", "disable")

	// Kafka configuration
	config.KafkaBrokers = getEnvSlice("KAFKA_BROKERS", []string{"localhost:9092"})
	config.KafkaTopic = getEnv("KAFKA_TOPIC", "email_requests")

	// SMTP configuration
	config.SMTPHost = getEnv("SMTP_HOST", "localhost")
	config.SMTPPort = getEnvInt("SMTP_PORT", 587)
	config.SMTPUsername = getEnv("SMTP_USERNAME", "")
	config.SMTPPassword = getEnv("SMTP_PASSWORD", "")
	config.SMTPFrom = getEnv("SMTP_FROM", "noreply@example.com")

	// Jaeger configuration
	config.JaegerEndpoint = getEnv("JAEGER_ENDPOINT", "http://localhost:14268/api/traces")
	config.JaegerService = getEnv("JAEGER_SERVICE", "email-service")

	// Rate limiting
	config.RateLimit = getEnvInt("RATE_LIMIT", 100)
	config.RateLimitBurst = getEnvInt("RATE_LIMIT_BURST", 200)

	return config, nil
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

// GetSMTPAddr returns the SMTP server address
func (c *Config) GetSMTPAddr() string {
	return fmt.Sprintf("%s:%d", c.SMTPHost, c.SMTPPort)
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
