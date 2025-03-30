package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LogConfig holds logging middleware configuration
type LogConfig struct {
	ServiceName     string
	LogLevel       string
	ElasticURL     string
	IndexPrefix    string
	EnableElastic  bool
	EnableConsole  bool
}

// DefaultLogConfig returns default logging configuration
func DefaultLogConfig() LogConfig {
	return LogConfig{
		ServiceName:    "kyc-service",
		LogLevel:       "info",
		ElasticURL:     "http://elasticsearch:9200",
		IndexPrefix:    "kyc-service",
		EnableElastic:  true,
		EnableConsole:  true,
	}
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp   time.Time              `json:"timestamp"`
	Service     string                 `json:"service"`
	Level       string                 `json:"level"`
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	Status      int                    `json:"status"`
	Duration    int64                  `json:"duration_ms"`
	ClientIP    string                 `json:"client_ip"`
	UserAgent   string                 `json:"user_agent"`
	RequestID   string                 `json:"request_id"`
	Error       string                 `json:"error,omitempty"`
	Extra       map[string]interface{} `json:"extra,omitempty"`
}

// LoggingMiddleware provides structured logging with ELK stack integration
func LoggingMiddleware(config LogConfig) gin.HandlerFunc {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Set log level
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(level)
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Read request body
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Create log entry
		entry := LogEntry{
			Timestamp: time.Now(),
			Service:   config.ServiceName,
			Level:     getLogLevel(c.Writer.Status()),
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Status:    c.Writer.Status(),
			Duration:  duration.Milliseconds(),
			ClientIP:  c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			RequestID: c.GetString("request_id"),
			Extra:     make(map[string]interface{}),
		}

		// Add request body for non-GET requests
		if c.Request.Method != "GET" && len(bodyBytes) > 0 {
			entry.Extra["request_body"] = string(bodyBytes)
		}

		// Add response body for errors
		if c.Writer.Status() >= 400 {
			entry.Error = c.Errors.String()
		}

		// Log to console if enabled
		if config.EnableConsole {
			logger.WithFields(logrus.Fields{
				"timestamp":   entry.Timestamp,
				"service":     entry.Service,
				"level":       entry.Level,
				"method":      entry.Method,
				"path":        entry.Path,
				"status":      entry.Status,
				"duration_ms": entry.Duration,
				"client_ip":   entry.ClientIP,
				"user_agent":  entry.UserAgent,
				"request_id":  entry.RequestID,
				"error":       entry.Error,
				"extra":       entry.Extra,
			}).Info("HTTP Request")
		}

		// Send to Elasticsearch if enabled
		if config.EnableElastic {
			go sendToElasticsearch(config, entry)
		}
	}
}

// getLogLevel returns the appropriate log level based on HTTP status code
func getLogLevel(status int) string {
	switch {
	case status >= 500:
		return "error"
	case status >= 400:
		return "warn"
	default:
		return "info"
	}
}

// sendToElasticsearch sends the log entry to Elasticsearch
func sendToElasticsearch(config LogConfig, entry LogEntry) {
	// Create index name with date
	indexName := fmt.Sprintf("%s-%s", config.IndexPrefix, entry.Timestamp.Format("2006.01.02"))

	// Convert entry to JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal log entry")
		return
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/%s/_doc", config.ElasticURL, indexName)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.WithError(err).Error("Failed to create Elasticsearch request")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("Failed to send log to Elasticsearch")
		return
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusCreated {
		logrus.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
		}).Error("Failed to send log to Elasticsearch")
	}
} 