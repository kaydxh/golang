/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
package sync

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// Cond 条件变量不保证顺序性，即signal只是通知wait去获取数据，wait拿到的数据是不是
// signal当时给的数据，不能保证
type Cond struct {
	L       sync.Locker
	ch      chan struct{}
	checker copyChecker
	noCopy  noCopy
}

func NewCond(l sync.Locker) *Cond {
	c := &Cond{
		L:  l,
		ch: make(chan struct{}),
	}
	return c
}

func (c *Cond) wait() {
	c.checker.check()
	c.L.Unlock()
	defer c.L.Lock()
	<-c.ch

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
	c.checker.check()
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
	c.checker.check()
	c.WaitUntilDo(pred, func() error { return nil })
}

func (c *Cond) WaitUntilDo(pred func() bool, do func() error) {
	c.checker.check()
	for {
		err := c.WaitForDo(5*time.Second, pred, do)
		if err == nil {
			break
		}
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
	c.checker.check()
	go func() {
		select {
		case c.ch <- struct{}{}:
		default:
		}
	}()
}

func (c *Cond) Broadcast() {
	c.checker.check()
	go func() {
		close(c.ch)
		c.ch = make(chan struct{})
	}()
}

// copyChecker holds back pointer to itself to detect object copying.
type copyChecker uintptr

func (c *copyChecker) check() {
	if uintptr(*c) != uintptr(unsafe.Pointer(c)) &&
		!atomic.CompareAndSwapUintptr((*uintptr)(c), 0, uintptr(unsafe.Pointer(c))) &&
		uintptr(*c) != uintptr(unsafe.Pointer(c)) {
		panic("sync.Cond is copied")
	}
}

// noCopy may be embedded into structs which must not be copied
// after the first use.
//
// See https://golang.org/issues/8005#issuecomment-190753527
// for details.
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
