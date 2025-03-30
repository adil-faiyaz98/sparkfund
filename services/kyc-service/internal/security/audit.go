package security

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// AuditEvent represents a security audit event
type AuditEvent struct {
	Timestamp   time.Time              `json:"timestamp"`
	EventType   string                 `json:"event_type"`
	UserID      string                 `json:"user_id"`
	IPAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	Resource    string                 `json:"resource"`
	Action      string                 `json:"action"`
	Status      string                 `json:"status"`
	Details     map[string]interface{} `json:"details"`
	RequestID   string                 `json:"request_id"`
	SessionID   string                 `json:"session_id"`
}

// AuditConfig holds audit logging configuration
type AuditConfig struct {
	LogDir      string
	MaxFileSize int64
	MaxFiles    int
	LogLevel    string
}

// DefaultAuditConfig returns default audit configuration
func DefaultAuditConfig() AuditConfig {
	return AuditConfig{
		LogDir:      "logs/audit",
		MaxFileSize: 100 * 1024 * 1024, // 100MB
		MaxFiles:    10,
		LogLevel:    "info",
	}
}

// AuditLogger handles security audit logging
type AuditLogger struct {
	config AuditConfig
	file   *os.File
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(config AuditConfig) (*AuditLogger, error) {
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return nil, err
	}

	logger := &AuditLogger{
		config: config,
	}

	if err := logger.rotateLog(); err != nil {
		return nil, err
	}

	return logger, nil
}

// LogEvent logs a security audit event
func (l *AuditLogger) LogEvent(event AuditEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if _, err := l.file.Write(append(data, '\n')); err != nil {
		return err
	}

	// Check if we need to rotate the log file
	if info, err := l.file.Stat(); err == nil && info.Size() >= l.config.MaxFileSize {
		return l.rotateLog()
	}

	return nil
}

// rotateLog rotates the log file
func (l *AuditLogger) rotateLog() error {
	if l.file != nil {
		l.file.Close()
	}

	timestamp := time.Now().Format("2006-01-02-15-04-05")
	filename := filepath.Join(l.config.LogDir, fmt.Sprintf("audit-%s.log", timestamp))

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	l.file = file

	// Clean up old log files
	return l.cleanupOldLogs()
}

// cleanupOldLogs removes old log files
func (l *AuditLogger) cleanupOldLogs() error {
	files, err := filepath.Glob(filepath.Join(l.config.LogDir, "audit-*.log"))
	if err != nil {
		return err
	}

	if len(files) > l.config.MaxFiles {
		// Sort files by modification time
		// Remove oldest files
		for i := 0; i < len(files)-l.config.MaxFiles; i++ {
			if err := os.Remove(files[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

// AuditMiddleware provides audit logging middleware
func AuditMiddleware(logger *AuditLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Create audit event
		event := AuditEvent{
			Timestamp: time.Now(),
			EventType: "http_request",
			IPAddress: c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
			Resource:  c.Request.URL.Path,
			Action:    c.Request.Method,
			Status:    fmt.Sprintf("%d", c.Writer.Status()),
			Details: map[string]interface{}{
				"duration_ms": time.Since(start).Milliseconds(),
				"headers":     c.Request.Header,
			},
			RequestID: c.GetString("request_id"),
			SessionID: c.GetString("session_id"),
		}

		// Get user ID if available
		if userID, exists := c.Get("user_id"); exists {
			event.UserID = userID.(string)
		}

		// Log the event
		if err := logger.LogEvent(event); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Failed to log audit event: %v\n", err)
		}
	}
}

// LogSecurityEvent logs a security-related event
func (l *AuditLogger) LogSecurityEvent(eventType string, userID string, details map[string]interface{}) error {
	event := AuditEvent{
		Timestamp: time.Now(),
		EventType: eventType,
		UserID:    userID,
		Details:   details,
	}

	return l.LogEvent(event)
}

// LogAuthenticationEvent logs an authentication event
func (l *AuditLogger) LogAuthenticationEvent(userID string, success bool, details map[string]interface{}) error {
	event := AuditEvent{
		Timestamp: time.Now(),
		EventType: "authentication",
		UserID:    userID,
		Status:    map[bool]string{true: "success", false: "failure"}[success],
		Details:   details,
	}

	return l.LogEvent(event)
}

// LogAuthorizationEvent logs an authorization event
func (l *AuditLogger) LogAuthorizationEvent(userID string, resource string, action string, success bool, details map[string]interface{}) error {
	event := AuditEvent{
		Timestamp: time.Now(),
		EventType: "authorization",
		UserID:    userID,
		Resource:  resource,
		Action:    action,
		Status:    map[bool]string{true: "success", false: "failure"}[success],
		Details:   details,
	}

	return l.LogEvent(event)
}

// LogDataAccessEvent logs a data access event
func (l *AuditLogger) LogDataAccessEvent(userID string, resource string, action string, details map[string]interface{}) error {
	event := AuditEvent{
		Timestamp: time.Now(),
		EventType: "data_access",
		UserID:    userID,
		Resource:  resource,
		Action:    action,
		Status:    "success",
		Details:   details,
	}

	return l.LogEvent(event)
}

// LogConfigurationChangeEvent logs a configuration change event
func (l *AuditLogger) LogConfigurationChangeEvent(userID string, configKey string, oldValue interface{}, newValue interface{}, details map[string]interface{}) error {
	if details == nil {
		details = make(map[string]interface{})
	}
	details["config_key"] = configKey
	details["old_value"] = oldValue
	details["new_value"] = newValue

	event := AuditEvent{
		Timestamp: time.Now(),
		EventType: "configuration_change",
		UserID:    userID,
		Resource:  "configuration",
		Action:    "update",
		Status:    "success",
		Details:   details,
	}

	return l.LogEvent(event)
} 