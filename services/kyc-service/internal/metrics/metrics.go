package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "endpoint"},
	)

	// Business metrics
	KYCVerifications = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kyc_verifications_total",
			Help: "Total number of KYC verifications",
		},
		[]string{"status", "type"},
	)

	DocumentUploads = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "document_uploads_total",
			Help: "Total number of document uploads",
		},
		[]string{"type", "status"},
	)

	VerificationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "verification_duration_seconds",
			Help:    "Time taken to complete verifications",
			Buckets: []float64{1, 5, 10, 30, 60, 120, 300, 600},
		},
		[]string{"type"},
	)

	ErrorCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors",
		},
		[]string{"type", "code"},
	)
)

// MetricsService handles all metrics operations
type MetricsService struct{}

func NewMetricsService() *MetricsService {
	return &MetricsService{}
}

// HTTP metrics
func (s *MetricsService) RecordRequest(method, endpoint, status string, duration float64) {
	RequestCounter.WithLabelValues(method, endpoint, status).Inc()
	RequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

// Business metrics
func (s *MetricsService) RecordKYCVerification(status, verType string) {
	KYCVerifications.WithLabelValues(status, verType).Inc()
}

func (s *MetricsService) RecordDocumentUpload(docType, status string) {
	DocumentUploads.WithLabelValues(docType, status).Inc()
}

func (s *MetricsService) RecordVerificationDuration(verType string, duration time.Duration) {
	VerificationDuration.WithLabelValues(verType).Observe(duration.Seconds())
}

func (s *MetricsService) RecordError(errorType, code string) {
	ErrorCounter.WithLabelValues(errorType, code).Inc()
}
