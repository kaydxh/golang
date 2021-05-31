package sync

import "sync"

type Cond struct {
	cond *sync.Cond
	mu   sync.Mutex
}

func NewCond() *Cond {
	c := &Cond{}
	c.cond = sync.NewCond(&c.mu)
	return c
}

func (c *Cond) WaitCond(pred func() bool, do func()) {
	c.cond.L.Lock()
	for !pred() {
		c.cond.Wait()
	}
	do()
	c.cond.L.Unlock()
}

func (c *Cond) SignalDo(do func()) {
	c.cond.L.Lock()
	do()
	c.Signal()
	c.cond.L.Unlock()
}

func (c *Cond) Signal() {
	c.cond.Signal()
}

func (c *Cond) Broadcast() {
	c.cond.Broadcast()
}
