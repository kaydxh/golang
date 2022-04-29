package pool

import (
	"context"
	"sync"

	sync_ "github.com/kaydxh/golang/go/sync"
)

type TaskHandler func(task interface{}) error

type Pool struct {
	Burst    int32
	TaskFunc TaskHandler

	taskChan chan interface{}
	err      error
	cond     *sync_.Cond

	wg      sync.WaitGroup
	burstMu sync.Mutex
	errMu   sync.Mutex
	errOnce sync.Once
}

func New(burst int32, taskFunc TaskHandler) *Pool {
	p := &Pool{
		Burst:    burst,
		TaskFunc: taskFunc,
		taskChan: make(chan interface{}, 0),
	}
	p.cond = sync_.NewCond(&p.burstMu)
	go p.run(context.Background())

	return p
}

func (p *Pool) Put(task interface{}) {
	p.taskChan <- task
}

func (p *Pool) Wait() error {
	p.wg.Wait()
	return p.err
}

func (p *Pool) Stop() {
	close(p.taskChan)
}

func (p *Pool) Error() error {
	p.errMu.Lock()
	defer p.errMu.Unlock()

	return p.err
}

func (p *Pool) trySetError(err error) {
	p.errOnce.Do(func() {
		p.errMu.Lock()
		defer p.errMu.Unlock()
		p.err = err
	})
}

func (p *Pool) run(ctx context.Context) (doneC <-chan struct{}) {
	done := make(chan struct{})
	if p.Burst <= 0 {
		p.Burst = 1
	}

	burst := p.Burst

	p.wg.Add(1)
	//	context.Background().Done()
	go func() {
		defer close(done)
		defer p.wg.Done()

		ctx, cancel := context.WithCancel(ctx)
		_ = cancel

		for {

			p.cond.WaitUntilDo(func() bool {
				return burst > 0
			}, func() error {
				burst--
				return nil
			})

			select {
			case task, ok := <-p.taskChan:
				if !ok {
					return
				}
				p.wg.Add(1)

				go func(t interface{}) {

					defer p.wg.Done()
					defer p.cond.SignalDo(func() error {
						burst++
						return nil
					})

					if err := p.TaskFunc(t); err != nil {
						//cancel()
						return
					}

				}(task)

			case <-ctx.Done():
				return
			}

		}
	}()

	return done

}
