package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
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
	logger     *zap.Logger
	config     *config.Config
	producer   *kafka.Producer
	dlqTopic   string        //DLQ Topic
	maxRetries int           //Max Retries
	retryDelay time.Duration //Retry Delay
	ctx        context.Context
}

// NewDLQService creates a new DLQ service instance
func NewDLQService(logger *zap.Logger, cfg *config.Config) (*DLQService, error) {
	producer, err := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic, logger) //Pass logger
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &DLQService{
		logger:     logger,
		config:     cfg,
		producer:   producer,
		dlqTopic:   cfg.Kafka.Topic + "-dlq", //Configure DLQ topic
		maxRetries: 3,                        //Configure max retries
		retryDelay: 5 * time.Second,          //Configure retry delay
		ctx:        context.Background(),
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
	if err := s.producer.SendMessage(context.Background(), s.dlqTopic, data); err != nil {
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

	// Process messages
	handler := func(ctx context.Context, msg *sarama.ConsumerMessage) error {
		var failedEmail FailedEmail
		if err := json.Unmarshal(msg.Value, &failedEmail); err != nil {
			s.logger.Error("Failed to unmarshal failed email", zap.Error(err))
			return nil //Don't retry if unmarshal fails
		}

		if failedEmail.Attempts >= s.maxRetries {
			s.logger.Error("Max retries reached for failed email",
				zap.Strings("to", failedEmail.To),
				zap.String("subject", failedEmail.Subject))
			return nil //Don't retry if max retries reached
		}

		//Retry the email
		failedEmail.Attempts++
		failedEmail.LastAttempt = time.Now()
		err := smtpService.SendEmail(&failedEmail) //Adapt to your send email function
		if err != nil {
			s.logger.Error("Failed to retry email",
				zap.Strings("to", failedEmail.To),
				zap.String("subject", failedEmail.Subject),
				zap.Int("attempts", failedEmail.Attempts),
				zap.Error(err))
			//Resend to DLQ
			s.HandleFailedEmail(&failedEmail) //Handle resending to DLQ
			return err
		}

		s.logger.Info("Successfully retried email",
			zap.Strings("to", failedEmail.To),
			zap.String("subject", failedEmail.Subject),
			zap.Int("attempts", failedEmail.Attempts))

		return nil
	}

	//Start consuming messages
	if err := consumer.Start(); err != nil {
		return fmt.Errorf("failed to start consumer: %v", err)
	}
	defer consumer.Stop()

	//Process messages in a loop
	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			time.Sleep(time.Second) //Avoid busy waiting
		}
	}
}
