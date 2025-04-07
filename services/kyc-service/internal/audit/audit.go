package audit

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// ActionType represents the type of action being audited
type ActionType string

const (
	// ActionCreate represents a create action
	ActionCreate ActionType = "CREATE"
	
	// ActionRead represents a read action
	ActionRead ActionType = "READ"
	
	// ActionUpdate represents an update action
	ActionUpdate ActionType = "UPDATE"
	
	// ActionDelete represents a delete action
	ActionDelete ActionType = "DELETE"
	
	// ActionLogin represents a login action
	ActionLogin ActionType = "LOGIN"
	
	// ActionLogout represents a logout action
	ActionLogout ActionType = "LOGOUT"
	
	// ActionVerify represents a verification action
	ActionVerify ActionType = "VERIFY"
	
	// ActionApprove represents an approval action
	ActionApprove ActionType = "APPROVE"
	
	// ActionReject represents a rejection action
	ActionReject ActionType = "REJECT"
	
	// ActionUpload represents an upload action
	ActionUpload ActionType = "UPLOAD"
	
	// ActionDownload represents a download action
	ActionDownload ActionType = "DOWNLOAD"
	
	// ActionAnalyze represents an analysis action
	ActionAnalyze ActionType = "ANALYZE"
)

// ResourceType represents the type of resource being audited
type ResourceType string

const (
	// ResourceUser represents a user resource
	ResourceUser ResourceType = "USER"
	
	// ResourceDocument represents a document resource
	ResourceDocument ResourceType = "DOCUMENT"
	
	// ResourceVerification represents a verification resource
	ResourceVerification ResourceType = "VERIFICATION"
	
	// ResourceKYC represents a KYC resource
	ResourceKYC ResourceType = "KYC"
	
	// ResourceAnalysis represents an analysis resource
	ResourceAnalysis ResourceType = "ANALYSIS"
	
	// ResourceFaceMatch represents a face match resource
	ResourceFaceMatch ResourceType = "FACE_MATCH"
	
	// ResourceRiskAnalysis represents a risk analysis resource
	ResourceRiskAnalysis ResourceType = "RISK_ANALYSIS"
	
	// ResourceAnomalyDetection represents an anomaly detection resource
	ResourceAnomalyDetection ResourceType = "ANOMALY_DETECTION"
)

// StatusType represents the status of the audited action
type StatusType string

const (
	// StatusSuccess represents a successful action
	StatusSuccess StatusType = "SUCCESS"
	
	// StatusFailure represents a failed action
	StatusFailure StatusType = "FAILURE"
	
	// StatusError represents an action that resulted in an error
	StatusError StatusType = "ERROR"
	
	// StatusDenied represents an action that was denied
	StatusDenied StatusType = "DENIED"
)

// Event represents an audit event
type Event struct {
	ID            uuid.UUID    `json:"id"`
	Timestamp     time.Time    `json:"timestamp"`
	UserID        string       `json:"user_id"`
	Action        ActionType   `json:"action"`
	Resource      ResourceType `json:"resource"`
	ResourceID    string       `json:"resource_id"`
	Status        StatusType   `json:"status"`
	ClientIP      string       `json:"client_ip"`
	UserAgent     string       `json:"user_agent"`
	RequestID     string       `json:"request_id"`
	RequestMethod string       `json:"request_method"`
	RequestPath   string       `json:"request_path"`
	RequestParams interface{}  `json:"request_params,omitempty"`
	ResponseCode  int          `json:"response_code"`
	ErrorMessage  string       `json:"error_message,omitempty"`
	Changes       interface{}  `json:"changes,omitempty"`
	Metadata      interface{}  `json:"metadata,omitempty"`
}

// Logger is the interface for audit logging
type Logger interface {
	// Log logs an audit event
	Log(ctx context.Context, event *Event) error
	
	// LogWithChanges logs an audit event with changes
	LogWithChanges(ctx context.Context, event *Event, before, after interface{}) error
	
	// Search searches for audit events
	Search(ctx context.Context, query map[string]interface{}, limit, offset int) ([]*Event, int, error)
	
	// GetByID gets an audit event by ID
	GetByID(ctx context.Context, id uuid.UUID) (*Event, error)
	
	// GetByResourceID gets audit events by resource ID
	GetByResourceID(ctx context.Context, resourceType ResourceType, resourceID string, limit, offset int) ([]*Event, int, error)
	
	// GetByUserID gets audit events by user ID
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*Event, int, error)
}

// Config holds configuration for the audit logger
type Config struct {
	Enabled       bool
	LogToDatabase bool
	LogToFile     bool
	LogToConsole  bool
	FilePath      string
	DatabaseURL   string
	LogLevel      string
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Enabled:       true,
		LogToDatabase: true,
		LogToFile:     false,
		LogToConsole:  true,
		FilePath:      "audit.log",
		DatabaseURL:   "",
		LogLevel:      "info",
	}
}

// logrusLogger implements the Logger interface using logrus
type logrusLogger struct {
	logger *logrus.Logger
	config Config
	repo   Repository
}

// NewLogrusLogger creates a new logrus-based audit logger
func NewLogrusLogger(logger *logrus.Logger, config Config, repo Repository) Logger {
	return &logrusLogger{
		logger: logger,
		config: config,
		repo:   repo,
	}
}

