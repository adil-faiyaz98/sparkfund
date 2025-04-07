package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
)

// CircuitBreakerConfig holds configuration for the circuit breaker
type CircuitBreakerConfig struct {
	Enabled             bool
	Timeout             time.Duration
	MaxConcurrentReqs   uint32
	ErrorThresholdPerc  int
	RequestVolumeThresh uint64
	SleepWindow         time.Duration
}

// DefaultCircuitBreakerConfig returns default circuit breaker configuration
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		Enabled:             true,
		Timeout:             30 * time.Second,
		MaxConcurrentReqs:   100,
		ErrorThresholdPerc:  50,
		RequestVolumeThresh: 20,
		SleepWindow:         5 * time.Second,
	}
}

var (
	circuitBreakersMutex sync.Mutex
	circuitBreakers      = make(map[string]*gobreaker.CircuitBreaker)
)

// CircuitBreakerMiddleware provides circuit breaker functionality for endpoints
func CircuitBreakerMiddleware(cfg CircuitBreakerConfig) gin.HandlerFunc {
	if !cfg.Enabled {
		// Return a pass-through middleware if circuit breakers are disabled
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// Get or create circuit breaker for this path
		cb := getCircuitBreaker(path, cfg)

		// Execute the request through the circuit breaker
		_, err := cb.Execute(func() (interface{}, error) {
			// Store the response writer
			originalWriter := c.Writer

			// Create a custom response writer to capture the status
			blw := &bodyLogWriter{ResponseWriter: originalWriter, status: 200}
			c.Writer = blw

			// Process request
			c.Next()

			// Check if there are any errors
			if len(c.Errors) > 0 {
				return nil, c.Errors.Last()
			}

			// If status code is an error, return an error to the circuit breaker
			if blw.status >= 500 {
				return nil, &ErrorResponse{Error: "Service error"}
			}

			return nil, nil
		})

		// If circuit is open, return service unavailable
		if err != nil {
			if cb.State() == gobreaker.StateOpen {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, ErrorResponse{
					Error: "Service temporarily unavailable, please try again later",
				})
				return
			}
		}
	}
}

// getCircuitBreaker gets or creates a circuit breaker for a specific path
func getCircuitBreaker(path string, cfg CircuitBreakerConfig) *gobreaker.CircuitBreaker {
	circuitBreakersMutex.Lock()
	defer circuitBreakersMutex.Unlock()

	if cb, found := circuitBreakers[path]; found {
		return cb
	}

	settings := gobreaker.Settings{
		Name:        path,
		MaxRequests: uint32(cfg.MaxConcurrentReqs),
		Interval:    cfg.SleepWindow,
		Timeout:     cfg.SleepWindow,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= uint32(cfg.RequestVolumeThresh) &&
				failureRatio >= float64(cfg.ErrorThresholdPerc)/100
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			// This could log the state change or emit metrics
		},
	}

	cb := gobreaker.NewCircuitBreaker(settings)
	circuitBreakers[path] = cb
	return cb
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// bodyLogWriter is a custom response writer that captures the status code
type bodyLogWriter struct {
	gin.ResponseWriter
	status int
}

// WriteHeader captures the status code
func (w *bodyLogWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
