package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the service
type Config struct {
	App            AppConfig            `mapstructure:"app"`
	Server         ServerConfig         `mapstructure:"server"`
	Database       DatabaseConfig       `mapstructure:"database"`
	JWT            JWTConfig            `mapstructure:"jwt"`
	RateLimit      RateLimitConfig      `mapstructure:"rate_limit"`
	Metrics        MetricsConfig        `mapstructure:"metrics"`
	Log            LogConfig            `mapstructure:"log"`
	CircuitBreaker CircuitBreakerConfig `mapstructure:"circuit_breaker"`
	Security       SecurityConfig       `mapstructure:"security"`
	Feature        FeatureConfig        `mapstructure:"feature"`
	FeatureFlags   map[string]bool      `mapstructure:"feature_flags"`
	Tracing        TracingConfig        `mapstructure:"tracing"`
	Cache          CacheConfig          `mapstructure:"cache"`
	TLS            TLSConfig            `mapstructure:"tls"`
	AI             AIConfig             `mapstructure:"ai"`
	Storage        StorageConfig        `mapstructure:"storage"`
	Validation     ValidationConfig     `mapstructure:"validation"`
	Notifications  NotificationConfig   `mapstructure:"notifications"`
	Monitoring     MonitoringConfig     `mapstructure:"monitoring"`
	Events         EventsConfig         `mapstructure:"events"`
}

// AppConfig holds application configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	Timeout         time.Duration `mapstructure:"timeout"`
	TrustedProxies  []string      `mapstructure:"trusted_proxies"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret  string        `mapstructure:"secret"`
	Expiry  time.Duration `mapstructure:"expiry"`
	Refresh time.Duration `mapstructure:"refresh"`
	Issuer  string        `mapstructure:"issuer"`
	Enabled bool          `mapstructure:"enabled"`
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled  bool          `mapstructure:"enabled"`
	Requests int           `mapstructure:"requests"`
	Window   time.Duration `mapstructure:"window"`
	Burst    int           `mapstructure:"burst"`
}

// MetricsConfig holds metrics configuration
type MetricsConfig struct {
	Enabled      bool          `mapstructure:"enabled"`
	Path         string        `mapstructure:"path"`
	Port         int           `mapstructure:"port"`
	PushInterval time.Duration `mapstructure:"push_interval"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	RequestLog bool   `mapstructure:"request_log"`
}

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	Enabled             bool          `mapstructure:"enabled"`
	Timeout             time.Duration `mapstructure:"timeout"`
	MaxConcurrentReqs   uint32        `mapstructure:"max_concurrent_requests"`
	ErrorThresholdPerc  int           `mapstructure:"error_threshold_percentage"`
	RequestVolumeThresh uint64        `mapstructure:"request_volume_threshold"`
	SleepWindow         time.Duration `mapstructure:"sleep_window"`
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	AllowedOrigins []string      `mapstructure:"allowed_origins"`
	AllowedMethods []string      `mapstructure:"allowed_methods"`
	AllowedHeaders []string      `mapstructure:"allowed_headers"`
	TrustedProxies []string      `mapstructure:"trusted_proxies"`
	EnableCSRF     bool          `mapstructure:"enable_csrf"`
	JWTSecret      string        `mapstructure:"jwt_secret"`
	JWTExpiry      time.Duration `mapstructure:"jwt_expiry"`
	RateLimit      int           `mapstructure:"rate_limit"`
	RateWindow     int           `mapstructure:"rate_window"`
	AuditLogging   struct {
		Enabled        bool   `mapstructure:"enabled"`
		LogPredictions bool   `mapstructure:"log_predictions"`
		LogRetention   string `mapstructure:"log_retention"`
	} `mapstructure:"audit_logging"`
	AccessControl struct {
		RoleBased     bool     `mapstructure:"role_based"`
		RequiredRoles []string `mapstructure:"required_roles"`
	} `mapstructure:"access_control"`
}

// FeatureConfig holds feature configuration
type FeatureConfig struct {
	EnableSwagger bool `mapstructure:"enable_swagger"`
	EnableAuth    bool `mapstructure:"enable_auth"`
	EnableMetrics bool `mapstructure:"enable_metrics"`
}

