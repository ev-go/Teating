package tracer

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	enabled = true
	level   = 0
	tp      *tracesdk.TracerProvider //nolint:gochecknoglobals
	tracer  trace.Tracer             //nolint:gochecknoglobals
)

func Init(url string, traceLevel int, name, environment string) error {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}

	tp = tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
			attribute.String("environment", environment),
		)),
	)
	tracer = tp.Tracer(name)
	level = traceLevel

	return nil
}

func Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	for _, opt := range opts {
		cfg := trace.NewSpanStartConfig(opt)
		attributes := cfg.Attributes()
		for _, attr := range attributes {
			if attr.Key == "traceLevel" {
				if attr.Value.AsInt64() < int64(level) {
					return ctx, nil
				}
			}
		}

	}

	if tp == nil || !enabled {
		return ctx, nil
	}

	return tracer.Start(ctx, spanName, opts...)
}

func End(span trace.Span) {
	if tp == nil || !enabled {
		return
	}

	span.End()
}

func Error(span trace.Span, err error) {
	if tp == nil || !enabled {
		return
	}

	span.SetStatus(codes.Error, "")
	span.RecordError(err)
}

func SetEnabled(v bool) {
	enabled = v
}

func SetTraceLevel(l int) {
	level = l
}
