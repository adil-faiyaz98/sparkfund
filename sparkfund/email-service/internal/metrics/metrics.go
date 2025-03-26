package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics defines the Prometheus metrics for the email service
type Metrics struct {
	EmailSentCounter           *prometheus.CounterVec
	EmailProcessingDuration    *prometheus.HistogramVec
	TemplateOperationsCounter  *prometheus.CounterVec
	DatabaseOperationsDuration *prometheus.HistogramVec
	DatabaseErrorsCounter      *prometheus.CounterVec
	KafkaMessagesProcessed     *prometheus.CounterVec
	KafkaProcessingDuration    *prometheus.HistogramVec
	SMTPOperationsDuration     *prometheus.HistogramVec
	SMTPErrorsCounter          *prometheus.CounterVec
}

// NewMetrics creates a new Metrics instance
func NewMetrics(registerer prometheus.Registerer) *Metrics {
	m := &Metrics{
		EmailSentCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "email_sent_total",
				Help: "Total number of emails sent",
			},
			[]string{"status"},
		),
		EmailProcessingDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "email_processing_duration_seconds",
				Help:    "Duration of email processing in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"status"},
		),
		TemplateOperationsCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "template_operations_total",
				Help: "Total number of template operations",
			},
			[]string{"operation", "status"},
		),
		DatabaseOperationsDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "database_operations_duration_seconds",
				Help:    "Duration of database operations in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "status"},
		),
		DatabaseErrorsCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "database_errors_total",
				Help: "Total number of database errors",
			},
			[]string{"operation", "error_type"},
		),
		KafkaMessagesProcessed: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "kafka_messages_processed_total",
				Help: "Total number of Kafka messages processed",
			},
			[]string{"topic", "status"},
		),
		KafkaProcessingDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "kafka_processing_duration_seconds",
				Help:    "Duration of Kafka message processing in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"topic"},
		),
		SMTPOperationsDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "smtp_operations_duration_seconds",
				Help:    "Duration of SMTP operations in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "status"},
		),
		SMTPErrorsCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "smtp_errors_total",
				Help: "Total number of SMTP errors",
			},
			[]string{"operation", "error_type"},
		),
	}

	// Register metrics with the provided registerer
	registerer.MustRegister(
		m.EmailSentCounter,
		m.EmailProcessingDuration,
		m.TemplateOperationsCounter,
		m.DatabaseOperationsDuration,
		m.DatabaseErrorsCounter,
		m.KafkaMessagesProcessed,
		m.KafkaProcessingDuration,
		m.SMTPOperationsDuration,
		m.SMTPErrorsCounter,
	)

	return m
}

// NewPrometheusRegistry creates a new Prometheus registry
func NewPrometheusRegistry() *prometheus.Registry {
	return prometheus.NewRegistry()
}

// NewAutoRegisterMetrics creates a new Metrics instance and registers it with the promauto global registry
func NewAutoRegisterMetrics() *Metrics {
	return NewMetrics(prometheus.DefaultRegisterer)
}
