package validation

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/adil-faiyaz98/sparkfund/email-service/internal/models"
)

var (
	// Email validation regex
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	// Template name validation regex (alphanumeric, hyphens, underscores)
	templateNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

// ValidateEmailRequest validates a SendEmailRequest
func ValidateEmailRequest(req *models.SendEmailRequest) error {
	if req == nil {
		return &ValidationError{
			Field:   "request",
			Message: "request cannot be nil",
		}
	}

	// Validate recipients
	if len(req.To) == 0 {
		return &ValidationError{
			Field:   "to",
			Message: "at least one recipient is required",
		}
	}
	for i, to := range req.To {
		if err := validateEmail(to); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("to[%d]", i),
				Message: err.Error(),
			}
		}
	}

	// Validate CC recipients
	for i, cc := range req.Cc {
		if err := validateEmail(cc); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("cc[%d]", i),
				Message: err.Error(),
			}
		}
	}

	// Validate BCC recipients
	for i, bcc := range req.Bcc {
		if err := validateEmail(bcc); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("bcc[%d]", i),
				Message: err.Error(),
			}
		}
	}

	// Validate sender
	if err := validateEmail(req.From); err != nil {
		return &ValidationError{
			Field:   "from",
			Message: err.Error(),
		}
	}

	// Validate subject
	if strings.TrimSpace(req.Subject) == "" {
		return &ValidationError{
			Field:   "subject",
			Message: "subject cannot be empty",
		}
	}
	if len(req.Subject) > 255 {
		return &ValidationError{
			Field:   "subject",
			Message: "subject cannot exceed 255 characters",
		}
	}

	// Validate body
	if strings.TrimSpace(req.Body) == "" {
		return &ValidationError{
			Field:   "body",
			Message: "body cannot be empty",
		}
	}

	// Validate content type
	if req.ContentType != "" {
		validContentTypes := map[string]bool{
			"text/plain": true,
			"text/html":  true,
		}
		if !validContentTypes[req.ContentType] {
			return &ValidationError{
				Field:   "content_type",
				Message: "invalid content type. Must be text/plain or text/html",
			}
		}
	}

	// Validate attachments
	if len(req.Attachments) > 0 {
		for i, attachment := range req.Attachments {
			if err := validateAttachment(attachment); err != nil {
				return &ValidationError{
					Field:   fmt.Sprintf("attachments[%d]", i),
					Message: err.Error(),
				}
			}
		}
	}

	return nil
}

// ValidateTemplate validates an email template
func ValidateTemplate(template *models.Template) error {
	if template == nil {
		return &ValidationError{
			Field:   "template",
			Message: "template cannot be nil",
		}
	}

	// Validate name
	if strings.TrimSpace(template.Name) == "" {
		return &ValidationError{
			Field:   "name",
			Message: "template name cannot be empty",
		}
	}
	if !templateNameRegex.MatchString(template.Name) {
		return &ValidationError{
			Field:   "name",
			Message: "template name can only contain alphanumeric characters, hyphens, and underscores",
		}
	}
	if len(template.Name) > 255 {
		return &ValidationError{
			Field:   "name",
			Message: "template name cannot exceed 255 characters",
		}
	}

	// Validate subject
	if strings.TrimSpace(template.Subject) == "" {
		return &ValidationError{
			Field:   "subject",
			Message: "template subject cannot be empty",
		}
	}
	if len(template.Subject) > 255 {
		return &ValidationError{
			Field:   "subject",
			Message: "template subject cannot exceed 255 characters",
		}
	}

	// Validate body
	if strings.TrimSpace(template.Body) == "" {
		return &ValidationError{
			Field:   "body",
			Message: "template body cannot be empty",
		}
	}

	// Validate variables
	if len(template.Variables) == 0 {
		return &ValidationError{
			Field:   "variables",
			Message: "template must have at least one variable",
		}
	}
	for i, variable := range template.Variables {
		if strings.TrimSpace(variable) == "" {
			return &ValidationError{
				Field:   fmt.Sprintf("variables[%d]", i),
				Message: "variable name cannot be empty",
			}
		}
		if !templateNameRegex.MatchString(variable) {
			return &ValidationError{
				Field:   fmt.Sprintf("variables[%d]", i),
				Message: "variable name can only contain alphanumeric characters, hyphens, and underscores",
			}
		}
	}

	return nil
}

