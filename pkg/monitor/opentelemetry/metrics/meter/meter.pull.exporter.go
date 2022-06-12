package meter

import (
	"context"

	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
)

type PullExporterBuilder interface {
	Build(ctx context.Context, c *controller.Controller) (aggregation.TemporalitySelector, error)
}
