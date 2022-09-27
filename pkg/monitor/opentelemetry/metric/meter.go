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

	"go.opentelemetry.io/otel/metric/global"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

type MeterOptinos struct {
	PushExporterBuilder PushExporterBuilder
	PullExporterBuilder PullExporterBuilder
	collectPeriod       time.Duration
}

type Meter struct {
	Controller controller.Controller
	opts       MeterOptinos
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

//https://github.com/open-telemetry/opentelemetry-go/blob/example/prometheus/v0.30.0/example/prometheus/main.go
//https://github.com/kaydxh/newrelic-opentelemetry-examples/blob/main/go/metrics.go
func (m *Meter) Install(ctx context.Context) (err error) {

	var metricControllerOptions []controller.Option
	if m.opts.collectPeriod > 0 {
		metricControllerOptions = append(
			metricControllerOptions,
			controller.WithCollectPeriod(m.opts.collectPeriod),
		)
	}

	if m.opts.PushExporterBuilder != nil {
		exporter, err := m.createPushExporter(ctx)
		if err != nil {
			return err
		}
		metricControllerOptions = append(metricControllerOptions, controller.WithExporter(exporter))

	}

	c := controller.New(
		processor.NewFactory(
			simple.NewWithHistogramDistribution(),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
		metricControllerOptions...,
	)
	if m.opts.PullExporterBuilder != nil {
		_, err = m.createPullExporter(ctx, c)
		if err != nil {
			return err
		}
	}

	err = c.Start(ctx)
	if err != nil {
		return err
	}

	global.SetMeterProvider(c)

	return nil
}

func (m *Meter) createPushExporter(ctx context.Context) (export.Exporter, error) {
	if m.opts.PushExporterBuilder == nil {
		return nil, fmt.Errorf("push metric exporter builder is nil")
	}

	return m.opts.PushExporterBuilder.Build(ctx)
}

// Pull Exporter supports Prometheus pulls.  It does not implement the
// sdk/export/metric.Exporter interface--instead it creates a pull
// controller and reads the latest checkpointed data on-scrape.
func (m *Meter) createPullExporter(ctx context.Context, c *controller.Controller,
) (aggregation.TemporalitySelector, error) {
	if m.opts.PullExporterBuilder == nil {
		return nil, fmt.Errorf("pull metric exporter builder is nil")
	}

	return m.opts.PullExporterBuilder.Build(ctx, c)
}
