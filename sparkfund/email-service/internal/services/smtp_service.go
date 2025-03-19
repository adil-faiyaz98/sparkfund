package services

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/smtp"
	"net/textproto"
	"path/filepath"
	"strings"
	"time"

	"github.com/adil-faiyaz98/sparkfund/email-service/internal/config"
	"go.uber.org/zap"
)

// SMTPError represents an SMTP-specific error
type SMTPError struct {
	Code    int
	Message string
	Err     error
}

func (e *SMTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("SMTP error %d: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("SMTP error %d: %s", e.Code, e.Message)
}

// SMTPService handles email sending with retries
type SMTPService struct {
	logger *zap.Logger
	config *config.Config
	client *smtp.Client
	conn   net.Conn
}

// NewSMTPService creates a new SMTP service instance
func NewSMTPService(logger *zap.Logger, cfg *config.Config) *SMTPService {
	return &SMTPService{
		logger: logger,
		config: cfg,
	}
}

// SendEmail sends an email with retries
func (s *SMTPService) SendEmail(to []string, subject, body string, attachments map[string][]byte) error {
	var lastErr error
	interval := s.config.Retry.InitialInterval

	for attempt := 1; attempt <= s.config.Retry.MaxRetries; attempt++ {
		err := s.sendEmailWithRetry(to, subject, body, attachments)
		if err == nil {
			s.logger.Info("Email sent successfully",
				zap.Strings("to", to),
				zap.String("subject", subject),
				zap.Int("attempt", attempt))
			return nil
		}

		lastErr = err
		s.logger.Error("Failed to send email, retrying",
			zap.Strings("to", to),
			zap.String("subject", subject),
			zap.Int("attempt", attempt),
			zap.Error(err),
			zap.Duration("next_retry", interval))

		time.Sleep(interval)
		interval = time.Duration(float64(interval) * s.config.Retry.Multiplier)
		if interval > s.config.Retry.MaxInterval {
			interval = s.config.Retry.MaxInterval
		}
	}

	return fmt.Errorf("failed to send email after %d attempts: %v", s.config.Retry.MaxRetries, lastErr)
}

// sendEmailWithRetry attempts to send an email with a single retry
func (s *SMTPService) sendEmailWithRetry(to []string, subject, body string, attachments map[string][]byte) error {
	// Connect to SMTP server
	if err := s.connect(); err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer s.disconnect()

	// Authenticate
	if err := s.authenticate(); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	// Send email
	if err := s.send(to, subject, body, attachments); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// connect establishes a connection to the SMTP server
func (s *SMTPService) connect() error {
	addr := fmt.Sprintf("%s:%d", s.config.SMTP.Host, s.config.SMTP.Port)
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return &SMTPError{Code: 0, Message: "connection failed", Err: err}
	}

	client, err := smtp.NewClient(conn, s.config.SMTP.Host)
	if err != nil {
		conn.Close()
		return &SMTPError{Code: 0, Message: "client creation failed", Err: err}
	}

	s.conn = conn
	s.client = client
	return nil
}

// disconnect closes the SMTP connection
func (s *SMTPService) disconnect() {
	if s.client != nil {
		s.client.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
}

// authenticate authenticates with the SMTP server
func (s *SMTPService) authenticate() error {
	auth := smtp.PlainAuth("", s.config.SMTP.Username, s.config.SMTP.Password, s.config.SMTP.Host)
	if err := s.client.Auth(auth); err != nil {
		return &SMTPError{Code: 0, Message: "authentication failed", Err: err}
	}
	return nil
}

// send sends the email with attachments
func (s *SMTPService) send(to []string, subject, body string, attachments map[string][]byte) error {
	// Set sender
	if err := s.client.Mail(s.config.SMTP.From); err != nil {
		return &SMTPError{Code: 0, Message: "sender rejected", Err: err}
	}

	// Set recipients
	for _, recipient := range to {
		if err := s.client.Rcpt(recipient); err != nil {
			return &SMTPError{Code: 0, Message: "recipient rejected", Err: err}
		}
	}

	// Create writer for email data
	w, err := s.client.Data()
	if err != nil {
		return &SMTPError{Code: 0, Message: "data command failed", Err: err}
	}

	// Write email headers and content
	if err := s.writeEmail(w, to, subject, body, attachments); err != nil {
		return err
	}

	// Close writer
	if err := w.Close(); err != nil {
		return &SMTPError{Code: 0, Message: "failed to close writer", Err: err}
	}

	return nil
}

// writeEmail writes the email content including headers and attachments
func (s *SMTPService) writeEmail(w io.Writer, to []string, subject, body string, attachments map[string][]byte) error {
	// Create multipart writer
	writer := multipart.NewWriter(w)
	boundary := writer.Boundary()

	// Write headers
	headers := []string{
		"From: " + s.config.SMTP.From,
		"To: " + strings.Join(to, ", "),
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: multipart/mixed; boundary=" + boundary,
		"",
	}
	if _, err := w.Write([]byte(strings.Join(headers, "\r\n"))); err != nil {
		return err
	}

	// Write text part
	if err := s.writeTextPart(writer, body); err != nil {
		return err
	}

	// Write attachments
	for filename, content := range attachments {
		if err := s.writeAttachment(writer, filename, content); err != nil {
			return err
		}
	}

	return writer.Close()
}

// writeTextPart writes the text part of the email
func (s *SMTPService) writeTextPart(writer *multipart.Writer, body string) error {
	part, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": []string{"text/plain; charset=UTF-8"},
	})
	if err != nil {
		return err
	}

	_, err = part.Write([]byte(body))
	return err
}

// writeAttachment writes an attachment to the email
func (s *SMTPService) writeAttachment(writer *multipart.Writer, filename string, content []byte) error {
	ext := filepath.Ext(filename)
	mimeType := "application/octet-stream"
	switch ext {
	case ".txt":
		mimeType = "text/plain"
	case ".pdf":
		mimeType = "application/pdf"
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".png":
		mimeType = "image/png"
	}

	part, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":              []string{mimeType},
		"Content-Disposition":       []string{fmt.Sprintf("attachment; filename=%s", filename)},
		"Content-Transfer-Encoding": []string{"base64"},
	})
	if err != nil {
		return err
	}

	encoder := base64.NewEncoder(base64.StdEncoding, part)
	if _, err := encoder.Write(content); err != nil {
		return err
	}
	return encoder.Close()
}
