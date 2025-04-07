package ai

import (
    "context"
    "time"
    "github.com/prometheus/client_golang/prometheus"
)

type ModelMonitor struct {
    metrics    *prometheus.Registry
    predictions *prometheus.CounterVec
    latency     *prometheus.HistogramVec
    accuracy    *prometheus.GaugeVec
    drift       *prometheus.GaugeVec
}

func NewModelMonitor() *ModelMonitor {
    registry := prometheus.NewRegistry()
    
    mm := &ModelMonitor{
        metrics: registry,
        predictions: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "model_predictions_total",
                Help: "Total number of model predictions",
            },
            []string{"model", "version", "outcome"},
        ),
        latency: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "model_prediction_latency_seconds",
                Help: "Model prediction latency in seconds",
            },
            []string{"model", "version"},
        ),
        accuracy: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Name: "model_accuracy",
                Help: "Model accuracy metrics",
            },
            []string{"model", "version", "metric"},
        ),
        drift: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Name: "model_drift",
                Help: "Model drift metrics",
            },
            []string{"model", "version", "feature"},
        ),
    }

    registry.MustRegister(mm.predictions)
    registry.MustRegister(mm.latency)
    registry.MustRegister(mm.accuracy)
    registry.MustRegister(mm.drift)

    return mm
}

func (mm *ModelMonitor) RecordPrediction(ctx context.Context, model, version, outcome string, duration time.Duration) {
    mm.predictions.WithLabelValues(model, version, outcome).Inc()
    mm.latency.WithLabelValues(model, version).Observe(duration.Seconds())
}

func (mm *ModelMonitor) UpdateAccuracy(model, version string, metrics map[string]float64) {
    for metric, value := range metrics {
        mm.accuracy.WithLabelValues(model, version, metric).Set(value)
    }
}

func (mm *ModelMonitor) UpdateDrift(model, version string, driftMetrics map[string]float64) {
    for feature, value := range driftMetrics {
        mm.drift.WithLabelValues(model, version, feature).Set(value)
    }
}