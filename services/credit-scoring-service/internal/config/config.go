package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sparkfund/credit-scoring-service/internal/errors"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Security SecurityConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxHeaderBytes  int
	GracefulTimeout time.Duration
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	URL            string
	Password       string
	DB             int
	PoolSize       int
	MinIdleConns   int
	MaxConnAge     time.Duration
	RequestTimeout time.Duration
}

type JWTConfig struct {
	SecretKey     string
	TokenExpiry   time.Duration
	Issuer        string
	Audience      string
	AllowedScopes []string
}

type SecurityConfig struct {
	RateLimit struct {
		RequestsPerMinute int
		BurstSize        int
	}
	CORS struct {
		AllowedOrigins []string
		AllowedMethods []string
		AllowedHeaders []string
	}
	Headers struct {
		EnableHSTS      bool
		EnableXSSFilter bool
		EnableNoSniff   bool
		FrameOptions    string
	}
}

type LoggingConfig struct {
	Level      string
	Format     string
	Output     string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

func Load() (*Config, error) {
	config := &Config{}

	// Server config
	config.Server = ServerConfig{
		Port:            getEnvOrDefault("PORT", "8080"),
		ReadTimeout:     getDurationEnvOrDefault("SERVER_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:    getDurationEnvOrDefault("SERVER_WRITE_TIMEOUT", 10*time.Second),
		MaxHeaderBytes:  getIntEnvOrDefault("SERVER_MAX_HEADER_BYTES", 1<<20), // 1MB
		GracefulTimeout: getDurationEnvOrDefault("SERVER_GRACEFUL_TIMEOUT", 5*time.Second),
	}

	// Database config
	config.Database = DatabaseConfig{
		Host:            getEnvOrDefault("DB_HOST", "localhost"),
		Port:            getEnvOrDefault("DB_PORT", "5432"),
		User:            getEnvOrDefault("DB_USER", "postgres"),
		Password:        getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:          getEnvOrDefault("DB_NAME", "sparkfund"),
		SSLMode:         getEnvOrDefault("DB_SSL_MODE", "disable"),
		MaxIdleConns:    getIntEnvOrDefault("DB_MAX_IDLE_CONNS", 10),
		MaxOpenConns:    getIntEnvOrDefault("DB_MAX_OPEN_CONNS", 100),
		ConnMaxLifetime: getDurationEnvOrDefault("DB_CONN_MAX_LIFETIME", time.Hour),
	}

	// Redis config
	config.Redis = RedisConfig{
		URL:            getEnvOrDefault("REDIS_URL", "localhost:6379"),
		Password:       getEnvOrDefault("REDIS_PASSWORD", ""),
		DB:             getIntEnvOrDefault("REDIS_DB", 0),
		PoolSize:       getIntEnvOrDefault("REDIS_POOL_SIZE", 10),
		MinIdleConns:   getIntEnvOrDefault("REDIS_MIN_IDLE_CONNS", 5),
		MaxConnAge:     getDurationEnvOrDefault("REDIS_MAX_CONN_AGE", 30*time.Minute),
		RequestTimeout: getDurationEnvOrDefault("REDIS_REQUEST_TIMEOUT", 3*time.Second),
	}

	// JWT config
	config.JWT = JWTConfig{
		SecretKey:     getEnvOrDefault("JWT_SECRET_KEY", ""),
		TokenExpiry:   getDurationEnvOrDefault("JWT_TOKEN_EXPIRY", 24*time.Hour),
		Issuer:        getEnvOrDefault("JWT_ISSUER", "sparkfund-credit-service"),
		Audience:      getEnvOrDefault("JWT_AUDIENCE", "sparkfund-api"),
		AllowedScopes: getStringSliceEnvOrDefault("JWT_ALLOWED_SCOPES", []string{"credit:read", "credit:write"}),
	}

	// Security config
	config.Security.RateLimit.RequestsPerMinute = getIntEnvOrDefault("RATE_LIMIT_REQUESTS_PER_MINUTE", 100)
	config.Security.RateLimit.BurstSize = getIntEnvOrDefault("RATE_LIMIT_BURST_SIZE", 10)
	config.Security.CORS.AllowedOrigins = getStringSliceEnvOrDefault("CORS_ALLOWED_ORIGINS", []string{"*"})
	config.Security.CORS.AllowedMethods = getStringSliceEnvOrDefault("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	config.Security.CORS.AllowedHeaders = getStringSliceEnvOrDefault("CORS_ALLOWED_HEADERS", []string{"Authorization", "Content-Type"})
	config.Security.Headers.EnableHSTS = getBoolEnvOrDefault("SECURITY_ENABLE_HSTS", true)
	config.Security.Headers.EnableXSSFilter = getBoolEnvOrDefault("SECURITY_ENABLE_XSS_FILTER", true)
	config.Security.Headers.EnableNoSniff = getBoolEnvOrDefault("SECURITY_ENABLE_NO_SNIFF", true)
	config.Security.Headers.FrameOptions = getEnvOrDefault("SECURITY_FRAME_OPTIONS", "DENY")

	// Logging config
	config.Logging = LoggingConfig{
		Level:      getEnvOrDefault("LOG_LEVEL", "info"),
		Format:     getEnvOrDefault("LOG_FORMAT", "json"),
		Output:     getEnvOrDefault("LOG_OUTPUT", "stdout"),
		MaxSize:    getIntEnvOrDefault("LOG_MAX_SIZE", 100), // MB
		MaxBackups: getIntEnvOrDefault("LOG_MAX_BACKUPS", 3),
		MaxAge:     getIntEnvOrDefault("LOG_MAX_AGE", 28),
		Compress:   getBoolEnvOrDefault("LOG_COMPRESS", true),
	}

	// Validate required fields
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.JWT.SecretKey == "" {
		return errors.NewAPIError(errors.ErrValidation, "JWT_SECRET_KEY is required")
	}

	if c.Database.Password == "" {
		return errors.NewAPIError(errors.ErrValidation, "DB_PASSWORD is required")
	}

	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnvOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnvOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getBoolEnvOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getStringSliceEnvOrDefault(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
} 