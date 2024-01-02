package resource

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kaydxh/golang/go/errors"
	net_ "github.com/kaydxh/golang/go/net"
	syscall_ "github.com/kaydxh/golang/go/syscall"
	time_ "github.com/kaydxh/golang/go/time"
	resource_ "github.com/kaydxh/golang/pkg/middleware/resource"
	app_ "github.com/kaydxh/golang/pkg/webserver/app"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
)

const (
	DefaultCheckInterval time.Duration = time.Minute
)

type ResourceStatsOptions struct {
	checkInterval time.Duration
}

type ResourceStatsMetrics struct {
	MemoryTotalHistogram syncfloat64.Histogram
	MemoryUsageHistogram syncfloat64.Histogram
	MemoryFreeHistogram  syncfloat64.Histogram
}

type ResourceStatsService struct {
	inShutdown atomic.Bool // true when when server is in shutdown

	opts ResourceStatsOptions

	metrics *ResourceStatsMetrics

	mu     sync.Mutex
	cancel func()
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

func NewResourceStatsService(opts ...ResourceStatsServiceOption) (*ResourceStatsService, error) {
	var err error
	r := &ResourceStatsService{
		metrics: &ResourceStatsMetrics{},
	}

	call := func(f func()) {
		if err != nil {
			return
		}
		f()
	}
	call(func() {
		r.metrics.MemoryTotalHistogram, err = resource_.GlobalMeter().SyncFloat64().Histogram("memory_total")
	})
	call(func() {
		r.metrics.MemoryUsageHistogram, err = resource_.GlobalMeter().SyncFloat64().Histogram("memory_usage")
	})
	call(func() {
		r.metrics.MemoryFreeHistogram, err = resource_.GlobalMeter().SyncFloat64().Histogram("memory_free")
	})
	if err != nil {
		otel.Handle(err)
	}

	r.ApplyOptions(opts...)

	return r, nil
}

// Run will initialize the backend. It must not block, but may run go routines in the background.
func (s *ResourceStatsService) Run(ctx context.Context) error {
	logger := s.getLogger()
	logger.Infoln("ResourceStatsService Run")
	if s.inShutdown.Load() {
		logger.Infoln("ResourceStatsService Shutdown")
		return fmt.Errorf("server closed")
	}
	go func() {
		errors.HandleError(s.Serve(ctx))
	}()
	return nil
}

func (s *ResourceStatsService) getLogger() *logrus.Entry {
	return logrus.WithField("module", "ResourceStatsService")
}

func (s *ResourceStatsService) ReportMetric(ctx context.Context) {
	attrs := Attrs()

	s.metrics.MemoryTotalHistogram.Record(ctx, float64(syscall_.MemoryUsage{}.SysTotalMemory()), attrs...)
	s.metrics.MemoryFreeHistogram.Record(ctx, float64(syscall_.MemoryUsage{}.SysFreeMemory()), attrs...)
	s.metrics.MemoryUsageHistogram.Record(ctx, float64(syscall_.MemoryUsage{}.SysUsageMemory()), attrs...)

}

// Serve ...
func (s *ResourceStatsService) Serve(ctx context.Context) error {
	logger := s.getLogger()
	logger.Infoln("ResourceStats Serve")

	if s.inShutdown.Load() {
		logger.Infoln("ResourceStats Shutdown")
		return fmt.Errorf("server closed")
	}

	defer s.inShutdown.Store(true)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	s.mu.Lock()
	s.cancel = cancel
	s.mu.Unlock()

	time_.UntilWithContxt(ctx, func(ctx context.Context) error {
		s.ReportMetric(ctx)
		return nil
	}, s.opts.checkInterval)
	if err := ctx.Err(); err != nil {
		logger.WithError(err).Errorf("stopped checking")
		return err
	}
	logger.Info("stopped checking")
	return nil
}

// Shutdown ...
func (s *ResourceStatsService) Shutdown() {
	s.inShutdown.Store(true)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cancel != nil {
		s.cancel()
	}
}
