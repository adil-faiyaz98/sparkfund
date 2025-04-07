package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// TracingConfig holds configuration for tracing
type TracingConfig struct {
	ServiceName string
	Enabled     bool
}

// DefaultTracingConfig returns default tracing configuration
func DefaultTracingConfig() TracingConfig {
	return TracingConfig{
		ServiceName: "kyc-service",
		Enabled:     true,
	}
}

// Metrics for monitoring
var (
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(RequestCounter)
	prometheus.MustRegister(RequestDuration)
}

// InitTracing initializes the global tracer
func InitTracing(cfg TracingConfig) (opentracing.Tracer, error) {
	if !cfg.Enabled {
		// Return a no-op tracer if tracing is disabled
		return opentracing.NoopTracer{}, nil
	}

	jaegerCfg := jaegercfg.Configuration{
		ServiceName: cfg.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	tracer, _, err := jaegerCfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)
	return tracer, nil
}

// TracingMiddleware adds distributed tracing capabilities
func TracingMiddleware(cfg TracingConfig) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		// Extract the parent span from the HTTP headers
		wireCtx, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header),
		)

		// Create a new span
		var span opentracing.Span
		if err != nil {
			// Create a new root span if no parent context
			span = opentracing.StartSpan(c.Request.URL.Path)
		} else {
			// Create a child span
			span = opentracing.StartSpan(
				c.Request.URL.Path,
				opentracing.ChildOf(wireCtx),
			)
		}
		defer span.Finish()

		// Set span tags
		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.HTTPUrl.Set(span, c.Request.URL.String())
		ext.Component.Set(span, cfg.ServiceName)

		// Set trace ID header for correlation
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			traceID := sc.TraceID().String()
			c.Header("X-Trace-ID", traceID)
		}

		// Create a new context with the span and pass it to the next handler
		ctx := opentracing.ContextWithSpan(c.Request.Context(), span)
		c.Request = c.Request.WithContext(ctx)

		// Start timer for request duration
		start := time.Now()

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := fmt.Sprintf("%d", c.Writer.Status())
		RequestCounter.WithLabelValues(c.Request.Method, c.FullPath(), status).Inc()
		RequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)

		// Set additional span tags after response
		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
		if c.Writer.Status() >= 400 {
			ext.Error.Set(span, true)
			if len(c.Errors) > 0 {
				span.SetTag("error.message", c.Errors.Last().Error())
			}
		}
	}
}

// ExtractTraceContext extracts the current tracing context from a gin context
func ExtractTraceContext(c *gin.Context) context.Context {
	return c.Request.Context()
}

// StartSpanFromContext starts a new span as a child of the parent context
func StartSpanFromContext(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
	return span, ctx
}
