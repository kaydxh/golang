package tracer

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

type TracerOptions struct {
	builer TracerExporterBuilder
}

type Tracer struct {
	opts TracerOptions
}

func NewTracer(opts ...TracerOption) *Tracer {
	t := &Tracer{}
	t.ApplyOptions(opts...)

	return t
}

//https://github.com/open-telemetry/opentelemetry-go/blob/main/example/jaeger/main.go
func (t *Tracer) Install(ctx context.Context) (err error) {
	exp, err := t.createExporter(ctx)
	if err != nil {
		return err
	}

	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
		)),
	)

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	return nil
}

func (t *Tracer) createExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	if t.opts.builer == nil {
		return nil, fmt.Errorf("exporter builder is nil")
	}

	return t.opts.builer.Build(ctx)
}
