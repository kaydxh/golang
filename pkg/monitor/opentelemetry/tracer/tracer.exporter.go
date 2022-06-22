package tracer

import (
	"context"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type TracerExporterBuilder interface {
	Build(ctx context.Context) (sdktrace.SpanExporter, error)
}
