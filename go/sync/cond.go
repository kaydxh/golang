package sync

import (
	"fmt"
	"sync"
	"time"
)

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

func (c *Cond) waitFor(timeout int) error {
	c.L.Unlock()
	defer c.L.Lock()

	select {
	case <-c.ch:
		return nil
	case <-time.After(time.Duration(timeout) * time.Second):
		return fmt.Errorf("wait timeout: %vs\n", timeout)
	}

}

func (c *Cond) WaitForDo(timeout int, pred func() bool, do func() error) error {
	c.L.Lock()
	defer c.L.Unlock()

	for !pred() {
		err := c.waitFor(timeout)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return err
		}
	}
	if do != nil {
		do()
	}
	return nil
}

func (c *Cond) WaitUntilDo(timeout int, pred func() bool, do func() error) error {
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
