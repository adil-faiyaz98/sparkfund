package ai

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"time"

	"github.com/sparkfund/security-monitoring/internal/models"
)

// SecurityPattern represents a sophisticated security detection pattern
type SecurityPattern struct {
	ID          string
	Name        string
	Description string
	Type        PatternType
	Weight      float64
	Threshold   float64
	Rules       []PatternRule
	Metadata    map[string]interface{}
}

// PatternType defines the type of security pattern
type PatternType string

const (
	PatternTypeBehavioral PatternType = "behavioral"
	PatternTypeNetwork    PatternType = "network"
	PatternTypeSystem     PatternType = "system"
	PatternTypeUser       PatternType = "user"
	PatternTypeData       PatternType = "data"
)

// PatternRule defines a rule within a security pattern
type PatternRule struct {
	ID         string
	Condition  string
	Expression string
	Weight     float64
	Parameters map[string]interface{}
}

// PatternEngine handles sophisticated pattern detection
type PatternEngine struct {
	patterns map[string]*SecurityPattern
	rules    map[string]*regexp.Regexp
}

// NewPatternEngine creates a new pattern engine
func NewPatternEngine() *PatternEngine {
	return &PatternEngine{
		patterns: make(map[string]*SecurityPattern),
		rules:    make(map[string]*regexp.Regexp),
	}
}

// RegisterPattern registers a new security pattern
func (e *PatternEngine) RegisterPattern(pattern *SecurityPattern) error {
	// Validate pattern
	if err := e.validatePattern(pattern); err != nil {
		return err
	}

	// Compile rules
	for _, rule := range pattern.Rules {
		re, err := regexp.Compile(rule.Expression)
		if err != nil {
			return err
		}
		e.rules[rule.ID] = re
	}

	e.patterns[pattern.ID] = pattern
	return nil
}

// AnalyzeSecurityData analyzes security data against registered patterns
func (e *PatternEngine) AnalyzeSecurityData(ctx context.Context, data models.SecurityData) ([]models.SecurityEvent, error) {
	var events []models.SecurityEvent

	for _, pattern := range e.patterns {
		score := e.evaluatePattern(pattern, data)
		if score >= pattern.Threshold {
			event := e.createSecurityEvent(pattern, score, data)
			events = append(events, event)
		}
	}

	return events, nil
}

// evaluatePattern evaluates a pattern against security data
func (e *PatternEngine) evaluatePattern(pattern *SecurityPattern, data models.SecurityData) float64 {
	var totalScore float64
	var totalWeight float64

	for _, rule := range pattern.Rules {
		score := e.evaluateRule(rule, data)
		totalScore += score * rule.Weight
		totalWeight += rule.Weight
	}

	if totalWeight == 0 {
		return 0
	}

	return totalScore / totalWeight
}

// evaluateRule evaluates a single rule against security data
func (e *PatternEngine) evaluateRule(rule PatternRule, data models.SecurityData) float64 {
	re := e.rules[rule.ID]
	if re == nil {
		return 0
	}

	switch rule.Condition {
	case "regex_match":
		if re.MatchString(data.RawData) {
			return 1.0
		}
	case "frequency":
		return e.evaluateFrequency(rule, data)
	case "threshold":
		return e.evaluateThreshold(rule, data)
	case "behavior":
		return e.evaluateBehavior(rule, data)
	default:
		return 0
	}

	return 0
}

// evaluateFrequency evaluates frequency-based rules
func (e *PatternEngine) evaluateFrequency(rule PatternRule, data models.SecurityData) float64 {
	// Implement frequency analysis
	// This could include:
	// - Event frequency over time
	// - Pattern frequency
	// - Anomaly detection
	return 0.5 // Placeholder
}

// evaluateThreshold evaluates threshold-based rules
func (e *PatternEngine) evaluateThreshold(rule PatternRule, data models.SecurityData) float64 {
	// Implement threshold analysis
	// This could include:
	// - Resource usage thresholds
	// - Connection thresholds
	// - Error rate thresholds
	return 0.5 // Placeholder
}

