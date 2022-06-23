package instance

import (
	"context"
	"sync"
	"time"

	errors_ "github.com/kaydxh/golang/go/errors"
)

type GlobalInitFunc func() error
type GlobalReleaseFunc func() error
type LocalInitFunc func(holder interface{}) error
type LocalReleaseFunc func(instance interface{}) error
type DeleteFunc func(instance interface{})

type PoolOptions struct {
	resevePoolSizePerGpu   int64
	capacityPoolSizePerGpu int64
	idleTimeout            time.Duration
	waitTimeout            time.Duration

	GpuIDs []int64

	globalInitFunc    GlobalInitFunc
	globalReleaseFunc GlobalReleaseFunc
	localInitFunc     LocalInitFunc
	localReleaseFunc  LocalReleaseFunc
	deleteFunc        DeleteFunc
}

type Pool struct {
	holders map[int64]chan *CoreInstanceHolder
	opts    PoolOptions
	mu      sync.RWMutex
}

func defaultPoolOptions() PoolOptions {
	return PoolOptions{
		resevePoolSizePerGpu:   1,
		capacityPoolSizePerGpu: 1,
		idleTimeout:            10 * time.Second,
		waitTimeout:            10 * time.Millisecond,
	}
}

func NewPool(opts ...PoolOption) *Pool {
	p := &Pool{}
	p.ApplyOptions(opts...)

	return p
}

func (p *Pool) GlobalInit(ctx context.Context) error {
	if p.opts.globalInitFunc == nil {
		return nil
	}

	return p.opts.globalInitFunc()
}

func (p *Pool) GlobalRelease(ctx context.Context) error {
	if p.opts.globalReleaseFunc == nil {
		return nil
	}

	return p.opts.globalReleaseFunc()
}

func (p *Pool) LocalInit(ctx context.Context, holder CoreInstanceHolder) (err error) {
	if p.opts.localInitFunc == nil {
		return nil
	}

	var errs []error

	for _, id := range p.opts.GpuIDs {
		if id >= 0 {

			var processErr error
			err = holder.Do(ctx, func() {
				processErr = p.opts.localInitFunc(holder.Instance)
			})
			if processErr != nil {
				errs = append(errs, processErr)
			}

			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errors_.NewAggregate(errs)
}
