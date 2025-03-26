package services

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/smtp"
	"net/textproto"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sparkfund/email-service/internal/config"
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
	logger         *zap.Logger
	config         *config.Config
	connectionPool chan *smtp.Client // Connection pool
	maxConnections int               // Max connections in the pool
	emailChan      chan emailTask    // Channel for asynchronous email sending
	metrics        *SmtpMetrics      //Metrics
	stopChan       chan struct{}
}

type SmtpMetrics struct {
	emailsSent       *prometheus.CounterVec
	emailsRetried    *prometheus.CounterVec
	emailsFailed     *prometheus.CounterVec
	connectionErrors *prometheus.CounterVec
}

type emailTask struct {
	to          []string
	subject     string
	body        string
	attachments map[string][]byte
}

// NewSMTPService creates a new SMTP service instance
func NewSMTPService(logger *zap.Logger, cfg *config.Config, metrics *SmtpMetrics) (*SMTPService, error) {
	maxConnections := 10 // Example - configure this
	connectionPool := make(chan *smtp.Client, maxConnections)
	stopChan := make(chan struct{})

	smtpService := &SMTPService{
		logger:         logger,
		config:         cfg,
		connectionPool: connectionPool,
		maxConnections: maxConnections,
		emailChan:      make(chan emailTask, 100), // Buffered channel
		metrics:        metrics,
		stopChan:       stopChan,
	}

	//Initialize connection pool
	for i := 0; i < maxConnections; i++ {
		client, err := smtpService.createClient()
		if err != nil {
			return nil, fmt.Errorf("failed to create initial SMTP client: %w", err)
		}
		connectionPool <- client
	}

	go smtpService.startWorkerPool()

	return smtpService, nil
}

func (s *SMTPService) createClient() (*smtp.Client, error) {
	addr := fmt.Sprintf("%s:%d", s.config.SMTP.Host, s.config.SMTP.Port)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true, // TODO: In production, get a real certificate
		ServerName:         s.config.SMTP.Host,
	}

	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		s.metrics.connectionErrors.With(prometheus.Labels{"error_type": "dial"}).Inc()
		return nil, &SMTPError{Code: 0, Message: "connection failed", Err: err}
	}

	client, err := smtp.NewClient(conn, s.config.SMTP.Host)
	if err != nil {
		conn.Close()
		s.metrics.connectionErrors.With(prometheus.Labels{"error_type": "client_creation"}).Inc()
		return nil, &SMTPError{Code: 0, Message: "client creation failed", Err: err}
	}

	if err = client.StartTLS(tlsconfig); err != nil {
		client.Close()
		s.metrics.connectionErrors.With(prometheus.Labels{"error_type": "start_tls"}).Inc()
		return nil, &SMTPError{Code: 0, Message: "start tls failed", Err: err}
	}

	// Authenticate
	auth := smtp.PlainAuth("", s.config.SMTP.Username, s.config.SMTP.Password, s.config.SMTP.Host)
	if err := client.Auth(auth); err != nil {
		client.Close()
		s.metrics.connectionErrors.With(prometheus.Labels{"error_type": "authentication"}).Inc()
		return nil, &SMTPError{Code: 0, Message: "authentication failed", Err: err}
	}

	return client, nil
}

// startWorkerPool launches worker goroutines to handle email sending
func (s *SMTPService) startWorkerPool() {
	numWorkers := 5 // Example - configure this
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.emailWorker()
		}()
	}

	go func() {
		wg.Wait()
		close(s.connectionPool) // Close connection pool after workers exit
	}()
}

// emailWorker retrieves email tasks from the channel and attempts to send them
func (s *SMTPService) emailWorker() {
	for {
		select {
		case task := <-s.emailChan:
			err := s.sendEmailWithRetry(task.to, task.subject, task.body, task.attachments)
			if err != nil {
				s.logger.Error("Failed to send email",
					zap.Strings("to", task.to),
					zap.String("subject", task.subject),
					zap.Error(err))
				s.metrics.emailsFailed.With(prometheus.Labels{"status": "failed"}).Inc()
			} else {
				s.metrics.emailsSent.With(prometheus.Labels{"status": "success"}).Inc()
			}
		case <-s.stopChan:
			return // Exit worker
		}
	}
}

