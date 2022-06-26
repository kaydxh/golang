package instance

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	runtime_ "github.com/kaydxh/golang/go/runtime"
	time_ "github.com/kaydxh/golang/go/time"
	"github.com/sirupsen/logrus"
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
	// the wait time to try get instance again, 0 means wait forever,
	// execute the block without any timeout
	waitTimeoutOnce time.Duration

	// the total wait time to try get instance, 0 means wait forever,
	// execute the block without any timeout
	waitTimeoutTotal time.Duration
	loadBalanceMode  LoadBalanceMode

	gpuIDs []int64

	modelPaths []string
	batchSize  int64

	globalInitFunc    GlobalInitFunc
	globalReleaseFunc GlobalReleaseFunc
	localInitFunc     LocalInitFunc
	localReleaseFunc  LocalReleaseFunc
	deleteFunc        DeleteFunc

	enabledPrintCostTime bool
}

type Pool struct {
	newFunc NewFunc

	holders map[int64]chan *CoreInstanceHolder
	opts    PoolOptions
	mu      sync.RWMutex
}

func defaultPoolOptions() PoolOptions {
	return PoolOptions{
		resevePoolSizePerGpu:   0,
		capacityPoolSizePerGpu: 1,
		idleTimeout:            10 * time.Second,
		waitTimeoutOnce:        10 * time.Millisecond,
	}
}

func NewPool(newFunc NewFunc, opts ...PoolOption) (*Pool, error) {
	if newFunc == nil {
		return nil, fmt.Errorf("new func is nil")

	}
	p := &Pool{
		newFunc: newFunc,
		holders: make(map[int64]chan *CoreInstanceHolder),
		opts:    defaultPoolOptions(),
	}
	p.ApplyOptions(opts...)

	if p.opts.resevePoolSizePerGpu < 0 {
		p.opts.resevePoolSizePerGpu = 0
	}

	if p.opts.capacityPoolSizePerGpu < p.opts.resevePoolSizePerGpu {
		p.opts.capacityPoolSizePerGpu = p.opts.resevePoolSizePerGpu
	}

	return p, nil
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
				holder.Instance = p.newFunc()
			})
			if err != nil {
				return err
			}
			if holder.Instance == nil {
				return fmt.Errorf("instance the result of new func is nil")
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

	// until get instance, unless context canceled
	//for {
	if p.opts.waitTimeoutOnce == 0 {
		select {
		case holder := <-p.holders[gpuID]:
			logrus.Infof("get a instance in gpu id: %v", gpuID)
			return holder, nil
		case <-ctx.Done():
			return nil, ContextCanceledError{Message: ctx.Err().Error()}
		}
	}

	/*
		if p.opts.waitTimeoutOnce == 0 {
			msg := fmt.Sprintf("no avaliable instance in gpu id: %v right now, try again.", gpuID)
			logrus.Warnf(msg)
			//logrus.Warnf("no avaliable instance in gpu id: %v right now, try again.", gpuID)
			//	continue
			return nil, TimeoutError{Message: msg}
		}
	*/

	timer := time.NewTimer(p.opts.waitTimeoutOnce)
	defer timer.Stop()

	select {
	case holder := <-p.holders[gpuID]:
		logrus.Infof("get a instance in gpu id: %v", gpuID)
		return holder, nil
	case <-ctx.Done():
		logrus.WithError(ctx.Err()).Errorf("get a instance in gpu id: %v context canceled", gpuID)
		return nil, ContextCanceledError{Message: ctx.Err().Error()}
	case <-timer.C:
		msg := fmt.Sprintf("get instance timeout: %v, try again.", p.opts.waitTimeoutOnce)
		logrus.Warnf(msg)
		return nil, TimeoutError{Message: msg}
		//logrus.Warnf("get instance timeout: %v, try again.", p.opts.waitTimeout)
	}
	//	}

}

func (p *Pool) Get(ctx context.Context) (*CoreInstanceHolder, error) {

	switch p.opts.loadBalanceMode {
	case RoundRobinBalanceMode:
		return p.GetWithRoundRobinMode(ctx)
	default:
		return nil, fmt.Errorf("not support the loadBalance mode: %v", p.opts.loadBalanceMode)
	}
}

// until get instance, unless context canceled
func (p *Pool) GetWithRoundRobinMode(ctx context.Context) (*CoreInstanceHolder, error) {

	remain := p.opts.waitTimeoutTotal

	for {
		for _, id := range p.opts.gpuIDs {
			tc := time_.New(true)

			holder, err := p.GetByGpuId(ctx, id)
			if err == nil {
				return holder, nil
			}

			if p.opts.waitTimeoutTotal > 0 {
				remain -= tc.Elapse()
				if remain <= 0 {
					msg := fmt.Sprintf("get instance total timeout: %v, remain: %v", p.opts.waitTimeoutTotal, remain)
					logrus.Errorf(msg)
					return nil, TimeoutError{Message: msg}
				}
			}

			switch err.(type) {
			case ContextCanceledError:
				return nil, err
			default:
			}
		}
	}
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

	holder := &CoreInstanceHolder{}
	tc := time_.New(p.opts.enabledPrintCostTime)
	summary := func() {
		var gpuID int64 = -1
		if holder != nil {
			gpuID = holder.GpuID
		}
		tc.Tick(fmt.Sprintf("Invoke %v no gpuID: %v", filepath.Base(runtime_.NameOfFunction(f)), gpuID))
		logrus.Infof(tc.String())
	}
	defer summary()

	if f == nil {
		return nil, nil
	}

	holder, err = p.Get(ctx)
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
