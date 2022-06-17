package prometheus

import (
	"context"
	"fmt"
	"net/http"

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
	exporter, err := prometheusmetric.New(prometheusmetric.Config{
		Registerer: prometheus.DefaultRegisterer,
		Gatherer:   prometheus.DefaultGatherer,
	}, c)
	if err != nil {
		return nil, fmt.Errorf("new prometheusmetric err: %v", err)
	}

	http.HandleFunc("/", exporter.ServeHTTP)
	return exporter, nil
}
