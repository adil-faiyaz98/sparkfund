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
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Cache    CacheConfig    `mapstructure:"cache"`
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
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret  string        `mapstructure:"secret"`
	Expiry  time.Duration `mapstructure:"expiry"`
	Refresh time.Duration `mapstructure:"refresh"`
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Enabled bool          `mapstructure:"enabled"`
	Type    string        `mapstructure:"type"`
	TTL     time.Duration `mapstructure:"ttl"`
	Redis   struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
}

// Load loads configuration from files and environment variables
func Load() (*Config, error) {
	// Create a new Viper instance
	v := viper.New()

	// Set default configuration file paths
	v.SetConfigName("config.base")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")
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

	return &config, nil
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

// GetDatabaseURL returns the database connection string
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}
