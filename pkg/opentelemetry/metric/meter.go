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
package metric

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
)

type MeterOptinos struct {
	PushExporterBuilder PushExporterBuilder
	PullExporterBuilder PullExporterBuilder
	collectPeriod       time.Duration
}

type Meter struct {
	opts MeterOptinos
}

func defaultMeterOptions() MeterOptinos {
	return MeterOptinos{
		collectPeriod: time.Minute,
	}
}

func NewMeter(opts ...MeterOption) *Meter {
	m := &Meter{
		opts: defaultMeterOptions(),
	}
	m.ApplyOptions(opts...)
	return m
}

// https://github.com/open-telemetry/opentelemetry-go/blob/example/prometheus/v0.30.0/example/prometheus/main.go
// https://github.com/kaydxh/newrelic-opentelemetry-examples/blob/main/go/metrics.go
func (m *Meter) Install(ctx context.Context) (err error) {
	var readers []metric.Reader

	if m.opts.PushExporterBuilder != nil {
		exporter, err := m.createPushExporter(ctx)
		if err != nil {
			return err
		}
		if exporter != nil { // such as prometheus, that's a puller actually
			var opts []metric.PeriodicReaderOption
			if m.opts.collectPeriod > 0 {
				opts = append(opts, metric.WithInterval(m.opts.collectPeriod))
			}
			reader := metric.NewPeriodicReader(exporter, opts...)
			readers = append(readers, reader)
		}
	}
	if m.opts.PullExporterBuilder != nil {
		reader, err := m.createPullReader(ctx)
		if err != nil {
			return err
		}
		readers = append(readers, reader)
	}
	var metricProviderOptions []metric.Option

	for _, reader := range readers {
		metricProviderOptions = append(metricProviderOptions, metric.WithReader(reader)) // 默认cumulative
	}

	sdk := metric.NewMeterProvider(metricProviderOptions...)
	otel.SetMeterProvider(sdk)
	return nil
}

func (m *Meter) createPushExporter(ctx context.Context) (metric.Exporter, error) {
	if m.opts.PushExporterBuilder == nil {
		return nil, fmt.Errorf("push metric reader builder is nil")
	}

	return m.opts.PushExporterBuilder.Build(ctx)
}

// Pull Exporter supports Prometheus pulls.  It does not implement the
// sdk/export/metric.Exporter interface--instead it creates a pull
// controller and reads the latest checkpointed data on-scrape.
func (m *Meter) createPullReader(ctx context.Context,
) (metric.Reader, error) {
	if m.opts.PullExporterBuilder == nil {
		return nil, fmt.Errorf("pull metric exporter builder is nil")
	}

	return m.opts.PullExporterBuilder.Build(ctx)
}
