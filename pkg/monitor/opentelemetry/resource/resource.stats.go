package resource

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kaydxh/golang/go/errors"
	time_ "github.com/kaydxh/golang/go/time"
	resource_ "github.com/kaydxh/golang/pkg/middleware/resource"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

const (
	DefaultCheckInterval time.Duration = time.Minute
)

type ResourceStatsOptions struct {
	checkInterval time.Duration
}

type ResourceStatsService struct {
	inShutdown atomic.Bool // true when when server is in shutdown

	opts ResourceStatsOptions

	metrics *ResourceStatsMetrics

	mu     sync.Mutex
	cancel func()
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

	if r.opts.checkInterval <= 0 {
		r.opts.checkInterval = DefaultCheckInterval
	}

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
		s.metrics.ReportMetric(ctx)
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
