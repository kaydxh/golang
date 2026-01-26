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
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	prometheusmetric "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

const (
	defaultMetricsUrl = "/metrics"
)

// Global registry and handler for Prometheus metrics
var (
	globalRegistry     *prometheus.Registry
	globalHandler      http.Handler
	globalRegistryOnce sync.Once
)

// GetGlobalRegistry returns the global Prometheus registry.
// Creates it on first call.
func GetGlobalRegistry() *prometheus.Registry {
	globalRegistryOnce.Do(func() {
		globalRegistry = prometheus.NewRegistry()
		// Register Go runtime metrics and process metrics
		globalRegistry.MustRegister(prometheus.NewGoCollector())
		globalRegistry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
		globalHandler = promhttp.HandlerFor(globalRegistry, promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		})
	})
	return globalRegistry
}

// GetMetricsHandler returns the HTTP handler for /metrics endpoint.
// Must be called after GetGlobalRegistry() or Build().
func GetMetricsHandler() http.Handler {
	GetGlobalRegistry() // Ensure initialized
	return globalHandler
}

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

func (p *PrometheusExporterBuilder) Build(ctx context.Context) (metric.Reader, error) {
	// Use global registry to ensure HTTP handler can access metrics
	return NewPrometheusReader(ctx, GetGlobalRegistry())
}

// GetUrlPath returns the configured URL path for metrics endpoint.
func (p *PrometheusExporterBuilder) GetUrlPath() string {
	return p.opts.Url
}

// NewPrometheusReader creates a new Prometheus exporter with the given registry.
func NewPrometheusReader(ctx context.Context, reg prometheus.Registerer) (*prometheusmetric.Exporter, error) {
	return prometheusmetric.New(prometheusmetric.WithRegisterer(reg))
}
