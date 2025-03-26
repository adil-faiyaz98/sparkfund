package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics represents the metrics collector
type Metrics struct {
	RequestCounter   *prometheus.CounterVec
	RequestDuration  *prometheus.HistogramVec
	ResponseSize     *prometheus.HistogramVec
	RequestsInFlight prometheus.Gauge
}

// NewAutoRegisterMetrics creates and registers metrics with Prometheus
func NewAutoRegisterMetrics() *Metrics {
	metrics := &Metrics{
		RequestCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "security_http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "security_http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
		ResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "security_http_response_size_bytes",
				Help:    "HTTP response size in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "endpoint"},
		),
		RequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "security_http_requests_in_flight",
				Help: "Current number of HTTP requests in flight",
			},
		),
	}

	return metrics
}

// MetricsHandler returns an HTTP handler for metrics endpoint
func (m *Metrics) MetricsHandler() gin.HandlerFunc {
	handler := promhttp.Handler()
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
