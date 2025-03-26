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

// Config holds all configuration for the API Gateway
type Config struct {
	Server         ServerConfig
	Proxy          ProxyConfig
	Authentication AuthConfig
	RateLimit      RateLimitConfig
	Metrics        MetricsConfig
	Tracing        TracingConfig
	Secrets        SecretsConfig
	Cache          CacheConfig
	Cors           CorsConfig
}

// ServerConfig holds the configuration for the server
type ServerConfig struct {
	Port            string
	Host            string
	ShutdownTimeout time.Duration
	Environment     string
	Debug           bool
	TLS             bool
	TLSCert         string
	TLSKey          string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
}

// ProxyConfig holds the configuration for service proxying
type ProxyConfig struct {
	Services            map[string]ServiceConfig
	Timeout             time.Duration
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration
	ResponseBufferSize  int
	EnableRetries       bool
	MaxRetries          int
	RetryWaitMin        time.Duration
	RetryWaitMax        time.Duration
}

// ServiceConfig holds configuration for each proxied service
type ServiceConfig struct {
	URL                   string
	HealthCheckPath       string
	Timeout               time.Duration
	CircuitBreakEnabled   bool
	CircuitBreakThreshold int
	CircuitBreakTimeout   time.Duration
	AuthRequired          bool
	LoadBalanceStrategy   string // "round-robin", "least-conn", "ip-hash"
	Endpoints             []string
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret       string
	JWTIssuer       string
	JWTAudience     string
	TokenExpiration time.Duration
	PublicRoutes    []string
	AuthServiceURL  string
	AuthTimeout     time.Duration
	CacheEnabled    bool
	CacheTTL        time.Duration
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled        bool
	RequestsPerMin int
	Burst          int
	Strategy       string // "token-bucket", "fixed-window", "sliding-window"
	IPHeaders      []string
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
}

// SecretsConfig holds configuration for external secrets management
type SecretsConfig struct {
	Provider  string
	Address   string
	Path      string
	TokenPath string
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Enabled         bool
	Type            string // "memory", "redis", "memcached"
	Address         string
	Password        string
	DB              int
	TTL             time.Duration
	CleanupInterval time.Duration
	MaxItems        int
}

