package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
	"investment-service/internal/config"
	"investment-service/internal/models"
)

var (
	circuitBreakersMutex sync.Mutex
	circuitBreakers      = make(map[string]*gobreaker.CircuitBreaker)
)

// CircuitBreakerMiddleware provides circuit breaker functionality for endpoints
func CircuitBreakerMiddleware() gin.HandlerFunc {
	cfg := config.Get().CircuitBreaker

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
				return nil, models.ErrorResponse{Error: "Service error"}
			}

			return nil, nil
		})

		// If circuit is open, return service unavailable
		if err != nil {
			if cb.State() == gobreaker.StateOpen {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, models.ErrorResponse{
					Error: "Service temporarily unavailable, please try again later",
				})
				return
			}
		}
	}
}

// getCircuitBreaker gets or creates a circuit breaker for a specific path
func getCircuitBreaker(path string, cfg config.CircuitBreaker) *gobreaker.CircuitBreaker {
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
