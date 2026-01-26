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
package rate

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Limiter 基于 channel 信号量的并发控制限流器
type Limiter struct {
	mu    sync.Mutex
	burst int
	sem   chan struct{}
}

// NewLimiter returns a new Limiter that permits bursts of at most b tokens.
func NewLimiter(b int) *Limiter {
	if b <= 0 {
		b = 1
	}
	l := &Limiter{
		burst: b,
		sem:   make(chan struct{}, b),
	}
	// 预填充令牌
	for i := 0; i < b; i++ {
		l.sem <- struct{}{}
	}
	return l
}

// Burst returns the maximum burst size.
func (lim *Limiter) Burst() int {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.burst
}

// Tokens returns the current available tokens.
func (lim *Limiter) Tokens() int {
	return len(lim.sem)
}

// Bursting returns the number of tokens currently in use.
func (lim *Limiter) Bursting() int {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.burst - len(lim.sem)
}

// Allow attempts to acquire a token immediately without waiting.
func (lim *Limiter) Allow() bool {
	select {
	case <-lim.sem:
		return true
	default:
		return false
	}
}

// AllowN attempts to acquire n tokens immediately without waiting.
func (lim *Limiter) AllowN(n int) bool {
	if n <= 0 {
		return true
	}
	if n > lim.burst {
		return false
	}

	acquired := 0
	for i := 0; i < n; i++ {
		select {
		case <-lim.sem:
			acquired++
		default:
			// 获取失败，归还已获取的令牌
			for j := 0; j < acquired; j++ {
				lim.sem <- struct{}{}
			}
			return false
		}
	}
	return true
}

// AllowWaitUntil waits indefinitely until a token is available.
func (lim *Limiter) AllowWaitUntil() bool {
	<-lim.sem
	return true
}

// AllowFor attempts to acquire a token with timeout.
// If timeout < 0, it waits indefinitely.
// If timeout == 0, it returns immediately.
func (lim *Limiter) AllowFor(timeout time.Duration) bool {
	if timeout == 0 {
		return lim.Allow()
	}
	if timeout < 0 {
		return lim.AllowWaitUntil()
	}

	select {
	case <-lim.sem:
		return true
	case <-time.After(timeout):
		return false
	}
}

// AllowContext attempts to acquire a token with context support.
func (lim *Limiter) AllowContext(ctx context.Context) error {
	select {
	case <-lim.sem:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Put returns one token to the limiter.
func (lim *Limiter) Put() {
	lim.PutN(1)
}

// PutN returns n tokens to the limiter.
func (lim *Limiter) PutN(n int) {
	for i := 0; i < n; i++ {
		select {
		case lim.sem <- struct{}{}:
		default:
			// channel is full, ignore
		}
	}
}

// WaitFor waits for one token with timeout.
func (lim *Limiter) WaitFor(timeout time.Duration) error {
	return lim.WaitN(timeout, 1)
}

// WaitN waits for n tokens with timeout.
func (lim *Limiter) WaitN(timeout time.Duration, n int) error {
	if n <= 0 {
		return nil
	}
	if n > lim.burst {
		return fmt.Errorf("rate: Wait(n=%d) exceeds limiter's burst %d", n, lim.burst)
	}

	acquired := 0
	// 失败时归还已获取的令牌
	defer func() {
		if acquired < n {
			lim.PutN(acquired)
		}
	}()

	// 无限等待
	if timeout < 0 {
		for acquired < n {
			<-lim.sem
			acquired++
		}
		return nil
	}

	// 带超时等待
	deadline := time.Now().Add(timeout)
	for acquired < n {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return fmt.Errorf("rate: Wait timeout after acquiring %d/%d tokens", acquired, n)
		}

		select {
		case <-lim.sem:
			acquired++
		case <-time.After(remaining):
			return fmt.Errorf("rate: Wait timeout after acquiring %d/%d tokens", acquired, n)
		}
	}
	return nil
}

// WaitContext waits for one token with context support.
func (lim *Limiter) WaitContext(ctx context.Context) error {
	return lim.WaitNContext(ctx, 1)
}

// WaitNContext waits for n tokens with context support.
func (lim *Limiter) WaitNContext(ctx context.Context, n int) error {
	if n <= 0 {
		return nil
	}
	if n > lim.burst {
		return fmt.Errorf("rate: Wait(n=%d) exceeds limiter's burst %d", n, lim.burst)
	}

	acquired := 0
	defer func() {
		if acquired < n {
			lim.PutN(acquired)
		}
	}()

	for acquired < n {
		select {
		case <-lim.sem:
			acquired++
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

// SetBurst dynamically adjusts the burst size.
func (lim *Limiter) SetBurst(newBurst int) {
	if newBurst <= 0 {
		newBurst = 1
	}

	lim.mu.Lock()
	defer lim.mu.Unlock()

	oldBurst := lim.burst
	if newBurst == oldBurst {
		return
	}

	// 创建新的 channel
	newSem := make(chan struct{}, newBurst)

	// 转移现有令牌
	currentTokens := len(lim.sem)
	transferTokens := currentTokens
	if transferTokens > newBurst {
		transferTokens = newBurst
	}

	for i := 0; i < transferTokens; i++ {
		newSem <- struct{}{}
	}

	// 消耗旧 channel 中多余的令牌
	for i := 0; i < currentTokens-transferTokens; i++ {
		<-lim.sem
	}

	lim.sem = newSem
	lim.burst = newBurst
}
