package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sparkfund/email-service/internal/config"
	"github.com/sparkfund/email-service/internal/kafka"
	"github.com/sparkfund/email-service/internal/models"
	"github.com/sparkfund/email-service/internal/repository"
	"go.uber.org/zap"
)

// EmailService defines the interface for email operations
type EmailService interface {
	SendEmail(ctx context.Context, req models.SendEmailRequest) error
	GetEmailLogs(ctx context.Context) ([]models.EmailLog, error)
	CreateTemplate(ctx context.Context, req models.CreateTemplateRequest) (*models.Template, error)
	GetTemplate(ctx context.Context, id uuid.UUID) (*models.Template, error)
	UpdateTemplate(ctx context.Context, id uuid.UUID, req models.UpdateTemplateRequest) (*models.Template, error)
	DeleteTemplate(ctx context.Context, id uuid.UUID) error
}

// Service implements the EmailService interface
type Service struct {
	logger      *zap.Logger
	config      *config.Config
	repo        repository.Repository
	producer    *kafka.Producer
	authService AuthService // Add auth service
}

type AuthService struct {
	logger *zap.Logger
	config *config.Config
}

// ValidateToken validates a token and returns the user information
func (a *AuthService) ValidateToken(token string) (*models.User, error) {
	// TODO: Implement token validation
	return nil, nil
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(logger *zap.Logger, cfg *config.Config) *AuthService {
	return &AuthService{
		logger: logger,
		config: cfg,
	}
}

// NewService creates a new email service instance
func NewService(
	logger *zap.Logger,
	cfg *config.Config,
	repo repository.Repository,
	authService AuthService, // Add auth service
) (*Service, error) {
	producer, err := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic, logger) // Pass logger
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &Service{
		logger:      logger,
		config:      cfg,
		repo:        repo,
		producer:    producer,
		authService: authService, // Add auth service
	}, nil
}

// SendEmail queues an email for sending
func (s *Service) SendEmail(ctx context.Context, req models.SendEmailRequest) error {
	// Create email log entry
	id := uuid.New()
	// Convert EmailAddress slices to string slices
	recipients := make([]string, len(req.To))
	for i, addr := range req.To {
		recipients[i] = addr.Validate().Error()
	}

	log := models.EmailLog{
		ID:         id,
		Recipients: recipients,
		Cc:         []string{},
		Bcc:        []string{},
		From:       s.config.SMTP.From,
		Subject:    req.Subject,
		Body:       req.Body,
		Status:     models.EmailStatusQueued,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save to database
	if err := s.repo.CreateEmailLog(ctx, log); err != nil {
		s.logger.Error("Failed to create email log", zap.Error(err))
		return fmt.Errorf("failed to create email log: %w", err)
	}

	// If template is specified, process it
	if req.TemplateID != "" {
		templateID, err := uuid.Parse(req.TemplateID)
		if err != nil {
			return fmt.Errorf("invalid template ID: %w", err)
		}
		template, err := s.repo.GetTemplate(ctx, templateID)
		if err != nil {
			return fmt.Errorf("failed to get template: %w", err)
		}

		// Process template variables
		subject, body, err := s.processTemplate(template, req.Data)
		if err != nil {
			return fmt.Errorf("failed to process template: %w", err)
		}

		req.Subject = subject
		req.Body = body
	}

	// Create Kafka message
	toAddresses := make([]string, 0, len(req.To)) // Use make with 0 length
	for _, addr := range req.To {
		if err := addr.Validate(); err != nil {
			s.logger.Error("Invalid To address", zap.Error(err), zap.String("email", string(addr)))
			return err // Or handle the error differently (e.g., skip the address)
		}
		toAddresses = append(toAddresses, string(addr)) // Append valid address
	}

	ccAddresses := make([]string, 0, len(req.Cc)) // Use make with 0 length
	for _, addr := range req.Cc {
		if err := addr.Validate(); err != nil {
			s.logger.Warn("Invalid CC address", zap.Error(err), zap.String("email", string(addr)))
			// Decide how to handle the error: skip the address or return an error
			continue // Skip the invalid CC address
		}
		ccAddresses = append(ccAddresses, string(addr)) // Append valid address
	}

	bccAddresses := make([]string, 0, len(req.Bcc)) // Use make with 0 length
	for _, addr := range req.Bcc {
		if err := addr.Validate(); err != nil {
			s.logger.Warn("Invalid BCC address", zap.Error(err), zap.String("email", string(addr)))
			// Decide how to handle the error: skip the address or return an error
			continue // Skip the invalid BCC address
		}
		bccAddresses = append(bccAddresses, string(addr)) // Append valid address
	}

	message := models.EmailMessage{
		ID:          id.String(),
		ToAddresses: toAddresses,
		FromAddress: s.config.SMTP.From,
		Subject:     req.Subject,
		Body:        req.Body,
		ContentType: "text/html", // Assuming HTML content
	}

	// Marshal message to JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Send to Kafka queue for processing
	if err := s.producer.SendMessage(ctx, s.config.Kafka.Topic, id.String(), messageBytes); err != nil {
		s.logger.Error("Failed to send message to Kafka", zap.Error(err))
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	s.logger.Info("Email queued for sending", zap.String("email_id", id.String()))
	return nil
}

// GetEmailLogs retrieves all email logs
func (s *Service) GetEmailLogs(ctx context.Context) ([]models.EmailLog, error) {
	return s.repo.GetEmailLogs(ctx)
}

// CreateTemplate creates a new email template
func (s *Service) CreateTemplate(ctx context.Context, req models.CreateTemplateRequest) (*models.Template, error) {
	id := uuid.New()
	template := &models.Template{
		ID:          id,
		Name:        req.Name,
		Subject:     req.Subject,
		Body:        req.Body,
		Variables:   req.Variables,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateTemplate(ctx, template); err != nil {
		s.logger.Error("Failed to create template", zap.Error(err))
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	return template, nil
}

// GetTemplate retrieves a template by ID
func (s *Service) GetTemplate(ctx context.Context, id uuid.UUID) (*models.Template, error) {
	return s.repo.GetTemplate(ctx, id)
}

// UpdateTemplate updates an existing template
func (s *Service) UpdateTemplate(ctx context.Context, id uuid.UUID, req models.UpdateTemplateRequest) (*models.Template, error) {
	template, err := s.repo.GetTemplate(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		template.Name = req.Name
	}
	if req.Subject != "" {
		template.Subject = req.Subject
	}
	if req.Body != "" {
		template.Body = req.Body
	}
	if len(req.Variables) > 0 {
		template.Variables = req.Variables
	}
	if req.Description != "" {
		template.Description = req.Description
	}
	template.UpdatedAt = time.Now()

	if err := s.repo.UpdateTemplate(ctx, template); err != nil {
		s.logger.Error("Failed to update template", zap.Error(err))
		return nil, fmt.Errorf("failed to update template: %w", err)
	}

	return template, nil
}

// DeleteTemplate deletes a template by ID
func (s *Service) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteTemplate(ctx, id)
}

// processTemplate processes a template with the given data
func (s *Service) processTemplate(template *models.Template, data map[string]string) (string, string, error) {
	// TODO: Implement template processing logic
	// This would replace variables in the template with actual values
	return template.Subject, template.Body, nil
}
