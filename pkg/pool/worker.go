package pool

import (
	"context"
	"fmt"
	"sync"

	errors_ "github.com/kaydxh/golang/go/errors"
	rate_ "github.com/kaydxh/golang/go/time/rate"
)

type WorkFunc func(work interface{}) error

type Worker struct {
	Burst int32

	cancel context.CancelFunc
	wg     sync.WaitGroup
	errs   []error
	errMu  sync.Mutex
}

func (w *Worker) Work(
	ctx context.Context,
	workChan <-chan interface{},
	workFn WorkFunc,
	errStop bool,
) (doneC <-chan struct{}) {
	done := make(chan struct{})
	if w.Burst == 0 {
		w.Burst = 1
	}

	go func() {
		defer close(done)

		ctx, w.cancel = context.WithCancel(ctx)
		limiter := rate_.NewLimiter(int(w.Burst))
		for {

			//util the condition is met, need one token, or will be blocked
			limiter.AllowWaitUntil()

			select {
			case work, ok := <-workChan:
				if !ok {
					//workChan is closed
					return
				}
				w.wg.Add(1)

				go func(wk interface{}) {
					defer limiter.Put()
					defer w.wg.Done()

					if err := workFn(wk); err != nil {
						w.TrySetError(err)
						if errStop {
							w.cancel()
						}
						return
					}

				}(work)

			case <-ctx.Done():
				//  err: context canceled
				w.TrySetError(ctx.Err())
				fmt.Println("===cancel")
				return
			}

		}
	}()

	return done

}

func (w *Worker) Error() error {
	w.errMu.Lock()
	defer w.errMu.Unlock()

	return errors_.NewAggregate(w.errs)
}

func (w *Worker) TrySetError(err error) {
	w.errMu.Lock()
	defer w.errMu.Unlock()
	w.errs = append(w.errs, err)
}

func (w *Worker) Wait() {
	w.wg.Wait()
}

func (w *Worker) Stop() {
	w.cancel()
}
