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
package report

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

const (
	// Meter names
	ServerReportMeterName = "server_report" // 被调上报
	ClientReportMeterName = "client_report" // 主调上报

	// Metric names
	RequestsMetricName  = "requests"  // 请求数
	SuccessMetricName   = "success"   // 成功数
	TimeoutMetricName   = "timeout"   // 超时数
	AbnormalMetricName  = "abnormal"  // 异常数
	CostMetricName      = "cost"      // 耗时(ms)
	DurationMetricName  = "duration"  // 耗时分布 Histogram
)

// Attribute keys for callee (被调方)
var (
	RetCodeKey    = attribute.Key("ret_code")    // 返回码
	PIpKey        = attribute.Key("p_ip")        // 被调 IP
	PAppKey       = attribute.Key("p_app")       // 被调应用名
	PServerKey    = attribute.Key("p_server")    // 被调服务名
	PServiceKey   = attribute.Key("p_service")   // 被调 Service
	PInterfaceKey = attribute.Key("p_interface") // 被调接口
)

// Attribute keys for caller (主调方)
var (
	AIpKey        = attribute.Key("a_ip")        // 主调 IP
	AAppKey       = attribute.Key("a_app")       // 主调应用名
	AServerKey    = attribute.Key("a_server")    // 主调服务名
	AServiceKey   = attribute.Key("a_service")   // 主调 Service
	AInterfaceKey = attribute.Key("a_interface") // 主调接口
)

// Default histogram bounds for duration (in milliseconds)
var DefaultDurationBounds = []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}

// MetricReporter handles metric reporting
type MetricReporter struct {
	meterProvider otelmetric.MeterProvider

	counters   map[string]otelmetric.Int64Counter
	histograms map[string]otelmetric.Float64Histogram
	mu         sync.RWMutex
}

var (
	globalReporter     *MetricReporter
	globalReporterOnce sync.Once
)

// GetGlobalReporter returns the global metric reporter
func GetGlobalReporter() *MetricReporter {
	globalReporterOnce.Do(func() {
		globalReporter = NewMetricReporter(otel.GetMeterProvider())
	})
	return globalReporter
}

// NewMetricReporter creates a new metric reporter
func NewMetricReporter(provider otelmetric.MeterProvider) *MetricReporter {
	if provider == nil {
		provider = otel.GetMeterProvider()
	}
	return &MetricReporter{
		meterProvider: provider,
		counters:      make(map[string]otelmetric.Int64Counter),
		histograms:    make(map[string]otelmetric.Float64Histogram),
	}
}

func (r *MetricReporter) getCounter(meterName, metricName string) otelmetric.Int64Counter {
	key := fmt.Sprintf("%s_%s", meterName, metricName)

	r.mu.RLock()
	if counter, ok := r.counters[key]; ok {
		r.mu.RUnlock()
		return counter
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()

	// Double check
	if counter, ok := r.counters[key]; ok {
		return counter
	}

	meter := r.meterProvider.Meter(meterName)
	counter, err := meter.Int64Counter(metricName)
	if err != nil {
		otel.Handle(err)
		return noop.Int64Counter{}
	}
	r.counters[key] = counter
	return counter
}

func (r *MetricReporter) getHistogram(meterName, metricName string, bounds []float64) otelmetric.Float64Histogram {
	key := fmt.Sprintf("%s_%s", meterName, metricName)

	r.mu.RLock()
	if histogram, ok := r.histograms[key]; ok {
		r.mu.RUnlock()
		return histogram
	}
	r.mu.RUnlock()

	r.mu.Lock()
	defer r.mu.Unlock()

	// Double check
	if histogram, ok := r.histograms[key]; ok {
		return histogram
	}

	meter := r.meterProvider.Meter(meterName)
	var opts []otelmetric.Float64HistogramOption
	if len(bounds) > 0 {
		opts = append(opts, otelmetric.WithExplicitBucketBoundaries(bounds...))
	}
	histogram, err := meter.Float64Histogram(metricName, opts...)
	if err != nil {
		otel.Handle(err)
		return noop.Float64Histogram{}
	}
	r.histograms[key] = histogram
	return histogram
}

// ReportServerMetric reports server-side (被调) metrics
func (r *MetricReporter) ReportServerMetric(ctx context.Context, dim *ServerDimension, costMs float64) {
	attrs := dim.ToAttributes()

	// requests
	r.getCounter(ServerReportMeterName, RequestsMetricName).Add(ctx, 1, otelmetric.WithAttributes(attrs...))

	// success / timeout / abnormal
	if dim.IsSuccess() {
		r.getCounter(ServerReportMeterName, SuccessMetricName).Add(ctx, 1, otelmetric.WithAttributes(attrs...))
	} else if dim.IsTimeout() {
		r.getCounter(ServerReportMeterName, TimeoutMetricName).Add(ctx, 1, otelmetric.WithAttributes(attrs...))
	} else {
		r.getCounter(ServerReportMeterName, AbnormalMetricName).Add(ctx, 1, otelmetric.WithAttributes(attrs...))
	}

	// cost
	r.getCounter(ServerReportMeterName, CostMetricName).Add(ctx, int64(costMs), otelmetric.WithAttributes(attrs...))

	// duration histogram
	r.getHistogram(ServerReportMeterName, DurationMetricName, DefaultDurationBounds).Record(ctx, costMs, otelmetric.WithAttributes(attrs...))
}

// ReportClientMetric reports client-side (主调) metrics
func (r *MetricReporter) ReportClientMetric(ctx context.Context, dim *ClientDimension, costMs float64) {
	attrs := dim.ToAttributes()

	// requests
	r.getCounter(ClientReportMeterName, RequestsMetricName).Add(ctx, 1, otelmetric.WithAttributes(attrs...))

	// success / timeout / abnormal
	if dim.IsSuccess() {
		r.getCounter(ClientReportMeterName, SuccessMetricName).Add(ctx, 1, otelmetric.WithAttributes(attrs...))
	} else if dim.IsTimeout() {
		r.getCounter(ClientReportMeterName, TimeoutMetricName).Add(ctx, 1, otelmetric.WithAttributes(attrs...))
	} else {
		r.getCounter(ClientReportMeterName, AbnormalMetricName).Add(ctx, 1, otelmetric.WithAttributes(attrs...))
	}

	// cost
	r.getCounter(ClientReportMeterName, CostMetricName).Add(ctx, int64(costMs), otelmetric.WithAttributes(attrs...))

	// duration histogram
	r.getHistogram(ClientReportMeterName, DurationMetricName, DefaultDurationBounds).Record(ctx, costMs, otelmetric.WithAttributes(attrs...))
}

// ReportServerMetric is a convenience function using the global reporter
func ReportServerMetric(ctx context.Context, dim *ServerDimension, costMs float64) {
	GetGlobalReporter().ReportServerMetric(ctx, dim, costMs)
}

// ReportClientMetric is a convenience function using the global reporter
func ReportClientMetric(ctx context.Context, dim *ClientDimension, costMs float64) {
	GetGlobalReporter().ReportClientMetric(ctx, dim, costMs)
}
