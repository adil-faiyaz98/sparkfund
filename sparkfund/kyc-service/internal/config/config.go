package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the KYC service
type Config struct {
	Server       ServerConfig
	Database     DatabaseConfig
	DocumentAPI  DocumentAPIConfig
	Metrics      MetricsConfig
	Tracing      TracingConfig
	Secrets      SecretsConfig
	Verification VerificationConfig
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
	Port           string
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

// DocumentAPIConfig holds configuration for document processing APIs
type DocumentAPIConfig struct {
	Provider     string // "aws-textract", "google-vision", "azure-ocr", "custom"
	Endpoint     string
	ApiKey       string
	Timeout      time.Duration
	MaxFileSize  int64    // in bytes
	AllowedTypes []string // e.g., "image/jpeg", "application/pdf"
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
	Provider      string
	Address       string
	Path          string
	TokenPath     string
	RotationCheck time.Duration
}

// VerificationConfig holds KYC verification specific configuration
type VerificationConfig struct {
	RequiredDocuments     []string
	RetentionPeriod       time.Duration
	PEPCheckEnabled       bool
	SanctionsCheckEnabled bool
	AutomaticApproval     bool
	MaxAttempts           int
	CooldownPeriod        time.Duration
	CountryRestrictions   []string
	VerificationTimeout   time.Duration
	BiometricEnabled      bool
	MinDocumentQuality    float64 // 0.0-1.0 quality threshold
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
		serviceEnv := filepath.Join("config", "kyc-service.env")
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
			RateLimit:       getEnvIntOrDefault("RATE_LIMIT", 50), // Lower limit due to document processing
			RateLimitBurst:  getEnvIntOrDefault("RATE_LIMIT_BURST", 100),
			Environment:     getEnvOrDefault("ENVIRONMENT", "development"),
			Debug:           getEnvOrDefault("DEBUG", "false") == "true",
			AllowOrigins:    getEnvSliceOrDefault("ALLOW_ORIGINS", []string{"*"}),
			TLS:             getEnvOrDefault("TLS_ENABLED", "false") == "true",
			TLSCert:         getEnvOrDefault("TLS_CERT_PATH", ""),
			TLSKey:          getEnvOrDefault("TLS_KEY_PATH", ""),
		},
		Database: DatabaseConfig{
			Host:           getEnvOrDefault("DB_HOST", "localhost"),
			Port:           getEnvOrDefault("DB_PORT", "5432"),
			User:           getEnvOrDefault("DB_USER", "postgres"),
			Password:       getEnvOrDefault("DB_PASSWORD", "postgres"),
			Name:           getEnvOrDefault("DB_NAME", "kyc_service"),
			SSLMode:        getEnvOrDefault("DB_SSL_MODE", "disable"),
			MaxConnections: getEnvIntOrDefault("DB_MAX_CONNECTIONS", 15),
			MaxIdleTime:    getEnvDurationOrDefault("DB_MAX_IDLE_TIME", 5*time.Minute),
			MigrationPath:  getEnvOrDefault("DB_MIGRATION_PATH", "migrations"),
			ConnectRetries: getEnvIntOrDefault("DB_CONNECT_RETRIES", 3),
			ConnectBackoff: getEnvDurationOrDefault("DB_CONNECT_BACKOFF", 5*time.Second),
		},
		DocumentAPI: DocumentAPIConfig{
			Provider:     getEnvOrDefault("DOC_API_PROVIDER", "aws-textract"),
			Endpoint:     getEnvOrDefault("DOC_API_ENDPOINT", "https://textract.us-east-1.amazonaws.com"),
			ApiKey:       getEnvOrDefault("DOC_API_KEY", ""),
			Timeout:      getEnvDurationOrDefault("DOC_API_TIMEOUT", 30*time.Second),
			MaxFileSize:  getEnvInt64OrDefault("DOC_MAX_FILE_SIZE", 10*1024*1024), // 10MB default
			AllowedTypes: getEnvSliceOrDefault("DOC_ALLOWED_TYPES", []string{"image/jpeg", "image/png", "application/pdf"}),
		},
		Metrics: MetricsConfig{
			Enabled:      getEnvOrDefault("METRICS_ENABLED", "true") == "true",
			Port:         getEnvOrDefault("METRICS_PORT", "9090"),
			Path:         getEnvOrDefault("METRICS_PATH", "/metrics"),
			ServiceLabel: getEnvOrDefault("METRICS_SERVICE_LABEL", "kyc-service"),
		},
		Tracing: TracingConfig{
			Enabled:        getEnvOrDefault("TRACING_ENABLED", "true") == "true",
			JaegerEndpoint: getEnvOrDefault("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
			ServiceName:    getEnvOrDefault("JAEGER_SERVICE_NAME", "kyc-service"),
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
		Verification: VerificationConfig{
			RequiredDocuments:     getEnvSliceOrDefault("KYC_REQUIRED_DOCUMENTS", []string{"identity", "address"}),
			RetentionPeriod:       getEnvDurationOrDefault("KYC_RETENTION_PERIOD", 7*24*365*time.Hour), // 7 years default
			PEPCheckEnabled:       getEnvOrDefault("KYC_PEP_CHECK_ENABLED", "true") == "true",
			SanctionsCheckEnabled: getEnvOrDefault("KYC_SANCTIONS_CHECK_ENABLED", "true") == "true",
			AutomaticApproval:     getEnvOrDefault("KYC_AUTOMATIC_APPROVAL", "false") == "true",
			MaxAttempts:           getEnvIntOrDefault("KYC_MAX_ATTEMPTS", 3),
			CooldownPeriod:        getEnvDurationOrDefault("KYC_COOLDOWN_PERIOD", 24*time.Hour),
			CountryRestrictions:   getEnvSliceOrDefault("KYC_COUNTRY_RESTRICTIONS", []string{}),
			VerificationTimeout:   getEnvDurationOrDefault("KYC_VERIFICATION_TIMEOUT", 72*time.Hour),
			BiometricEnabled:      getEnvOrDefault("KYC_BIOMETRIC_ENABLED", "false") == "true",
			MinDocumentQuality:    getEnvFloatOrDefault("KYC_MIN_DOCUMENT_QUALITY", 0.8),
		},
	}

	// Validate configuration
	if err := validateKYCConfiguration(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// validateKYCConfiguration performs KYC-specific validations
func validateKYCConfiguration(cfg *Config) error {
	if cfg.Server.Environment == "production" {
		// Check for secure configuration in production
		if cfg.DocumentAPI.ApiKey == "" {
			return errors.New("document API key must be set in production environment")
		}

		if !cfg.Server.TLS {
			return errors.New("TLS must be enabled in production environment")
		}

		if !cfg.PEPCheckEnabled || !cfg.SanctionsCheckEnabled {
			return errors.New("PEP and sanctions checks must be enabled in production")
		}
	}

	// Check reasonable document size limits
	if cfg.DocumentAPI.MaxFileSize <= 0 || cfg.DocumentAPI.MaxFileSize > 100*1024*1024 {
		return errors.New("document max file size must be between 1 byte and 100MB")
	}

	return nil
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
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

func getEnvInt64OrDefault(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
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
