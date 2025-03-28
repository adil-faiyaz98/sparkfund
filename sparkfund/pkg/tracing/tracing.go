package tracing

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"github.com/sparkfund/pkg/errors"
)

// Config represents tracing configuration
type Config struct {
	ServiceName string
	JaegerURL   string
	Environment string
}

// Init initializes tracing
func Init(cfg *Config) error {
	// Create Jaeger exporter
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.JaegerURL)))
	if err != nil {
		return errors.ErrInternalServer(err)
	}

	// Create resource with service information
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			attribute.String("environment", cfg.Environment),
		),
	)
	if err != nil {
		return errors.ErrInternalServer(err)
	}

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	return nil
}

// StartSpan starts a new span
func StartSpan(ctx context.Context, name string, opts ...sdktrace.SpanStartOption) (context.Context, sdktrace.Span) {
	tracer := otel.Tracer("")
	return tracer.Start(ctx, name, opts...)
}

// EndSpan ends a span
func EndSpan(span sdktrace.Span) {
	span.End()
}

// AddEvent adds an event to a span
func AddEvent(span sdktrace.Span, name string, attrs ...attribute.KeyValue) {
	span.AddEvent(name, sdktrace.WithAttributes(attrs...))
}

// SetAttributes sets attributes on a span
func SetAttributes(span sdktrace.Span, attrs ...attribute.KeyValue) {
	span.SetAttributes(attrs...)
}

// RecordError records an error on a span
func RecordError(span sdktrace.Span, err error, opts ...sdktrace.EventOption) {
	span.RecordError(err, opts...)
	span.SetStatus(sdktrace.Error, err.Error())
}

// WithTimeout creates a context with timeout
func WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}

// WithCancel creates a context with cancel
func WithCancel(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(ctx)
}

// WithDeadline creates a context with deadline
func WithDeadline(ctx context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(ctx, deadline)
}

// WithValue creates a context with value
func WithValue(ctx context.Context, key, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

// FromContext gets a span from context
func FromContext(ctx context.Context) (sdktrace.Span, bool) {
	return sdktrace.SpanFromContext(ctx)
}

// SpanContext gets span context from context
func SpanContext(ctx context.Context) (sdktrace.SpanContext, bool) {
	span, ok := FromContext(ctx)
	if !ok {
		return sdktrace.SpanContext{}, false
	}
	return span.SpanContext(), true
}

// TraceID gets trace ID from context
func TraceID(ctx context.Context) (string, bool) {
	span, ok := FromContext(ctx)
	if !ok {
		return "", false
	}
	return span.SpanContext().TraceID().String(), true
}

// SpanID gets span ID from context
func SpanID(ctx context.Context) (string, bool) {
	span, ok := FromContext(ctx)
	if !ok {
		return "", false
	}
	return span.SpanContext().SpanID().String(), true
}

// ParentSpanID gets parent span ID from context
func ParentSpanID(ctx context.Context) (string, bool) {
	span, ok := FromContext(ctx)
	if !ok {
		return "", false
	}
	return span.Parent().SpanID().String(), true
}

// IsSampled checks if span is sampled
func IsSampled(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.SpanContext().IsSampled()
}

// IsRemote checks if span is remote
func IsRemote(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.SpanContext().IsRemote()
}

// IsValid checks if span is valid
func IsValid(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.SpanContext().IsValid()
}

// IsDebug checks if span is debug
func IsDebug(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.SpanContext().IsDebug()
}

// IsDeferred checks if span is deferred
func IsDeferred(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.SpanContext().IsDeferred()
}

// IsTraceValid checks if trace is valid
func IsTraceValid(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.SpanContext().TraceID().IsValid()
}

// IsSpanValid checks if span is valid
func IsSpanValid(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.SpanContext().SpanID().IsValid()
}

// IsParentValid checks if parent is valid
func IsParentValid(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.Parent().SpanID().IsValid()
}

// IsTraceSampled checks if trace is sampled
func IsTraceSampled(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.SpanContext().TraceFlags().IsSampled()
}

// IsTraceDebug checks if trace is debug
func IsTraceDebug(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.SpanContext().TraceFlags().IsDebug()
}

// IsTraceDeferred checks if trace is deferred
func IsTraceDeferred(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.SpanContext().TraceFlags().IsDeferred()
}

// IsTraceRemote checks if trace is remote
func IsTraceRemote(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.SpanContext().IsRemote()
}

// IsTraceValid checks if trace is valid
func IsTraceValid(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.SpanContext().TraceID().IsValid()
}

// IsSpanValid checks if span is valid
func IsSpanValid(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.SpanContext().SpanID().IsValid()
}

// IsParentValid checks if parent is valid
func IsParentValid(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.Parent().SpanID().IsValid()
}

// IsTraceSampled checks if trace is sampled
func IsTraceSampled(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.SpanContext().TraceFlags().IsSampled()
}

// IsTraceDebug checks if trace is debug
func IsTraceDebug(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.SpanContext().TraceFlags().IsDebug()
}

// IsTraceDeferred checks if trace is deferred
func IsTraceDeferred(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.SpanContext().TraceFlags().IsDeferred()
}

// IsTraceRemote checks if trace is remote
func IsTraceRemote(ctx context.Context) bool {
	span, ok := FromContext(ctx)
	if !ok {
		return false
	}
	return span.SpanContext().IsRemote()
} 