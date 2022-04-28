package sync

import (
	"fmt"
	"sync"
	"time"
)

// Cond 条件变量不保证顺序性，即signal只是通知wait去获取数据，wait拿到的数据是不是
// signal当时给的数据，不能保证
type Cond struct {
	L  sync.Locker
	ch chan struct{}
}

func NewCond(l sync.Locker) *Cond {
	c := &Cond{
		L:  l,
		ch: make(chan struct{}),
	}
	return c
}

func (c *Cond) wait() error {
	c.L.Unlock()
	defer c.L.Lock()
	<-c.ch
	return nil
}

func (c *Cond) waitFor(timeout time.Duration) error {
	c.L.Unlock()
	defer c.L.Lock()

	select {
	case <-c.ch:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("wait timeout: %v\n", timeout)
	}

}

func (c *Cond) WaitForDo(timeout time.Duration, pred func() bool, do func() error) error {
	c.L.Lock()
	defer c.L.Unlock()

	for !pred() {
		err := c.waitFor(timeout)
		if err != nil {
			return err
		}

	}
	if do != nil {
		do()
	}
	return nil
}

func (c *Cond) WaitUntilDo(timeout time.Duration, pred func() bool, do func() error) error {
	c.L.Lock()
	defer c.L.Unlock()

	for !pred() {
		c.wait()
	}
	if do != nil {
		do()
	}
	return nil
}

func (c *Cond) SignalDo(do func() error) {
	c.L.Lock()
	if do != nil {
		do()
	}
	c.L.Unlock()

	c.Signal()
}

func (c *Cond) BroadcastDo(do func() error) {
	c.L.Lock()
	if do != nil {
		do()
	}
	c.L.Unlock()

	c.Broadcast()
}

func (c *Cond) Signal() {
	go func() {
		c.ch <- struct{}{}
	}()
}

func (c *Cond) Broadcast() {
	go func() {
		close(c.ch)
		c.ch = make(chan struct{})
	}()
}
