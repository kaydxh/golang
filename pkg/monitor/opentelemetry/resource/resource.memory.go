package resource

import (
	"context"

	net_ "github.com/kaydxh/golang/go/net"
	resource_ "github.com/kaydxh/golang/pkg/middleware/resource"
	app_ "github.com/kaydxh/golang/pkg/webserver/app"
	"github.com/shirou/gopsutil/mem"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	MemoryTotalKey     = "memory_total"
	MemoryUsageKey     = "memory_usage"
	MemoryAvailableKey = "memory_available"
)

type ResourceStatsMetrics struct {
	MemoryTotalHistogram     metric.Float64Histogram
	MemoryUsageHistogram     metric.Float64Histogram
	MemoryAvailableHistogram metric.Float64Histogram
}

/*
func RecordOptions() []metric.RecordOption {
	attrs := Attrs()
	opts := make([]metric.RecordOption, 0, len(attrs))
	for _, attr := range attrs {
		opts = append(opts, metric.WithAttribute(attr))
	}
	return opts
}
*/

func Attrs() []attribute.KeyValue {
	var attrs []attribute.KeyValue
	hostIP, err := net_.GetHostIP()
	if err == nil && hostIP.String() != "" {
		attrs = append(attrs, resource_.PodIpKey.String(hostIP.String()))
	}
	appName := app_.GetVersion().AppName
	if appName != "" {
		attrs = append(attrs, resource_.ServerNameKey.String(appName))
	}

	return attrs
}

func NewResourceStatsMetrics() (*ResourceStatsMetrics, error) {
	/*
		var err error
		r := &ResourceStatsMetrics{}
		call := func(f func()) {
			if err != nil {
				return
			}
			f()
		}
		call(func() {
			r.MemoryTotalHistogram, err = resource_.GlobalMeter().Float64Histogram(MemoryTotalKey)
		})
		call(func() {
			r.MemoryUsageHistogram, err = resource_.GlobalMeter().Float64Histogram(MemoryUsageKey)
		})
		call(func() {
			r.MemoryAvaliableHistogram, err = resource_.GlobalMeter().Float64Histogram(MemoryAvaliableKey)
		})
		if err != nil {
			otel.Handle(err)
		}

		return r, nil
	*/

	var err error
	r := &ResourceStatsMetrics{}

	// 获取全局 meter
	meter := resource_.GlobalMeter()

	// 创建 histogram 的辅助函数
	createHistogram := func(name string) (metric.Float64Histogram, error) {
		return meter.Float64Histogram(name,
			metric.WithDescription("Resource "+name+" metrics"),
			metric.WithUnit("bytes"), // 根据需要调整单位
		)
	}

	// 创建内存总量 histogram
	r.MemoryTotalHistogram, err = createHistogram(MemoryTotalKey)
	if err != nil {
		otel.Handle(err)
		return nil, err
	}

	// 创建内存使用量 histogram
	r.MemoryUsageHistogram, err = meter.Float64Histogram(MemoryUsageKey,
		metric.WithDescription("Memory usage percentage"),
		metric.WithUnit("%"),
	)
	if err != nil {
		otel.Handle(err)
		return nil, err
	}

	// 创建内存可用量 histogram
	r.MemoryAvailableHistogram, err = createHistogram(MemoryAvailableKey)
	if err != nil {
		otel.Handle(err)
		return nil, err
	}

	return r, nil

}

func (r *ResourceStatsMetrics) ReportMetric(ctx context.Context) (total, available, usage float64) {
	attrs := Attrs()

	v, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, 0
	}

	total = float64(v.Total)
	available = float64(v.Available)
	usage = v.UsedPercent

	r.MemoryTotalHistogram.Record(ctx, total, metric.WithAttributes(attrs...))
	r.MemoryAvailableHistogram.Record(ctx, available, metric.WithAttributes(attrs...))
	r.MemoryUsageHistogram.Record(ctx, usage, metric.WithAttributes(attrs...))

	// r.MemoryTotalHistogram.Record(ctx, float64(v.Total), attrs...)
	// r.MemoryAvaliableHistogram.Record(ctx, float64(v.Available), attrs...)
	// r.MemoryUsageHistogram.Record(ctx, v.UsedPercent, attrs...)

	return total, available, usage
}