// ValidateEmailLog validates an email log entry
func ValidateEmailLog(log *models.EmailLog) error {
	if log == nil {
		return &ValidationError{
			Field:   "log",
			Message: "log cannot be nil",
		}
	}

	// Validate ID
	if strings.TrimSpace(log.ID) == "" {
		return &ValidationError{
			Field:   "id",
			Message: "ID cannot be empty",
		}
	}

	// Validate recipients
	if len(log.To) == 0 {
		return &ValidationError{
			Field:   "to",
			Message: "at least one recipient is required",
		}
	}
	for i, to := range log.To {
		if err := validateEmail(to); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("to[%d]", i),
				Message: err.Error(),
			}
		}
	}

	// Validate sender
	if err := validateEmail(log.From); err != nil {
		return &ValidationError{
			Field:   "from",
			Message: err.Error(),
		}
	}

	// Validate subject
	if strings.TrimSpace(log.Subject) == "" {
		return &ValidationError{
			Field:   "subject",
			Message: "subject cannot be empty",
		}
	}
	if len(log.Subject) > 255 {
		return &ValidationError{
			Field:   "subject",
			Message: "subject cannot exceed 255 characters",
		}
	}

	// Validate body
	if strings.TrimSpace(log.Body) == "" {
		return &ValidationError{
			Field:   "body",
			Message: "body cannot be empty",
		}
	}

	// Validate content type
	if log.ContentType != "" {
		validContentTypes := map[string]bool{
			"text/plain": true,
			"text/html":  true,
		}
		if !validContentTypes[log.ContentType] {
			return &ValidationError{
				Field:   "content_type",
				Message: "invalid content type. Must be text/plain or text/html",
			}
		}
	}

	// Validate status
	if log.Status == "" {
		return &ValidationError{
			Field:   "status",
			Message: "status cannot be empty",
		}
	}

	validStatuses := map[string]bool{
		string(models.EmailStatusPending):   true,
		string(models.EmailStatusQueued):    true,
		string(models.EmailStatusSent):      true,
		string(models.EmailStatusFailed):    true,
		string(models.EmailStatusDelivered): true,
	}

	if !validStatuses[log.Status] {
		return &ValidationError{
			Field:   "status",
			Message: "invalid status. Must be one of: pending, queued, sent, failed, delivered",
		}
	}

	return nil
}

// validateEmail validates an email address
func validateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return &ValidationError{
			Field:   "email",
			Message: "email address cannot be empty",
		}
	}

	if !emailRegex.MatchString(email) {
		return &ValidationError{
			Field:   "email",
			Message: "invalid email address format",
		}
	}

	// Additional validation using net/mail
	if _, err := mail.ParseAddress(email); err != nil {
		return &ValidationError{
			Field:   "email",
			Message: "invalid email address",
		}
	}

	return nil
}

// validateAttachment validates an email attachment
func validateAttachment(attachment models.Attachment) error {
	if strings.TrimSpace(attachment.Filename) == "" {
		return &ValidationError{
			Field:   "filename",
			Message: "filename cannot be empty",
		}
	}

	if strings.TrimSpace(attachment.ContentType) == "" {
		return &ValidationError{
			Field:   "content_type",
			Message: "content type cannot be empty",
		}
	}

	if strings.TrimSpace(attachment.Content) == "" {
		return &ValidationError{
			Field:   "content",
			Message: "content cannot be empty",
		}
	}

	return nil
}
