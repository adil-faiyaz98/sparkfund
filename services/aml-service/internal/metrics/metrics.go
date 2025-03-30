package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all the prometheus metrics for the service
type Metrics struct {
	transactionProcessTime prometheus.Histogram
	errors                 *prometheus.CounterVec
}

// New creates a new Metrics instance
func New() *Metrics {
	return &Metrics{
		transactionProcessTime: promauto.NewHistogram(prometheus.HistogramOpts{
			Name: "aml_transaction_process_time_seconds",
			Help: "Time taken to process an AML transaction",
			Buckets: []float64{
				0.1, 0.5, 1, 2.5, 5, 10,
			},
		}),
		errors: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "aml_errors_total",
			Help: "Total number of errors by type",
		}, []string{"type"}),
	}
}

// RecordTransactionProcessTime records the time taken to process a transaction
func (m *Metrics) RecordTransactionProcessTime(seconds float64) {
	m.transactionProcessTime.Observe(seconds)
}

// RecordError increments the error counter for a given error type
func (m *Metrics) RecordError(errorType string) {
	m.errors.WithLabelValues(errorType).Inc()
}

func NewMetrics() *Metrics {
	return &Metrics{
		TransactionsProcessed: promauto.NewCounter(prometheus.CounterOpts{
			Name: "aml_transactions_processed_total",
			Help: "The total number of processed transactions",
		}),
		TransactionProcessTime: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "aml_transaction_process_duration_seconds",
			Help:    "Time taken to process a transaction",
			Buckets: prometheus.DefBuckets,
		}),
		HighRiskTransactions: promauto.NewCounter(prometheus.CounterOpts{
			Name: "aml_high_risk_transactions_total",
			Help: "The total number of high-risk transactions detected",
		}),
		AlertsGenerated: promauto.NewCounter(prometheus.CounterOpts{
			Name: "aml_alerts_generated_total",
			Help: "The total number of alerts generated",
		}),
		CircuitBreakerState: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "aml_circuit_breaker_state",
			Help: "The current state of circuit breakers (0: Open, 1: Half-Open, 2: Closed)",
		}, []string{"name"}),
		DatabaseOperationTime: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "aml_database_operation_duration_seconds",
			Help:    "Time taken for database operations",
			Buckets: prometheus.DefBuckets,
		}),
		ExternalServiceLatency: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "aml_external_service_latency_seconds",
			Help:    "Latency of external service calls",
			Buckets: prometheus.DefBuckets,
		}),
		ErrorCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "aml_errors_total",
			Help: "The total number of errors by type",
		}, []string{"type"}),
	}
}

// RecordTransactionProcessed increments the transactions processed counter
func (m *Metrics) RecordTransactionProcessed() {
	m.TransactionsProcessed.Inc()
}

// RecordTransactionProcessTime records the time taken to process a transaction
func (m *Metrics) RecordTransactionProcessTime(duration float64) {
	m.TransactionProcessTime.Observe(duration)
}

// RecordHighRiskTransaction increments the high-risk transactions counter
func (m *Metrics) RecordHighRiskTransaction() {
	m.HighRiskTransactions.Inc()
}

// RecordAlertGenerated increments the alerts generated counter
func (m *Metrics) RecordAlertGenerated() {
	m.AlertsGenerated.Inc()
}

// SetCircuitBreakerState sets the state of a circuit breaker
func (m *Metrics) SetCircuitBreakerState(name string, state float64) {
	m.CircuitBreakerState.WithLabelValues(name).Set(state)
}

// RecordDatabaseOperationTime records the time taken for a database operation
func (m *Metrics) RecordDatabaseOperationTime(duration float64) {
	m.DatabaseOperationTime.Observe(duration)
}

// RecordExternalServiceLatency records the latency of an external service call
func (m *Metrics) RecordExternalServiceLatency(duration float64) {
	m.ExternalServiceLatency.Observe(duration)
}

// RecordError increments the error counter for a specific error type
func (m *Metrics) RecordError(errorType string) {
	m.ErrorCounter.WithLabelValues(errorType).Inc()
}
