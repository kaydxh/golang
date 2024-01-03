package resource

import (
	"context"

	net_ "github.com/kaydxh/golang/go/net"
	resource_ "github.com/kaydxh/golang/pkg/middleware/resource"
	app_ "github.com/kaydxh/golang/pkg/webserver/app"
	"github.com/shirou/gopsutil/mem"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
)

const (
	MemoryTotalKey     = "memory_total"
	MemoryUsageKey     = "memory_usage"
	MemoryAvaliableKey = "memory_avaliable"
)

type ResourceStatsMetrics struct {
	MemoryTotalHistogram     syncfloat64.Histogram
	MemoryUsageHistogram     syncfloat64.Histogram
	MemoryAvaliableHistogram syncfloat64.Histogram
}

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
	var err error
	r := &ResourceStatsMetrics{}
	call := func(f func()) {
		if err != nil {
			return
		}
		f()
	}
	call(func() {
		r.MemoryTotalHistogram, err = resource_.GlobalMeter().SyncFloat64().Histogram(MemoryTotalKey)
	})
	call(func() {
		r.MemoryUsageHistogram, err = resource_.GlobalMeter().SyncFloat64().Histogram(MemoryUsageKey)
	})
	call(func() {
		r.MemoryAvaliableHistogram, err = resource_.GlobalMeter().SyncFloat64().Histogram(MemoryAvaliableKey)
	})
	if err != nil {
		otel.Handle(err)
	}

	return r, nil
}

func (r *ResourceStatsMetrics) ReportMetric(ctx context.Context) (total, avaiable, usage float64) {
	attrs := Attrs()

	v, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, 0
	}

	total = float64(v.Total)
	avaiable = float64(v.Available)
	usage = v.UsedPercent
	r.MemoryTotalHistogram.Record(ctx, total, attrs...)
	r.MemoryAvaliableHistogram.Record(ctx, avaiable, attrs...)
	r.MemoryUsageHistogram.Record(ctx, usage, attrs...)

	return total, avaiable, usage
}
