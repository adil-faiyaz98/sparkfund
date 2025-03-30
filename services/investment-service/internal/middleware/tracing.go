package middleware

import (
	"context"
	"fmt"
	"time"

	"investment-service/internal/metrics"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// InitTracing initializes the global tracer
func InitTracing(serviceName string) (opentracing.Tracer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	tracer, _, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)
	return tracer, nil
}

// TracingMiddleware adds distributed tracing capabilities
func TracingMiddleware() gin.HandlerFunc {
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
		ext.Component.Set(span, "investment-service")

		// Set trace ID header for correlation
		traceID := span.Context().(jaeger.SpanContext).TraceID().String()
		c.Header("X-Trace-ID", traceID)

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
		metrics.RequestCounter.WithLabelValues(c.Request.Method, c.FullPath(), status).Inc()
		metrics.RequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)

		// Set additional span tags after response
		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
		if c.Writer.Status() >= 400 {
			ext.Error.Set(span, true)
			span.SetTag("error.message", c.Errors.Last().Error())
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
