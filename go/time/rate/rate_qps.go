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
  "context"
  "sync"
  "time"
)

// QPSLimiter is a rate limiter based on token bucket algorithm.
// It controls the rate of requests per second (QPS).
type QPSLimiter struct {
  mu         sync.Mutex
  qps        float64   // tokens generated per second
  burst      int       // max tokens (bucket size)
  tokens     float64   // current available tokens
  lastUpdate time.Time // last token update time
}

// NewQPSLimiter creates a new QPS-based rate limiter.
// qps: the rate of tokens generated per second (queries per second)
// burst: the maximum number of tokens that can be stored (allows burst traffic)
func NewQPSLimiter(qps float64, burst int) *QPSLimiter {
  return &QPSLimiter{
    qps:        qps,
    burst:      burst,
    tokens:     float64(burst),
    lastUpdate: time.Now(),
  }
}

// SetQPS dynamically updates the QPS rate.
func (l *QPSLimiter) SetQPS(qps float64) {
  l.mu.Lock()
  defer l.mu.Unlock()
  l.qps = qps
}

// SetBurst dynamically updates the burst size.
func (l *QPSLimiter) SetBurst(burst int) {
  l.mu.Lock()
  defer l.mu.Unlock()
  l.burst = burst
  if l.tokens > float64(burst) {
    l.tokens = float64(burst)
  }
}

// QPS returns the current QPS rate.
func (l *QPSLimiter) QPS() float64 {
  l.mu.Lock()
  defer l.mu.Unlock()
  return l.qps
}

// Burst returns the current burst size.
func (l *QPSLimiter) Burst() int {
  l.mu.Lock()
  defer l.mu.Unlock()
  return l.burst
}

// Tokens returns the current available tokens.
func (l *QPSLimiter) Tokens() float64 {
  l.mu.Lock()
  defer l.mu.Unlock()
  l.refillTokens()
  return l.tokens
}

// refillTokens adds tokens based on elapsed time since last update.
// Must be called with lock held.
func (l *QPSLimiter) refillTokens() {
  now := time.Now()
  elapsed := now.Sub(l.lastUpdate).Seconds()
  l.lastUpdate = now

  // Add tokens based on QPS and elapsed time
  l.tokens += elapsed * l.qps
  if l.tokens > float64(l.burst) {
    l.tokens = float64(l.burst)
  }
}

// Allow reports whether one event may happen now.
// It returns true if the event is allowed, false otherwise.
// This method does not block.
func (l *QPSLimiter) Allow() bool {
  return l.AllowN(1)
}

// AllowN reports whether n events may happen now.
// It returns true if the events are allowed, false otherwise.
// This method does not block.
func (l *QPSLimiter) AllowN(n int) bool {
  l.mu.Lock()
  defer l.mu.Unlock()

  l.refillTokens()

  if l.tokens >= float64(n) {
    l.tokens -= float64(n)
    return true
  }
  return false
}

// Wait blocks until one event is allowed or the context is canceled.
func (l *QPSLimiter) Wait(ctx context.Context) error {
  return l.WaitN(ctx, 1)
}

// WaitN blocks until n events are allowed or the context is canceled.
func (l *QPSLimiter) WaitN(ctx context.Context, n int) error {
  l.mu.Lock()
  l.refillTokens()

  if l.tokens >= float64(n) {
    l.tokens -= float64(n)
    l.mu.Unlock()
    return nil
  }

  // Calculate wait time
  tokensNeeded := float64(n) - l.tokens
  waitDuration := time.Duration(tokensNeeded / l.qps * float64(time.Second))
  l.mu.Unlock()

  // Check if context has deadline and if we can wait that long
  if deadline, ok := ctx.Deadline(); ok {
    if time.Until(deadline) < waitDuration {
      return context.DeadlineExceeded
    }
  }

  timer := time.NewTimer(waitDuration)
  defer timer.Stop()

  select {
  case <-timer.C:
    l.mu.Lock()
    l.refillTokens()
    if l.tokens >= float64(n) {
      l.tokens -= float64(n)
      l.mu.Unlock()
      return nil
    }
    l.mu.Unlock()
    // Retry if not enough tokens (rare case due to timing)
    return l.WaitN(ctx, n)
  case <-ctx.Done():
    return ctx.Err()
  }
}

// AllowFor tries to get a token within the specified timeout.
// Returns true if a token is obtained, false if timeout expires.
func (l *QPSLimiter) AllowFor(timeout time.Duration) bool {
  if timeout <= 0 {
    return l.Allow()
  }

  ctx, cancel := context.WithTimeout(context.Background(), timeout)
  defer cancel()

  err := l.Wait(ctx)
  return err == nil
}

// Put is a no-op for QPS limiter (tokens are time-based, not returned).
// This method exists to satisfy the Limiter interface compatibility.
func (l *QPSLimiter) Put() {
  // QPS limiter uses time-based token refill, no need to return tokens
}

// Reserve returns a Reservation that can be used to wait for the token.
func (l *QPSLimiter) Reserve() *QPSReservation {
  return l.ReserveN(1)
}

// ReserveN returns a Reservation for n tokens.
func (l *QPSLimiter) ReserveN(n int) *QPSReservation {
  l.mu.Lock()
  defer l.mu.Unlock()

  l.refillTokens()

  r := &QPSReservation{
    limiter: l,
    tokens:  n,
  }

  if l.tokens >= float64(n) {
    l.tokens -= float64(n)
    r.ok = true
    r.timeToAct = time.Now()
    return r
  }

  // Calculate when we'll have enough tokens
  tokensNeeded := float64(n) - l.tokens
  waitDuration := time.Duration(tokensNeeded / l.qps * float64(time.Second))

  r.ok = true
  r.timeToAct = time.Now().Add(waitDuration)
  l.tokens = 0 // Reserve all current tokens

  return r
}

// QPSReservation holds information about a rate-limited event.
type QPSReservation struct {
  limiter   *QPSLimiter
  ok        bool
  tokens    int
  timeToAct time.Time
}

// OK returns whether the reservation is valid.
func (r *QPSReservation) OK() bool {
  return r.ok
}

// Delay returns the duration to wait before the event can happen.
func (r *QPSReservation) Delay() time.Duration {
  delay := time.Until(r.timeToAct)
  if delay < 0 {
    return 0
  }
  return delay
}

// Cancel cancels the reservation.
func (r *QPSReservation) Cancel() {
  if !r.ok {
    return
  }

  r.limiter.mu.Lock()
  defer r.limiter.mu.Unlock()

  // Return tokens if we haven't acted yet
  if time.Now().Before(r.timeToAct) {
    r.limiter.tokens += float64(r.tokens)
    if r.limiter.tokens > float64(r.limiter.burst) {
      r.limiter.tokens = float64(r.limiter.burst)
    }
  }
  r.ok = false
}
