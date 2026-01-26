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
  "fmt"
  "sync"
  "time"
)

// MethodQPSLimiter supports different QPS limits for different API methods.
type MethodQPSLimiter struct {
  mu       sync.RWMutex
  limiters map[string]*QPSLimiter // method -> limiter mapping
  global   *QPSLimiter            // default global limiter
}

// MethodQPSConfig defines QPS configuration for a specific method.
type MethodQPSConfig struct {
  Method string  // API method name (e.g., "/api/v1/users", "/service.Method")
  QPS    float64 // queries per second
  Burst  int     // max burst size
}

// NewMethodQPSLimiter creates a new method-level QPS limiter.
// defaultQPS and defaultBurst are used for methods without specific configuration.
func NewMethodQPSLimiter(defaultQPS float64, defaultBurst int) *MethodQPSLimiter {
  return &MethodQPSLimiter{
    limiters: make(map[string]*QPSLimiter),
    global:   NewQPSLimiter(defaultQPS, defaultBurst),
  }
}

// NewMethodQPSLimiterWithConfigs creates a limiter with predefined method configurations.
func NewMethodQPSLimiterWithConfigs(defaultQPS float64, defaultBurst int, configs []MethodQPSConfig) (*MethodQPSLimiter, error) {
  m := NewMethodQPSLimiter(defaultQPS, defaultBurst)
  for _, cfg := range configs {
    if err := m.AddMethod(cfg.Method, cfg.QPS, cfg.Burst); err != nil {
      return nil, err
    }
  }
  return m, nil
}

// AddMethod adds a QPS limit for a specific method.
func (m *MethodQPSLimiter) AddMethod(method string, qps float64, burst int) error {
  if method == "" {
    return fmt.Errorf("method cannot be empty")
  }
  if qps <= 0 {
    return fmt.Errorf("qps must be positive, got %v", qps)
  }
  if burst <= 0 {
    return fmt.Errorf("burst must be positive, got %v", burst)
  }

  m.mu.Lock()
  defer m.mu.Unlock()

  if _, exists := m.limiters[method]; exists {
    return fmt.Errorf("method %q already has a limiter configured", method)
  }

  m.limiters[method] = NewQPSLimiter(qps, burst)
  return nil
}

// SetMethodQPS dynamically updates the QPS for a specific method.
// If the method doesn't exist, it creates a new limiter.
func (m *MethodQPSLimiter) SetMethodQPS(method string, qps float64, burst int) error {
  if method == "" {
    return fmt.Errorf("method cannot be empty")
  }
  if qps <= 0 {
    return fmt.Errorf("qps must be positive, got %v", qps)
  }
  if burst <= 0 {
    return fmt.Errorf("burst must be positive, got %v", burst)
  }

  m.mu.Lock()
  defer m.mu.Unlock()

  if limiter, exists := m.limiters[method]; exists {
    limiter.SetQPS(qps)
    limiter.SetBurst(burst)
  } else {
    m.limiters[method] = NewQPSLimiter(qps, burst)
  }
  return nil
}

// RemoveMethod removes the QPS limit for a specific method.
// After removal, the method will use the global limiter.
func (m *MethodQPSLimiter) RemoveMethod(method string) {
  m.mu.Lock()
  defer m.mu.Unlock()
  delete(m.limiters, method)
}

// SetGlobalQPS updates the global (default) QPS limit.
func (m *MethodQPSLimiter) SetGlobalQPS(qps float64, burst int) {
  m.global.SetQPS(qps)
  m.global.SetBurst(burst)
}

// getLimiter returns the limiter for the given method.
// Falls back to global limiter if no specific limiter is configured.
func (m *MethodQPSLimiter) getLimiter(method string) *QPSLimiter {
  m.mu.RLock()
  defer m.mu.RUnlock()

  if limiter, exists := m.limiters[method]; exists {
    return limiter
  }
  return m.global
}

// Allow checks if a request for the given method is allowed.
// Returns true if allowed, false if rate limited.
func (m *MethodQPSLimiter) Allow(method string) bool {
  return m.getLimiter(method).Allow()
}

// AllowN checks if n requests for the given method are allowed.
func (m *MethodQPSLimiter) AllowN(method string, n int) bool {
  return m.getLimiter(method).AllowN(n)
}

// AllowFor tries to get permission within the specified timeout.
func (m *MethodQPSLimiter) AllowFor(method string, timeout time.Duration) bool {
  return m.getLimiter(method).AllowFor(timeout)
}

// Wait blocks until a request for the given method is allowed.
func (m *MethodQPSLimiter) Wait(ctx context.Context, method string) error {
  return m.getLimiter(method).Wait(ctx)
}

// WaitN blocks until n requests for the given method are allowed.
func (m *MethodQPSLimiter) WaitN(ctx context.Context, method string, n int) error {
  return m.getLimiter(method).WaitN(ctx, n)
}

// Put is a no-op for QPS limiter (for interface compatibility).
func (m *MethodQPSLimiter) Put(method string) {
  // QPS limiter uses time-based token refill, no need to return tokens
}

// GetMethodQPS returns the QPS configuration for a specific method.
// Returns (qps, burst, exists).
func (m *MethodQPSLimiter) GetMethodQPS(method string) (qps float64, burst int, exists bool) {
  m.mu.RLock()
  defer m.mu.RUnlock()

  if limiter, ok := m.limiters[method]; ok {
    return limiter.QPS(), limiter.Burst(), true
  }
  return m.global.QPS(), m.global.Burst(), false
}

// ListMethods returns all methods with specific QPS limits configured.
func (m *MethodQPSLimiter) ListMethods() []string {
  m.mu.RLock()
  defer m.mu.RUnlock()

  methods := make([]string, 0, len(m.limiters))
  for method := range m.limiters {
    methods = append(methods, method)
  }
  return methods
}

// Stats returns current QPS stats for all configured methods.
type MethodQPSStats struct {
  Method string  `json:"method"`
  QPS    float64 `json:"qps"`
  Burst  int     `json:"burst"`
  Tokens float64 `json:"tokens"` // current available tokens
}

func (m *MethodQPSLimiter) Stats() []MethodQPSStats {
  m.mu.RLock()
  defer m.mu.RUnlock()

  stats := make([]MethodQPSStats, 0, len(m.limiters)+1)

  // Global limiter stats
  stats = append(stats, MethodQPSStats{
    Method: "*",
    QPS:    m.global.QPS(),
    Burst:  m.global.Burst(),
    Tokens: m.global.Tokens(),
  })

  // Per-method stats
  for method, limiter := range m.limiters {
    stats = append(stats, MethodQPSStats{
      Method: method,
      QPS:    limiter.QPS(),
      Burst:  limiter.Burst(),
      Tokens: limiter.Tokens(),
    })
  }

  return stats
}
