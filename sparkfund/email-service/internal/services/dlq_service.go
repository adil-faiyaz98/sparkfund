package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/adil-faiyaz98/sparkfund/email-service/internal/config"
	"github.com/adil-faiyaz98/sparkfund/email-service/internal/kafka"
	"go.uber.org/zap"
)

// FailedEmail represents an email that failed to send
type FailedEmail struct {
	To          []string          `json:"to"`
	Subject     string            `json:"subject"`
	Body        string            `json:"body"`
	Attachments map[string][]byte `json:"attachments,omitempty"`
	Error       string            `json:"error"`
	Attempts    int               `json:"attempts"`
	LastAttempt time.Time         `json:"last_attempt"`
}

// DLQService handles failed email processing
type DLQService struct {
	logger   *zap.Logger
	config   *config.Config
	producer *kafka.Producer
}

// NewDLQService creates a new DLQ service instance
func NewDLQService(logger *zap.Logger, cfg *config.Config) (*DLQService, error) {
	producer, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &DLQService{
		logger:   logger,
		config:   cfg,
		producer: producer,
	}, nil
}

// HandleFailedEmail processes a failed email and sends it to the DLQ
func (s *DLQService) HandleFailedEmail(email *FailedEmail) error {
	// Log the failure
	s.logger.Error("Email failed to send",
		zap.Strings("to", email.To),
		zap.String("subject", email.Subject),
		zap.Int("attempts", email.Attempts),
		zap.String("error", email.Error),
		zap.Time("last_attempt", email.LastAttempt))

	// Marshal the failed email
	data, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("failed to marshal failed email: %w", err)
	}

	// Send to DLQ topic
	if err := s.producer.SendMessage(context.Background(), "email-dlq", email.Subject, data); err != nil {
		return fmt.Errorf("failed to publish to DLQ: %w", err)
	}

	// Log successful DLQ processing
	s.logger.Info("Failed email sent to DLQ",
		zap.Strings("to", email.To),
		zap.String("subject", email.Subject))

	return nil
}

// RetryFailedEmails attempts to retry emails from the DLQ
func (s *DLQService) RetryFailedEmails(consumer *kafka.Consumer, smtpService *SMTPService) error {
	// Subscribe to DLQ topic
	if err := consumer.Subscribe("email-dlq"); err != nil {
		return fmt.Errorf("failed to subscribe to DLQ topic: %v", err)
	}

	// Process messages
	for {
		msg, err := consumer.Consume()
		if err != nil {
			s.logger.Error("Failed to consume DLQ message", zap.Error(err))
			continue
		}

		// Unmarshal failed email
		var failedEmail FailedEmail
		if err := json.Unmarshal(msg.Value, &failedEmail); err != nil {
			s.logger.Error("Failed to unmarshal DLQ message", zap.Error(err))
			continue
		}

		// Check if we should retry
		if failedEmail.Attempts >= s.config.Retry.MaxRetries {
			s.logger.Warn("Email exceeded maximum retry attempts",
				zap.Strings("to", failedEmail.To),
				zap.String("subject", failedEmail.Subject))
			continue
		}

		// Attempt to send the email
		err = smtpService.SendEmail(
			failedEmail.To,
			failedEmail.Subject,
			failedEmail.Body,
			failedEmail.Attachments,
		)

		if err != nil {
			// Update failed email and send back to DLQ
			failedEmail.Attempts++
			failedEmail.LastAttempt = time.Now()
			failedEmail.Error = err.Error()

			if err := s.HandleFailedEmail(&failedEmail); err != nil {
				s.logger.Error("Failed to handle retry failure", zap.Error(err))
			}
		} else {
			s.logger.Info("Successfully retried failed email",
				zap.Strings("to", failedEmail.To),
				zap.String("subject", failedEmail.Subject))
		}
	}
}
