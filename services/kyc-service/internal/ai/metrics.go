package ai

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    modelPredictionLatency = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "kyc_ai_model_prediction_latency_seconds",
            Help: "Latency of AI model predictions",
            Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
        },
        []string{"model_name", "model_version"},
    )

    modelPredictionErrors = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "kyc_ai_model_prediction_errors_total",
            Help: "Total number of AI model prediction errors",
        },
        []string{"model_name", "error_type"},
    )
)