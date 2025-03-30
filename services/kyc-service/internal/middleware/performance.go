package middleware

import (
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CacheConfig holds cache configuration
type CacheConfig struct {
	MaxAge         time.Duration
	NoStore        bool
	NoCache        bool
	MustRevalidate bool
	Private        bool
}

// DefaultCacheConfig returns default cache configuration
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		MaxAge:         24 * time.Hour,
		NoStore:        false,
		NoCache:        false,
		MustRevalidate: false,
		Private:        false,
	}
}

// CacheMiddleware adds cache control headers
func CacheMiddleware(config CacheConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip caching for non-GET requests
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Build cache control header
		var directives []string
		if config.NoStore {
			directives = append(directives, "no-store")
		}
		if config.NoCache {
			directives = append(directives, "no-cache")
		}
		if config.MustRevalidate {
			directives = append(directives, "must-revalidate")
		}
		if config.Private {
			directives = append(directives, "private")
		}
		if config.MaxAge > 0 {
			directives = append(directives, "max-age="+string(config.MaxAge.Seconds()))
		}

		if len(directives) > 0 {
			c.Header("Cache-Control", strings.Join(directives, ", "))
		}

		c.Next()
	}
}

// GzipMiddleware compresses responses using gzip
func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip compression for small responses
		if c.Request.Header.Get("Accept-Encoding") == "" {
			c.Next()
			return
		}

		// Check if client supports gzip
		if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		// Create gzip writer
		gz := gzip.NewWriter(c.Writer)
		defer gz.Close()

		// Create custom response writer
		writer := &gzipWriter{
			ResponseWriter: c.Writer,
			writer:         gz,
		}

		// Set headers
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		// Use custom writer
		c.Writer = writer

		c.Next()
	}
}

// gzipWriter is a custom response writer that compresses the response
type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

// TimeoutMiddleware adds a timeout to the request context
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create timeout context
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Update request context
		c.Request = c.Request.WithContext(ctx)

		// Create done channel
		done := make(chan struct{})
		go func() {
			c.Next()
			close(done)
		}()

		// Wait for request to complete or timeout
		select {
		case <-done:
			// Request completed successfully
		case <-ctx.Done():
			// Request timed out
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"error": "Request timeout",
			})
		}
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()

		// Add request ID to context
		c.Set("RequestID", requestID)

		// Add request ID to response header
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// ResponseSizeMiddleware limits the size of responses
func ResponseSizeMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create custom response writer
		writer := &sizeLimitWriter{
			ResponseWriter: c.Writer,
			maxSize:        maxSize,
		}

		// Use custom writer
		c.Writer = writer

		c.Next()
	}
}

// sizeLimitWriter is a custom response writer that limits response size
type sizeLimitWriter struct {
	gin.ResponseWriter
	maxSize int64
	written int64
}

func (s *sizeLimitWriter) Write(data []byte) (int, error) {
	if s.written+int64(len(data)) > s.maxSize {
		return 0, fmt.Errorf("response size limit exceeded")
	}
	n, err := s.ResponseWriter.Write(data)
	s.written += int64(n)
	return n, err
}
