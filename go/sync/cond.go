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
		start := time.Now()
		err := c.waitFor(timeout)
		if err != nil {
			return err
		}

		timeout -= time.Now().Sub(start)
	}
	if do != nil {
		do()
	}
	return nil
}

func (c *Cond) WaitUntil(pred func() bool) {
	c.WaitUntilDo(pred, func() error { return nil })
}

func (c *Cond) WaitUntilDo(pred func() bool, do func() error) {
	for {
		err := c.WaitForDo(5*time.Second, pred, do)
		if err == nil {
			break
		}
		c.Signal()
	}

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
