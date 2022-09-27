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
package rate_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	rate_ "github.com/kaydxh/golang/go/time/rate"
	"gotest.tools/assert"
)

func TestLimitAllow(t *testing.T) {
	// The test runs for a few seconds executing many requests and then checks
	// that overall number of requests is reasonable.
	const (
		limit = 100
		burst = 100
	)

	var (
		numAll    int32 = 0
		numOK     int32 = 0
		numFailed int32 = 0
	)

	lim := rate_.NewLimiter(burst)

	var wg sync.WaitGroup
	f := func() {
		if ok := lim.Allow(); ok {
			defer lim.Put()
			atomic.AddInt32(&numOK, 1)
			time.Sleep(10 * time.Millisecond)

		} else {
			atomic.AddInt32(&numFailed, 1)
			//			t.Logf("not Allowed, bursting: %v", lim.Bursting())
		}

		atomic.AddInt32(&numAll, 1)
		wg.Done()
	}

	start := time.Now()
	end := start.Add(5 * time.Second)

	for time.Now().Before(end) {
		wg.Add(1)
		go f()
	}
	wg.Wait()

	assert.Equal(t, numAll, numOK+numFailed)

	// numOK should get very close to the number of requests allowed  50000.
	t.Logf("%v request ==> %v requset Allowed, %v request Unallowed ", numAll, numOK, numFailed)

}

func TestLimitAllowFor(t *testing.T) {
	// The test runs for a few seconds executing many requests and then checks
	// that overall number of requests is reasonable.
	const (
		limit = 100
		burst = 100
	)

	var (
		numAll    int32 = 0
		numOK     int32 = 0
		numFailed int32 = 0
	)

	//timout > 5s, so must process 50001 times
	timeout := 6 * time.Second

	lim := rate_.NewLimiter(burst)

	var wg sync.WaitGroup
	f := func() {
		if ok := lim.AllowFor(timeout); ok {
			defer lim.Put()
			atomic.AddInt32(&numOK, 1)
			//about process 50000 times in 5s
			time.Sleep(10 * time.Millisecond)

		} else {
			atomic.AddInt32(&numFailed, 1)
			//			t.Logf("not Allowed, bursting: %v", lim.Bursting())
		}

		atomic.AddInt32(&numAll, 1)
		wg.Done()
	}

	for i := 0; i < 50001; i++ {
		wg.Add(1)
		go f()
	}
	wg.Wait()

	assert.Equal(t, numAll, numOK+numFailed)

	// numOK should get very close to the number of requests allowed  50000.
	t.Logf("%v request ==> %v requset Allowed, %v request Unallowed ", numAll, numOK, numFailed)

}