// Log logs an audit event
func (l *logrusLogger) Log(ctx context.Context, event *Event) error {
	if !l.config.Enabled {
		return nil
	}

	// Generate ID if not provided
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}

	// Set timestamp if not provided
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	// Log to database
	if l.config.LogToDatabase && l.repo != nil {
		if err := l.repo.Create(ctx, event); err != nil {
			l.logger.WithError(err).Error("Failed to log audit event to database")
		}
	}

	// Log to console/file
	if l.config.LogToConsole || l.config.LogToFile {
		entry := l.logger.WithFields(logrus.Fields{
			"audit_id":        event.ID.String(),
			"timestamp":       event.Timestamp,
			"user_id":         event.UserID,
			"action":          event.Action,
			"resource":        event.Resource,
			"resource_id":     event.ResourceID,
			"status":          event.Status,
			"client_ip":       event.ClientIP,
			"user_agent":      event.UserAgent,
			"request_id":      event.RequestID,
			"request_method":  event.RequestMethod,
			"request_path":    event.RequestPath,
			"response_code":   event.ResponseCode,
			"error_message":   event.ErrorMessage,
		})

		// Add request params if provided
		if event.RequestParams != nil {
			params, err := json.Marshal(event.RequestParams)
			if err == nil {
				entry = entry.WithField("request_params", string(params))
			}
		}

		// Add metadata if provided
		if event.Metadata != nil {
			metadata, err := json.Marshal(event.Metadata)
			if err == nil {
				entry = entry.WithField("metadata", string(metadata))
			}
		}

		// Log with appropriate level based on status
		switch event.Status {
		case StatusSuccess:
			entry.Info("Audit event")
		case StatusFailure:
			entry.Warn("Audit event")
		case StatusError:
			entry.Error("Audit event")
		case StatusDenied:
			entry.Warn("Audit event")
		default:
			entry.Info("Audit event")
		}
	}

	return nil
}

// LogWithChanges logs an audit event with changes
func (l *logrusLogger) LogWithChanges(ctx context.Context, event *Event, before, after interface{}) error {
	if !l.config.Enabled {
		return nil
	}

	// Calculate changes
	changes, err := calculateChanges(before, after)
	if err != nil {
		l.logger.WithError(err).Error("Failed to calculate changes for audit event")
	} else {
		event.Changes = changes
	}

	return l.Log(ctx, event)
}

// Search searches for audit events
func (l *logrusLogger) Search(ctx context.Context, query map[string]interface{}, limit, offset int) ([]*Event, int, error) {
	if !l.config.Enabled || l.repo == nil {
		return nil, 0, nil
	}

	return l.repo.Search(ctx, query, limit, offset)
}

// GetByID gets an audit event by ID
func (l *logrusLogger) GetByID(ctx context.Context, id uuid.UUID) (*Event, error) {
	if !l.config.Enabled || l.repo == nil {
		return nil, nil
	}

	return l.repo.GetByID(ctx, id)
}

// GetByResourceID gets audit events by resource ID
func (l *logrusLogger) GetByResourceID(ctx context.Context, resourceType ResourceType, resourceID string, limit, offset int) ([]*Event, int, error) {
	if !l.config.Enabled || l.repo == nil {
		return nil, 0, nil
	}

	return l.repo.GetByResourceID(ctx, resourceType, resourceID, limit, offset)
}

// GetByUserID gets audit events by user ID
func (l *logrusLogger) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*Event, int, error) {
	if !l.config.Enabled || l.repo == nil {
		return nil, 0, nil
	}

	return l.repo.GetByUserID(ctx, userID, limit, offset)
}

// calculateChanges calculates the changes between two objects
func calculateChanges(before, after interface{}) (map[string]interface{}, error) {
	changes := make(map[string]interface{})

	// Convert to JSON and back to map to normalize the objects
	beforeJSON, err := json.Marshal(before)
	if err != nil {
		return nil, err
	}

	afterJSON, err := json.Marshal(after)
	if err != nil {
		return nil, err
	}

	var beforeMap map[string]interface{}
	var afterMap map[string]interface{}

	if err := json.Unmarshal(beforeJSON, &beforeMap); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(afterJSON, &afterMap); err != nil {
		return nil, err
	}

	// Compare the maps
	for key, afterValue := range afterMap {
		beforeValue, exists := beforeMap[key]
		if !exists {
			// New field
			changes[key] = map[string]interface{}{
				"old": nil,
				"new": afterValue,
			}
		} else if !jsonEqual(beforeValue, afterValue) {
			// Changed field
			changes[key] = map[string]interface{}{
				"old": beforeValue,
				"new": afterValue,
			}
		}
	}

	// Check for deleted fields
	for key, beforeValue := range beforeMap {
		if _, exists := afterMap[key]; !exists {
			changes[key] = map[string]interface{}{
				"old": beforeValue,
				"new": nil,
			}
		}
	}

	return changes, nil
}

// jsonEqual compares two JSON-serializable values for equality
func jsonEqual(a, b interface{}) bool {
	aJSON, err := json.Marshal(a)
	if err != nil {
		return false
	}

	bJSON, err := json.Marshal(b)
	if err != nil {
		return false
	}

	return string(aJSON) == string(bJSON)
}
