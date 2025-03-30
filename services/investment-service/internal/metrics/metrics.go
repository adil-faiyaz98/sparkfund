package metrics

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// RequestCounter counts HTTP requests
	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "investment_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// RequestDuration measures request duration
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "investment_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// DatabaseQueryCounter counts database queries
	DatabaseQueryCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "investment_db_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"operation"},
	)

	// DatabaseQueryDuration measures database query duration
	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "investment_db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
		},
		[]string{"operation"},
	)

	// BusinessErrorCounter counts business logic errors
	BusinessErrorCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "investment_business_errors_total",
			Help: "Total number of business logic errors",
		},
		[]string{"type"},
	)

	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "investment_service_active_connections",
			Help: "Number of active connections",
		},
	)
)

func init() {
	prometheus.MustRegister(RequestCounter, RequestDuration, DatabaseQueryCounter, DatabaseQueryDuration, BusinessErrorCounter, ActiveConnections)
}

// RegisterMetrics registers the metrics endpoint
func RegisterMetrics(router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// MetricsMiddleware collects metrics for each request
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Increment active connections
		ActiveConnections.Inc()
		defer ActiveConnections.Dec()

		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime).Seconds()

		// Record metrics
		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = "unknown"
		}

		status := c.Writer.Status()
		method := c.Request.Method

		RequestCounter.WithLabelValues(method, endpoint, fmt.Sprintf("%d", status)).Inc()
		RequestDuration.WithLabelValues(method, endpoint).Observe(duration)
	}
}

// TrackRequestDuration returns a function to track request duration
func TrackRequestDuration(method, endpoint string) func() {
	startTime := time.Now()
	return func() {
		duration := time.Since(startTime).Seconds()
		RequestDuration.WithLabelValues(method, endpoint).Observe(duration)
	}
}

// TrackDBQuery returns a function to track database query duration
func TrackDBQuery(operation string) func() {
	startTime := time.Now()
	return func() {
		duration := time.Since(startTime).Seconds()
		DatabaseQueryDuration.WithLabelValues(operation).Observe(duration)
		DatabaseQueryCounter.WithLabelValues(operation).Inc()
	}
}

// RecordBusinessError records a business logic error
func RecordBusinessError(errorType string) {
	BusinessErrorCounter.WithLabelValues(errorType).Inc()
}
