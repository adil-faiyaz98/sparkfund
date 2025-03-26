package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the service
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Auth     AuthConfig
	Metrics  MetricsConfig
	Tracing  TracingConfig
	Secrets  SecretsConfig
	Port     string
	LogLevel string
}

// ServerConfig holds the configuration for the server
type ServerConfig struct {
	Port            string
	ShutdownTimeout time.Duration
	RateLimit       int
	RateLimitBurst  int
	Environment     string
	Debug           bool
	AllowOrigins    []string
	TLS             bool
	TLSCert         string
	TLSKey          string
}

// DatabaseConfig holds the configuration for the database
type DatabaseConfig struct {
	Host           string
	Port           int
	User           string
	Password       string
	Name           string
	SSLMode        string
	MaxConnections int
	MaxIdleTime    time.Duration
	MigrationPath  string
	ConnectRetries int
	ConnectBackoff time.Duration
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret            string
	ExpirationMinutes int
	RefreshExpiration int
	Issuer            string
	Audience          string
	RotationInterval  time.Duration
	Algorithm         string
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	MaxFailedAttempts    int
	LockoutDuration      time.Duration
	PasswordMinLength    int
	RequireSpecialChar   bool
	RequireNumber        bool
	RequireUppercase     bool
	PasswordHistoryCount int
	SessionTimeout       time.Duration
	MFAEnabled           bool
	MFAIssuer            string
}

// MetricsConfig holds metrics configuration
type MetricsConfig struct {
	Enabled      bool
	Port         string
	Path         string
	ServiceLabel string
}

// TracingConfig holds tracing configuration
type TracingConfig struct {
	Enabled        bool
	JaegerEndpoint string
	ServiceName    string
	SamplingRate   float64
	MaxExportBatch int
	BatchTimeout   time.Duration
	ExportTimeout  time.Duration
}

// SecretsConfig holds configuration for external secrets management
type SecretsConfig struct {
	Provider      string // "vault", "aws-ssm", "azure-keyvault"
	Address       string
	Path          string
	TokenPath     string
	RotationCheck time.Duration
}

