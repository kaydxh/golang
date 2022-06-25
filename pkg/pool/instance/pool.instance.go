package instance

import (
	"context"
	"fmt"
	"sync"
	"time"

	errors_ "github.com/kaydxh/golang/go/errors"
)

type GlobalInitFunc func() error
type GlobalReleaseFunc func() error
type LocalInitFunc func(instance interface{}) error
type LocalReleaseFunc func(instance interface{}) error
type NewFunc func() interface{}
type DeleteFunc func(instance interface{})

type LoadBalanceMode int

const (
	RoundRobinBalanceMode LoadBalanceMode = 0
	RandomLoadBalanceMode LoadBalanceMode = 1
	LeastLoadBalanceMode  LoadBalanceMode = 2
)

type PoolOptions struct {
	name                   string
	resevePoolSizePerGpu   int64
	capacityPoolSizePerGpu int64
	idleTimeout            time.Duration
	waitTimeout            time.Duration
	loadBalanceMode        LoadBalanceMode

	gpuIDs []int64

	modelPaths []string
	batchSize  int64

	globalInitFunc    GlobalInitFunc
	globalReleaseFunc GlobalReleaseFunc
	localInitFunc     LocalInitFunc
	localReleaseFunc  LocalReleaseFunc
	newFunc           NewFunc
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
	p := &Pool{
		holders: make(map[int64]chan *CoreInstanceHolder),
	}
	p.ApplyOptions(opts...)

	if p.opts.resevePoolSizePerGpu < 0 {
		p.opts.resevePoolSizePerGpu = 0
	}

	if p.opts.capacityPoolSizePerGpu < p.opts.resevePoolSizePerGpu {
		p.opts.capacityPoolSizePerGpu = p.opts.resevePoolSizePerGpu
	}

	return p
}

func (p *Pool) init(ctx context.Context) error {

	for _, id := range p.opts.gpuIDs {
		if id < 0 {
			continue
		}

		p.holders[id] = make(chan *CoreInstanceHolder, p.opts.capacityPoolSizePerGpu)
		for i := int64(0); i < p.opts.resevePoolSizePerGpu; i++ {
			holder := &CoreInstanceHolder{
				Name:       p.opts.name,
				GpuID:      id,
				ModelPaths: p.opts.modelPaths,
				BatchSize:  p.opts.batchSize,
			}
			err := holder.Do(ctx, func() {
				holder.Instance = p.opts.newFunc()
			})
			if err != nil {
				return err
			}

			var processErr error
			err = holder.Do(ctx, func() {
				processErr = p.opts.localInitFunc(holder.Instance)
			})

			if processErr != nil {
				return processErr
			}
			if err != nil {
				return err
			}

			p.holders[id] <- holder
		}
	}

	return nil
}

func (p *Pool) GlobalInit(ctx context.Context) error {
	err := p.globalInit(ctx)
	if err != nil {
		return err
	}

	return p.init(ctx)
}

func (p *Pool) globalInit(ctx context.Context) error {
	if p.opts.globalInitFunc == nil {
		return nil
	}

	return p.opts.globalInitFunc()
}

func (p *Pool) GlobalRelease(ctx context.Context) error {
	for _, ch := range p.holders {
		select {
		case holder := <-ch:
			holder.Do(ctx, func() {
				if p.opts.localReleaseFunc != nil {
					p.opts.localReleaseFunc(holder.Instance)
				}

				if p.opts.deleteFunc != nil {
					p.opts.deleteFunc(holder.Instance)
				}
			})

		case <-ctx.Done():
			return ctx.Err()

		default:
		}
	}

	return p.opts.globalReleaseFunc()
}

func (p *Pool) globalRelease(ctx context.Context) error {
	if p.opts.globalReleaseFunc == nil {
		return nil
	}

	return p.opts.globalReleaseFunc()
}

func (p *Pool) GetByGpuId(ctx context.Context, gpuID int64) (*CoreInstanceHolder, error) {
	if p.opts.capacityPoolSizePerGpu <= 0 {
		return nil, fmt.Errorf("no instance, capacity pool size per gpu is <= 0")
	}

	select {
	case holder := <-p.holders[gpuID]:
		fmt.Println("get")
		return holder, nil

	case <-ctx.Done():
		return nil, ctx.Err()

	default:
	}

	if p.opts.waitTimeout == 0 {
		return nil, fmt.Errorf("no avaliable instance in gpu id: %v right now", gpuID)
	}

	timer := time.NewTimer(p.opts.waitTimeout)
	defer timer.Stop()

	select {
	case holder := <-p.holders[gpuID]:
		fmt.Println("get")
		return holder, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timer.C:
		return nil, fmt.Errorf("get instance timeout: %v", p.opts.waitTimeout)
	}

}

func (p *Pool) Get(ctx context.Context) (*CoreInstanceHolder, error) {

	switch p.opts.loadBalanceMode {
	case RoundRobinBalanceMode:
		return p.GetWithRoundRobinMode(ctx)
	default:
		return nil, fmt.Errorf("not support the loadBalance mode: %v", p.opts.loadBalanceMode)
	}
}

func (p *Pool) GetWithRoundRobinMode(ctx context.Context) (*CoreInstanceHolder, error) {

	var errs []error

	for _, id := range p.opts.gpuIDs {
		holder, err := p.GetByGpuId(ctx, id)
		if err == nil {
			return holder, nil
		}
		errs = append(errs, err)
	}

	return nil, errors_.NewAggregate(errs)
}

func (p *Pool) Put(ctx context.Context, holder *CoreInstanceHolder) error {
	select {
	case p.holders[holder.GpuID] <- holder:
		fmt.Println("put")
		return nil

	case <-ctx.Done():
		return ctx.Err()

	default:
	}

	// if get this line, it means put instance error, then we delete this instance
	if p.opts.localReleaseFunc != nil {
		holder.Do(ctx, func() {
			p.opts.localReleaseFunc(holder.Instance)
		})

	}
	if p.opts.deleteFunc != nil {
		holder.Do(ctx, func() {
			p.opts.deleteFunc(holder.Instance)
		})
	}

	return nil
}

func (p *Pool) Invoke(
	ctx context.Context,
	f func(ctx context.Context, instance interface{}) (interface{}, error),
) (response interface{}, err error) {

	if f == nil {
		return nil, nil
	}

	holder, err := p.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer p.Put(ctx, holder)

	var processErr error
	err = holder.Do(ctx, func() {
		response, processErr = f(ctx, holder.Instance)
	})

	if processErr != nil {
		err = processErr
	}

	return response, err
}
