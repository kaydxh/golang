package meter

import (
	"context"

	"go.opentelemetry.io/otel/sdk/metric/export"
)

type ExporterBuilder interface {
	Build(ctx context.Context) (export.Exporter, error)
}
