package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sparkfund/pkg/errors"
)

// Config represents metrics configuration
type Config struct {
	Enabled bool
	Port    int
}

// Init initializes metrics
func Init(cfg *Config) error {
	if !cfg.Enabled {
		return nil
	}

	// Start metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != nil {
			errors.ErrInternalServer(err)
		}
	}()

	return nil
}

// HTTP metrics
var (
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path"},
	)

	HTTPRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Size of HTTP requests in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path"},
	)

	HTTPResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path"},
	)

	HTTPRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being served",
		},
	)
)

// Database metrics
var (
	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query_type"},
	)

	DatabaseConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections",
			Help: "Current number of database connections",
		},
	)

	DatabaseErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_errors_total",
			Help: "Total number of database errors",
		},
		[]string{"error_type"},
	)
)

// Cache metrics
var (
	CacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	CacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
	)

	CacheSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cache_size_bytes",
			Help: "Current size of cache in bytes",
		},
	)
)

// Business metrics
var (
	ActiveUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users",
			Help: "Current number of active users",
		},
	)

	TotalTransactions = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "transactions_total",
			Help: "Total number of transactions",
		},
	)

	TransactionAmount = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "transaction_amount",
			Help:    "Amount of transactions",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
	)
)

// System metrics
var (
	CPUUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percent",
			Help: "Current CPU usage percentage",
		},
	)

	MemoryUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Current memory usage in bytes",
		},
	)

	Goroutines = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "goroutines",
			Help: "Current number of goroutines",
		},
	)
)

// TrackHTTPRequest tracks an HTTP request
func TrackHTTPRequest(method, path string, status int, duration time.Duration) {
	HTTPRequestDuration.WithLabelValues(method, path, string(status)).Observe(duration.Seconds())
	HTTPRequestsTotal.WithLabelValues(method, path).Inc()
}

// TrackDatabaseQuery tracks a database query
func TrackDatabaseQuery(queryType string, duration time.Duration) {
	DatabaseQueryDuration.WithLabelValues(queryType).Observe(duration.Seconds())
}

// TrackDatabaseError tracks a database error
func TrackDatabaseError(errorType string) {
	DatabaseErrors.WithLabelValues(errorType).Inc()
}

// TrackCacheHit tracks a cache hit
func TrackCacheHit() {
	CacheHits.Inc()
}

// TrackCacheMiss tracks a cache miss
func TrackCacheMiss() {
	CacheMisses.Inc()
}

// SetActiveUsers sets the number of active users
func SetActiveUsers(count int) {
	ActiveUsers.Set(float64(count))
}

// TrackTransaction tracks a transaction
func TrackTransaction(amount float64) {
	TotalTransactions.Inc()
	TransactionAmount.Observe(amount)
}

// SetSystemMetrics sets system metrics
func SetSystemMetrics(cpuPercent, memoryBytes, goroutineCount float64) {
	CPUUsage.Set(cpuPercent)
	MemoryUsage.Set(memoryBytes)
	Goroutines.Set(goroutineCount)
} 