package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Config represents the main configuration
type Config struct {
	Server     ServerConfig     `json:"server"`
	Security   SecurityConfig   `json:"security"`
	AI         AIConfig         `json:"ai"`
	Alerts     AlertsConfig     `json:"alerts"`
	RateLimit  RateLimitConfig  `json:"rate_limit"`
	Monitoring MonitoringConfig `json:"monitoring"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Address      string        `json:"address"`
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	BatchSize      int           `json:"batch_size"`
	UpdateInterval time.Duration `json:"update_interval"`
	AlertThreshold float64       `json:"alert_threshold"`
	MaxAlerts      int           `json:"max_alerts"`
	RetentionDays  int           `json:"retention_days"`
	LogLevel       string        `json:"log_level"`
	EnableTLS      bool          `json:"enable_tls"`
	CertFile       string        `json:"cert_file"`
	KeyFile        string        `json:"key_file"`
}

// AIConfig represents AI engine configuration
type AIConfig struct {
	ModelPath      string            `json:"model_path"`
	BatchSize      int               `json:"batch_size"`
	Threshold      float64           `json:"threshold"`
	ModelConfigs   map[string]Config `json:"model_configs"`
	UpdateInterval time.Duration     `json:"update_interval"`
	EnableGPU      bool              `json:"enable_gpu"`
	MaxWorkers     int               `json:"max_workers"`
	CacheSize      int               `json:"cache_size"`
}

// AlertsConfig represents alert configuration
type AlertsConfig struct {
	EnableEmail     bool     `json:"enable_email"`
	EnableSlack     bool     `json:"enable_slack"`
	EnablePagerDuty bool     `json:"enable_pagerduty"`
	EmailRecipients []string `json:"email_recipients"`
	SlackWebhook    string   `json:"slack_webhook"`
	PagerDutyKey    string   `json:"pagerduty_key"`
	MinSeverity     string   `json:"min_severity"`
	AlertTemplate   string   `json:"alert_template"`
}

// RateLimitConfig represents rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int           `json:"requests_per_minute"`
	BurstSize         int           `json:"burst_size"`
	WindowSize        time.Duration `json:"window_size"`
	EnableIPBased     bool          `json:"enable_ip_based"`
	EnableUserBased   bool          `json:"enable_user_based"`
}

// MonitoringConfig represents monitoring configuration
type MonitoringConfig struct {
	EnablePrometheus bool     `json:"enable_prometheus"`
	EnableJaeger     bool     `json:"enable_jaeger"`
	EnableELK        bool     `json:"enable_elk"`
	PrometheusPort   int      `json:"prometheus_port"`
	JaegerEndpoint   string   `json:"jaeger_endpoint"`
	ELKEndpoint      string   `json:"elk_endpoint"`
	MetricsPath      string   `json:"metrics_path"`
	LogLevel         string   `json:"log_level"`
	Tags             []string `json:"tags"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	// Default configuration
	cfg := &Config{
		Server: ServerConfig{
			Address:      ":8080",
			Port:         8080,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		Security: SecurityConfig{
			BatchSize:      100,
			UpdateInterval: 5 * time.Second,
			AlertThreshold: 0.8,
			MaxAlerts:      1000,
			RetentionDays:  30,
			LogLevel:       "info",
			EnableTLS:      false,
		},
		AI: AIConfig{
			ModelPath:      "./models",
			BatchSize:      32,
			Threshold:      0.7,
			UpdateInterval: 1 * time.Hour,
			EnableGPU:      false,
			MaxWorkers:     4,
			CacheSize:      1000,
		},
		Alerts: AlertsConfig{
			EnableEmail:     true,
			EnableSlack:     true,
			EnablePagerDuty: true,
			MinSeverity:     "high",
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: 100,
			BurstSize:         200,
			WindowSize:        1 * time.Minute,
			EnableIPBased:     true,
			EnableUserBased:   true,
		},
		Monitoring: MonitoringConfig{
			EnablePrometheus: true,
			EnableJaeger:     true,
			EnableELK:        true,
			PrometheusPort:   9090,
			MetricsPath:      "/metrics",
			LogLevel:         "info",
		},
	}

	// Load from file if exists
	if err := loadFromFile(cfg); err != nil {
		return nil, err
	}

	// Load from environment variables
	if err := loadFromEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// loadFromFile loads configuration from file
func loadFromFile(cfg *Config) error {
	file, err := os.Open("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(cfg)
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(cfg *Config) error {
	// Server configuration
	if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
		cfg.Server.Address = addr
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.Server.Port = parseInt(port)
	}

	// Security configuration
	if batchSize := os.Getenv("SECURITY_BATCH_SIZE"); batchSize != "" {
		cfg.Security.BatchSize = parseInt(batchSize)
	}
	if threshold := os.Getenv("SECURITY_ALERT_THRESHOLD"); threshold != "" {
		cfg.Security.AlertThreshold = parseFloat(threshold)
	}

	// AI configuration
	if modelPath := os.Getenv("AI_MODEL_PATH"); modelPath != "" {
		cfg.AI.ModelPath = modelPath
	}
	if batchSize := os.Getenv("AI_BATCH_SIZE"); batchSize != "" {
		cfg.AI.BatchSize = parseInt(batchSize)
	}

	// Alerts configuration
	if enableEmail := os.Getenv("ALERTS_ENABLE_EMAIL"); enableEmail != "" {
		cfg.Alerts.EnableEmail = parseBool(enableEmail)
	}
	if enableSlack := os.Getenv("ALERTS_ENABLE_SLACK"); enableSlack != "" {
		cfg.Alerts.EnableSlack = parseBool(enableSlack)
	}

	// Rate limit configuration
	if requests := os.Getenv("RATE_LIMIT_REQUESTS"); requests != "" {
		cfg.RateLimit.RequestsPerMinute = parseInt(requests)
	}
	if burst := os.Getenv("RATE_LIMIT_BURST"); burst != "" {
		cfg.RateLimit.BurstSize = parseInt(burst)
	}

	// Monitoring configuration
	if enablePrometheus := os.Getenv("MONITORING_ENABLE_PROMETHEUS"); enablePrometheus != "" {
		cfg.Monitoring.EnablePrometheus = parseBool(enablePrometheus)
	}
	if enableJaeger := os.Getenv("MONITORING_ENABLE_JAEGER"); enableJaeger != "" {
		cfg.Monitoring.EnableJaeger = parseBool(enableJaeger)
	}

	return nil
}

// Helper functions
func parseInt(s string) int {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	if err != nil {
		return 0
	}
	return i
}

func parseFloat(s string) float64 {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	if err != nil {
		return 0
	}
	return f
}

func parseBool(s string) bool {
	return strings.ToLower(s) == "true"
}