// LoadConfig loads the configuration from environment variables
func LoadConfig(envFile ...string) (*Config, error) {
	// Load specified .env file or default to .env
	envPath := ".env"
	if len(envFile) > 0 && envFile[0] != "" {
		envPath = envFile[0]
	}

	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			log.Printf("Warning: Error loading %s file: %v", envPath, err)
		} else {
			log.Printf("Loaded configuration from %s", envPath)
		}
	}

	// Load service-specific .env file if in development
	if getEnvOrDefault("ENVIRONMENT", "development") == "development" {
		serviceEnv := filepath.Join("config", "security-service.env")
		if _, err := os.Stat(serviceEnv); err == nil {
			if err := godotenv.Load(serviceEnv); err != nil {
				log.Printf("Warning: Error loading %s file: %v", serviceEnv, err)
			} else {
				log.Printf("Loaded service configuration from %s", serviceEnv)
			}
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:            getEnvOrDefault("PORT", "8080"),
			ShutdownTimeout: getEnvDurationOrDefault("SHUTDOWN_TIMEOUT", 30*time.Second),
			RateLimit:       getEnvIntOrDefault("RATE_LIMIT", 100),
			RateLimitBurst:  getEnvIntOrDefault("RATE_LIMIT_BURST", 200),
			Environment:     getEnvOrDefault("ENVIRONMENT", "development"),
			Debug:           getEnvOrDefault("DEBUG", "false") == "true",
			AllowOrigins:    getEnvSliceOrDefault("ALLOW_ORIGINS", []string{"*"}),
			TLS:             getEnvOrDefault("TLS_ENABLED", "false") == "true",
			TLSCert:         getEnvOrDefault("TLS_CERT_PATH", ""),
			TLSKey:          getEnvOrDefault("TLS_KEY_PATH", ""),
		},
		Database: DatabaseConfig{
			Host:           getEnvOrDefault("DB_HOST", "localhost"),
			Port:           getEnvAsInt("DB_PORT", 5432),
			User:           getEnvOrDefault("DB_USER", "postgres"),
			Password:       getEnvOrDefault("DB_PASSWORD", "postgres"),
			Name:           getEnvOrDefault("DB_NAME", "security_service"),
			SSLMode:        getEnvOrDefault("DB_SSL_MODE", "disable"),
			MaxConnections: getEnvIntOrDefault("DB_MAX_CONNECTIONS", 10),
			MaxIdleTime:    getEnvDurationOrDefault("DB_MAX_IDLE_TIME", 15*time.Minute),
			MigrationPath:  getEnvOrDefault("DB_MIGRATION_PATH", "migrations"),
			ConnectRetries: getEnvIntOrDefault("DB_CONNECT_RETRIES", 3),
			ConnectBackoff: getEnvDurationOrDefault("DB_CONNECT_BACKOFF", 5*time.Second),
		},
		JWT: JWTConfig{
			Secret:            getEnvOrDefault("JWT_SECRET", "your-default-secret-change-in-production"),
			ExpirationMinutes: getEnvIntOrDefault("JWT_EXPIRATION_MINUTES", 60),
			RefreshExpiration: getEnvIntOrDefault("JWT_REFRESH_EXPIRATION", 1440),
			Issuer:            getEnvOrDefault("JWT_ISSUER", "sparkfund-security-service"),
			Audience:          getEnvOrDefault("JWT_AUDIENCE", "sparkfund"),
			RotationInterval:  getEnvDurationOrDefault("JWT_ROTATION_INTERVAL", 24*time.Hour),
			Algorithm:         getEnvOrDefault("JWT_ALGORITHM", "HS256"),
		},
		Auth: AuthConfig{
			MaxFailedAttempts:    getEnvIntOrDefault("AUTH_MAX_FAILED_ATTEMPTS", 5),
			LockoutDuration:      getEnvDurationOrDefault("AUTH_LOCKOUT_DURATION", 15*time.Minute),
			PasswordMinLength:    getEnvIntOrDefault("AUTH_PASSWORD_MIN_LENGTH", 8),
			RequireSpecialChar:   getEnvOrDefault("AUTH_REQUIRE_SPECIAL_CHAR", "true") == "true",
			RequireNumber:        getEnvOrDefault("AUTH_REQUIRE_NUMBER", "true") == "true",
			RequireUppercase:     getEnvOrDefault("AUTH_REQUIRE_UPPERCASE", "true") == "true",
			PasswordHistoryCount: getEnvIntOrDefault("AUTH_PASSWORD_HISTORY_COUNT", 5),
			SessionTimeout:       getEnvDurationOrDefault("AUTH_SESSION_TIMEOUT", 30*time.Minute),
			MFAEnabled:           getEnvOrDefault("AUTH_MFA_ENABLED", "false") == "true",
			MFAIssuer:            getEnvOrDefault("AUTH_MFA_ISSUER", "sparkfund"),
		},
		Metrics: MetricsConfig{
			Enabled:      getEnvOrDefault("METRICS_ENABLED", "true") == "true",
			Port:         getEnvOrDefault("METRICS_PORT", "9090"),
			Path:         getEnvOrDefault("METRICS_PATH", "/metrics"),
			ServiceLabel: getEnvOrDefault("METRICS_SERVICE_LABEL", "security-service"),
		},
		Tracing: TracingConfig{
			Enabled:        getEnvOrDefault("TRACING_ENABLED", "true") == "true",
			JaegerEndpoint: getEnvOrDefault("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
			ServiceName:    getEnvOrDefault("JAEGER_SERVICE_NAME", "security-service"),
			SamplingRate:   getEnvFloatOrDefault("TRACING_SAMPLING_RATE", 1.0),
			MaxExportBatch: getEnvIntOrDefault("TRACING_MAX_EXPORT_BATCH", 100),
			BatchTimeout:   getEnvDurationOrDefault("TRACING_BATCH_TIMEOUT", 5*time.Second),
			ExportTimeout:  getEnvDurationOrDefault("TRACING_EXPORT_TIMEOUT", 30*time.Second),
		},
		Secrets: SecretsConfig{
			Provider:      getEnvOrDefault("SECRETS_PROVIDER", "vault"),
			Address:       getEnvOrDefault("SECRETS_ADDRESS", "http://localhost:8200"),
			Path:          getEnvOrDefault("SECRETS_PATH", "secret/data"),
			TokenPath:     getEnvOrDefault("SECRETS_TOKEN_PATH", "/var/run/secrets/kubernetes.io/serviceaccount/token"),
			RotationCheck: getEnvDurationOrDefault("SECRETS_ROTATION_CHECK", 1*time.Hour),
		},
		Port:     getEnvOrDefault("PORT", "8080"),
		LogLevel: getEnvOrDefault("LOG_LEVEL", "info"),
	}

	return cfg, nil
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
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
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
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

func getEnvFloatOrDefault(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// String returns the string representation of the config
func (c *Config) String() string {
	bytes, _ := json.MarshalIndent(struct {
		Server   ServerConfig
		Database struct {
			Host     string
			Port     int
			User     string
			DBName   string
			SSLMode  string
			Password string `json:"-"` // Hide password in logs
		}
		JWT struct {
			Secret            string
			ExpirationMinutes int
			RefreshExpiration int
			Issuer            string
			Audience          string
			RotationInterval  time.Duration
			Algorithm         string
		}
		Auth     AuthConfig
		Metrics  MetricsConfig
		Tracing  TracingConfig
		Secrets  SecretsConfig
		Port     string
		LogLevel string
	}{
		Server: c.Server,
		Database: struct {
			Host     string
			Port     int
			User     string
			DBName   string
			SSLMode  string
			Password string `json:"-"`
		}{
			Host:     c.Database.Host,
			Port:     c.Database.Port,
			User:     c.Database.User,
			DBName:   c.Database.Name,
			SSLMode:  c.Database.SSLMode,
			Password: "********",
		},
		JWT: struct {
			Secret            string
			ExpirationMinutes int
			RefreshExpiration int
			Issuer            string
			Audience          string
			RotationInterval  time.Duration
			Algorithm         string
		}{
			Secret:            c.JWT.Secret,
			ExpirationMinutes: c.JWT.ExpirationMinutes,
			RefreshExpiration: c.JWT.RefreshExpiration,
			Issuer:            c.JWT.Issuer,
			Audience:          c.JWT.Audience,
			RotationInterval:  c.JWT.RotationInterval,
			Algorithm:         c.JWT.Algorithm,
		},
		Auth:     c.Auth,
		Metrics:  c.Metrics,
		Tracing:  c.Tracing,
		Secrets:  c.Secrets,
		Port:     c.Port,
		LogLevel: c.LogLevel,
	}, "", "  ")

	return string(bytes)
}
