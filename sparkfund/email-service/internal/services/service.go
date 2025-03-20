package services

import (
	"context"
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
	GetTemplate(ctx context.Context, id string) (*models.Template, error)
	UpdateTemplate(ctx context.Context, id string, req models.UpdateTemplateRequest) (*models.Template, error)
	DeleteTemplate(ctx context.Context, id string) error
}

// Service implements the EmailService interface
type Service struct {
	logger      *zap.Logger
	config      *config.Config
	repo        repository.Repository
	producer    *kafka.Producer
	authService AuthService
}

type AuthService interface {
	ValidateToken(token string) (*models.User, error)
}

// NewService creates a new email service instance
func NewService(
	logger *zap.Logger,
	cfg *config.Config,
	repo repository.Repository,
) (*Service, error) {
	producer, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &Service{
		logger:   logger,
		config:   cfg,
		repo:     repo,
		producer: producer,
	}, nil
}

// SendEmail queues an email for sending
func (s *Service) SendEmail(ctx context.Context, req models.SendEmailRequest) error {
	// Create email log entry
	log := models.EmailLog{
		ID:         uuid.New().String(),
		Recipients: req.To,
		Cc:         req.Cc,
		Bcc:        req.Bcc,
		From:       s.config.SMTP.From,
		Subject:    req.Subject,
		Body:       req.Body,
		Status:     models.EmailStatusQueued,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save to database
	if err := s.repo.CreateEmailLog(ctx, log); err != nil {
		return fmt.Errorf("failed to create email log: %w", err)
	}

	// If template is specified, process it
	if req.TemplateID != "" {
		template, err := s.repo.GetTemplate(ctx, req.TemplateID)
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

	// TODO: Send to Kafka queue for processing
	// This would be implemented in a separate worker service

	return nil
}

// GetEmailLogs retrieves all email logs
func (s *Service) GetEmailLogs(ctx context.Context) ([]models.EmailLog, error) {
	return s.repo.GetEmailLogs(ctx)
}

// CreateTemplate creates a new email template
func (s *Service) CreateTemplate(ctx context.Context, req models.CreateTemplateRequest) (*models.Template, error) {
	template := &models.Template{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Subject:     req.Subject,
		Body:        req.Body,
		Variables:   req.Variables,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateTemplate(ctx, template); err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	return template, nil
}

// GetTemplate retrieves a template by ID
func (s *Service) GetTemplate(ctx context.Context, id string) (*models.Template, error) {
	return s.repo.GetTemplate(ctx, id)
}

// UpdateTemplate updates an existing template
func (s *Service) UpdateTemplate(ctx context.Context, id string, req models.UpdateTemplateRequest) (*models.Template, error) {
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
		return nil, fmt.Errorf("failed to update template: %w", err)
	}

	return template, nil
}

// DeleteTemplate deletes a template by ID
func (s *Service) DeleteTemplate(ctx context.Context, id string) error {
	return s.repo.DeleteTemplate(ctx, id)
}

// processTemplate processes a template with the given data
func (s *Service) processTemplate(template *models.Template, data map[string]string) (string, string, error) {
	// TODO: Implement template processing logic
	// This would replace variables in the template with actual values
	return template.Subject, template.Body, nil
}