// CorsConfig holds CORS configuration
type CorsConfig struct {
	Enabled          bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
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
		serviceEnv := filepath.Join("config", "api-gateway.env")
		if _, err := os.Stat(serviceEnv); err == nil {
			if err := godotenv.Load(serviceEnv); err != nil {
				log.Printf("Warning: Error loading %s file: %v", serviceEnv, err)
			} else {
				log.Printf("Loaded service configuration from %s", serviceEnv)
			}
		}
	}

	// Load service definitions
	serviceConfigs, err := loadServiceConfigs()
	if err != nil {
		log.Printf("Warning: Error loading service configs: %v", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:            getEnvOrDefault("PORT", "8000"),
			Host:            getEnvOrDefault("HOST", "0.0.0.0"),
			ShutdownTimeout: getEnvDurationOrDefault("SHUTDOWN_TIMEOUT", 30*time.Second),
			Environment:     getEnvOrDefault("ENVIRONMENT", "development"),
			Debug:           getEnvOrDefault("DEBUG", "false") == "true",
			TLS:             getEnvOrDefault("TLS_ENABLED", "false") == "true",
			TLSCert:         getEnvOrDefault("TLS_CERT_PATH", ""),
			TLSKey:          getEnvOrDefault("TLS_KEY_PATH", ""),
			ReadTimeout:     getEnvDurationOrDefault("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout:    getEnvDurationOrDefault("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:     getEnvDurationOrDefault("SERVER_IDLE_TIMEOUT", 120*time.Second),
		},
		Proxy: ProxyConfig{
			Services:            serviceConfigs,
			Timeout:             getEnvDurationOrDefault("PROXY_TIMEOUT", 30*time.Second),
			MaxIdleConns:        getEnvIntOrDefault("PROXY_MAX_IDLE_CONNS", 100),
			MaxIdleConnsPerHost: getEnvIntOrDefault("PROXY_MAX_IDLE_CONNS_PER_HOST", 10),
			IdleConnTimeout:     getEnvDurationOrDefault("PROXY_IDLE_CONN_TIMEOUT", 90*time.Second),
			ResponseBufferSize:  getEnvIntOrDefault("PROXY_RESPONSE_BUFFER_SIZE", 32*1024),
			EnableRetries:       getEnvOrDefault("PROXY_ENABLE_RETRIES", "true") == "true",
			MaxRetries:          getEnvIntOrDefault("PROXY_MAX_RETRIES", 3),
			RetryWaitMin:        getEnvDurationOrDefault("PROXY_RETRY_WAIT_MIN", 1*time.Second),
			RetryWaitMax:        getEnvDurationOrDefault("PROXY_RETRY_WAIT_MAX", 5*time.Second),
		},
		Authentication: AuthConfig{
			JWTSecret:       getEnvOrDefault("JWT_SECRET", "your-default-secret-change-in-production"),
			JWTIssuer:       getEnvOrDefault("JWT_ISSUER", "sparkfund"),
			JWTAudience:     getEnvOrDefault("JWT_AUDIENCE", "sparkfund-api"),
			TokenExpiration: getEnvDurationOrDefault("TOKEN_EXPIRATION", 60*time.Minute),
			PublicRoutes:    getEnvSliceOrDefault("PUBLIC_ROUTES", []string{"/health", "/api/v1/auth/login"}),
			AuthServiceURL:  getEnvOrDefault("AUTH_SERVICE_URL", "http://auth-service:8080"),
			AuthTimeout:     getEnvDurationOrDefault("AUTH_TIMEOUT", 5*time.Second),
			CacheEnabled:    getEnvOrDefault("AUTH_CACHE_ENABLED", "true") == "true",
			CacheTTL:        getEnvDurationOrDefault("AUTH_CACHE_TTL", 5*time.Minute),
		},
		RateLimit: RateLimitConfig{
			Enabled:        getEnvOrDefault("RATE_LIMIT_ENABLED", "true") == "true",
			RequestsPerMin: getEnvIntOrDefault("RATE_LIMIT_REQUESTS_PER_MIN", 60),
			Burst:          getEnvIntOrDefault("RATE_LIMIT_BURST", 100),
			Strategy:       getEnvOrDefault("RATE_LIMIT_STRATEGY", "token-bucket"),
			IPHeaders:      getEnvSliceOrDefault("RATE_LIMIT_IP_HEADERS", []string{"X-Forwarded-For", "X-Real-IP"}),
		},
		Metrics: MetricsConfig{
			Enabled:      getEnvOrDefault("METRICS_ENABLED", "true") == "true",
			Port:         getEnvOrDefault("METRICS_PORT", "9090"),
			Path:         getEnvOrDefault("METRICS_PATH", "/metrics"),
			ServiceLabel: getEnvOrDefault("METRICS_SERVICE_LABEL", "api-gateway"),
		},
		Tracing: TracingConfig{
			Enabled:        getEnvOrDefault("TRACING_ENABLED", "true") == "true",
			JaegerEndpoint: getEnvOrDefault("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
			ServiceName:    getEnvOrDefault("JAEGER_SERVICE_NAME", "api-gateway"),
			SamplingRate:   getEnvFloatOrDefault("TRACING_SAMPLING_RATE", 1.0),
		},
		Secrets: SecretsConfig{
			Provider:  getEnvOrDefault("SECRETS_PROVIDER", "vault"),
			Address:   getEnvOrDefault("SECRETS_ADDRESS", "http://localhost:8200"),
			Path:      getEnvOrDefault("SECRETS_PATH", "secret/data"),
			TokenPath: getEnvOrDefault("SECRETS_TOKEN_PATH", "/var/run/secrets/kubernetes.io/serviceaccount/token"),
		},
		Cache: CacheConfig{
			Enabled:         getEnvOrDefault("CACHE_ENABLED", "true") == "true",
			Type:            getEnvOrDefault("CACHE_TYPE", "memory"),
			Address:         getEnvOrDefault("CACHE_ADDRESS", "localhost:6379"),
			Password:        getEnvOrDefault("CACHE_PASSWORD", ""),
			DB:              getEnvIntOrDefault("CACHE_DB", 0),
			TTL:             getEnvDurationOrDefault("CACHE_TTL", 5*time.Minute),
			CleanupInterval: getEnvDurationOrDefault("CACHE_CLEANUP_INTERVAL", 10*time.Minute),
			MaxItems:        getEnvIntOrDefault("CACHE_MAX_ITEMS", 10000),
		},
		Cors: CorsConfig{
			Enabled:          getEnvOrDefault("CORS_ENABLED", "true") == "true",
			AllowOrigins:     getEnvSliceOrDefault("CORS_ALLOW_ORIGINS", []string{"*"}),
			AllowMethods:     getEnvSliceOrDefault("CORS_ALLOW_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowHeaders:     getEnvSliceOrDefault("CORS_ALLOW_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization"}),
			ExposeHeaders:    getEnvSliceOrDefault("CORS_EXPOSE_HEADERS", []string{}),
			AllowCredentials: getEnvOrDefault("CORS_ALLOW_CREDENTIALS", "false") == "true",
			MaxAge:           getEnvDurationOrDefault("CORS_MAX_AGE", 12*time.Hour),
		},
	}

	// Validate configuration
	if err := validateGatewayConfig(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// loadServiceConfigs loads service configurations from environment variables
func loadServiceConfigs() (map[string]ServiceConfig, error) {
	services := make(map[string]ServiceConfig)

	// Get the list of services to configure
	servicesList := getEnvSliceOrDefault("GATEWAY_SERVICES", []string{
		"auth-service",
		"kyc-service",
		"aml-service",
		"fraud-detection",
		"credit-scoring",
		"risk-management",
		"notification",
	})

	for _, service := range servicesList {
		serviceName := strings.ReplaceAll(service, "-", "_")
		urlEnvVar := fmt.Sprintf("%s_URL", strings.ToUpper(serviceName))
		healthCheckEnvVar := fmt.Sprintf("%s_HEALTH_CHECK_PATH", strings.ToUpper(serviceName))
		timeoutEnvVar := fmt.Sprintf("%s_TIMEOUT", strings.ToUpper(serviceName))
		circuitBreakerEnvVar := fmt.Sprintf("%s_CIRCUIT_BREAKER_ENABLED", strings.ToUpper(serviceName))
		circuitBreakerThresholdEnvVar := fmt.Sprintf("%s_CIRCUIT_BREAKER_THRESHOLD", strings.ToUpper(serviceName))
		circuitBreakerTimeoutEnvVar := fmt.Sprintf("%s_CIRCUIT_BREAKER_TIMEOUT", strings.ToUpper(serviceName))
		authRequiredEnvVar := fmt.Sprintf("%s_AUTH_REQUIRED", strings.ToUpper(serviceName))
		loadBalanceStrategyEnvVar := fmt.Sprintf("%s_LOAD_BALANCE_STRATEGY", strings.ToUpper(serviceName))
		endpointsEnvVar := fmt.Sprintf("%s_ENDPOINTS", strings.ToUpper(serviceName))

		services[service] = ServiceConfig{
			URL:                   getEnvOrDefault(urlEnvVar, fmt.Sprintf("http://%s:8080", service)),
			HealthCheckPath:       getEnvOrDefault(healthCheckEnvVar, "/health"),
			Timeout:               getEnvDurationOrDefault(timeoutEnvVar, 30*time.Second),
			CircuitBreakEnabled:   getEnvOrDefault(circuitBreakerEnvVar, "true") == "true",
			CircuitBreakThreshold: getEnvIntOrDefault(circuitBreakerThresholdEnvVar, 5),
			CircuitBreakTimeout:   getEnvDurationOrDefault(circuitBreakerTimeoutEnvVar, 10*time.Second),
			AuthRequired:          getEnvOrDefault(authRequiredEnvVar, "true") == "true",
			LoadBalanceStrategy:   getEnvOrDefault(loadBalanceStrategyEnvVar, "round-robin"),
			Endpoints:             getEnvSliceOrDefault(endpointsEnvVar, []string{}),
		}
	}

	return services, nil
}

// validateGatewayConfig performs API Gateway-specific validations
func validateGatewayConfig(cfg *Config) error {
	// Check for JWT secret in production
	if cfg.Server.Environment == "production" {
		if cfg.Authentication.JWTSecret == "your-default-secret-change-in-production" {
			return errors.New("production environment detected but using default JWT secret")
		}

		if cfg.Server.TLS && (cfg.Server.TLSCert == "" || cfg.Server.TLSKey == "") {
			return errors.New("TLS is enabled but certificate or key path is missing")
		}
	}

	// Ensure we have at least one service defined if not in development mode
	if cfg.Server.Environment != "development" && len(cfg.Proxy.Services) == 0 {
		return errors.New("no API services configured")
	}

	return nil
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
