package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Collector handles specific KYC-related metrics collection
type Collector struct {
	requestCounter   *prometheus.CounterVec
	requestDuration  *prometheus.HistogramVec
	mlPredictions    *prometheus.CounterVec
	mlLatency        *prometheus.HistogramVec
	documentVerified *prometheus.CounterVec
}

func NewCollector() *Collector {
	return &Collector{
		requestCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "kyc_requests_total",
				Help: "Total number of KYC requests",
			},
			[]string{"endpoint", "status"},
		),
		requestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "kyc_request_duration_seconds",
				Help:    "Request duration in seconds",
				Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"endpoint"},
		),
		mlPredictions: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "kyc_ml_predictions_total",
				Help: "Total number of ML predictions",
			},
			[]string{"model", "outcome"},
		),
		mlLatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "kyc_ml_prediction_duration_seconds",
				Help:    "ML prediction duration in seconds",
				Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
			},
			[]string{"model"},
		),
		documentVerified: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "kyc_documents_verified_total",
				Help: "Total number of verified documents",
			},
			[]string{"type", "status"},
		),
	}
}

// RecordRequest records KYC-specific request metrics
func (c *Collector) RecordRequest(endpoint, status string) {
	c.requestCounter.WithLabelValues(endpoint, status).Inc()
}

// RecordRequestDuration records KYC-specific request duration
func (c *Collector) RecordRequestDuration(endpoint string, duration float64) {
	c.requestDuration.WithLabelValues(endpoint).Observe(duration)
}

// RecordMLPrediction records ML model prediction metrics
func (c *Collector) RecordMLPrediction(model, outcome string) {
	c.mlPredictions.WithLabelValues(model, outcome).Inc()
}

// RecordMLLatency records ML model prediction latency
func (c *Collector) RecordMLLatency(model string, duration float64) {
	c.mlLatency.WithLabelValues(model).Observe(duration)
}

// RecordDocumentVerification records document verification metrics
func (c *Collector) RecordDocumentVerification(docType, status string) {
	c.documentVerified.WithLabelValues(docType, status).Inc()
}

// GetCollector returns a singleton instance of the Collector
var defaultCollector *Collector

func GetCollector() *Collector {
	if defaultCollector == nil {
		defaultCollector = NewCollector()
	}
	return defaultCollector
}
