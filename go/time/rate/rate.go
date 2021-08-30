package rate

import (
	"sync"
	"time"
)

type Limiter struct {
	mu     sync.Mutex
	burst  int
	tokens int
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
	return &Limiter{
		burst:  b,
		tokens: b,
	}
}

func (lim *Limiter) Bursting() int {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	return lim.burst - lim.tokens
}

// Allow is shorthand for AllowN(time.Now(), 1).
func (lim *Limiter) Allow() bool {
	return lim.AllowN(time.Now(), 1)
}

// AllowN reports whether n events may happen at time now.
// Use this method if you intend to drop / skip events that exceed the rate limit.
// Otherwise use Reserve or Wait.
func (lim *Limiter) AllowN(now time.Time, n int) bool {
	return lim.reserveN(now, n, 0).ok
}

func (lim *Limiter) Put() bool {
	return lim.PutN(1)
}

func (lim *Limiter) PutN(n int) bool {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	lim.tokens += n
	if lim.tokens > lim.burst {
		lim.tokens = lim.burst
	}

	return true
}

// reserveN is a helper method for AllowN, ReserveN, and WaitN.
// maxFutureReserve specifies the maximum reservation wait duration allowed.
// reserveN returns Reservation, not *Reservation, to avoid allocation in AllowN and WaitN.
func (lim *Limiter) reserveN(now time.Time, n int, maxFutureReserve time.Duration) Reservation {
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

// A Reservation holds information about events that are permitted by a Limiter to happen after a delay.
// A Reservation may be canceled, which may enable the Limiter to permit additional events.
type Reservation struct {
	ok     bool
	lim    *Limiter
	tokens int
}
