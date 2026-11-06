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

func TestKeyedMutexTryLockReportsHolder(t *testing.T) {
	var m sync_.KeyedMutex[string, string]
	release, _, ok := m.TryLock("a", "alice")
	if !ok {
		t.Fatal("first TryLock failed")
	}
	if _, holder, ok := m.TryLock("a", "bob"); ok || holder != "alice" {
		t.Fatalf("second TryLock = %q, %v; want alice, false", holder, ok)
	}
	if other, _, ok := m.TryLock("b", "bob"); !ok {
		t.Fatal("another key was blocked")
	} else {
		other()
	}
	if holder, held := m.Holder("a"); !held || holder != "alice" {
		t.Fatalf("Holder = %q, %v", holder, held)
	}
	release()
	release() // idempotent
	if _, held := m.Holder("a"); held {
		t.Fatal("key still held after release")
	}
	if again, _, ok := m.TryLock("a", "bob"); !ok {
		t.Fatal("released key not free")
	} else {
		again()
	}
}

func TestKeyedMutexLockWaitsForRelease(t *testing.T) {
	var m sync_.KeyedMutex[string, int]
	release, _, _ := m.TryLock("a", 1)
	got := make(chan func(), 1)
	go func() {
		next, err := m.Lock(context.Background(), "a", 2)
		if err != nil {
			t.Error(err)
		}
		got <- next
	}()
	select {
	case <-got:
		t.Fatal("Lock returned while key was held")
	case <-time.After(20 * time.Millisecond):
	}
	release()
	select {
	case next := <-got:
		if holder, _ := m.Holder("a"); holder != 2 {
			t.Fatalf("holder = %d, want 2", holder)
		}
		next()
	case <-time.After(time.Second):
		t.Fatal("Lock did not take the released key")
	}
}

func TestKeyedMutexLockStopsWithContext(t *testing.T) {
	var m sync_.KeyedMutex[string, int]
	release, _, _ := m.TryLock("a", 1)
	defer release()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if _, err := m.Lock(ctx, "a", 2); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Lock = %v, want deadline exceeded", err)
	}
}
