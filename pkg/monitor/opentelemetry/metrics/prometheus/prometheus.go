package prometheus

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	prometheusmetric "go.opentelemetry.io/otel/exporters/prometheus"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
)

type PrometheusExporterBuiler struct{}

func (p *PrometheusExporterBuiler) Build(
	ctx context.Context,
	c *controller.Controller,
) (aggregation.TemporalitySelector, error) {
	return prometheusmetric.New(prometheusmetric.Config{
		Registerer: prometheus.DefaultRegisterer,
		Gatherer:   prometheus.DefaultGatherer,
	}, c)
}
