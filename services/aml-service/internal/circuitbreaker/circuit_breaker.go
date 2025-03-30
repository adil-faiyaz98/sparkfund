package circuitbreaker

import (
	"time"

	"github.com/sony/gobreaker"
)

// CircuitBreakerConfig holds configuration for circuit breaker
type CircuitBreakerConfig struct {
	Name          string
	MaxRequests   uint32
	Interval      time.Duration
	Timeout       time.Duration
	ReadyToTrip   func(counts gobreaker.Counts) bool
	OnStateChange func(name string, from gobreaker.State, to gobreaker.State)
}

// NewCircuitBreaker creates a new circuit breaker with the given configuration
func NewCircuitBreaker(config CircuitBreakerConfig) *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:          config.Name,
		MaxRequests:   config.MaxRequests,
		Interval:      config.Interval,
		Timeout:       config.Timeout,
		ReadyToTrip:   config.ReadyToTrip,
		OnStateChange: config.OnStateChange,
	})
}

// DefaultCircuitBreaker creates a circuit breaker with default settings
func DefaultCircuitBreaker(name string) *gobreaker.CircuitBreaker {
	return NewCircuitBreaker(CircuitBreakerConfig{
		Name:        name,
		MaxRequests: 100,
		Interval:    time.Minute,
		Timeout:     time.Minute * 2,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 10 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			// Log state changes or emit metrics
		},
	})
}

// ExternalServiceCircuitBreaker creates a circuit breaker for external service calls
func ExternalServiceCircuitBreaker(name string) *gobreaker.CircuitBreaker {
	return NewCircuitBreaker(CircuitBreakerConfig{
		Name:        name,
		MaxRequests: 50,
		Interval:    time.Minute * 5,
		Timeout:     time.Minute * 10,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio >= 0.5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			// Log state changes or emit metrics
		},
	})
}

// DatabaseCircuitBreaker creates a circuit breaker for database operations
func DatabaseCircuitBreaker(name string) *gobreaker.CircuitBreaker {
	return NewCircuitBreaker(CircuitBreakerConfig{
		Name:        name,
		MaxRequests: 200,
		Interval:    time.Minute * 2,
		Timeout:     time.Minute * 5,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 20 && failureRatio >= 0.7
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			// Log state changes or emit metrics
		},
	})
}
