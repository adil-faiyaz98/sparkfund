package monitoring

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sparkfund/security-monitoring/internal/models"
)

// AlertManager represents the alert management system
type AlertManager struct {
	config   AlertsConfig
	metrics  *Metrics
	alerts   []*models.SecurityAlert
	mu       sync.RWMutex
	stopChan chan struct{}
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
	RetentionDays   int      `json:"retention_days"`
	MaxAlerts       int      `json:"max_alerts"`
}

// NewAlertManager creates a new alert manager instance
func NewAlertManager(cfg AlertsConfig) *AlertManager {
	return &AlertManager{
		config:   cfg,
		alerts:   make([]*models.SecurityAlert, 0),
		stopChan: make(chan struct{}),
	}
}

// Start begins the alert management process
func (am *AlertManager) Start(ctx context.Context) error {
	// Start alert processing goroutine
	go am.processAlerts(ctx)

	// Start alert cleanup goroutine
	go am.cleanupAlerts(ctx)

	return nil
}

// Stop gracefully stops the alert management process
func (am *AlertManager) Stop() {
	close(am.stopChan)
}

// AddAlert adds a new security alert
func (am *AlertManager) AddAlert(alert *models.SecurityAlert) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	// Check severity threshold
	if !am.meetsSeverityThreshold(alert.Severity) {
		return nil
	}

	// Add alert to list
	am.alerts = append(am.alerts, alert)

	// Trim alerts if exceeding max
	if len(am.alerts) > am.config.MaxAlerts {
		am.alerts = am.alerts[len(am.alerts)-am.config.MaxAlerts:]
	}

	// Update metrics
	if am.metrics != nil {
		am.metrics.alertCount.WithLabelValues(alert.Type, alert.Severity).Inc()
	}

	return nil
}

// GetAlerts returns all current alerts
func (am *AlertManager) GetAlerts() []*models.SecurityAlert {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return am.alerts
}

// GetAlertsBySeverity returns alerts filtered by severity
func (am *AlertManager) GetAlertsBySeverity(severity string) []*models.SecurityAlert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var filtered []*models.SecurityAlert
	for _, alert := range am.alerts {
		if alert.Severity == severity {
			filtered = append(filtered, alert)
		}
	}
	return filtered
}

// GetAlertsByType returns alerts filtered by type
func (am *AlertManager) GetAlertsByType(alertType string) []*models.SecurityAlert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var filtered []*models.SecurityAlert
	for _, alert := range am.alerts {
		if alert.Type == alertType {
			filtered = append(filtered, alert)
		}
	}
	return filtered
}

// UpdateAlertStatus updates the status of an alert
func (am *AlertManager) UpdateAlertStatus(alertID string, status string, resolution string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	for _, alert := range am.alerts {
		if alert.ID == alertID {
			alert.Status = status
			if status == "resolved" {
				now := time.Now()
				alert.ResolvedAt = &now
				alert.Resolution = resolution
			}
			return nil
		}
	}
	return fmt.Errorf("alert not found: %s", alertID)
}

// processAlerts processes alerts and sends notifications
func (am *AlertManager) processAlerts(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-am.stopChan:
			return
		default:
			am.mu.RLock()
			for _, alert := range am.alerts {
				if alert.Status == "new" {
					// Send notifications
					if err := am.sendNotifications(alert); err != nil {
						// Log error but continue processing
					}
					// Update alert status
					alert.Status = "notified"
				}
			}
			am.mu.RUnlock()
			time.Sleep(time.Second)
		}
	}
}

// cleanupAlerts removes old alerts
func (am *AlertManager) cleanupAlerts(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			am.mu.Lock()
			var active []*models.SecurityAlert
			cutoff := time.Now().AddDate(0, 0, -am.config.RetentionDays)

			for _, alert := range am.alerts {
				if alert.Timestamp.After(cutoff) {
					active = append(active, alert)
				}
			}
			am.alerts = active
			am.mu.Unlock()

		case <-ctx.Done():
			return
		case <-am.stopChan:
			return
		}
	}
}

// sendNotifications sends alert notifications through configured channels
func (am *AlertManager) sendNotifications(alert *models.SecurityAlert) error {
	// Prepare notification content
	content := am.prepareNotificationContent(alert)

	// Send email notification
	if am.config.EnableEmail {
		if err := am.sendEmailNotification(alert, content); err != nil {
			// Log error but continue with other notifications
		}
	}

	// Send Slack notification
	if am.config.EnableSlack {
		if err := am.sendSlackNotification(alert, content); err != nil {
			// Log error but continue with other notifications
		}
	}

	// Send PagerDuty notification
	if am.config.EnablePagerDuty {
		if err := am.sendPagerDutyNotification(alert, content); err != nil {
			// Log error but continue with other notifications
		}
	}

	return nil
}

// prepareNotificationContent prepares the content for notifications
func (am *AlertManager) prepareNotificationContent(alert *models.SecurityAlert) string {
	if am.config.AlertTemplate != "" {
		// Use custom template if provided
		return fmt.Sprintf(am.config.AlertTemplate,
			alert.Title,
			alert.Description,
			alert.Severity,
			alert.Type,
			alert.Source,
			alert.Timestamp.Format(time.RFC3339),
		)
	}

	// Default template
	return fmt.Sprintf("Security Alert: %s\nSeverity: %s\nType: %s\nSource: %s\nTime: %s\nDescription: %s",
		alert.Title,
		alert.Severity,
		alert.Type,
		alert.Source,
		alert.Timestamp.Format(time.RFC3339),
		alert.Description,
	)
}

// sendEmailNotification sends an email notification
func (am *AlertManager) sendEmailNotification(alert *models.SecurityAlert, content string) error {
	// Implementation for sending email notifications
	return nil
}

// sendSlackNotification sends a Slack notification
func (am *AlertManager) sendSlackNotification(alert *models.SecurityAlert, content string) error {
	// Implementation for sending Slack notifications
	return nil
}

// sendPagerDutyNotification sends a PagerDuty notification
func (am *AlertManager) sendPagerDutyNotification(alert *models.SecurityAlert, content string) error {
	// Implementation for sending PagerDuty notifications
	return nil
}

// meetsSeverityThreshold checks if the alert meets the minimum severity threshold
func (am *AlertManager) meetsSeverityThreshold(severity string) bool {
	severityLevels := map[string]int{
		"critical": 4,
		"high":     3,
		"medium":   2,
		"low":      1,
	}

	alertLevel := severityLevels[severity]
	thresholdLevel := severityLevels[am.config.MinSeverity]

	return alertLevel >= thresholdLevel
}
