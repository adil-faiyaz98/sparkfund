package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxRetries     int           `json:"max_retries"`
	InitialBackoff time.Duration `json:"initial_backoff"`
	MaxBackoff     time.Duration `json:"max_backoff"`
}

// Config holds all configuration for the email service
type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Kafka       KafkaConfig
	SMTP        SMTPConfig
	Jaeger      JaegerConfig
	Retry       RetryConfig
	Environment string
	Port             int
	ShutdownTimeout  time.Duration
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds PostgreSQL configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Brokers []string
	Topic   string
}

// SMTPConfig holds SMTP server configuration
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// JaegerConfig holds Jaeger configuration
type JaegerConfig struct {
	Endpoint    string
	ServiceName string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &Config{
		Port:            getEnvIntOrDefault("SERVER_PORT", 8080),
		ShutdownTimeout: getEnvDurationOrDefault("SHUTDOWN_TIMEOUT", 30*time.Second),
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
			DBName:   getEnvOrDefault("DB_NAME", "email_service"),
		},
		Jaeger: JaegerConfig{
			Endpoint:    getEnvOrDefault("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
			ServiceName: getEnvOrDefault("JAEGER_SERVICE_NAME", "email-service"),
		},
	}

	// Server config
	cfg.Server = ServerConfig{
		Port:         cfg.Port,
		ReadTimeout:  getEnvDurationOrDefault("SERVER_READ_TIMEOUT", 10*time.Second),
		WriteTimeout: getEnvDurationOrDefault("SERVER_WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:  getEnvDurationOrDefault("SERVER_IDLE_TIMEOUT", 120*time.Second),
	}

	// Kafka config
	cfg.Kafka = KafkaConfig{
		Brokers: getEnvSliceOrDefault("KAFKA_BROKERS", []string{"localhost:9092"}),
		Topic:   getEnvOrDefault("KAFKA_TOPIC", "email-queue"),
	}

	// SMTP config
	cfg.SMTP = SMTPConfig{
		Host:     getEnvOrDefault("SMTP_HOST", "localhost"),
		Port:     getEnvIntOrDefault("SMTP_PORT", 587),
		Username: getEnvOrDefault("SMTP_USERNAME", ""),
		Password: getEnvOrDefault("SMTP_PASSWORD", ""),
		From:     getEnvOrDefault("SMTP_FROM", "noreply@sparkfund.com"),
	}

	// Retry config
	cfg.Retry = RetryConfig{
		MaxRetries:     getEnvIntOrDefault("RETRY_MAX_RETRIES", 3),
		InitialBackoff: getEnvDurationOrDefault("RETRY_INITIAL_BACKOFF", 1*time.Second),
		MaxBackoff:     getEnvDurationOrDefault("RETRY_MAX_BACKOFF", 30*time.Second),
	}

	return cfg, nil
}

// Helper functions for environment variables
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvSliceOrDefault(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.DBName)
}
