package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config holds all configuration for the service
type Config struct {
	Environment string `mapstructure:"environment"`

	Server struct {
		Port            string        `mapstructure:"port"`
		ReadTimeout     time.Duration `mapstructure:"read_timeout"`
		WriteTimeout    time.Duration `mapstructure:"write_timeout"`
		IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
		ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
		TrustedProxies  []string      `mapstructure:"trusted_proxies"`
	} `mapstructure:"server"`

	Database struct {
		Host            string        `mapstructure:"host"`
		Port            string        `mapstructure:"port"`
		User            string        `mapstructure:"user"`
		Password        string        `mapstructure:"password"`
		Name            string        `mapstructure:"name"`
		SSLMode         string        `mapstructure:"sslmode"`
		MaxIdleConns    int           `mapstructure:"max_idle_conns"`
		MaxOpenConns    int           `mapstructure:"max_open_conns"`
		ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
		ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	} `mapstructure:"database"`

	JWT struct {
		Secret  string        `mapstructure:"secret"`
		Expiry  time.Duration `mapstructure:"expiry"`
		Refresh time.Duration `mapstructure:"refresh"`
		Issuer  string        `mapstructure:"issuer"`
		Enabled bool          `mapstructure:"enabled"`
	} `mapstructure:"jwt"`

	RateLimit struct {
		Enabled  bool          `mapstructure:"enabled"`
		Requests int           `mapstructure:"requests"`
		Window   time.Duration `mapstructure:"window"`
		Burst    int           `mapstructure:"burst"`
	} `mapstructure:"rate_limit"`

	Metrics struct {
		Enabled bool   `mapstructure:"enabled"`
		Path    string `mapstructure:"path"`
	} `mapstructure:"metrics"`

	Log struct {
		Level      string `mapstructure:"level"`
		Format     string `mapstructure:"format"`
		Output     string `mapstructure:"output"`
		RequestLog bool   `mapstructure:"request_log"`
	} `mapstructure:"log"`

	CircuitBreaker struct {
		Enabled             bool          `mapstructure:"enabled"`
		Timeout             time.Duration `mapstructure:"timeout"`
		MaxConcurrentReqs   uint32        `mapstructure:"max_concurrent_requests"`
		ErrorThresholdPerc  int           `mapstructure:"error_threshold_percentage"`
		RequestVolumeThresh uint64        `mapstructure:"request_volume_threshold"`
		SleepWindow         time.Duration `mapstructure:"sleep_window"`
	} `mapstructure:"circuit_breaker"`

	Security struct {
		AllowedOrigins []string `mapstructure:"allowed_origins"`
		AllowedMethods []string `mapstructure:"allowed_methods"`
		AllowedHeaders []string `mapstructure:"allowed_headers"`
		TrustedProxies []string `mapstructure:"trusted_proxies"`
		EnableCSRF     bool     `mapstructure:"enable_csrf"`
	} `mapstructure:"security"`

	Feature struct {
		EnableSwagger bool `mapstructure:"enable_swagger"`
		EnableAuth    bool `mapstructure:"enable_auth"`
		EnableMetrics bool `mapstructure:"enable_metrics"`
	} `mapstructure:"feature"`

	Tracing struct {
		Enabled     bool   `mapstructure:"enabled"`
		ServiceName string `mapstructure:"service_name"`
	} `mapstructure:"tracing"`

	Cache struct {
		Enabled         bool          `mapstructure:"enabled"`
		TTL             time.Duration `mapstructure:"ttl"`
		CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
	} `mapstructure:"cache"`

	TLS struct {
		Enabled   bool   `mapstructure:"enabled"`
		CertFile  string `mapstructure:"cert_file"`
		KeyFile   string `mapstructure:"key_file"`
		MinVersion string `mapstructure:"min_version"`
	} `mapstructure:"tls"`
}

var (
	config  Config
	once    sync.Once
	logger  *logrus.Logger
	initErr error
)

// Initialize sets up the logger for config package
func Initialize(l *logrus.Logger) {
	logger = l
}

// Load loads configuration from file and environment
func Load(configPath string) error {
	once.Do(func() {
		// If logger wasn't initialized, create a default one
		if logger == nil {
			logger = logrus.New()
		}

		// Set defaults
		setDefaults()

		// Get environment
		env := os.Getenv("APP_ENV")
		if env == "" {
			env = "development"
		}

		// Initialize viper
		v := viper.New()

		// Load base configuration that applies to all environments
		v.SetConfigName("config.base")
		v.SetConfigType("yaml")
		v.AddConfigPath(configPath)
		v.AddConfigPath("./config")
		v.AddConfigPath(".")

		// Read base config file (optional)
		if err := v.ReadInConfig(); err != nil {
			if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
				logger.Warnf("Error reading base config: %v", err)
			}
		}

		// Now load environment specific config to override base config
		v.SetConfigName(fmt.Sprintf("config.%s", env))

		// Attempt to read the environment-specific config
		if err := v.MergeInConfig(); err != nil {
			if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
				logger.Warnf("Error reading %s config: %v", env, err)
			}
		}

		// Enable environment variable overrides
		v.SetEnvPrefix("APP")
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		v.AutomaticEnv()

		// Handle secret injection via files (e.g., Kubernetes secrets)
		loadSecretsFromFiles(v)

		// Unmarshal config
		if err := v.Unmarshal(&config); err != nil {
			logger.Errorf("Failed to unmarshal configuration: %v", err)
			initErr = err
			return
		}

		// Set environment in config
		config.Environment = env

		// Validate critical configuration
		initErr = validateConfig(&config)
	})

	return initErr
}

