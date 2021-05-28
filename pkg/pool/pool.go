package pool

import (
	"context"
	"fmt"
	"sync"
)

type Pool struct {
	Burst    int32
	TaskFunc TaskHandler

	taskChan chan interface{}
	err      error

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
	go p.run(context.Background())

	return p
}

type TaskHandler func(task interface{}) error

func (p *Pool) Put(task interface{}) {
	p.taskChan <- task
}

func (p *Pool) Wait() error {
	p.wg.Wait()
	return p.err
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

	c := sync.NewCond(&p.burstMu)
	go func() {
		defer close(done)
		defer p.wg.Done()

		ctx, cancel := context.WithCancel(ctx)

		for {

			c.L.Lock()
			for burst <= 0 {
				c.Wait()
			}
			burst--
			//fmt.Printf("burst: %v\n", burst)
			c.L.Unlock()

			select {
			case task, ok := <-p.taskChan:
				if !ok {
					return
				}
				p.wg.Add(1)

				go func(t interface{}) {
					cleanfunc := func() {
						c.L.Lock()
						burst++
						c.Signal()
						c.L.Unlock()

						p.wg.Done()

					}
					defer cleanfunc()

					if err := p.TaskFunc(t); err != nil {
						cancel()
						return
					}

				}(task)

			case <-ctx.Done():
				fmt.Println("done")
				return
			}

		}
	}()

	return done

}
