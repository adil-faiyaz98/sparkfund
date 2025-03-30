package resilience

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sparkfund/services/investment-service/internal/config"
	"github.com/sparkfund/services/investment-service/internal/metrics"
)

// CircuitState represents the state of a circuit breaker
type CircuitState int

const (
	// Closed means the circuit breaker allows requests
	Closed CircuitState = iota
	// Open means the circuit breaker blocks requests
	Open
	// HalfOpen means the circuit breaker is testing if the service is healthy
	HalfOpen
)

var (
	// ErrCircuitOpen is returned when the circuit is open
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	name            string
	state           CircuitState
	failures        int
	successes       int
	lastStateChange time.Time
	mutex           sync.RWMutex
	config          config.Config
	logger          *logrus.Logger
	metrics         *metrics.MetricsCollector
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(name string, cfg config.Config, logger *logrus.Logger, metrics *metrics.MetricsCollector) *CircuitBreaker {
	return &CircuitBreaker{
		name:            name,
		state:           Closed,
		failures:        0,
		successes:       0,
		lastStateChange: time.Now(),
		config:          cfg,
		logger:          logger,
		metrics:         metrics,
	}
}

// Execute executes a function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, operation string, fn func(context.Context) error) error {
	// Check if circuit is open
	if !cb.AllowRequest() {
		cb.metrics.CircuitBreakerReject(cb.name, operation)
		return ErrCircuitOpen
	}

	// Track time
	start := time.Now()

	// Execute the function
	err := fn(ctx)

	// Record result
	if err != nil {
		cb.RecordFailure()
		cb.metrics.CircuitBreakerFailure(cb.name, operation, time.Since(start).Seconds())
	} else {
		cb.RecordSuccess()
		cb.metrics.CircuitBreakerSuccess(cb.name, operation, time.Since(start).Seconds())
	}

	return err
}

// AllowRequest checks if a request should be allowed
func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	switch cb.state {
	case Closed:
		return true
	case Open:
		// Check if timeout has elapsed to move to half-open
		if time.Since(cb.lastStateChange) > cb.config.Resilience.CircuitBreaker.Timeout {
			cb.mutex.RUnlock()
			cb.toHalfOpen()
			cb.mutex.RLock()
			return true
		}
		return false
	case HalfOpen:
		// In half-open, we allow a limited number of requests
		return true
	default:
		return true
	}
}

// RecordSuccess records a successful request
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if cb.state == HalfOpen {
		cb.successes++
		if cb.successes >= int(cb.config.Resilience.CircuitBreaker.SuccessThreshold) {
			cb.toClose()
		}
	}
}

// RecordFailure records a failed request
func (cb *CircuitBreaker) RecordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if cb.state == Closed {
		cb.failures++
		if cb.failures >= int(cb.config.Resilience.CircuitBreaker.ErrorThreshold) {
			cb.toOpen()
		}
	} else if cb.state == HalfOpen {
		cb.toOpen()
	}
}

// toClose changes the state to closed
func (cb *CircuitBreaker) toClose() {
	cb.state = Closed
	cb.failures = 0
	cb.successes = 0
	cb.lastStateChange = time.Now()
	cb.logger.WithField("circuit", cb.name).Info("Circuit breaker state changed to CLOSED")
}

// toOpen changes the state to open
func (cb *CircuitBreaker) toOpen() {
	cb.state = Open
	cb.failures = 0
	cb.successes = 0
	cb.lastStateChange = time.Now()
	cb.logger.WithField("circuit", cb.name).Warn("Circuit breaker state changed to OPEN")
}

// toHalfOpen changes the state to half-open
func (cb *CircuitBreaker) toHalfOpen() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.state = HalfOpen
	cb.failures = 0
	cb.successes = 0
	cb.lastStateChange = time.Now()
	cb.logger.WithField("circuit", cb.name).Info("Circuit breaker state changed to HALF-OPEN")
}
