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
package rate_test

import (
  "context"
  "sync"
  "sync/atomic"
  "testing"
  "time"

  "github.com/kaydxh/golang/go/time/rate"
)

func TestQPSLimiter_Allow(t *testing.T) {
  // Create a limiter with 10 QPS and burst of 5
  limiter := rate.NewQPSLimiter(10, 5)

  // Should allow burst requests immediately
  for i := 0; i < 5; i++ {
    if !limiter.Allow() {
      t.Errorf("Allow() should return true for burst request %d", i)
    }
  }

  // 6th request should be rejected (burst exhausted)
  if limiter.Allow() {
    t.Error("Allow() should return false when burst is exhausted")
  }

  // Wait for token refill (100ms = 1 token at 10 QPS)
  time.Sleep(110 * time.Millisecond)

  // Should allow one more request
  if !limiter.Allow() {
    t.Error("Allow() should return true after token refill")
  }
}

func TestQPSLimiter_Wait(t *testing.T) {
  limiter := rate.NewQPSLimiter(100, 1) // 100 QPS, burst 1

  // Exhaust the burst
  limiter.Allow()

  // Wait with timeout
  ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
  defer cancel()

  start := time.Now()
  err := limiter.Wait(ctx)
  elapsed := time.Since(start)

  if err != nil {
    t.Errorf("Wait() error = %v", err)
  }

  // Should wait approximately 10ms (1/100 second)
  if elapsed < 5*time.Millisecond || elapsed > 30*time.Millisecond {
    t.Errorf("Wait() took %v, expected ~10ms", elapsed)
  }
}

func TestQPSLimiter_WaitTimeout(t *testing.T) {
  limiter := rate.NewQPSLimiter(1, 1) // 1 QPS, burst 1

  // Exhaust the burst
  limiter.Allow()

  // Wait with very short timeout (should fail)
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
  defer cancel()

  err := limiter.Wait(ctx)
  if err != context.DeadlineExceeded {
    t.Errorf("Wait() error = %v, want context.DeadlineExceeded", err)
  }
}

func TestQPSLimiter_Concurrent(t *testing.T) {
  limiter := rate.NewQPSLimiter(1000, 100) // 1000 QPS, burst 100

  var allowed int64
  var denied int64
  var wg sync.WaitGroup

  // Launch 200 concurrent requests
  for i := 0; i < 200; i++ {
    wg.Add(1)
    go func() {
      defer wg.Done()
      if limiter.Allow() {
        atomic.AddInt64(&allowed, 1)
      } else {
        atomic.AddInt64(&denied, 1)
      }
    }()
  }

  wg.Wait()

  // Should allow exactly burst size immediately
  if allowed != 100 {
    t.Logf("Allowed: %d, Denied: %d (expected 100 allowed)", allowed, denied)
  }
}

func TestMethodQPSLimiter_DifferentMethods(t *testing.T) {
  limiter := rate.NewMethodQPSLimiter(10, 5) // Default: 10 QPS, burst 5

  // Configure different QPS for different methods
  err := limiter.AddMethod("/api/v1/users", 100, 20)
  if err != nil {
    t.Fatalf("AddMethod() error = %v", err)
  }

  err = limiter.AddMethod("/api/v1/orders", 50, 10)
  if err != nil {
    t.Fatalf("AddMethod() error = %v", err)
  }

  // Test method-specific limits
  // /api/v1/users has burst of 20
  for i := 0; i < 20; i++ {
    if !limiter.Allow("/api/v1/users") {
      t.Errorf("Allow() should return true for /api/v1/users request %d", i)
    }
  }
  if limiter.Allow("/api/v1/users") {
    t.Error("Allow() should return false when /api/v1/users burst is exhausted")
  }

  // /api/v1/orders has burst of 10
  for i := 0; i < 10; i++ {
    if !limiter.Allow("/api/v1/orders") {
      t.Errorf("Allow() should return true for /api/v1/orders request %d", i)
    }
  }
  if limiter.Allow("/api/v1/orders") {
    t.Error("Allow() should return false when /api/v1/orders burst is exhausted")
  }

  // Unknown method uses global limiter (burst 5)
  for i := 0; i < 5; i++ {
    if !limiter.Allow("/api/v1/unknown") {
      t.Errorf("Allow() should return true for unknown method request %d", i)
    }
  }
  if limiter.Allow("/api/v1/unknown") {
    t.Error("Allow() should return false when global burst is exhausted")
  }
}

