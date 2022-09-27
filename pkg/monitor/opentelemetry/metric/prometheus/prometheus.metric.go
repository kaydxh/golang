/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
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
