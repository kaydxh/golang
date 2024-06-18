/*
 *Copyright (c) 2024, kaydxh
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
package rate

import (
	"fmt"
	"sync"
	"time"
)

// Limiters support for differenet methods

type MethodLimiter struct {
	mu       sync.RWMutex
	Limiters map[string]*Limiter
	*Limiter
}

func NewMethodLimiter(burst int) *MethodLimiter {
	ml := &MethodLimiter{
		Limiters: make(map[string]*Limiter, 0),
		Limiter:  NewLimiter(burst),
	}

	return ml
}

func (m *MethodLimiter) AddLimiter(method string, limiter *Limiter) error {
	if method == "" {
		return fmt.Errorf("the method can not be empty")
	}
	if limiter == nil {
		return fmt.Errorf("limiter can not be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.Limiters[method]; ok {
		return fmt.Errorf("method: %s already exist", method)
	}
	m.Limiters[method] = limiter
	return nil
}

func (m *MethodLimiter) Allow(method string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	limiter, ok := m.Limiters[method]
	if !ok {
		return m.Limiter.Allow()
	}

	return limiter.Allow()
}

func (m *MethodLimiter) AllowFor(method string, timeout time.Duration) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	limiter, ok := m.Limiters[method]
	if !ok {
		return m.Limiter.AllowFor(timeout)
	}

	return limiter.AllowFor(timeout)
}

func (m *MethodLimiter) Put(method string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	limiter, ok := m.Limiters[method]
	if !ok {
		m.Limiter.Put()
		return
	}

	limiter.Put()
}
