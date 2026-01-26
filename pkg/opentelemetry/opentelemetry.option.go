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
package opentelemetry

import (
	"time"

	"github.com/kaydxh/golang/pkg/opentelemetry/metric"
	"github.com/kaydxh/golang/pkg/opentelemetry/tracer"
	"go.opentelemetry.io/otel/sdk/resource"
)

func WithMeterPushExporter(pushExporterBuilder metric.PushExporterBuilder) OpenTelemetryServiceOption {
	return OpenTelemetryServiceOptionFunc(func(o *OpenTelemetryService) {
		o.opts.meterOptions = append(o.opts.meterOptions, metric.WithPushExporter(pushExporterBuilder))

	})
}

func WithMeterPullExporter(pullExporterBuilder metric.PullExporterBuilder) OpenTelemetryServiceOption {
	return OpenTelemetryServiceOptionFunc(func(o *OpenTelemetryService) {
		o.opts.meterOptions = append(o.opts.meterOptions, metric.WithPullExporter(pullExporterBuilder))
	})
}

func WithMetricCollectDuration(period time.Duration) OpenTelemetryServiceOption {
	return OpenTelemetryServiceOptionFunc(func(o *OpenTelemetryService) {
		o.opts.meterOptions = append(o.opts.meterOptions, metric.WithCollectPeriod(period))
	})
}

func WithTracerExporter(exporterBuilder tracer.TracerExporterBuilder) OpenTelemetryServiceOption {
	return OpenTelemetryServiceOptionFunc(func(o *OpenTelemetryService) {
		o.opts.tracerOptions = append(o.opts.tracerOptions, tracer.WithExporterBuilder(exporterBuilder))
	})
}

// WithResource sets a custom resource for both tracer and meter
func WithResource(res *resource.Resource) OpenTelemetryServiceOption {
	return OpenTelemetryServiceOptionFunc(func(o *OpenTelemetryService) {
		o.opts.meterOptions = append(o.opts.meterOptions, metric.WithResource(res))
		o.opts.tracerOptions = append(o.opts.tracerOptions, tracer.WithResource(res))
	})
}

// ========================================
// App MeterProvider Options (双 Provider 支持)
// ========================================

// WithAppMeterPushExporter sets the push exporter for App MeterProvider
func WithAppMeterPushExporter(pushExporterBuilder metric.PushExporterBuilder) OpenTelemetryServiceOption {
	return OpenTelemetryServiceOptionFunc(func(o *OpenTelemetryService) {
		o.opts.appMeterOptions = append(o.opts.appMeterOptions, metric.WithPushExporter(pushExporterBuilder))
	})
}

// WithAppMeterPullExporter sets the pull exporter for App MeterProvider
func WithAppMeterPullExporter(pullExporterBuilder metric.PullExporterBuilder) OpenTelemetryServiceOption {
	return OpenTelemetryServiceOptionFunc(func(o *OpenTelemetryService) {
		o.opts.appMeterOptions = append(o.opts.appMeterOptions, metric.WithPullExporter(pullExporterBuilder))
	})
}

// WithAppMetricCollectDuration sets the collect duration for App MeterProvider
func WithAppMetricCollectDuration(period time.Duration) OpenTelemetryServiceOption {
	return OpenTelemetryServiceOptionFunc(func(o *OpenTelemetryService) {
		o.opts.appMeterOptions = append(o.opts.appMeterOptions, metric.WithCollectPeriod(period))
	})
}

// WithAppMeterResource sets a custom resource for App MeterProvider
func WithAppMeterResource(res *resource.Resource) OpenTelemetryServiceOption {
	return OpenTelemetryServiceOptionFunc(func(o *OpenTelemetryService) {
		o.opts.appMeterOptions = append(o.opts.appMeterOptions, metric.WithResource(res))
	})
}
