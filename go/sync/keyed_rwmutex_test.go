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
package sync_test

import (
	"context"
	"errors"
	"testing"
	"time"

	sync_ "github.com/kaydxh/golang/go/sync"
)

func returnsWithin(t *testing.T, d time.Duration, ch <-chan func()) func() {
	t.Helper()
	select {
	case release := <-ch:
		return release
	case <-time.After(d):
		return nil
	}
}

func lockAsync(m *sync_.KeyedRWMutex[string], key string, writer bool) <-chan func() {
	got := make(chan func(), 1)
	go func() {
		var release func()
		var err error
		if writer {
			release, err = m.Lock(context.Background(), key)
		} else {
			release, err = m.RLock(context.Background(), key)
		}
		if err == nil {
			got <- release
		}
	}()
	return got
}

func TestKeyedRWMutexReadersShare(t *testing.T) {
	var m sync_.KeyedRWMutex[string]
	first, err := m.RLock(context.Background(), "a")
	if err != nil {
		t.Fatal(err)
	}
	second := returnsWithin(t, time.Second, lockAsync(&m, "a", false))
	if second == nil {
		t.Fatal("second reader waited")
	}
	if _, ok := m.TryLock("a"); ok {
		t.Fatal("writer took a key readers hold")
	}
	first()
	second()
	release, ok := m.TryLock("a")
	if !ok {
		t.Fatal("writer could not take the released key")
	}
	release()
}

func TestKeyedRWMutexWriterWaitsForReadersAndHoldsNewOnesBack(t *testing.T) {
	var m sync_.KeyedRWMutex[string]
	reader, _ := m.RLock(context.Background(), "a")
	writer := lockAsync(&m, "a", true)
	if returnsWithin(t, 30*time.Millisecond, writer) != nil {
		t.Fatal("writer took the key from a reader")
	}
	late := lockAsync(&m, "a", false)
	if returnsWithin(t, 30*time.Millisecond, late) != nil {
		t.Fatal("a new reader passed a waiting writer")
	}
	reader()
	release := returnsWithin(t, time.Second, writer)
	if release == nil {
		t.Fatal("writer did not get the key")
	}
	if returnsWithin(t, 30*time.Millisecond, late) != nil {
		t.Fatal("reader joined a writer")
	}
	release()
	if next := returnsWithin(t, time.Second, late); next == nil {
		t.Fatal("reader did not follow the writer")
	} else {
		next()
	}
}

func TestKeyedRWMutexKeysAreIndependent(t *testing.T) {
	var m sync_.KeyedRWMutex[string]
	a, _ := m.Lock(context.Background(), "a")
	defer a()
	b, ok := m.TryLock("b")
	if !ok {
		t.Fatal("another key was blocked")
	}
	b()
}

func TestKeyedRWMutexCancelledWriterLetsReadersGo(t *testing.T) {
	var m sync_.KeyedRWMutex[string]
	reader, _ := m.RLock(context.Background(), "a")
	defer reader()
	ctx, cancel := context.WithCancel(context.Background())
	failed := make(chan error, 1)
	go func() {
		_, err := m.Lock(ctx, "a")
		failed <- err
	}()
	time.Sleep(20 * time.Millisecond)
	late := lockAsync(&m, "a", false)
	cancel()
	if err := <-failed; !errors.Is(err, context.Canceled) {
		t.Fatalf("Lock = %v, want canceled", err)
	}
	if next := returnsWithin(t, time.Second, late); next == nil {
		t.Fatal("reader stayed blocked by a cancelled writer")
	} else {
		next()
	}
}

func TestKeyedRWMutexReleaseIsIdempotent(t *testing.T) {
	var m sync_.KeyedRWMutex[string]
	reader, _ := m.RLock(context.Background(), "a")
	reader()
	reader()
	release, ok := m.TryLock("a")
	if !ok {
		t.Fatal("double release left the key held")
	}
	release()
	release()
	if again, ok := m.TryLock("a"); !ok {
		t.Fatal("double release of the writer left the key held")
	} else {
		again()
	}
}
