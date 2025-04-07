package validation

import (
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics handles validation-related metrics
type Metrics struct {
    validationDuration *prometheus.HistogramVec
    validationErrors   *prometheus.CounterVec
    validationScores   *prometheus.HistogramVec
    validationResults  *prometheus.CounterVec
}

// NewMetrics creates a new metrics instance
func NewMetrics() *Metrics {
    return &Metrics{
        validationDuration: promauto.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "kyc_validation_duration_seconds",
                Help: "Duration of validation operations",
                Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
            },
            []string{"validator_type"},
        ),
        validationErrors: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Name: "kyc_validation_errors_total",
                Help: "Total number of validation errors",
            },
            []string{"validator_type", "error_type"},
        ),
        validationScores: promauto.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "kyc_validation_scores",
                Help: "Distribution of validation scores",
                Buckets: []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0},
            },
            []string{"validator_type"},
        ),
        validationResults: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Name: "kyc_validation_results_total",
                Help: "Total number of validation results by outcome",
            },
            []string{"validator_type", "result"},
        ),
    }
}

// RecordValidationDuration records the duration of a validation operation
func (m *Metrics) RecordValidationDuration(validatorType string, duration time.Duration) {
    m.validationDuration.WithLabelValues(validatorType).Observe(duration.Seconds())
}

// RecordValidationError records a validation error
func (m *Metrics) RecordValidationError(validatorType, errorType string) {
    m.validationErrors.WithLabelValues(validatorType, errorType).Inc()
}

// RecordValidationScore records a validation score
func (m *Metrics) RecordValidationScore(validatorType string, score float64) {
    m.validationScores.WithLabelValues(validatorType).Observe(score)
}

// RecordValidationResult records a validation result
func (m *Metrics) RecordValidationResult(validatorType string, isValid bool) {
    result := "invalid"
    if isValid {
        result = "valid"
    }
    m.validationResults.WithLabelValues(validatorType, result).Inc()
}