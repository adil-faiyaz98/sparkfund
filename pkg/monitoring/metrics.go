package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RequestDuration tracks HTTP request duration
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests in seconds",
			Buckets: []float64{
				0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10,
			},
		},
		[]string{"handler", "method", "status"},
	)

	// RequestsTotal tracks total number of HTTP requests
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"handler", "method", "status"},
	)

	// DatabaseOperationDuration tracks database operation duration
	DatabaseOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "database_operation_duration_seconds",
			Help: "Duration of database operations in seconds",
			Buckets: []float64{
				0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5,
			},
		},
		[]string{"operation", "status"},
	)

	// DatabaseConnectionsOpen tracks number of open database connections
	DatabaseConnectionsOpen = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_open",
			Help: "Number of open database connections",
		},
	)

	// AuthenticationAttempts tracks authentication attempts
	AuthenticationAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "authentication_attempts_total",
			Help: "Total number of authentication attempts",
		},
		[]string{"status"},
	)

	// BusinessMetrics tracks various business-related metrics
	BusinessMetrics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "business_metrics",
			Help: "Various business-related metrics",
		},
		[]string{"metric_name"},
	)
)

// RecordHTTPRequest records metrics for an HTTP request
func RecordHTTPRequest(handler, method, status string, duration float64) {
	RequestDuration.WithLabelValues(handler, method, status).Observe(duration)
	RequestsTotal.WithLabelValues(handler, method, status).Inc()
}

// RecordDatabaseOperation records metrics for a database operation
func RecordDatabaseOperation(operation string, status string, duration float64) {
	DatabaseOperationDuration.WithLabelValues(operation, status).Observe(duration)
}

// SetDatabaseConnections sets the current number of open database connections
func SetDatabaseConnections(count float64) {
	DatabaseConnectionsOpen.Set(count)
}

// RecordAuthenticationAttempt records an authentication attempt
func RecordAuthenticationAttempt(status string) {
	AuthenticationAttempts.WithLabelValues(status).Inc()
}

// RecordBusinessMetric records a business metric
func RecordBusinessMetric(metricName string, value float64) {
	BusinessMetrics.WithLabelValues(metricName).Set(value)
}
