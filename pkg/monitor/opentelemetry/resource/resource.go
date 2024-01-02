package resource

import (
	"context"

	net_ "github.com/kaydxh/golang/go/net"
	syscall_ "github.com/kaydxh/golang/go/syscall"
	resource_ "github.com/kaydxh/golang/pkg/middleware/resource"
	app_ "github.com/kaydxh/golang/pkg/webserver/app"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
)

type ResourceStatsMetrics struct {
	MemoryTotalHistogram syncfloat64.Histogram
	MemoryUsageHistogram syncfloat64.Histogram
	MemoryFreeHistogram  syncfloat64.Histogram
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

func (r *ResourceStatsMetrics) ReportMetric(ctx context.Context) {
	attrs := Attrs()

	total := float64(syscall_.MemoryUsage{}.SysTotalMemory())
	free := float64(syscall_.MemoryUsage{}.SysFreeMemory())
	usage := float64(syscall_.MemoryUsage{}.SysUsageMemory())

	r.MemoryTotalHistogram.Record(ctx, total, attrs...)
	r.MemoryFreeHistogram.Record(ctx, free, attrs...)
	r.MemoryUsageHistogram.Record(ctx, usage, attrs...)

}
