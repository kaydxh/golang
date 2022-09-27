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
	name                    string
	resevePoolSizePerCore   int64
	capacityPoolSizePerCore int64
	idleTimeout             time.Duration
	// the wait time to try get instance again, 0 means wait forever,
	// execute the block without any timeout
	waitTimeoutOnce time.Duration

	// the total wait time to try get instance, 0 means wait forever,
	// execute the block without any timeout
	waitTimeoutTotal time.Duration
	loadBalanceMode  LoadBalanceMode

	coreIDs []int64

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
		resevePoolSizePerCore:   0,
		capacityPoolSizePerCore: 1,
		idleTimeout:             10 * time.Second,
		waitTimeoutOnce:         10 * time.Millisecond,
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

	if p.opts.resevePoolSizePerCore < 0 {
		p.opts.resevePoolSizePerCore = 0
	}

	if p.opts.capacityPoolSizePerCore < p.opts.resevePoolSizePerCore {
		p.opts.capacityPoolSizePerCore = p.opts.resevePoolSizePerCore
	}

	return p, nil
}

func (p *Pool) init(ctx context.Context) error {

	for _, id := range p.opts.coreIDs {
		if id < 0 {
			continue
		}

		p.holders[id] = make(chan *CoreInstanceHolder, p.opts.capacityPoolSizePerCore)
		for i := int64(0); i < p.opts.resevePoolSizePerCore; i++ {
			holder := &CoreInstanceHolder{
				Name:       p.opts.name,
				CoreID:     id,
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

func (p *Pool) GetByCoreId(ctx context.Context, coreID int64) (*CoreInstanceHolder, error) {
	if p.opts.capacityPoolSizePerCore <= 0 {
		return nil, fmt.Errorf("no instance, capacity pool size per core is <= 0")
	}

	if p.opts.waitTimeoutOnce == 0 {
		select {
		case holder := <-p.holders[coreID]:
			logrus.Infof("get a instance in core id: %v", coreID)
			return holder, nil
		case <-ctx.Done():
			return nil, ContextCanceledError{Message: ctx.Err().Error()}
		}
	}

	timer := time.NewTimer(p.opts.waitTimeoutOnce)
	defer timer.Stop()

	select {
	case holder := <-p.holders[coreID]:
		logrus.Infof("get a instance in core id: %v", coreID)
		return holder, nil
	case <-ctx.Done():
		logrus.WithError(ctx.Err()).Errorf("get a instance in core id: %v context canceled", coreID)
		return nil, ContextCanceledError{Message: ctx.Err().Error()}
	case <-timer.C:
		msg := fmt.Sprintf("get instance timeout: %v, try again.", p.opts.waitTimeoutOnce)
		logrus.Warnf(msg)
		return nil, TimeoutError{Message: msg}
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

// until get instance, unless context canceled
func (p *Pool) GetWithRoundRobinMode(ctx context.Context) (*CoreInstanceHolder, error) {

	remain := p.opts.waitTimeoutTotal

	for {
		for _, id := range p.opts.coreIDs {
			tc := time_.New(true)

			holder, err := p.GetByCoreId(ctx, id)
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
	case p.holders[holder.CoreID] <- holder:
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
		var coreID int64 = -1
		if holder != nil {
			coreID = holder.CoreID
		}
		tc.Tick(fmt.Sprintf("Invoke %v no coreID: %v", filepath.Base(runtime_.NameOfFunction(f)), coreID))
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
