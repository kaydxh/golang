package sync

import (
	"fmt"
	"sync"
	"time"
)

type CondTimeout struct {
	L  sync.Locker
	ch chan struct{}
}

func NewCondTimeout(l sync.Locker) *CondTimeout {
	c := &CondTimeout{
		L:  l,
		ch: make(chan struct{}),
	}
	return c
}

func (c *CondTimeout) wait() error {
	c.L.Unlock()
	defer c.L.Lock()
	<-c.ch
	return nil
}

func (c *CondTimeout) waitFor(timeout int) error {
	c.L.Unlock()
	defer c.L.Lock()

	select {
	case <-c.ch:
		return nil
	case <-time.After(time.Duration(timeout) * time.Second):
		return fmt.Errorf("wait timeout: %vs\n", timeout)
	}

}

func (c *CondTimeout) WaitFoDo(timeout int, pred func() bool, do func() error) error {
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

func (c *CondTimeout) WaitUntilDo(timeout int, pred func() bool, do func() error) error {
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

func (c *CondTimeout) SignalDo(do func() error) {
	c.L.Lock()
	if do != nil {
		do()
	}
	defer c.L.Unlock()

	c.Signal()
}

func (c *CondTimeout) BroadcastDo(do func() error) {
	c.L.Lock()
	if do != nil {
		do()
	}
	defer c.L.Unlock()

	c.Broadcast()
}

func (c *CondTimeout) Signal() {
	go func() {
		c.ch <- struct{}{}
	}()
}

func (c *CondTimeout) Broadcast() {
	go func() {
		close(c.ch)
		c.ch = make(chan struct{})
	}()
}
