package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// KYCOperations tracks the number of KYC operations
	KYCOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kyc_operations_total",
			Help: "Total number of KYC operations by type",
		},
		[]string{"operation", "status"},
	)

	// KYCOperationDuration tracks the duration of KYC operations
	KYCOperationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "kyc_operation_duration_seconds",
			Help:    "Duration of KYC operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// KYCStatusCount tracks the number of KYC records by status
	KYCStatusCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kyc_status_count",
			Help: "Number of KYC records by status",
		},
		[]string{"status"},
	)

	// DatabaseErrors tracks database operation errors
	DatabaseErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kyc_database_errors_total",
			Help: "Total number of database errors",
		},
		[]string{"operation"},
	)

	// HTTPRequests tracks HTTP requests
	HTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kyc_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HTTPRequestDuration tracks HTTP request duration
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "kyc_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	prometheus.MustRegister(KYCOperations)
	prometheus.MustRegister(KYCOperationDuration)
	prometheus.MustRegister(KYCStatusCount)
	prometheus.MustRegister(DatabaseErrors)
	prometheus.MustRegister(HTTPRequests)
	prometheus.MustRegister(HTTPRequestDuration)
} 