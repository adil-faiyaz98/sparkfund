package ai

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds configuration for AI components
type Config struct {
	// Model configuration
	ModelPath          string
	BatchSize          int
	UpdateInterval     time.Duration
	ValidationInterval time.Duration
	MinAccuracy        float64

	// Training configuration
	MaxRetries        int
	RetryDelay        time.Duration
	MaxConcurrentJobs int
	JobTimeout        time.Duration

	// Validation configuration
	ValidationBatchSize  int
	TestDataRatio        float64
	CrossValidationFolds int

	// Storage configuration
	StorageType      string
	StoragePath      string
	StorageRetention time.Duration

	// Monitoring configuration
	MetricsEnabled  bool
	MetricsInterval time.Duration
	LogLevel        string

	// Security configuration
	EncryptionEnabled bool
	EncryptionKey     string
	AccessControl     AccessControlConfig
}

// AccessControlConfig defines access control settings
type AccessControlConfig struct {
	Enabled           bool
	AllowedIPs        []string
	AllowedUsers      []string
	AllowedRoles      []string
	MaxFailedAttempts int
	LockoutDuration   time.Duration
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		ModelPath:          "./models",
		BatchSize:          100,
		UpdateInterval:     time.Minute * 5,
		ValidationInterval: time.Hour,
		MinAccuracy:        0.85,

		MaxRetries:        3,
		RetryDelay:        time.Second * 5,
		MaxConcurrentJobs: 5,
		JobTimeout:        time.Minute * 30,

		ValidationBatchSize:  1000,
		TestDataRatio:        0.2,
		CrossValidationFolds: 5,

		StorageType:      "file",
		StoragePath:      "./storage",
		StorageRetention: time.Hour * 24 * 30,

		MetricsEnabled:  true,
		MetricsInterval: time.Minute,
		LogLevel:        "info",

		EncryptionEnabled: true,
		EncryptionKey:     "", // Must be set in production
		AccessControl: AccessControlConfig{
			Enabled:           true,
			AllowedIPs:        []string{},
			AllowedUsers:      []string{},
			AllowedRoles:      []string{"admin", "analyst"},
			MaxFailedAttempts: 5,
			LockoutDuration:   time.Minute * 15,
		},
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.ModelPath == "" {
		return fmt.Errorf("model path is required")
	}
	if c.BatchSize <= 0 {
		return fmt.Errorf("batch size must be positive")
	}
	if c.UpdateInterval <= 0 {
		return fmt.Errorf("update interval must be positive")
	}
	if c.ValidationInterval <= 0 {
		return fmt.Errorf("validation interval must be positive")
	}
	if c.MinAccuracy < 0 || c.MinAccuracy > 1 {
		return fmt.Errorf("minimum accuracy must be between 0 and 1")
	}
	if c.MaxRetries < 0 {
		return fmt.Errorf("max retries must be non-negative")
	}
	if c.RetryDelay <= 0 {
		return fmt.Errorf("retry delay must be positive")
	}
	if c.MaxConcurrentJobs <= 0 {
		return fmt.Errorf("max concurrent jobs must be positive")
	}
	if c.JobTimeout <= 0 {
		return fmt.Errorf("job timeout must be positive")
	}
	if c.ValidationBatchSize <= 0 {
		return fmt.Errorf("validation batch size must be positive")
	}
	if c.TestDataRatio <= 0 || c.TestDataRatio >= 1 {
		return fmt.Errorf("test data ratio must be between 0 and 1")
	}
	if c.CrossValidationFolds < 2 {
		return fmt.Errorf("cross validation folds must be at least 2")
	}
	if c.StorageType == "" {
		return fmt.Errorf("storage type is required")
	}
	if c.StoragePath == "" {
		return fmt.Errorf("storage path is required")
	}
	if c.StorageRetention <= 0 {
		return fmt.Errorf("storage retention must be positive")
	}
	if c.MetricsInterval <= 0 {
		return fmt.Errorf("metrics interval must be positive")
	}
	if c.EncryptionEnabled && c.EncryptionKey == "" {
		return fmt.Errorf("encryption key is required when encryption is enabled")
	}
	if c.AccessControl.Enabled {
		if c.AccessControl.MaxFailedAttempts <= 0 {
			return fmt.Errorf("max failed attempts must be positive")
		}
		if c.AccessControl.LockoutDuration <= 0 {
			return fmt.Errorf("lockout duration must be positive")
		}
	}

	return nil
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := DefaultConfig()

	// Load from environment variables
	// This is a simplified example - in production, you would use a proper config loader
	if modelPath := os.Getenv("AI_MODEL_PATH"); modelPath != "" {
		config.ModelPath = modelPath
	}
	if batchSize := os.Getenv("AI_BATCH_SIZE"); batchSize != "" {
		if size, err := strconv.Atoi(batchSize); err == nil {
			config.BatchSize = size
		}
	}
	if updateInterval := os.Getenv("AI_UPDATE_INTERVAL"); updateInterval != "" {
		if interval, err := time.ParseDuration(updateInterval); err == nil {
			config.UpdateInterval = interval
		}
	}
	if validationInterval := os.Getenv("AI_VALIDATION_INTERVAL"); validationInterval != "" {
		if interval, err := time.ParseDuration(validationInterval); err == nil {
			config.ValidationInterval = interval
		}
	}
	if minAccuracy := os.Getenv("AI_MIN_ACCURACY"); minAccuracy != "" {
		if accuracy, err := strconv.ParseFloat(minAccuracy, 64); err == nil {
			config.MinAccuracy = accuracy
		}
	}
	if encryptionKey := os.Getenv("AI_ENCRYPTION_KEY"); encryptionKey != "" {
		config.EncryptionKey = encryptionKey
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %v", err)
	}

	return config, nil
}
