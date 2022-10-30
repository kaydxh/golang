package resource

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
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
	TotalReqCounter syncint64.Counter
	SuccCntCounter  syncint64.Counter
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
		m.SuccCntCounter, err = meter.SyncInt64().Counter("succ_cnt")
	})
	if err != nil {
		otel.Handle(err)
	}

	return m
}