func TestMethodQPSLimiter_SetMethodQPS(t *testing.T) {
  limiter := rate.NewMethodQPSLimiter(10, 5)

  // Add initial config
  limiter.AddMethod("/api/v1/test", 100, 10)

  // Update QPS
  err := limiter.SetMethodQPS("/api/v1/test", 200, 20)
  if err != nil {
    t.Fatalf("SetMethodQPS() error = %v", err)
  }

  qps, burst, exists := limiter.GetMethodQPS("/api/v1/test")
  if !exists {
    t.Error("Method should exist")
  }
  if qps != 200 {
    t.Errorf("QPS = %v, want 200", qps)
  }
  if burst != 20 {
    t.Errorf("Burst = %v, want 20", burst)
  }
}

func TestMethodQPSLimiter_Stats(t *testing.T) {
  limiter := rate.NewMethodQPSLimiter(10, 5)
  limiter.AddMethod("/api/v1/users", 100, 20)
  limiter.AddMethod("/api/v1/orders", 50, 10)

  stats := limiter.Stats()

  if len(stats) != 3 { // global + 2 methods
    t.Errorf("Stats() returned %d entries, want 3", len(stats))
  }

  // Check that global is included
  hasGlobal := false
  for _, s := range stats {
    if s.Method == "*" {
      hasGlobal = true
      if s.QPS != 10 {
        t.Errorf("Global QPS = %v, want 10", s.QPS)
      }
    }
  }
  if !hasGlobal {
    t.Error("Stats() should include global limiter")
  }
}

func TestMethodQPSLimiter_WithConfigs(t *testing.T) {
  configs := []rate.MethodQPSConfig{
    {Method: "/api/v1/users", QPS: 100, Burst: 20},
    {Method: "/api/v1/orders", QPS: 50, Burst: 10},
    {Method: "/api/v1/products", QPS: 200, Burst: 50},
  }

  limiter, err := rate.NewMethodQPSLimiterWithConfigs(10, 5, configs)
  if err != nil {
    t.Fatalf("NewMethodQPSLimiterWithConfigs() error = %v", err)
  }

  methods := limiter.ListMethods()
  if len(methods) != 3 {
    t.Errorf("ListMethods() returned %d methods, want 3", len(methods))
  }

  // Verify each method's config
  for _, cfg := range configs {
    qps, burst, exists := limiter.GetMethodQPS(cfg.Method)
    if !exists {
      t.Errorf("Method %q should exist", cfg.Method)
    }
    if qps != cfg.QPS {
      t.Errorf("Method %q QPS = %v, want %v", cfg.Method, qps, cfg.QPS)
    }
    if burst != cfg.Burst {
      t.Errorf("Method %q Burst = %v, want %v", cfg.Method, burst, cfg.Burst)
    }
  }
}

func BenchmarkQPSLimiter_Allow(b *testing.B) {
  limiter := rate.NewQPSLimiter(1000000, 10000)

  b.ResetTimer()
  b.RunParallel(func(pb *testing.PB) {
    for pb.Next() {
      limiter.Allow()
    }
  })
}

func BenchmarkMethodQPSLimiter_Allow(b *testing.B) {
  limiter := rate.NewMethodQPSLimiter(1000000, 10000)
  limiter.AddMethod("/api/v1/users", 1000000, 10000)

  b.ResetTimer()
  b.RunParallel(func(pb *testing.PB) {
    for pb.Next() {
      limiter.Allow("/api/v1/users")
    }
  })
}
