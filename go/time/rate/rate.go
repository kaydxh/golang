package rate

import (
	"fmt"
	"sync"
	"time"

	sync_ "github.com/kaydxh/golang/go/sync"
)

type Limiter struct {
	mu     sync.Mutex
	burst  int
	tokens int
	cond   *sync_.Cond
}

// Burst returns the maximum burst size. Burst is the maximum number of tokens
// that can be consumed in a single call to Allow, Reserve, or Wait, so higher
// Burst values allow more events to happen at once.
// A zero Burst allows no events, unless limit == Inf.
func (lim *Limiter) Burst() int {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.burst
}

// NewLimiter returns a new Limiter that allows events up to rate r and permits
// bursts of at most b tokens.
func NewLimiter(b int) *Limiter {
	l := &Limiter{
		burst:  b,
		tokens: b,
	}

	l.cond = sync_.NewCond(&l.mu)

	return l
}

func (lim *Limiter) Bursting() int {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	return lim.burst - lim.tokens
}

// Allow is shorthand for AllowN(time.Now(), 1).
func (lim *Limiter) Allow() bool {
	return lim.AllowN(time.Now(), 1, 0)
}

// AllowWaitUntil is shorthand for AllowFor(-1).
func (lim *Limiter) AllowWaitUntil() bool {
	return lim.AllowFor(-1 * time.Second)
}

// Allow is shorthand for AllowN(time.Now(), 1).
func (lim *Limiter) AllowFor(timeout time.Duration) bool {
	return lim.AllowN(time.Now(), 1, timeout)
}

// AllowN reports whether n events may happen at time now.
// Use this method if you intend to drop / skip events that exceed the rate limit.
// Otherwise use Reserve or Wait.
func (lim *Limiter) AllowN(now time.Time, n int, timeout time.Duration) bool {
	ok := lim.reserveN(now, n).ok
	if ok {
		return true
	}

	if timeout == 0 {
		return false
	}

	err := lim.WaitFor(timeout)
	if err != nil {
		return false
	}

	return true
}

func (lim *Limiter) Put() {
	lim.PutN(1)
}

func (lim *Limiter) PutN(n int) {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	lim.tokens += n
	if lim.tokens > lim.burst {
		lim.tokens = lim.burst
	}

	lim.cond.Signal()
}

// reserveN is a helper method for AllowN, ReserveN, and WaitN.
// maxFutureReserve specifies the maximum reservation wait duration allowed.
// reserveN returns Reservation, not *Reservation, to avoid allocation in AllowN and WaitN.
func (lim *Limiter) reserveN(now time.Time, n int) Reservation {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	if lim.tokens >= n {
		lim.tokens -= n
		return Reservation{
			ok:  true,
			lim: lim,
		}
	}

	return Reservation{
		ok:  false,
		lim: lim,
	}
}

// ReserveN returns a Reservation that indicates how long the caller must wait before n events happen.
// The Limiter takes this Reservation into account when allowing future events.
// The returned Reservation’s OK() method returns false if n exceeds the Limiter's burst size.
// Usage example:
//   r := lim.ReserveN(time.Now(), 1)
//   if !r.OK() {
//     // Not allowed to act! Did you remember to set lim.burst to be > 0 ?
//     return
//   }
//	err := lim.WaitFor(timeout)
//   Act()
// Use this method if you wish to wait and slow down in accordance with the rate limit without dropping events.
// If you need to respect a deadline or cancel the delay, use Wait instead.
// To drop or skip events exceeding rate limit, use Allow instead.

// A Reservation holds information about events that are permitted by a Limiter to happen after a delay.
// A Reservation may be canceled, which may enable the Limiter to permit additional events.
// Wait is shorthand for WaitN(ctx, 1).
func (lim *Limiter) WaitFor(timeout time.Duration) (err error) {
	return lim.WaitN(timeout, 1)
}

// WaitN blocks until lim permits n events to happen.
// It returns an error if n exceeds the Limiter's burst size, the Context is
// canceled, or the expected wait time exceeds the Context's Deadline.
// The burst limit is ignored if the rate limit is Inf.
func (lim *Limiter) WaitN(timeout time.Duration, n int) (err error) {
	lim.mu.Lock()
	burst := lim.burst
	lim.mu.Unlock()
	if n > burst {
		return fmt.Errorf("rate: Wait(n=%d) exceeds limiter's burst %d", n, burst)
	}

	pred := func() bool {
		return lim.tokens >= n
	}

	do := func() error {
		lim.tokens -= n
		return nil
	}

	if timeout >= 0 {
		return lim.cond.WaitForDo(timeout, pred, do)
	} else {
		lim.cond.WaitUntilDo(pred, do)
		return nil
	}
}

type Reservation struct {
	ok     bool
	lim    *Limiter
	tokens int
}
