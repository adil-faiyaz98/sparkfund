package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Document metrics
	documentUploadsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "document_uploads_total",
			Help: "Total number of document uploads",
		},
	)

	documentVerificationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "document_verifications_total",
			Help: "Total number of document verifications",
		},
		[]string{"status"},
	)

	// Verification metrics
	verificationDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "verification_duration_seconds",
			Help:    "Document verification duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	// Error metrics
	errorTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors",
		},
		[]string{"type"},
	)
)

// PrometheusMiddleware adds Prometheus metrics to the Gin router
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.Request.URL.Path,
			status,
		).Inc()

		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			c.Request.URL.Path,
		).Observe(duration)
	}
}

// RecordDocumentUpload increments the document upload counter
func RecordDocumentUpload() {
	documentUploadsTotal.Inc()
}

// RecordDocumentVerification records a document verification with its status
func RecordDocumentVerification(status string) {
	documentVerificationsTotal.WithLabelValues(status).Inc()
}

// RecordVerificationDuration records the duration of a verification process
func RecordVerificationDuration(duration time.Duration) {
	verificationDuration.Observe(duration.Seconds())
}

// RecordError records an error occurrence
func RecordError(errorType string) {
	errorTotal.WithLabelValues(errorType).Inc()
} 