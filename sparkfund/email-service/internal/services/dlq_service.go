package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sparkfund/email-service/internal/config"
	"github.com/sparkfund/email-service/internal/kafka"
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
	ctx      context.Context
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
		ctx:      context.Background(),
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
	// Start consuming messages
	if err := consumer.Start(); err != nil {
		return fmt.Errorf("failed to start consumer: %v", err)
	}
	defer consumer.Stop()

	// Process messages
	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			// Process messages in a loop
			// The actual message processing is handled by the MessageHandler function
			// that was passed to NewConsumer
			time.Sleep(time.Second) // Avoid busy waiting
		}
	}
}
