package models

import (
	"time"
)

// SecurityData represents the input data for security analysis
type SecurityData struct {
	Timestamp time.Time              `json:"timestamp"`
	EventType string                 `json:"event_type"`
	Source    string                 `json:"source"`
	Data      map[string]interface{} `json:"data"`
	Metadata  map[string]string      `json:"metadata"`
	RawData   []byte                 `json:"raw_data,omitempty"`
	Context   map[string]interface{} `json:"context"`
}

// SecurityAnalysis represents the comprehensive security analysis result
type SecurityAnalysis struct {
	Timestamp       time.Time   `json:"timestamp"`
	Threats         []Threat    `json:"threats"`
	Intrusions      []Intrusion `json:"intrusions"`
	Malware         []Malware   `json:"malware"`
	Patterns        []Pattern   `json:"patterns"`
	RiskScore       float64     `json:"risk_score"`
	Confidence      float64     `json:"confidence"`
	Recommendations []string    `json:"recommendations"`
}

// Threat represents a detected security threat
type Threat struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	Description string                 `json:"description"`
	Source      string                 `json:"source"`
	Timestamp   time.Time              `json:"timestamp"`
	Details     map[string]interface{} `json:"details"`
	Confidence  float64                `json:"confidence"`
	Status      string                 `json:"status"`
}

// Intrusion represents a detected intrusion attempt
type Intrusion struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Severity   string                 `json:"severity"`
	SourceIP   string                 `json:"source_ip"`
	Target     string                 `json:"target"`
	Timestamp  time.Time              `json:"timestamp"`
	Details    map[string]interface{} `json:"details"`
	Confidence float64                `json:"confidence"`
	Status     string                 `json:"status"`
}

// Malware represents detected malware
type Malware struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Severity   string                 `json:"severity"`
	Location   string                 `json:"location"`
	Timestamp  time.Time              `json:"timestamp"`
	Details    map[string]interface{} `json:"details"`
	Confidence float64                `json:"confidence"`
	Status     string                 `json:"status"`
}

// Pattern represents a detected security pattern
type Pattern struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Source      string                 `json:"source"`
	Timestamp   time.Time              `json:"timestamp"`
	Details     map[string]interface{} `json:"details"`
	Confidence  float64                `json:"confidence"`
	Status      string                 `json:"status"`
}

// SecurityEvent represents a security event for logging and analysis
type SecurityEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Severity  string                 `json:"severity"`
	Source    string                 `json:"source"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
	Metadata  map[string]string      `json:"metadata"`
	Context   map[string]interface{} `json:"context"`
	Analysis  *SecurityAnalysis      `json:"analysis,omitempty"`
}

// SecurityMetrics represents security monitoring metrics
type SecurityMetrics struct {
	TotalThreats      int64            `json:"total_threats"`
	TotalIntrusions   int64            `json:"total_intrusions"`
	TotalMalware      int64            `json:"total_malware"`
	TotalPatterns     int64            `json:"total_patterns"`
	AverageRiskScore  float64          `json:"average_risk_score"`
	LastUpdate        time.Time        `json:"last_update"`
	MetricsByType     map[string]int64 `json:"metrics_by_type"`
	MetricsBySeverity map[string]int64 `json:"metrics_by_severity"`
}

// SecurityAlert represents a security alert for notification
type SecurityAlert struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Source      string                 `json:"source"`
	Timestamp   time.Time              `json:"timestamp"`
	Details     map[string]interface{} `json:"details"`
	Actions     []string               `json:"actions"`
	Status      string                 `json:"status"`
	AssignedTo  string                 `json:"assigned_to,omitempty"`
	ResolvedAt  *time.Time             `json:"resolved_at,omitempty"`
	Resolution  string                 `json:"resolution,omitempty"`
}
