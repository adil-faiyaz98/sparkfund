package tracing

import (
    "context"
    "github.com/opentracing/opentracing-go"
    "github.com/uber/jaeger-client-go"
    "github.com/uber/jaeger-client-go/config"
)

type Tracer struct {
    tracer opentracing.Tracer
}

func NewTracer(serviceName string) (*Tracer, error) {
    cfg := &config.Configuration{
        ServiceName: serviceName,
        Sampler: &config.SamplerConfig{
            Type:  "const",
            Param: 1,
        },
        Reporter: &config.ReporterConfig{
            LogSpans: true,
        },
    }

    tracer, _, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
    if err != nil {
        return nil, err
    }

    opentracing.SetGlobalTracer(tracer)
    return &Tracer{tracer: tracer}, nil
}

func (t *Tracer) StartSpan(name string) opentracing.Span {
    return t.tracer.StartSpan(name)
}

func (t *Tracer) StartSpanFromContext(ctx context.Context, name string) (opentracing.Span, context.Context) {
    return opentracing.StartSpanFromContext(ctx, name)
}