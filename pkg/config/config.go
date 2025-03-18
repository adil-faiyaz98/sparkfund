package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration settings
type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Auth        AuthConfig
	Monitoring  MonitoringConfig
	Logging     LoggingConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	RateLimit       RateLimitConfig
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	Timeout  time.Duration
}

// AuthConfig holds authentication-related configuration
type AuthConfig struct {
	JWTSecret        string
	TokenExpiration  time.Duration
	RefreshSecret    string
	RefreshDuration  time.Duration
	PasswordSaltCost int
}

// MonitoringConfig holds monitoring-related configuration
type MonitoringConfig struct {
	PrometheusPort int
	MetricsPath    string
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level      string
	Format     string
	OutputPath string
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerSecond float64
	BurstSize         int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{
		Environment: getEnvString("ENVIRONMENT", "development"),
		Server: ServerConfig{
			Port:            getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:     getEnvDuration("SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout:    getEnvDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			ShutdownTimeout: getEnvDuration("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
			RateLimit: RateLimitConfig{
				RequestsPerSecond: getEnvFloat("RATE_LIMIT_RPS", 100),
				BurstSize:         getEnvInt("RATE_LIMIT_BURST", 200),
			},
		},
		Database: DatabaseConfig{
			Host:     getEnvString("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnvString("DB_USER", "postgres"),
			Password: getEnvString("DB_PASSWORD", ""),
			DBName:   getEnvString("DB_NAME", "moneypulse"),
			SSLMode:  getEnvString("DB_SSLMODE", "disable"),
			MaxConns: getEnvInt("DB_MAX_CONNECTIONS", 20),
			Timeout:  getEnvDuration("DB_TIMEOUT", 5*time.Second),
		},
		Auth: AuthConfig{
			JWTSecret:        getEnvString("JWT_SECRET", ""),
			TokenExpiration:  getEnvDuration("TOKEN_EXPIRATION", 24*time.Hour),
			RefreshSecret:    getEnvString("REFRESH_SECRET", ""),
			RefreshDuration:  getEnvDuration("REFRESH_DURATION", 7*24*time.Hour),
			PasswordSaltCost: getEnvInt("PASSWORD_SALT_COST", 10),
		},
		Monitoring: MonitoringConfig{
			PrometheusPort: getEnvInt("PROMETHEUS_PORT", 9090),
			MetricsPath:    getEnvString("METRICS_PATH", "/metrics"),
		},
		Logging: LoggingConfig{
			Level:      getEnvString("LOG_LEVEL", "info"),
			Format:     getEnvString("LOG_FORMAT", "json"),
			OutputPath: getEnvString("LOG_OUTPUT", "stdout"),
		},
	}

	// Validate required configurations
	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// validate checks if all required configurations are set
func (c *Config) validate() error {
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.Auth.RefreshSecret == "" {
		return fmt.Errorf("REFRESH_SECRET is required")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	return nil
}

// Helper functions to get environment variables with default values
func getEnvString(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	valueStr := getEnvString(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvFloat(key string, defaultVal float64) float64 {
	valueStr := getEnvString(key, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	valueStr := getEnvString(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
