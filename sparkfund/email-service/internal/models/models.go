package models

import (
	"time"
)

// EmailStatus represents the status of an email
type EmailStatus string

const (
	EmailStatusPending   EmailStatus = "pending"
	EmailStatusQueued    EmailStatus = "queued"
	EmailStatusSent      EmailStatus = "sent"
	EmailStatusFailed    EmailStatus = "failed"
	EmailStatusDelivered EmailStatus = "delivered"
)

// SendEmailRequest represents the request body for sending an email
type SendEmailRequest struct {
	To          []string          `json:"to" binding:"required"`
	Cc          []string          `json:"cc,omitempty"`
	Bcc         []string          `json:"bcc,omitempty"`
	Subject     string            `json:"subject" binding:"required"`
	Body        string            `json:"body" binding:"required"`
	TemplateID  string            `json:"template_id,omitempty"`
	Data        map[string]string `json:"data,omitempty"`
	Attachments []Attachment      `json:"attachments,omitempty"`
}

// Attachment represents an email attachment
type Attachment struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	Content     string `json:"content" binding:"required"` // Base64 encoded content
}

// EmailResponse represents the response for email operations
type EmailResponse struct {
	Message string `json:"message"`
}

// EmailLog represents a log entry for an email
type EmailLog struct {
	ID          string    `json:"id" db:"id"`
	To          []string  `json:"to" db:"to_addresses"`
	From        string    `json:"from" db:"from_address"`
	Subject     string    `json:"subject" db:"subject"`
	Body        string    `json:"body" db:"body"`
	ContentType string    `json:"content_type" db:"content_type"`
	Status      string    `json:"status" db:"status"`
	Error       string    `json:"error,omitempty" db:"error"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateTemplateRequest represents the request body for creating a template
type CreateTemplateRequest struct {
	Name        string   `json:"name" binding:"required"`
	Subject     string   `json:"subject" binding:"required"`
	Body        string   `json:"body" binding:"required"`
	Variables   []string `json:"variables" binding:"required"`
	Description string   `json:"description,omitempty"`
}

// UpdateTemplateRequest represents the request body for updating a template
type UpdateTemplateRequest struct {
	Name        string   `json:"name,omitempty"`
	Subject     string   `json:"subject,omitempty"`
	Body        string   `json:"body,omitempty"`
	Variables   []string `json:"variables,omitempty"`
	Description string   `json:"description,omitempty"`
}

// Template represents an email template
type Template struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Subject     string    `json:"subject" db:"subject"`
	Body        string    `json:"body" db:"body"`
	Variables   []string  `json:"variables" db:"variables"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// EmailMessage represents a message to be sent via Kafka
type EmailMessage struct {
	ID          string       `json:"id"`
	ToAddresses []string     `json:"to_addresses"`
	FromAddress string       `json:"from_address"`
	Subject     string       `json:"subject"`
	Body        string       `json:"body"`
	ContentType string       `json:"content_type"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// User represents a user in the system
type User struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
