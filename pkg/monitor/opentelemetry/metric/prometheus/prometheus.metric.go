package prometheus

import (
	"context"
	"fmt"
	"net/http"
	url_ "net/url"

	"github.com/prometheus/client_golang/prometheus"
	prometheusmetric "go.opentelemetry.io/otel/exporters/prometheus"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
)

const (
	defaultMetricsUrl = "noop://localhost/metrics"
)

type PrometheusExporterBuilderOptions struct {
	Url string
}

type PrometheusExporterBuilder struct {
	opts PrometheusExporterBuilderOptions
}

func defaultBuilderOptions() PrometheusExporterBuilderOptions {
	return PrometheusExporterBuilderOptions{
		Url: defaultMetricsUrl,
	}
}

func NewPrometheusExporterBuilder(opts ...PrometheusExporterBuilderOption) *PrometheusExporterBuilder {

	builder := &PrometheusExporterBuilder{
		opts: defaultBuilderOptions(),
	}
	builder.ApplyOptions(opts...)

	return builder
}

func (p *PrometheusExporterBuilder) Build(
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

	u, err := url_.Parse(p.opts.Url)
	if err != nil {
		return nil, fmt.Errorf("parse url: %v, err: %v", p.opts.Url, err)
	}

	http.HandleFunc(u.Path, exporter.ServeHTTP)
	return exporter, nil
}
