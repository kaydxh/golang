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
package task

import (
	"context"
	"fmt"
	"sync"

	errors_ "github.com/kaydxh/golang/go/errors"
	rate_ "github.com/kaydxh/golang/go/time/rate"
)

type TaskHandler func(task interface{}) error

type PoolConfig struct {
	burst int
	// errstop, only take effect to the tasks that excessed the number burst,
	errStop bool
}

func defaultPoolConfig() PoolConfig {
	return PoolConfig{
		burst:   1,
		errStop: false,
	}
}

type Pool struct {
	TaskFunc TaskHandler
	opts     PoolConfig
	ctx      context.Context

	taskChan     chan interface{}
	errs         []error
	cancel       context.CancelFunc
	workDoneChan chan struct{}

	wg    sync.WaitGroup
	errMu sync.Mutex
}

func New(taskFunc TaskHandler, opts ...PoolOptions) *Pool {
	p := &Pool{
		TaskFunc:     taskFunc,
		taskChan:     make(chan interface{}),
		workDoneChan: make(chan struct{}),
		opts:         defaultPoolConfig(),
	}
	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.ApplyOptions(opts...)

	go p.run()
	return p
}

func (p *Pool) Put(task interface{}) error {
	select {
	case p.taskChan <- task:
	case <-p.workDoneChan:
		return fmt.Errorf("workdone channel is closed, work goroutine is exit")
	case <-p.ctx.Done():
		return fmt.Errorf("work is stoped, work goroutine is exit")

	}

	return nil
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Stop() {
	p.cancel()
}

func (p *Pool) Error() error {
	p.errMu.Lock()
	defer p.errMu.Unlock()

	return errors_.NewAggregate(p.errs)
}

func (p *Pool) trySetError(err error) {
	p.errMu.Lock()
	defer p.errMu.Unlock()
	p.errs = append(p.errs, err)
}

func (p *Pool) run() (doneC <-chan struct{}) {

	go func() {
		defer close(p.workDoneChan)

		limiter := rate_.NewLimiter(int(p.opts.burst))
		for {
			//util the condition is met, need one token, or will be blocked
			limiter.AllowWaitUntil()

			select {
			case task, ok := <-p.taskChan:
				if !ok {
					return
				}
				p.wg.Add(1)

				go func(t interface{}) {
					defer limiter.Put()
					defer p.wg.Done()

					if err := p.TaskFunc(t); err != nil {
						p.trySetError(err)
						if p.opts.errStop {
							p.cancel()
						}
						return
					}

				}(task)

			case <-p.ctx.Done():
				//  err: context canceled
				p.trySetError(p.ctx.Err())
				return
			}

		}
	}()

	return p.workDoneChan
}
