package meter

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
	ExporterBuilder     ExporterBuilder
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

	if m.opts.ExporterBuilder != nil {
		exporter, err := m.createExporter(ctx)
		if err != nil {
			return err
		}
		metricControllerOptions = append(metricControllerOptions, controller.WithExporter(exporter))

		if m.opts.collectPeriod > 0 {
			metricControllerOptions = append(
				metricControllerOptions,
				controller.WithCollectPeriod(m.opts.collectPeriod),
			)
		}
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

func (m *Meter) createExporter(ctx context.Context) (export.Exporter, error) {
	if m.opts.ExporterBuilder == nil {
		return nil, fmt.Errorf("exporter builder is nil")
	}

	return m.opts.ExporterBuilder.Build(ctx)
}

// Pull Exporter supports Prometheus pulls.  It does not implement the
// sdk/export/metric.Exporter interface--instead it creates a pull
// controller and reads the latest checkpointed data on-scrape.
func (m *Meter) createPullExporter(ctx context.Context, c *controller.Controller,
) (aggregation.TemporalitySelector, error) {
	if m.opts.PullExporterBuilder == nil {
		return nil, fmt.Errorf("pull exporter builder is nil")
	}

	return m.opts.PullExporterBuilder.Build(ctx, c)
}
