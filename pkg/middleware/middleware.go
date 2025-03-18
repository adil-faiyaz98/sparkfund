package middleware

import (
	"context"
	"net/http"
	"time"

	"money-pulse/pkg/auth"
	"money-pulse/pkg/logger"
	"money-pulse/pkg/monitoring"

	"golang.org/x/time/rate"
)

// Middleware represents a collection of middleware components
type Middleware struct {
	log  *logger.Logger
	auth *auth.TokenManager
}

// NewMiddleware creates a new middleware instance
func NewMiddleware(log *logger.Logger, auth *auth.TokenManager) *Middleware {
	return &Middleware{
		log:  log,
		auth: auth,
	}
}

// LoggingMiddleware logs HTTP request details
func (m *Middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer that captures the status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Process the request
		next.ServeHTTP(rw, r)

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Log request details
		m.log.WithFields(map[string]interface{}{
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     rw.statusCode,
			"duration":   duration,
			"remote_ip":  r.RemoteAddr,
			"user_agent": r.UserAgent(),
		}).Info("HTTP request completed")

		// Record metrics
		monitoring.RecordHTTPRequest(
			r.URL.Path,
			r.Method,
			http.StatusText(rw.statusCode),
			duration,
		)
	})
}

// AuthMiddleware handles authentication and authorization
func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			monitoring.RecordAuthenticationAttempt("missing_token")
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		claims, err := m.auth.ValidateToken(token)
		if err != nil {
			monitoring.RecordAuthenticationAttempt("invalid_token")
			m.log.WithError(err).Error("Invalid authentication token")
			http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
			return
		}

		// Add claims to request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "roles", claims.Roles)

		monitoring.RecordAuthenticationAttempt("success")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RateLimitMiddleware implements rate limiting
func (m *Middleware) RateLimitMiddleware(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(100), 200) // 100 requests per second with burst of 200
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			m.log.WithFields(map[string]interface{}{
				"remote_ip": r.RemoteAddr,
				"path":      r.URL.Path,
			}).Warn("Rate limit exceeded")
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// responseWriter is a custom response writer that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Chain combines multiple middleware into a single middleware
func Chain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}
