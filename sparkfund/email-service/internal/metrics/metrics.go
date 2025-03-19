package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Email metrics
	EmailSentCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "email_sent_total",
			Help: "Total number of emails sent",
		},
		[]string{"status"},
	)

	EmailProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "email_processing_duration_seconds",
			Help:    "Duration of email processing in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status"},
	)

	// Template metrics
	TemplateOperationsCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "template_operations_total",
			Help: "Total number of template operations",
		},
		[]string{"operation", "status"},
	)

	// Database metrics
	DatabaseOperationsDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_operations_duration_seconds",
			Help:    "Duration of database operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "status"},
	)

	DatabaseErrorsCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_errors_total",
			Help: "Total number of database errors",
		},
		[]string{"operation", "error_type"},
	)

	// Kafka metrics
	KafkaMessagesProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_messages_processed_total",
			Help: "Total number of Kafka messages processed",
		},
		[]string{"topic", "status"},
	)

	KafkaProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "kafka_processing_duration_seconds",
			Help:    "Duration of Kafka message processing in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"topic"},
	)

	// SMTP metrics
	SMTPOperationsDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "smtp_operations_duration_seconds",
			Help:    "Duration of SMTP operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "status"},
	)

	SMTPErrorsCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "smtp_errors_total",
			Help: "Total number of SMTP errors",
		},
		[]string{"operation", "error_type"},
	)
)
