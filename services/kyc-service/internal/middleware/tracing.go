package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracingConfig holds tracing middleware configuration
type TracingConfig struct {
	ServiceName     string
	SamplingRate    float64
	ExporterEndpoint string
}

// DefaultTracingConfig returns default tracing configuration
func DefaultTracingConfig() TracingConfig {
	return TracingConfig{
		ServiceName:      "kyc-service",
		SamplingRate:     1.0,
		ExporterEndpoint: "http://jaeger:14268/api/traces",
	}
}

// TracingMiddleware provides distributed tracing
func TracingMiddleware(config TracingConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a new span
		ctx, span := otel.Tracer(config.ServiceName).Start(c.Request.Context(), fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path))
		defer span.End()

		// Add request attributes
		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.String("http.host", c.Request.Host),
			attribute.String("http.user_agent", c.Request.UserAgent()),
		)

		// Add request headers as span attributes
		for k, v := range c.Request.Header {
			span.SetAttributes(attribute.String(fmt.Sprintf("http.header.%s", k), v[0]))
		}

		// Store span in context
		c.Request = c.Request.WithContext(ctx)

		// Process request
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		// Add response attributes
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
			attribute.Int64("http.duration_ms", duration.Milliseconds()),
		)

		// Set span status
		if c.Writer.Status() >= 500 {
			span.SetStatus(codes.Error, "Internal Server Error")
		} else if c.Writer.Status() >= 400 {
			span.SetStatus(codes.Error, "Client Error")
		} else {
			span.SetStatus(codes.Ok, "Success")
		}

		// Add response headers as span attributes
		for k, v := range c.Writer.Header() {
			span.SetAttributes(attribute.String(fmt.Sprintf("http.response.header.%s", k), v[0]))
		}
	}
}

// GetTracer returns a tracer for the service
func GetTracer(serviceName string) trace.Tracer {
	return otel.Tracer(serviceName)
}

// GetSpanFromContext returns the current span from context
func GetSpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// AddSpanEvent adds an event to the current span
func AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// AddSpanError adds an error to the current span
func AddSpanError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
} 