package circuitbreaker

import (
	"errors"
	"sync"
	"time"

	"investment-service/internal/config"
)

// State represents the circuit breaker state
type State int

const (
	StateClosed   State = iota // Circuit is closed, requests are allowed
	StateOpen                  // Circuit is open, requests are not allowed
	StateHalfOpen              // Circuit is half-open, allowing test requests
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	name            string
	state           State
	failureCount    int
	lastFailureTime time.Time
	mutex           sync.RWMutex
	maxFailures     int
	timeout         time.Duration
	interval        time.Duration
}

// ErrCircuitOpen is returned when the circuit is open
var ErrCircuitOpen = errors.New("circuit breaker is open")

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(name string) *CircuitBreaker {
	cfg := config.Get().CircuitBreaker

	return &CircuitBreaker{
		name:        name,
		state:       StateClosed,
		maxFailures: cfg.MaxErrors,
		timeout:     cfg.Timeout,
		interval:    cfg.Interval,
	}
}

// Execute runs the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !config.Get().CircuitBreaker.Enabled {
		return fn() // Circuit breaker disabled, just execute the function
	}

	cb.mutex.RLock()
	if cb.state == StateOpen {
		// Check if timeout has elapsed
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.mutex.RUnlock()
			cb.mutex.Lock()
			cb.state = StateHalfOpen
			cb.mutex.Unlock()
		} else {
			cb.mutex.RUnlock()
			return ErrCircuitOpen
		}
	} else {
		cb.mutex.RUnlock()
	}

	err := fn()

	if err != nil {
		cb.mutex.Lock()
		defer cb.mutex.Unlock()

		// Record failure
		cb.failureCount++
		cb.lastFailureTime = time.Now()

		// Check if circuit should open
		if cb.state == StateHalfOpen || (cb.state == StateClosed && cb.failureCount >= cb.maxFailures) {
			cb.state = StateOpen
		}

		return err
	}

	// Success - reset if in half-open state
	if cb.state == StateHalfOpen {
		cb.mutex.Lock()
		cb.state = StateClosed
		cb.failureCount = 0
		cb.mutex.Unlock()
	}

	return nil
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.state = StateClosed
	cb.failureCount = 0
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() State {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	return cb.state
}
