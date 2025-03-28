package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	// Server configuration
	Server struct {
		Port         int
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}

	// Database configuration
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}

	// Redis configuration
	Redis struct {
		Host     string
		Port     int
		Password string
		DB       int
	}

	// JWT configuration
	JWT struct {
		SecretKey     string
		ExpirationTime time.Duration
	}

	// AWS configuration
	AWS struct {
		Region          string
		AccessKeyID     string
		SecretAccessKey string
		BucketName      string
	}

	// Logging configuration
	Logging struct {
		Level string
	}

	// Metrics configuration
	Metrics struct {
		Enabled bool
		Port    int
	}
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// It's okay if .env file doesn't exist
	}

	cfg := &Config{}

	// Server configuration
	cfg.Server.Port = getEnvInt("SERVER_PORT", 8080)
	cfg.Server.ReadTimeout = getEnvDuration("SERVER_READ_TIMEOUT", 5*time.Second)
	cfg.Server.WriteTimeout = getEnvDuration("SERVER_WRITE_TIMEOUT", 10*time.Second)
	cfg.Server.IdleTimeout = getEnvDuration("SERVER_IDLE_TIMEOUT", 120*time.Second)

	// Database configuration
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnvInt("DB_PORT", 5432)
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "postgres")
	cfg.Database.DBName = getEnv("DB_NAME", "sparkfund")
	cfg.Database.SSLMode = getEnv("DB_SSL_MODE", "disable")

	// Redis configuration
	cfg.Redis.Host = getEnv("REDIS_HOST", "localhost")
	cfg.Redis.Port = getEnvInt("REDIS_PORT", 6379)
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")
	cfg.Redis.DB = getEnvInt("REDIS_DB", 0)

	// JWT configuration
	cfg.JWT.SecretKey = getEnv("JWT_SECRET_KEY", "your-secret-key")
	cfg.JWT.ExpirationTime = getEnvDuration("JWT_EXPIRATION_TIME", 24*time.Hour)

	// AWS configuration
	cfg.AWS.Region = getEnv("AWS_REGION", "us-east-1")
	cfg.AWS.AccessKeyID = getEnv("AWS_ACCESS_KEY_ID", "")
	cfg.AWS.SecretAccessKey = getEnv("AWS_SECRET_ACCESS_KEY", "")
	cfg.AWS.BucketName = getEnv("AWS_BUCKET_NAME", "sparkfund")

	// Logging configuration
	cfg.Logging.Level = getEnv("LOG_LEVEL", "info")

	// Metrics configuration
	cfg.Metrics.Enabled = getEnvBool("METRICS_ENABLED", true)
	cfg.Metrics.Port = getEnvInt("METRICS_PORT", 9090)

	return cfg, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an environment variable as an integer with a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool gets an environment variable as a boolean with a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvDuration gets an environment variable as a duration with a default value
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
