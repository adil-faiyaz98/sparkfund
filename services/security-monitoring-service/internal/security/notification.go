package security

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// NotificationService handles transaction notifications
type NotificationService struct {
	config *NotificationConfig
	mu     sync.RWMutex
}

// NotificationConfig defines notification configuration
type NotificationConfig struct {
	// Notification channels
	Channels struct {
		Email struct {
			Enabled     bool
			FromAddress string
			SMTPConfig  SMTPConfig
		}
		SMS struct {
			Enabled   bool
			Provider  string
			APIKey    string
			APISecret string
		}
		Phone struct {
			Enabled   bool
			Provider  string
			APIKey    string
			APISecret string
		}
	}

	// Notification templates
	Templates struct {
		HighRiskEmail      string
		HighRiskSMS        string
		HighRiskPhone      string
		BlockedEmail       string
		BlockedSMS         string
		BlockedPhone       string
		UnusualEmail       string
		UnusualSMS         string
		UnusualPhone       string
		LimitExceededEmail string
		LimitExceededSMS   string
		LimitExceededPhone string
	}

	// Notification settings
	Settings struct {
		MaxRetries    int
		RetryInterval time.Duration
		MaxConcurrent int
		QueueSize     int
		BatchSize     int
		BatchInterval time.Duration
	}
}

// SMTPConfig defines SMTP configuration
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	UseTLS   bool
}

// NewNotificationService creates a new notification service
func NewNotificationService(config NotificationConfig) *NotificationService {
	return &NotificationService{
		config: &config,
	}
}

// NotifyTransactionAlert sends notifications for transaction alerts
func (n *NotificationService) NotifyTransactionAlert(ctx context.Context, alert *TransactionAlert) error {
	n.mu.RLock()
	defer n.mu.RUnlock()

	// Create notification channels
	channels := make([]chan error, 0)

	// Send email notification if enabled
	if n.config.Channels.Email.Enabled {
		emailChan := make(chan error, 1)
		go func() {
			emailChan <- n.sendEmailNotification(ctx, alert)
		}()
		channels = append(channels, emailChan)
	}

	// Send SMS notification if enabled
	if n.config.Channels.SMS.Enabled {
		smsChan := make(chan error, 1)
		go func() {
			smsChan <- n.sendSMSNotification(ctx, alert)
		}()
		channels = append(channels, smsChan)
	}

	// Send phone notification if enabled
	if n.config.Channels.Phone.Enabled {
		phoneChan := make(chan error, 1)
		go func() {
			phoneChan <- n.sendPhoneNotification(ctx, alert)
		}()
		channels = append(channels, phoneChan)
	}

	// Wait for all notifications to complete
	var errors []error
	for _, ch := range channels {
		if err := <-ch; err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to send notifications: %v", errors)
	}

	return nil
}

// sendEmailNotification sends an email notification
func (n *NotificationService) sendEmailNotification(ctx context.Context, alert *TransactionAlert) error {
	// Get appropriate template based on alert type
	template := n.getEmailTemplate(alert.Type)

	// Format message with alert details
	message := fmt.Sprintf(template,
		alert.UserName,
		alert.TransactionID,
		alert.Amount,
		alert.Currency,
		alert.Timestamp.Format(time.RFC3339),
		alert.Reason,
	)

	// Send email using configured SMTP settings
	return n.sendEmail(ctx, alert.UserEmail, message)
}

// sendSMSNotification sends an SMS notification
func (n *NotificationService) sendSMSNotification(ctx context.Context, alert *TransactionAlert) error {
	// Get appropriate template based on alert type
	template := n.getSMSTemplate(alert.Type)

	// Format message with alert details
	message := fmt.Sprintf(template,
		alert.UserName,
		alert.TransactionID,
		alert.Amount,
		alert.Currency,
		alert.Reason,
	)

	// Send SMS using configured provider
	return n.sendSMS(ctx, alert.UserPhone, message)
}

// sendPhoneNotification sends a phone notification
func (n *NotificationService) sendPhoneNotification(ctx context.Context, alert *TransactionAlert) error {
	// Get appropriate template based on alert type
	template := n.getPhoneTemplate(alert.Type)

	// Format message with alert details
	message := fmt.Sprintf(template,
		alert.UserName,
		alert.TransactionID,
		alert.Amount,
		alert.Currency,
		alert.Reason,
	)

	// Send phone notification using configured provider
	return n.sendPhone(ctx, alert.UserPhone, message)
}

// Helper functions
func (n *NotificationService) getEmailTemplate(alertType string) string {
	switch alertType {
	case "high_risk":
		return n.config.Templates.HighRiskEmail
	case "blocked":
		return n.config.Templates.BlockedEmail
	case "unusual":
		return n.config.Templates.UnusualEmail
	case "limit_exceeded":
		return n.config.Templates.LimitExceededEmail
	default:
		return n.config.Templates.HighRiskEmail
	}
}

func (n *NotificationService) getSMSTemplate(alertType string) string {
	switch alertType {
	case "high_risk":
		return n.config.Templates.HighRiskSMS
	case "blocked":
		return n.config.Templates.BlockedSMS
	case "unusual":
		return n.config.Templates.UnusualSMS
	case "limit_exceeded":
		return n.config.Templates.LimitExceededSMS
	default:
		return n.config.Templates.HighRiskSMS
	}
}

func (n *NotificationService) getPhoneTemplate(alertType string) string {
	switch alertType {
	case "high_risk":
		return n.config.Templates.HighRiskPhone
	case "blocked":
		return n.config.Templates.BlockedPhone
	case "unusual":
		return n.config.Templates.UnusualPhone
	case "limit_exceeded":
		return n.config.Templates.LimitExceededPhone
	default:
		return n.config.Templates.HighRiskPhone
	}
}

func (n *NotificationService) sendEmail(ctx context.Context, to, message string) error {
	// Implement email sending logic
	// This should use the configured SMTP settings
	return nil
}

func (n *NotificationService) sendSMS(ctx context.Context, to, message string) error {
	// Implement SMS sending logic
	// This should use the configured SMS provider
	return nil
}

func (n *NotificationService) sendPhone(ctx context.Context, to, message string) error {
	// Implement phone notification logic
	// This should use the configured phone provider
	return nil
}

// TransactionAlert represents a transaction alert
type TransactionAlert struct {
	Type          string
	UserID        string
	UserName      string
	UserEmail     string
	UserPhone     string
	TransactionID string
	Amount        float64
	Currency      string
	Timestamp     time.Time
	Reason        string
	Details       map[string]interface{}
}
