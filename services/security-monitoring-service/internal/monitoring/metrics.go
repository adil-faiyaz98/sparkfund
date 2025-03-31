package monitoring

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sparkfund/security-monitoring/internal/models"
)

// Metrics represents the monitoring metrics
type Metrics struct {
	mu sync.RWMutex

	// Security metrics
	threatsTotal    *prometheus.CounterVec
	intrusionsTotal *prometheus.CounterVec
	malwareTotal    *prometheus.CounterVec
	patternsTotal   *prometheus.CounterVec
	riskScore       *prometheus.GaugeVec
	confidenceScore *prometheus.GaugeVec
	alertCount      *prometheus.CounterVec
	eventCount      *prometheus.CounterVec

	// Performance metrics
	processingTime *prometheus.HistogramVec
	queueSize      *prometheus.GaugeVec
	errorCount     *prometheus.CounterVec
	modelLatency   *prometheus.HistogramVec

	// System metrics
	memoryUsage       *prometheus.GaugeVec
	cpuUsage          *prometheus.GaugeVec
	goroutineCount    *prometheus.GaugeVec
	activeConnections *prometheus.GaugeVec
}

// NewMetrics creates a new metrics instance
func NewMetrics() *Metrics {
	m := &Metrics{}

	// Initialize security metrics
	m.threatsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "security_threats_total",
			Help: "Total number of security threats detected",
		},
		[]string{"type", "severity"},
	)

	m.intrusionsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "security_intrusions_total",
			Help: "Total number of intrusion attempts detected",
		},
		[]string{"type", "severity"},
	)

	m.malwareTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "security_malware_total",
			Help: "Total number of malware instances detected",
		},
		[]string{"type", "severity"},
	)

	m.patternsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "security_patterns_total",
			Help: "Total number of security patterns detected",
		},
		[]string{"type", "severity"},
	)

	m.riskScore = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "security_risk_score",
			Help: "Current security risk score",
		},
		[]string{"source"},
	)

	m.confidenceScore = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "security_confidence_score",
			Help: "Current confidence score for security analysis",
		},
		[]string{"source"},
	)

	m.alertCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "security_alerts_total",
			Help: "Total number of security alerts generated",
		},
		[]string{"type", "severity"},
	)

	m.eventCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "security_events_total",
			Help: "Total number of security events processed",
		},
		[]string{"type", "source"},
	)

	// Initialize performance metrics
	m.processingTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "security_processing_time_seconds",
			Help:    "Time taken to process security events",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"type"},
	)

	m.queueSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "security_queue_size",
			Help: "Current size of the security event queue",
		},
		[]string{"type"},
	)

	m.errorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "security_errors_total",
			Help: "Total number of security processing errors",
		},
		[]string{"type", "error"},
	)

	m.modelLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "security_model_latency_seconds",
			Help:    "Latency of AI model predictions",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"model"},
	)

	// Initialize system metrics
	m.memoryUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "security_memory_usage_bytes",
			Help: "Current memory usage of the security service",
		},
		[]string{"type"},
	)

	m.cpuUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "security_cpu_usage_percent",
			Help: "Current CPU usage of the security service",
		},
		[]string{"type"},
	)

	m.goroutineCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "security_goroutines_total",
			Help: "Current number of goroutines",
		},
		[]string{"type"},
	)

	m.activeConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "security_active_connections",
			Help: "Current number of active connections",
		},
		[]string{"type"},
	)

	// Register all metrics
	prometheus.MustRegister(m.threatsTotal)
	prometheus.MustRegister(m.intrusionsTotal)
	prometheus.MustRegister(m.malwareTotal)
	prometheus.MustRegister(m.patternsTotal)
	prometheus.MustRegister(m.riskScore)
	prometheus.MustRegister(m.confidenceScore)
	prometheus.MustRegister(m.alertCount)
	prometheus.MustRegister(m.eventCount)
	prometheus.MustRegister(m.processingTime)
	prometheus.MustRegister(m.queueSize)
	prometheus.MustRegister(m.errorCount)
	prometheus.MustRegister(m.modelLatency)
	prometheus.MustRegister(m.memoryUsage)
	prometheus.MustRegister(m.cpuUsage)
	prometheus.MustRegister(m.goroutineCount)
	prometheus.MustRegister(m.activeConnections)

	return m
}

// UpdateMetrics updates metrics based on security analysis
func (m *Metrics) UpdateMetrics(analysis *models.SecurityAnalysis) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Update threat metrics
	for _, threat := range analysis.Threats {
		m.threatsTotal.WithLabelValues(threat.Type, threat.Severity).Inc()
	}

	// Update intrusion metrics
	for _, intrusion := range analysis.Intrusions {
		m.intrusionsTotal.WithLabelValues(intrusion.Type, intrusion.Severity).Inc()
	}

	// Update malware metrics
	for _, malware := range analysis.Malware {
		m.malwareTotal.WithLabelValues(malware.Type, malware.Severity).Inc()
	}

	// Update pattern metrics
	for _, pattern := range analysis.Patterns {
		m.patternsTotal.WithLabelValues(pattern.Type, "unknown").Inc()
	}

	// Update risk and confidence scores
	m.riskScore.WithLabelValues("overall").Set(analysis.RiskScore)
	m.confidenceScore.WithLabelValues("overall").Set(analysis.Confidence)
}

// RecordProcessingTime records the time taken to process an event
func (m *Metrics) RecordProcessingTime(eventType string, duration time.Duration) {
	m.processingTime.WithLabelValues(eventType).Observe(duration.Seconds())
}

// RecordModelLatency records the latency of AI model predictions
func (m *Metrics) RecordModelLatency(model string, duration time.Duration) {
	m.modelLatency.WithLabelValues(model).Observe(duration.Seconds())
}

// UpdateQueueSize updates the current queue size
func (m *Metrics) UpdateQueueSize(queueType string, size int) {
	m.queueSize.WithLabelValues(queueType).Set(float64(size))
}

// RecordError records a processing error
func (m *Metrics) RecordError(errorType, errorMessage string) {
	m.errorCount.WithLabelValues(errorType, errorMessage).Inc()
}

// UpdateSystemMetrics updates system-related metrics
func (m *Metrics) UpdateSystemMetrics(memoryBytes, cpuPercent float64, goroutines, connections int) {
	m.memoryUsage.WithLabelValues("total").Set(memoryBytes)
	m.cpuUsage.WithLabelValues("total").Set(cpuPercent)
	m.goroutineCount.WithLabelValues("total").Set(float64(goroutines))
	m.activeConnections.WithLabelValues("total").Set(float64(connections))
}

// Reset resets all metrics
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.threatsTotal.Reset()
	m.intrusionsTotal.Reset()
	m.malwareTotal.Reset()
	m.patternsTotal.Reset()
	m.riskScore.Reset()
	m.confidenceScore.Reset()
	m.alertCount.Reset()
	m.eventCount.Reset()
	m.processingTime.Reset()
	m.queueSize.Reset()
	m.errorCount.Reset()
	m.modelLatency.Reset()
	m.memoryUsage.Reset()
	m.cpuUsage.Reset()
	m.goroutineCount.Reset()
	m.activeConnections.Reset()
}
