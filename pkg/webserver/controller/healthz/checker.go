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
package healthz

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

// HealthChecker is the interface for health check implementations.
type HealthChecker interface {
	// Name returns the name of the health check.
	Name() string
	// Check performs the health check and returns an error if unhealthy.
	Check(ctx context.Context) error
}

// HealthCheckResult represents the result of a single health check.
type HealthCheckResult struct {
	Name    string `json:"name"`
	Healthy bool   `json:"healthy"`
	Error   string `json:"error,omitempty"`
	Latency string `json:"latency,omitempty"`
}

// HealthCheckResponse represents the overall health check response.
type HealthCheckResponse struct {
	Status    string              `json:"status"`
	Checks    []HealthCheckResult `json:"checks,omitempty"`
	Timestamp string              `json:"timestamp"`
}

// PingHealthChecker is a basic health check that always returns healthy.
type PingHealthChecker struct{}

func (p PingHealthChecker) Name() string {
	return "ping"
}

func (p PingHealthChecker) Check(ctx context.Context) error {
	return nil
}

// HTTPHealthChecker checks the health of an HTTP endpoint.
type HTTPHealthChecker struct {
	name    string
	url     string
	timeout time.Duration
	client  *http.Client
}

// NewHTTPHealthChecker creates a new HTTP health checker.
func NewHTTPHealthChecker(name, url string, timeout time.Duration) *HTTPHealthChecker {
	return &HTTPHealthChecker{
		name:    name,
		url:     url,
		timeout: timeout,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (h *HTTPHealthChecker) Name() string {
	return h.name
}

func (h *HTTPHealthChecker) Check(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unhealthy status code: %d", resp.StatusCode)
	}

	return nil
}

// TCPHealthChecker checks the health of a TCP endpoint.
type TCPHealthChecker struct {
	name    string
	addr    string
	timeout time.Duration
}

// NewTCPHealthChecker creates a new TCP health checker.
func NewTCPHealthChecker(name, addr string, timeout time.Duration) *TCPHealthChecker {
	return &TCPHealthChecker{
		name:    name,
		addr:    addr,
		timeout: timeout,
	}
}

func (t *TCPHealthChecker) Name() string {
	return t.name
}

func (t *TCPHealthChecker) Check(ctx context.Context) error {
	dialer := &net.Dialer{
		Timeout: t.timeout,
	}

	conn, err := dialer.DialContext(ctx, "tcp", t.addr)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", t.addr, err)
	}
	defer conn.Close()

	return nil
}

// FuncHealthChecker wraps a function as a health checker.
type FuncHealthChecker struct {
	name      string
	checkFunc func(ctx context.Context) error
}

// NewFuncHealthChecker creates a new function-based health checker.
func NewFuncHealthChecker(name string, checkFunc func(ctx context.Context) error) *FuncHealthChecker {
	return &FuncHealthChecker{
		name:      name,
		checkFunc: checkFunc,
	}
}

func (f *FuncHealthChecker) Name() string {
	return f.name
}

func (f *FuncHealthChecker) Check(ctx context.Context) error {
	if f.checkFunc == nil {
		return nil
	}
	return f.checkFunc(ctx)
}

// CompositeHealthChecker combines multiple health checkers.
type CompositeHealthChecker struct {
	mu       sync.RWMutex
	checkers []HealthChecker
}

// NewCompositeHealthChecker creates a new composite health checker.
func NewCompositeHealthChecker(checkers ...HealthChecker) *CompositeHealthChecker {
	return &CompositeHealthChecker{
		checkers: checkers,
	}
}

// AddChecker adds a health checker to the composite.
func (c *CompositeHealthChecker) AddChecker(checker HealthChecker) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checkers = append(c.checkers, checker)
}

// RemoveChecker removes a health checker by name.
func (c *CompositeHealthChecker) RemoveChecker(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, checker := range c.checkers {
		if checker.Name() == name {
			c.checkers = append(c.checkers[:i], c.checkers[i+1:]...)
			return
		}
	}
}

// Check performs all health checks and returns an error if any fails.
func (c *CompositeHealthChecker) Check(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, checker := range c.checkers {
		if err := checker.Check(ctx); err != nil {
			return fmt.Errorf("%s: %w", checker.Name(), err)
		}
	}
	return nil
}

// CheckAll performs all health checks and returns detailed results.
func (c *CompositeHealthChecker) CheckAll(ctx context.Context) ([]HealthCheckResult, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	results := make([]HealthCheckResult, 0, len(c.checkers))
	allHealthy := true

	for _, checker := range c.checkers {
		start := time.Now()
		err := checker.Check(ctx)
		latency := time.Since(start)

		result := HealthCheckResult{
			Name:    checker.Name(),
			Healthy: err == nil,
			Latency: latency.String(),
		}

		if err != nil {
			result.Error = err.Error()
			allHealthy = false
		}

		results = append(results, result)
	}

	return results, allHealthy
}

// Checkers returns a copy of the registered checkers.
func (c *CompositeHealthChecker) Checkers() []HealthChecker {
	c.mu.RLock()
	defer c.mu.RUnlock()

	checkers := make([]HealthChecker, len(c.checkers))
	copy(checkers, c.checkers)
	return checkers
}
