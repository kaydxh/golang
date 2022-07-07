package metric

import (
	"context"

	"go.opentelemetry.io/otel/sdk/metric/export"
)

type PushExporterBuilder interface {
	Build(ctx context.Context) (export.Exporter, error)
}