// SendEmail queues an email for sending
func (s *SMTPService) SendEmail(to []string, subject, body string, attachments map[string][]byte) error {
	task := emailTask{
		to:          to,
		subject:     subject,
		body:        body,
		attachments: attachments,
	}

	select {
	case s.emailChan <- task:
		return nil
	default:
		return fmt.Errorf("email channel is full, email queued failed") //Handle queue full
	}
}

// sendEmailWithRetry attempts to send an email with a single retry
func (s *SMTPService) sendEmailWithRetry(to []string, subject, body string, attachments map[string][]byte) error {
	var client *smtp.Client
	select {
	case client = <-s.connectionPool:
		defer func() { s.connectionPool <- client }() //Return connection to pool
	default:
		//No connection available in the pool, create a new one
		var err error
		client, err = s.createClient()
		if err != nil {
			return err
		}
		defer client.Close() //Close the client as it is not part of the pool
	}

	// Set sender
	if err := client.Mail(s.config.SMTP.From); err != nil {
		return &SMTPError{Code: 0, Message: "sender rejected", Err: err}
	}

	// Set recipients
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return &SMTPError{Code: 0, Message: "recipient rejected", Err: err}
		}
	}

	// Create writer for email data
	w, err := client.Data()
	if err != nil {
		return &SMTPError{Code: 0, Message: "data command failed", Err: err}
	}

	// Write email headers and content
	if err := s.writeEmail(w, s.config.SMTP.From, to, subject, body, attachments); err != nil {
		return err
	}

	// Close writer
	if err := w.Close(); err != nil {
		return &SMTPError{Code: 0, Message: "failed to close writer", Err: err}
	}

	return nil
}

// writeEmail writes the email content including headers and attachments
func (s *SMTPService) writeEmail(w io.Writer, from string, to []string, subject, body string, attachments map[string][]byte) error {
	// Create multipart writer
	writer := multipart.NewWriter(w)
	boundary := writer.Boundary()

	// Write headers
	headers := []string{
		"From: " + from,
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

func (s *SMTPService) isPermanentError(err error) bool {
	if smtpErr, ok := err.(*SMTPError); ok {
		// Check for specific SMTP error codes that indicate permanent failures
		return smtpErr.Code >= 500 && smtpErr.Code < 600 // 5xx errors are permanent
	}
	// If it's not an SMTPError, assume it's transient
	return false
}

// Stop gracefully shuts down the SMTP service
func (s *SMTPService) Stop() {
	s.logger.Info("Stopping SMTP service...")
	close(s.stopChan) // Signal workers to stop

	//Drain the email channel
	for i := 0; i < len(s.emailChan); i++ {
		<-s.emailChan
	}
}

func NewSmtpMetrics(registerer prometheus.Registerer) *SmtpMetrics {
	m := &SmtpMetrics{
		emailsSent: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "smtp_emails_sent_total",
				Help: "Total number of emails successfully sent via SMTP",
			},
			[]string{"status"},
		),
		emailsRetried: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "smtp_emails_retried_total",
				Help: "Total number of email retries",
			},
			[]string{"status"},
		),
		emailsFailed: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "smtp_emails_failed_total",
				Help: "Total number of emails failed to send via SMTP",
			},
			[]string{"status"},
		),
		connectionErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "smtp_connection_errors_total",
				Help: "Total number of SMTP connection errors",
			},
			[]string{"error_type"},
		),
	}

	if registerer != nil {
		registerer.MustRegister(
			m.emailsSent,
			m.emailsRetried,
			m.emailsFailed,
			m.connectionErrors,
		)
	}
	return m
}
