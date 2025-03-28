package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
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

	httpRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 2, 10),
		},
		[]string{"method", "path"},
	)

	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 2, 10),
		},
		[]string{"method", "path"},
	)

	activeRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Number of currently active HTTP requests",
		},
	)
)

// MonitoringMiddleware provides monitoring and metrics collection
type MonitoringMiddleware struct {
	logger *zap.Logger
}

// NewMonitoringMiddleware creates a new MonitoringMiddleware instance
func NewMonitoringMiddleware(logger *zap.Logger) *MonitoringMiddleware {
	return &MonitoringMiddleware{
		logger: logger,
	}
}

// Metrics returns a middleware that collects HTTP metrics
func (mm *MonitoringMiddleware) Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		activeRequests.Inc()

		// Record request size
		httpRequestSize.WithLabelValues(c.Request.Method, c.FullPath()).Observe(float64(c.Request.ContentLength))

		// Create a response writer that captures the response size
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			path:          c.FullPath(),
			method:        c.Request.Method,
		}

		// Replace the response writer
		c.Writer = writer

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()
		httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), http.StatusText(status)).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)
		httpResponseSize.WithLabelValues(c.Request.Method, c.FullPath()).Observe(float64(writer.size))

		activeRequests.Dec()

		// Log slow requests
		if duration > 1.0 { // Log requests taking more than 1 second
			mm.logger.Warn("slow request detected",
				zap.String("method", c.Request.Method),
				zap.String("path", c.FullPath()),
				zap.Duration("duration", time.Since(start)),
				zap.Int("status", status),
			)
		}
	}
}

// responseWriter is a custom response writer that captures response size
type responseWriter struct {
	gin.ResponseWriter
	path   string
	method string
	size   int64
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.size += int64(len(b))
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	w.size += int64(len(s))
	return w.ResponseWriter.WriteString(s)
}

// HealthCheck returns a middleware that checks system health
func (mm *MonitoringMiddleware) HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check system resources
		// TODO: Implement system resource checks (memory, CPU, disk space)

		// Check database connection
		// TODO: Implement database health check

		// Check Redis connection
		// TODO: Implement Redis health check

		c.Next()
	}
}

// RequestLogger returns a middleware that logs request details
func (mm *MonitoringMiddleware) RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Log request
		mm.logger.Info("incoming request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.FullPath()),
			zap.String("remote_addr", c.ClientIP()),
			zap.String("user_agent", c.GetHeader("User-Agent")),
		)

		c.Next()

		// Log response
		mm.logger.Info("request completed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.FullPath()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(start)),
		)
	}
} 