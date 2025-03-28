package metrics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
}

// RegisterMetrics registers the metrics endpoint
func RegisterMetrics(router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// RecordRequest records an HTTP request
func RecordRequest(method, endpoint string, status int) {
	httpRequestsTotal.WithLabelValues(method, endpoint, string(status)).Inc()
} 