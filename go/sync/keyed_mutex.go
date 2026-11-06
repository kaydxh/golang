/*
 *Copyright (c) 2026, kaydxh
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
	"context"
	"sync"
)

// KeyedMutex is a set of mutexes, one per key, each remembering who holds
// it so a caller that finds a key taken can say by whom. Lock waits are
// bounded by a context. The zero value is ready to use; it must not be
// copied after first use.
type KeyedMutex[K comparable, H any] struct {
	mu   sync.Mutex
	held map[K]*keyedHold[H]
}

type keyedHold[H any] struct {
	holder H
	done   chan struct{}
}

// TryLock takes key for holder without waiting. When key is held, it
// reports the current holder and false.
func (m *KeyedMutex[K, H]) TryLock(key K, holder H) (release func(), current H, ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if hold := m.held[key]; hold != nil {
		return nil, hold.holder, false
	}
	return m.take(key, holder), current, true
}

// Lock waits for key until ctx ends, then takes it for holder.
func (m *KeyedMutex[K, H]) Lock(ctx context.Context, key K, holder H) (func(), error) {
	for {
		m.mu.Lock()
		hold := m.held[key]
		if hold == nil {
			release := m.take(key, holder)
			m.mu.Unlock()
			return release, nil
		}
		m.mu.Unlock()
		select {
		case <-hold.done:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// Holder reports who holds key, if anyone.
func (m *KeyedMutex[K, H]) Holder(key K) (H, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if hold := m.held[key]; hold != nil {
		return hold.holder, true
	}
	var none H
	return none, false
}

// take records holder on key; m.mu is held. The returned release is
// idempotent.
func (m *KeyedMutex[K, H]) take(key K, holder H) func() {
	if m.held == nil {
		m.held = map[K]*keyedHold[H]{}
	}
	hold := &keyedHold[H]{holder: holder, done: make(chan struct{})}
	m.held[key] = hold
	var once sync.Once
	return func() {
		once.Do(func() {
			m.mu.Lock()
			if m.held[key] == hold {
				delete(m.held, key)
			}
			m.mu.Unlock()
			close(hold.done)
		})
	}
}
