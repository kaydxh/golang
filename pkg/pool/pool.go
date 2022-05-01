package pool

import (
	"context"
	"fmt"
	"sync"

	errors_ "github.com/kaydxh/golang/go/errors"
	rate_ "github.com/kaydxh/golang/go/time/rate"
)

type TaskHandler func(task interface{}) error

type Pool struct {
	Burst    int32
	TaskFunc TaskHandler

	ctx context.Context

	taskChan     chan interface{}
	errs         []error
	cancel       context.CancelFunc
	workDoneChan chan struct{}

	wg    sync.WaitGroup
	errMu sync.Mutex
}

func New(burst int32, taskFunc TaskHandler) *Pool {
	p := &Pool{
		Burst:        burst,
		TaskFunc:     taskFunc,
		taskChan:     make(chan interface{}),
		workDoneChan: make(chan struct{}),
	}
	p.ctx, p.cancel = context.WithCancel(context.Background())
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
	if p.Burst <= 0 {
		p.Burst = 1
	}

	go func() {
		defer close(p.workDoneChan)

		limiter := rate_.NewLimiter(int(p.Burst))
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
						//	p.cancel()
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
