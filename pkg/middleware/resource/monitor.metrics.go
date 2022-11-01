package resource

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
)

const (
	instrumentationName    = "github/kaydxh/golang/pkg/middleware/resource"
	instrumentationVersion = "v0.0.1"
)

var (
	meter = global.MeterProvider().Meter(
		instrumentationName,
		metric.WithInstrumentationVersion(instrumentationVersion),
	)
)

type MetricMonitor struct {
	TotalReqCounter   syncint64.Counter
	FailCntCounter    syncint64.Counter
	CostTimeHistogram syncfloat64.Histogram
}

var (
	DefaultMetricMonitor = NewMetricMonitor()
)

func NewMetricMonitor() *MetricMonitor {
	var err error
	m := &MetricMonitor{}
	call := func(f func()) {
		if err != nil {
			return
		}
		f()
	}
	call(func() {
		m.TotalReqCounter, err = meter.SyncInt64().Counter("total_req")
	})
	call(func() {
		m.FailCntCounter, err = meter.SyncInt64().Counter("fail_cnt")
	})
	call(func() {
		m.CostTimeHistogram, err = meter.SyncFloat64().Histogram("cost_time")
	})
	if err != nil {
		otel.Handle(err)
	}

	return m
}

func ReportMetric(ctx context.Context, dim Dimension, costTime time.Duration) {
	attrs := Attrs(dim)
	DefaultMetricMonitor.TotalReqCounter.Add(ctx, 1, attrs...)
	if dim.Error != nil {
		DefaultMetricMonitor.FailCntCounter.Add(ctx, 1, attrs...)
	}
	DefaultMetricMonitor.CostTimeHistogram.Record(ctx, float64(costTime.Milliseconds()), attrs...)
}
