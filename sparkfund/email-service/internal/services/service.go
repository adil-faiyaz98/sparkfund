package services

import (
	"context"
	"fmt"
	"time"

	"github.com/adil-faiyaz98/sparkfund/email-service/internal/config"
	"github.com/adil-faiyaz98/sparkfund/email-service/internal/kafka"
	"github.com/adil-faiyaz98/sparkfund/email-service/internal/models"
	"github.com/adil-faiyaz98/sparkfund/email-service/internal/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// EmailService defines the interface for email operations
type EmailService interface {
	SendEmail(req models.SendEmailRequest) error
	GetEmailLogs() ([]models.EmailLog, error)
	CreateTemplate(req models.CreateTemplateRequest) (*models.Template, error)
	GetTemplate(id string) (*models.Template, error)
	UpdateTemplate(id string, req models.UpdateTemplateRequest) (*models.Template, error)
	DeleteTemplate(id string) error
}

// Service implements the EmailService interface
type Service struct {
	logger      *zap.Logger
	config      *config.Config
	repo        repositories.Repository
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
	repo repositories.Repository,
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
func (s *Service) SendEmail(req models.SendEmailRequest) error {
	// Create email log entry
	log := models.EmailLog{
		ID:        uuid.New().String(),
		To:        req.To,
		Cc:        req.Cc,
		Bcc:       req.Bcc,
		Subject:   req.Subject,
		Body:      req.Body,
		Status:    "queued",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	if err := s.repo.CreateEmailLog(context.Background(), log); err != nil {
		return fmt.Errorf("failed to create email log: %w", err)
	}

	// If template is specified, process it
	if req.TemplateID != "" {
		template, err := s.repo.GetTemplate(context.Background(), req.TemplateID)
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
func (s *Service) GetEmailLogs() ([]models.EmailLog, error) {
	return s.repo.GetEmailLogs(context.Background())
}

// CreateTemplate creates a new email template
func (s *Service) CreateTemplate(req models.CreateTemplateRequest) (*models.Template, error) {
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

	if err := s.repo.CreateTemplate(context.Background(), template); err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	return template, nil
}

// GetTemplate retrieves a template by ID
func (s *Service) GetTemplate(id string) (*models.Template, error) {
	return s.repo.GetTemplate(context.Background(), id)
}

// UpdateTemplate updates an existing template
func (s *Service) UpdateTemplate(id string, req models.UpdateTemplateRequest) (*models.Template, error) {
	template, err := s.repo.GetTemplate(context.Background(), id)
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

	if err := s.repo.UpdateTemplate(context.Background(), template); err != nil {
		return nil, fmt.Errorf("failed to update template: %w", err)
	}

	return template, nil
}

// DeleteTemplate deletes a template by ID
func (s *Service) DeleteTemplate(id string) error {
	return s.repo.DeleteTemplate(context.Background(), id)
}

// processTemplate processes a template with the given data
func (s *Service) processTemplate(template *models.Template, data map[string]string) (string, string, error) {
	// TODO: Implement template processing logic
	// This would replace variables in the template with actual values
	return template.Subject, template.Body, nil
}
