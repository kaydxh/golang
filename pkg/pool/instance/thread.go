package instance

import (
	"context"
	"runtime"
	"sync"
	"time"
)

type Thread struct {
	enableOsThread bool

	ctx       context.Context
	cancel    context.CancelFunc
	handlerCh chan func()
	once      sync.Once
	mu        sync.Mutex
}

func NewThread(enableOsThread bool) *Thread {
	t := &Thread{
		enableOsThread: enableOsThread,
	}
	t.initOnce()
	return t
}

func (t *Thread) initOnce() {
	t.once.Do(func() {
		t.mu.Lock()
		defer t.mu.Unlock()
		t.ctx, t.cancel = context.WithCancel(context.Background())

		t.handlerCh = make(chan func())
		go t.DoInOSThread()
	})
}

func (t *Thread) Do(ctx context.Context, f func()) error {
	t.initOnce()

	time.Sleep(time.Second)
	select {
	case t.handlerCh <- f:
		return nil

	case <-ctx.Done():
		return ctx.Err()

	case <-t.ctx.Done():
		return t.ctx.Err()
	}
}

func (t *Thread) DoInOSThread() {

	if t.enableOsThread {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
	}

	for {
		select {
		case handler, ok := <-t.handlerCh:
			if !ok {
				return
			}

			if handler == nil {
				continue
			}
			handler()

		case <-t.ctx.Done():
			return
		}

	}
}
