package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Logging  LoggingConfig  `yaml:"logging"`
	Security SecurityConfig `yaml:"security"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port            string `yaml:"port"`
	ReadTimeout     int    `yaml:"read_timeout"`
	WriteTimeout    int    `yaml:"write_timeout"`
	MaxHeaderBytes  int    `yaml:"max_header_bytes"`
	MaxRequestSize  int64  `yaml:"max_request_size"`
	AllowedOrigins  []string `yaml:"allowed_origins"`
	AllowedMethods  []string `yaml:"allowed_methods"`
	AllowedHeaders  []string `yaml:"allowed_headers"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	JWTSecret     string `yaml:"jwt_secret"`
	JWTExpiry     int    `yaml:"jwt_expiry"`
	RateLimit     int    `yaml:"rate_limit"`
	RateWindow    int    `yaml:"rate_window"`
	BlockDuration int    `yaml:"block_duration"`
}

// Load loads configuration from a YAML file
func Load(path string) (*Config, error) {
	// Read configuration file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse configuration
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Set default values
	setDefaults(&config)

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %v", err)
	}

	return &config, nil
}

// setDefaults sets default values for configuration
func setDefaults(config *Config) {
	// Server defaults
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 10
	}
	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 10
	}
	if config.Server.MaxHeaderBytes == 0 {
		config.Server.MaxHeaderBytes = 1 << 20 // 1MB
	}
	if config.Server.MaxRequestSize == 0 {
		config.Server.MaxRequestSize = 10 << 20 // 10MB
	}
	if len(config.Server.AllowedOrigins) == 0 {
		config.Server.AllowedOrigins = []string{"*"}
	}
	if len(config.Server.AllowedMethods) == 0 {
		config.Server.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}
	if len(config.Server.AllowedHeaders) == 0 {
		config.Server.AllowedHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	}

	// Database defaults
	if config.Database.Host == "" {
		config.Database.Host = "localhost"
	}
	if config.Database.Port == 0 {
		config.Database.Port = 5432
	}
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}

	// Redis defaults
	if config.Redis.Host == "" {
		config.Redis.Host = "localhost"
	}
	if config.Redis.Port == 0 {
		config.Redis.Port = 6379
	}

	// Logging defaults
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.Format == "" {
		config.Logging.Format = "json"
	}
	if config.Logging.Output == "" {
		config.Logging.Output = "stdout"
	}
	if config.Logging.MaxSize == 0 {
		config.Logging.MaxSize = 100
	}
	if config.Logging.MaxBackups == 0 {
		config.Logging.MaxBackups = 3
	}
	if config.Logging.MaxAge == 0 {
		config.Logging.MaxAge = 28
	}

	// Security defaults
	if config.Security.JWTExpiry == 0 {
		config.Security.JWTExpiry = 24
	}
	if config.Security.RateLimit == 0 {
		config.Security.RateLimit = 100
	}
	if config.Security.RateWindow == 0 {
		config.Security.RateWindow = 60
	}
	if config.Security.BlockDuration == 0 {
		config.Security.BlockDuration = 300
	}
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	// Validate server configuration
	if config.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	// Validate database configuration
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if config.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if config.Database.Password == "" {
		return fmt.Errorf("database password is required")
	}
	if config.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	// Validate Redis configuration
	if config.Redis.Host == "" {
		return fmt.Errorf("redis host is required")
	}

	// Validate security configuration
	if config.Security.JWTSecret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	return nil
}