// TracingConfig holds tracing configuration
type TracingConfig struct {
	Enabled      bool    `mapstructure:"enabled"`
	ServiceName  string  `mapstructure:"service_name"`
	SamplingRate float64 `mapstructure:"sampling_rate"`
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Enabled         bool          `mapstructure:"enabled"`
	Type            string        `mapstructure:"type"`
	TTL             time.Duration `mapstructure:"ttl"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
	Redis           struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
		Prefix   string `mapstructure:"prefix"`
	} `mapstructure:"redis"`
}

// TLSConfig holds TLS configuration
type TLSConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	CertFile   string `mapstructure:"cert_file"`
	KeyFile    string `mapstructure:"key_file"`
	MinVersion string `mapstructure:"min_version"`
}

// AIConfig holds AI configuration
type AIConfig struct {
	ServiceURL string                 `mapstructure:"service_url"`
	APIKey     string                 `mapstructure:"api_key"`
	Timeout    time.Duration          `mapstructure:"timeout"`
	MaxRetries int                    `mapstructure:"max_retries"`
	RetryDelay time.Duration          `mapstructure:"retry_delay"`
	Models     map[string]ModelConfig `mapstructure:"models"`
}

// ModelConfig holds AI model configuration
type ModelConfig struct {
	Path          string  `mapstructure:"path"`
	Version       string  `mapstructure:"version"`
	BatchSize     int     `mapstructure:"batch_size"`
	Threshold     float64 `mapstructure:"threshold"`
	RealTime      bool    `mapstructure:"real_time"`
	AuditLogging  bool    `mapstructure:"audit_logging"`
	Preprocessing struct {
		ImageSize []int `mapstructure:"image_size"`
		Normalize bool  `mapstructure:"normalize"`
	} `mapstructure:"preprocessing"`
	Security struct {
		EncryptionEnabled bool   `mapstructure:"encryption_enabled"`
		KeyRotationPeriod string `mapstructure:"key_rotation_period"`
	} `mapstructure:"security"`
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	Type  string `mapstructure:"type"`
	Local struct {
		Path    string `mapstructure:"path"`
		TempDir string `mapstructure:"temp_dir"`
	} `mapstructure:"local"`
	S3 struct {
		Bucket string `mapstructure:"bucket"`
		Region string `mapstructure:"region"`
	} `mapstructure:"s3"`
	Retention struct {
		Documents           string `mapstructure:"documents"`
		VerificationResults string `mapstructure:"verification_results"`
	} `mapstructure:"retention"`
}

// ValidationConfig holds validation configuration
type ValidationConfig struct {
	Document struct {
		MaxSize       int64    `mapstructure:"max_size"`
		AllowedTypes  []string `mapstructure:"allowed_types"`
		MinResolution int      `mapstructure:"min_resolution"`
		MaxPages      int      `mapstructure:"max_pages"`
	} `mapstructure:"document"`
	Face struct {
		MinSize        int      `mapstructure:"min_size"`
		MaxSize        int      `mapstructure:"max_size"`
		AllowedFormats []string `mapstructure:"allowed_formats"`
	} `mapstructure:"face"`
}

// NotificationConfig holds notification configuration
type NotificationConfig struct {
	Enabled   bool `mapstructure:"enabled"`
	Providers struct {
		Email struct {
			SMTPHost    string `mapstructure:"smtp_host"`
			SMTPPort    int    `mapstructure:"smtp_port"`
			FromAddress string `mapstructure:"from_address"`
		} `mapstructure:"email"`
		SMS struct {
			Provider   string `mapstructure:"provider"`
			FromNumber string `mapstructure:"from_number"`
		} `mapstructure:"sms"`
	} `mapstructure:"providers"`
}

// MonitoringConfig holds monitoring configuration
type MonitoringConfig struct {
	Tracing struct {
		Enabled      bool    `mapstructure:"enabled"`
		SamplingRate float64 `mapstructure:"sampling_rate"`
	} `mapstructure:"tracing"`
	Metrics struct {
		Enabled      bool          `mapstructure:"enabled"`
		PushInterval time.Duration `mapstructure:"push_interval"`
	} `mapstructure:"metrics"`
	Alerts struct {
		Enabled  bool     `mapstructure:"enabled"`
		Channels []string `mapstructure:"channels"`
	} `mapstructure:"alerts"`
}

// EventsConfig holds events configuration
type EventsConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	BrokerType  string `mapstructure:"broker_type"`
	BrokerURL   string `mapstructure:"broker_url"`
	TopicPrefix string `mapstructure:"topic_prefix"`
}

// Global configuration instance
var cfg *Config

// Load loads configuration from files and environment variables
func Load() (*Config, error) {
	// Create a new Viper instance
	v := viper.New()

	// Set default configuration file paths
	v.SetConfigName("config.base")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	// Read base configuration
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read base config: %w", err)
	}

	// Get environment
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// Load environment-specific configuration
	v.SetConfigName(fmt.Sprintf("config.%s", env))
	if err := v.MergeInConfig(); err != nil {
		// It's okay if the environment-specific config doesn't exist
		if !strings.Contains(err.Error(), "Not Found") {
			return nil, fmt.Errorf("failed to read %s config: %w", env, err)
		}
	}

	// Enable environment variable overrides
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Unmarshal configuration
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set environment in config
	config.App.Environment = env

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Store configuration globally
	cfg = &config

	return &config, nil
}

// Get returns the loaded configuration
func Get() *Config {
	return cfg
}

// validateConfig validates the configuration
func validateConfig(cfg *Config) error {
	// Validate required fields
	if cfg.App.Name == "" {
		return fmt.Errorf("app name is required")
	}

	if cfg.Server.Port == 0 {
		return fmt.Errorf("server port is required")
	}

	// Validate database configuration in production
	if cfg.App.Environment == "production" {
		if cfg.Database.Host == "" {
			return fmt.Errorf("database host is required in production")
		}

		if cfg.Database.Password == "" || cfg.Database.Password == "postgres" {
			return fmt.Errorf("insecure database password in production")
		}

		if cfg.JWT.Secret == "" || cfg.JWT.Secret == "your-secret-key" {
			return fmt.Errorf("insecure JWT secret in production")
		}

		if cfg.Database.SSLMode == "disable" {
			return fmt.Errorf("database SSL should be enabled in production")
		}
	}

	return nil
}
