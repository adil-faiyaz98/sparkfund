package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all configuration for the email service
type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Kafka       KafkaConfig
	SMTP        SMTPConfig
	Jaeger      JaegerConfig
	Environment string
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
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	MaxIdle  int
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
	Endpoint string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
	}

	// Server config
	cfg.Server = ServerConfig{
		Port:         getEnvIntOrDefault("SERVER_PORT", 8080),
		ReadTimeout:  getEnvDurationOrDefault("SERVER_READ_TIMEOUT", 10*time.Second),
		WriteTimeout: getEnvDurationOrDefault("SERVER_WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:  getEnvDurationOrDefault("SERVER_IDLE_TIMEOUT", 120*time.Second),
	}

	// Database config
	cfg.Database = DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvIntOrDefault("DB_PORT", 5432),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:   getEnvOrDefault("DB_NAME", "email_service"),
		SSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
		MaxConns: getEnvIntOrDefault("DB_MAX_CONNS", 10),
		MaxIdle:  getEnvIntOrDefault("DB_MAX_IDLE", 5),
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

	// Jaeger config
	cfg.Jaeger = JaegerConfig{
		Endpoint: getEnvOrDefault("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
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
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
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
