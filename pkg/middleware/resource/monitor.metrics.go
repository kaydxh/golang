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
	// metricGroup is the scope name for meter, can be overridden by SetMetricGroup
	// For ZhiYan, common values: "server_report", "client_report", "default"
	metricGroup   = instrumentationName
	metricGroupMu sync.RWMutex

	meter     metric.Meter
	meterOnce sync.Once
)

// SetMetricGroup sets the metric group (scope name) for meter
// Should be called before any metrics are created
// For ZhiYan, use "server_report" or "client_report"
func SetMetricGroup(group string) {
	metricGroupMu.Lock()
	defer metricGroupMu.Unlock()
	if group != "" {
		metricGroup = group
	}
}

// GetMetricGroup returns the current metric group
func GetMetricGroup() string {
	metricGroupMu.RLock()
	defer metricGroupMu.RUnlock()
	return metricGroup
}

func getMeter() metric.Meter {
	meterOnce.Do(func() {
		metricGroupMu.RLock()
		group := metricGroup
		metricGroupMu.RUnlock()
		meter = otel.GetMeterProvider().Meter(
			group,
			metric.WithInstrumentationVersion(instrumentationVersion),
		)
	})
	return meter
}

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
	DefaultMetricMonitor *MetricMonitor
	defaultMonitorOnce   sync.Once
)

func GetDefaultMetricMonitor() *MetricMonitor {
	defaultMonitorOnce.Do(func() {
		DefaultMetricMonitor = NewMetricMonitor()
	})
	return DefaultMetricMonitor
}

func GlobalMeter() metric.Meter {
	return getMeter()
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
		m.TotalReqCounter, err = getMeter().Int64Counter("total_req")
	})
	call(func() {
		m.FailCntCounter, err = getMeter().Int64Counter("fail_cnt")
	})
	call(func() {
		m.CostTimeHistogram, err = getMeter().Float64Histogram("cost_time")
	})
	if err != nil {
		otel.Handle(err)
	}

	return m
}

func (m *MetricMonitor) GetOrNewBusinessCounter(key string) (metric.Int64Counter, error) {
	m.businessCountersMu.Lock()
	defer m.businessCountersMu.Unlock()
	counter, ok := m.BusinessCounters[key]
	if ok {
		return counter, nil
	}

	counter, err := getMeter().Int64Counter(key)
	if err != nil {
		return nil, err
	}
	m.BusinessCounters[key] = counter
	return counter, nil
}

func (m *MetricMonitor) GetOrNewBusinessHistogram(key string) (metric.Float64Histogram, error) {
	m.businessHistogramMu.Lock()
	defer m.businessHistogramMu.Unlock()
	histogram, ok := m.BusinessHistogram[key]
	if ok {
		return histogram, nil
	}

	histogram, err := getMeter().Float64Histogram(key)
	if err != nil {
		return nil, err
	}
	m.BusinessHistogram[key] = histogram
	return histogram, nil
}

func ReportMetric(ctx context.Context, dim Dimension, costTime time.Duration) {
	attrs := ExtractAttrsWithContext(ctx)
	attrs = append(attrs, Attrs(dim)...)

	// Use GetDefaultMetricMonitor() to ensure lazy initialization
	monitor := GetDefaultMetricMonitor()
	monitor.TotalReqCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	if dim.Error != nil {
		monitor.FailCntCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
	monitor.CostTimeHistogram.Record(ctx, float64(costTime.Milliseconds()), metric.WithAttributes(attrs...))
	ReportBusinessMetric(ctx, attrs)
}