// evaluateBehavior evaluates behavior-based rules
func (e *PatternEngine) evaluateBehavior(rule PatternRule, data models.SecurityData) float64 {
	// Implement behavior analysis
	// This could include:
	// - User behavior patterns
	// - System behavior patterns
	// - Network behavior patterns
	return 0.5 // Placeholder
}

// createSecurityEvent creates a security event from a pattern match
func (e *PatternEngine) createSecurityEvent(pattern *SecurityPattern, score float64, data models.SecurityData) models.SecurityEvent {
	return models.SecurityEvent{
		ID:          generateEventID(),
		Type:        string(pattern.Type),
		Severity:    determineSeverity(score),
		Description: pattern.Description,
		Timestamp:   time.Now(),
		Source:      data.Source,
		Details: map[string]interface{}{
			"pattern_id": pattern.ID,
			"score":      score,
			"metadata":   pattern.Metadata,
		},
	}
}

// validatePattern validates a security pattern
func (e *PatternEngine) validatePattern(pattern *SecurityPattern) error {
	if pattern.ID == "" {
		return fmt.Errorf("pattern ID is required")
	}
	if pattern.Name == "" {
		return fmt.Errorf("pattern name is required")
	}
	if pattern.Type == "" {
		return fmt.Errorf("pattern type is required")
	}
	if pattern.Weight <= 0 {
		return fmt.Errorf("pattern weight must be positive")
	}
	if pattern.Threshold < 0 || pattern.Threshold > 1 {
		return fmt.Errorf("pattern threshold must be between 0 and 1")
	}
	if len(pattern.Rules) == 0 {
		return fmt.Errorf("pattern must have at least one rule")
	}
	return nil
}

// generateEventID generates a unique event ID
func generateEventID() string {
	hash := sha256.New()
	hash.Write([]byte(time.Now().String()))
	return hex.EncodeToString(hash.Sum(nil))
}

// Zero Trust Pattern Definitions
var ZeroTrustPatterns = []*SecurityPattern{
	{
		ID:          "zt_identity_verification",
		Name:        "Identity Verification",
		Description: "Verifies user identity and access patterns",
		Type:        PatternTypeUser,
		Weight:      1.0,
		Threshold:   0.8,
		Rules: []PatternRule{
			{
				ID:         "multi_factor_auth",
				Condition:  "threshold",
				Expression: ".*",
				Weight:     0.4,
				Parameters: map[string]interface{}{
					"required_factors": 2,
				},
			},
			{
				ID:         "access_pattern",
				Condition:  "behavior",
				Expression: ".*",
				Weight:     0.6,
				Parameters: map[string]interface{}{
					"max_failed_attempts": 3,
					"time_window":         "5m",
				},
			},
		},
	},
	{
		ID:          "zt_network_segmentation",
		Name:        "Network Segmentation",
		Description: "Enforces network segmentation and access controls",
		Type:        PatternTypeNetwork,
		Weight:      1.0,
		Threshold:   0.9,
		Rules: []PatternRule{
			{
				ID:         "segment_access",
				Condition:  "threshold",
				Expression: ".*",
				Weight:     0.5,
				Parameters: map[string]interface{}{
					"allowed_segments": []string{"frontend", "backend", "database"},
				},
			},
			{
				ID:         "connection_monitoring",
				Condition:  "behavior",
				Expression: ".*",
				Weight:     0.5,
				Parameters: map[string]interface{}{
					"max_connections": 100,
					"time_window":     "1m",
				},
			},
		},
	},
	{
		ID:          "zt_data_protection",
		Name:        "Data Protection",
		Description: "Ensures data security and access controls",
		Type:        PatternTypeData,
		Weight:      1.0,
		Threshold:   0.95,
		Rules: []PatternRule{
			{
				ID:         "data_access",
				Condition:  "threshold",
				Expression: ".*",
				Weight:     0.4,
				Parameters: map[string]interface{}{
					"encryption_required": true,
					"access_logging":      true,
				},
			},
			{
				ID:         "data_transfer",
				Condition:  "behavior",
				Expression: ".*",
				Weight:     0.6,
				Parameters: map[string]interface{}{
					"allowed_protocols": []string{"HTTPS", "SFTP"},
					"max_transfer_size": "1GB",
				},
			},
		},
	},
}