// Get returns the loaded configuration
func Get() Config {
	return config
}

// setDefaults sets default configuration values
func setDefaults() {
	config.Environment = "development"

	config.Server.Port = "8081"
	config.Server.ReadTimeout = 5 * time.Second
	config.Server.WriteTimeout = 10 * time.Second
	config.Server.IdleTimeout = 120 * time.Second
	config.Server.ShutdownTimeout = 30 * time.Second
	config.Server.TrustedProxies = []string{"127.0.0.1", "172.16.0.0/12", "172.17.0.0/16", "192.168.0.0/16"}

	config.Database.Host = "postgres"
	config.Database.Port = "5432"
	config.Database.User = "postgres"
	config.Database.Name = "kyc_service"
	config.Database.SSLMode = "disable"
	config.Database.MaxIdleConns = 10
	config.Database.MaxOpenConns = 100
	config.Database.ConnMaxLifetime = time.Hour
	config.Database.ConnMaxIdleTime = 10 * time.Minute

	config.JWT.Expiry = 24 * time.Hour
	config.JWT.Refresh = 7 * 24 * time.Hour
	config.JWT.Issuer = "sparkfund"
	config.JWT.Enabled = true

	config.RateLimit.Enabled = true
	config.RateLimit.Requests = 60
	config.RateLimit.Window = time.Minute
	config.RateLimit.Burst = 10

	config.Metrics.Enabled = true
	config.Metrics.Path = "/metrics"

	config.Log.Level = "info"
	config.Log.Format = "json"
	config.Log.Output = "stdout"
	config.Log.RequestLog = true

	config.CircuitBreaker.Enabled = true
	config.CircuitBreaker.Timeout = 30 * time.Second
	config.CircuitBreaker.MaxConcurrentReqs = 100
	config.CircuitBreaker.ErrorThresholdPerc = 50
	config.CircuitBreaker.RequestVolumeThresh = 20
	config.CircuitBreaker.SleepWindow = 5 * time.Second

	config.Security.AllowedOrigins = []string{"*"}
	config.Security.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.Security.AllowedHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-CSRF-Token"}
	config.Security.TrustedProxies = []string{"127.0.0.1", "172.16.0.0/12", "192.168.0.0/16"}
	config.Security.EnableCSRF = true

	config.Feature.EnableSwagger = true
	config.Feature.EnableAuth = true
	config.Feature.EnableMetrics = true

	config.Tracing.Enabled = true
	config.Tracing.ServiceName = "kyc-service"

	config.Cache.Enabled = true
	config.Cache.TTL = 5 * time.Minute
	config.Cache.CleanupInterval = 10 * time.Minute

	config.TLS.Enabled = false
	config.TLS.MinVersion = "1.2"
}

// loadSecretsFromFiles loads secrets from mounted files (k8s secrets)
func loadSecretsFromFiles(v *viper.Viper) {
	// Check for JWT secret file
	if jwtSecretFile := os.Getenv("JWT_SECRET_FILE"); jwtSecretFile != "" {
		if data, err := os.ReadFile(jwtSecretFile); err == nil {
			v.Set("jwt.secret", string(data))
		}
	}

	// Check for database password file
	if dbPasswordFile := os.Getenv("DB_PASSWORD_FILE"); dbPasswordFile != "" {
		if data, err := os.ReadFile(dbPasswordFile); err == nil {
			v.Set("database.password", string(data))
		}
	}

	// Check for TLS certificate and key
	if tlsKeyFile := os.Getenv("TLS_KEY_FILE"); tlsKeyFile != "" {
		v.Set("tls.key_file", tlsKeyFile)
	}

	if tlsCertFile := os.Getenv("TLS_CERT_FILE"); tlsCertFile != "" {
		v.Set("tls.cert_file", tlsCertFile)
	}
}

// validateConfig checks that critical configuration is present
func validateConfig(cfg *Config) error {
	// In production, enforce certain security settings
	if os.Getenv("APP_ENV") == "production" {
		// Require JWT secret in production
		if cfg.JWT.Secret == "" {
			return errors.New("JWT authentication is not properly configured for production")
		}

		// Require non-default database password in production
		if cfg.Database.Password == "" || cfg.Database.Password == "postgres" {
			return errors.New("insecure database password in production")
		}

		// Require SSL for database in production
		if cfg.Database.SSLMode == "disable" {
			return errors.New("database SSL should be enabled in production")
		}
	}

	return nil
}

// Reload refreshes configuration at runtime
func Reload(configPath string) error {
	once = sync.Once{}
	return Load(configPath)
}

// IsProduction returns true if the current environment is production
func IsProduction() bool {
	return config.Environment == "production"
}
