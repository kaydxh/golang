package resource

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	instrumentationName    = "github/kaydxh/golang/pkg/middleware/resource"
	instrumentationVersion = "v0.0.1"
)

var (
	meter = otel.GetMeterProvider().Meter(
		instrumentationName,
		metric.WithInstrumentationVersion(instrumentationVersion),
	)
)

type MetricMonitor struct {
	TotalReqCounter   metric.Int64Counter
	FailCntCounter    metric.Int64Counter
	CostTimeHistogram metric.Float64Histogram

	BusinessCounters   map[string]metric.Int64Counter
	businessCountersMu sync.RWMutex

	BusinessHistogram   map[string]metric.Float64Histogram
	businessHistogramMu sync.RWMutex
}

var (
	DefaultMetricMonitor = NewMetricMonitor()
)

func GlobalMeter() metric.Meter {
	return meter
}

func NewMetricMonitor() *MetricMonitor {
	var err error
	m := &MetricMonitor{
		BusinessCounters:  make(map[string]metric.Int64Counter, 0),
		BusinessHistogram: make(map[string]metric.Float64Histogram, 0),
	}
	call := func(f func()) {
		if err != nil {
			return
		}
		f()
	}
	call(func() {
		m.TotalReqCounter, err = meter.Int64Counter("total_req")
	})
	call(func() {
		m.FailCntCounter, err = meter.Int64Counter("fail_cnt")
	})
	call(func() {
		m.CostTimeHistogram, err = meter.Float64Histogram("cost_time")
	})
	if err != nil {
		otel.Handle(err)
	}

	return m
}

func (m *MetricMonitor) GetOrNewBusinessCounter(key string) (metric.Int64Counter, error) {
	m.businessCountersMu.Lock()
	defer m.businessCountersMu.Unlock()
	counter, ok := DefaultMetricMonitor.BusinessCounters[key]
	if ok {
		return counter, nil
	}

	counter, err := meter.Int64Counter(key)
	if err != nil {
		return nil, err
	}
	DefaultMetricMonitor.BusinessCounters[key] = counter
	return counter, nil
}

func (m *MetricMonitor) GetOrNewBusinessHistogram(key string) (metric.Float64Histogram, error) {
	m.businessHistogramMu.Lock()
	defer m.businessHistogramMu.Unlock()
	histogram, ok := DefaultMetricMonitor.BusinessHistogram[key]
	if ok {
		return histogram, nil
	}

	histogram, err := meter.Float64Histogram(key)
	if err != nil {
		return nil, err
	}
	DefaultMetricMonitor.BusinessHistogram[key] = histogram
	return histogram, nil
}

func ReportMetric(ctx context.Context, dim Dimension, costTime time.Duration) {
	attrs := ExtractAttrsWithContext(ctx)
	attrs = append(attrs, Attrs(dim)...)

	DefaultMetricMonitor.TotalReqCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	if dim.Error != nil {
		DefaultMetricMonitor.FailCntCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
	DefaultMetricMonitor.CostTimeHistogram.Record(ctx, float64(costTime.Milliseconds()), metric.WithAttributes(attrs...))
	ReportBusinessMetric(ctx, attrs)
}
