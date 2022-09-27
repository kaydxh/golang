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

	"github.com/kaydxh/golang/pkg/monitor/opentelemetry/metric"
	"github.com/kaydxh/golang/pkg/monitor/opentelemetry/tracer"
)

func WithMeterPushExporter(pushExporterBuilder metric.PushExporterBuilder) OpenTelemetryOption {
	return OpenTelemetryOptionFunc(func(o *OpenTelemetry) {
		o.opts.meterOptions = append(o.opts.meterOptions, metric.WithPushExporter(pushExporterBuilder))

	})
}

func WithMeterPullExporter(pullExporterBuilder metric.PullExporterBuilder) OpenTelemetryOption {
	return OpenTelemetryOptionFunc(func(o *OpenTelemetry) {
		o.opts.meterOptions = append(o.opts.meterOptions, metric.WithPullExporter(pullExporterBuilder))
	})
}

func WithMetricCollectDuration(period time.Duration) OpenTelemetryOption {
	return OpenTelemetryOptionFunc(func(o *OpenTelemetry) {
		o.opts.meterOptions = append(o.opts.meterOptions, metric.WithCollectPeriod(period))
	})
}

func WithTracerExporter(exporterBuilder tracer.TracerExporterBuilder) OpenTelemetryOption {
	return OpenTelemetryOptionFunc(func(o *OpenTelemetry) {
		o.opts.tracerOptions = append(o.opts.tracerOptions, tracer.WithExporterBuilder(exporterBuilder))
	})
}
