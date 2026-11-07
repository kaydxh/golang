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

// KeyedRWMutex is a set of reader/writer locks, one per key. Many readers of
// a key may hold it together; a writer holds it alone. A waiting writer stops
// new readers, so a steady stream of readers cannot starve it. Waits are
// bounded by a context. The zero value is ready to use; it must not be copied
// after first use.
type KeyedRWMutex[K comparable] struct {
	mu   sync.Mutex
	held map[K]*rwHold
}

type rwHold struct {
	readers int
	writer  bool
	// writers counts the writers waiting for the key.
	writers int
	// changed is closed, and replaced, whenever the key is released.
	changed chan struct{}
}

// RLock waits until no writer holds or waits for key, then holds it as a
// reader.
func (m *KeyedRWMutex[K]) RLock(ctx context.Context, key K) (func(), error) {
	for {
		m.mu.Lock()
		hold := m.hold(key)
		if !hold.writer && hold.writers == 0 {
			hold.readers++
			m.mu.Unlock()
			return m.releaser(key, false), nil
		}
		changed := hold.changed
		m.mu.Unlock()
		select {
		case <-changed:
		case <-ctx.Done():
			m.mu.Lock()
			m.forget(key)
			m.mu.Unlock()
			return nil, ctx.Err()
		}
	}
}

// Lock waits until key is free, then holds it alone.
func (m *KeyedRWMutex[K]) Lock(ctx context.Context, key K) (func(), error) {
	m.mu.Lock()
	m.hold(key).writers++
	for {
		hold := m.hold(key)
		if !hold.writer && hold.readers == 0 {
			hold.writers--
			hold.writer = true
			m.mu.Unlock()
			return m.releaser(key, true), nil
		}
		changed := hold.changed
		m.mu.Unlock()
		select {
		case <-changed:
			m.mu.Lock()
		case <-ctx.Done():
			m.mu.Lock()
			hold := m.hold(key)
			hold.writers--
			// Readers held back by this writer may go now.
			m.broadcast(hold)
			m.forget(key)
			m.mu.Unlock()
			return nil, ctx.Err()
		}
	}
}

// TryLock holds key alone if it is free right now.
func (m *KeyedRWMutex[K]) TryLock(key K) (func(), bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	hold := m.hold(key)
	if hold.writer || hold.readers > 0 || hold.writers > 0 {
		m.forget(key)
		return nil, false
	}
	hold.writer = true
	return m.releaser(key, true), true
}

// hold returns key's state, creating it; m.mu is held.
func (m *KeyedRWMutex[K]) hold(key K) *rwHold {
	if m.held == nil {
		m.held = map[K]*rwHold{}
	}
	hold := m.held[key]
	if hold == nil {
		hold = &rwHold{changed: make(chan struct{})}
		m.held[key] = hold
	}
	return hold
}

// forget drops an idle key; m.mu is held.
func (m *KeyedRWMutex[K]) forget(key K) {
	if hold := m.held[key]; hold != nil && hold.readers == 0 && !hold.writer && hold.writers == 0 {
		delete(m.held, key)
	}
}

// broadcast wakes everyone waiting on hold; m.mu is held.
func (m *KeyedRWMutex[K]) broadcast(hold *rwHold) {
	close(hold.changed)
	hold.changed = make(chan struct{})
}

// releaser returns an idempotent release of one reader or the writer.
func (m *KeyedRWMutex[K]) releaser(key K, writer bool) func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			m.mu.Lock()
			defer m.mu.Unlock()
			hold := m.hold(key)
			if writer {
				hold.writer = false
			} else {
				hold.readers--
			}
			m.broadcast(hold)
			m.forget(key)
		})
	}
}
