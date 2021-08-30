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
