package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request duration histogram
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	// HTTP request counter
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path"},
	)

	// HTTP request size histogram
	httpRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Size of HTTP requests in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path"},
	)

	// HTTP response size histogram
	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "path"},
	)

	// Active requests gauge
	activeRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Number of currently active HTTP requests",
		},
	)

	// Error rate counter
	errorRate = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors",
		},
		[]string{"method", "path", "status"},
	)
)

// PrometheusMiddleware adds Prometheus metrics collection
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Track active requests
		activeRequests.Inc()
		defer activeRequests.Dec()

		// Track request size
		if c.Request.ContentLength > 0 {
			httpRequestSize.WithLabelValues(
				c.Request.Method,
				c.FullPath(),
			).Observe(float64(c.Request.ContentLength))
		}

		// Create custom response writer to track response size
		writer := &responseWriter{
			ResponseWriter: c.Writer,
		}

		// Use custom writer
		c.Writer = writer

		// Process request
		c.Next()

		// Track response size
		if writer.size > 0 {
			httpResponseSize.WithLabelValues(
				c.Request.Method,
				c.FullPath(),
			).Observe(float64(writer.size))
		}

		// Track request duration
		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			strconv.Itoa(c.Writer.Status()),
		).Observe(duration)

		// Track total requests
		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Inc()

		// Track errors
		if c.Writer.Status() >= 400 {
			errorRate.WithLabelValues(
				c.Request.Method,
				c.FullPath(),
				strconv.Itoa(c.Writer.Status()),
			).Inc()
		}
	}
}

// responseWriter is a custom response writer that tracks response size
type responseWriter struct {
	gin.ResponseWriter
	size int64
}

func (w *responseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.size += int64(n)
	return n, err
}

// HealthCheckMiddleware adds health check endpoint
func HealthCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.JSON(200, gin.H{
				"status": "ok",
				"time":   time.Now().Format(time.RFC3339),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// MetricsMiddleware adds metrics endpoint
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" {
			prometheus.Handler().ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}
		c.Next()
	}
}

// LoggingMiddleware adds structured logging
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
		fields := map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   time.Since(start).Seconds(),
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}

		// Add request ID if available
		if requestID, exists := c.Get("RequestID"); exists {
			fields["request_id"] = requestID
		}

		// Log with appropriate level based on status code
		if c.Writer.Status() >= 500 {
			// Error level for server errors
			c.Error(fmt.Errorf("server error: %d", c.Writer.Status()))
		} else if c.Writer.Status() >= 400 {
			// Warning level for client errors
			c.Warning(fmt.Errorf("client error: %d", c.Writer.Status()))
		} else {
			// Info level for successful requests
			c.Info("request completed", fields)
		}
	}
}
