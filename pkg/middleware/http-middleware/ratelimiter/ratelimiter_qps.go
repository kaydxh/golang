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
package ratelimit

import (
  "encoding/json"
  "fmt"
  "net/http"
  "time"

  rate_ "github.com/kaydxh/golang/go/time/rate"
  logs_ "github.com/kaydxh/golang/pkg/logs"
)

// QPSRateLimiter provides QPS-based rate limiting for HTTP handlers.
// It supports different QPS limits for different URL paths.
type QPSRateLimiter struct {
  limiter *rate_.MethodQPSLimiter
}

// NewQPSRateLimiter creates a new QPS-based rate limiter.
// defaultQPS: default queries per second for paths without specific config
// defaultBurst: default burst size for paths without specific config
func NewQPSRateLimiter(defaultQPS float64, defaultBurst int) *QPSRateLimiter {
  return &QPSRateLimiter{
    limiter: rate_.NewMethodQPSLimiter(defaultQPS, defaultBurst),
  }
}

// NewQPSRateLimiterWithConfigs creates a QPS limiter with predefined path configurations.
func NewQPSRateLimiterWithConfigs(
  defaultQPS float64,
  defaultBurst int,
  configs []rate_.MethodQPSConfig,
) (*QPSRateLimiter, error) {
  limiter, err := rate_.NewMethodQPSLimiterWithConfigs(defaultQPS, defaultBurst, configs)
  if err != nil {
    return nil, err
  }
  return &QPSRateLimiter{limiter: limiter}, nil
}

// AddPath adds a QPS limit for a specific URL path.
func (l *QPSRateLimiter) AddPath(path string, qps float64, burst int) error {
  return l.limiter.AddMethod(path, qps, burst)
}

// SetPathQPS dynamically updates the QPS for a specific path.
func (l *QPSRateLimiter) SetPathQPS(path string, qps float64, burst int) error {
  return l.limiter.SetMethodQPS(path, qps, burst)
}

// RemovePath removes the QPS limit for a specific path.
func (l *QPSRateLimiter) RemovePath(path string) {
  l.limiter.RemoveMethod(path)
}

// Handler returns an HTTP middleware that applies QPS rate limiting.
// The rate limit is based on the request URL path.
func (l *QPSRateLimiter) Handler(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path

    if !l.limiter.Allow(path) {
      logger := logs_.GetLogger(r.Context())
      logger.WithField("path", path).
        WithField("method", r.Method).
        Warn("request rejected by QPS rate limiter")

      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("Retry-After", "1") // Suggest retry after 1 second
      w.WriteHeader(http.StatusTooManyRequests)

      resp := map[string]interface{}{
        "error":   "rate_limit_exceeded",
        "message": fmt.Sprintf("%s %s is rejected by http_ratelimit_qps middleware, QPS limit exceeded", r.Method, path),
        "code":    http.StatusTooManyRequests,
      }
      json.NewEncoder(w).Encode(resp)
      return
    }

    next.ServeHTTP(w, r)
  })
}

// HandlerFunc returns an HTTP middleware function.
func (l *QPSRateLimiter) HandlerFunc(next http.HandlerFunc) http.HandlerFunc {
  return l.Handler(next).ServeHTTP
}

// HandlerWait returns an HTTP middleware that waits for rate limit token
// instead of rejecting immediately.
func (l *QPSRateLimiter) HandlerWait(next http.Handler, timeout time.Duration) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path

    if !l.limiter.AllowFor(path, timeout) {
      logger := logs_.GetLogger(r.Context())
      logger.WithField("path", path).
        WithField("method", r.Method).
        WithField("timeout", timeout).
        Warn("request rejected by QPS rate limiter after waiting")

      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("Retry-After", "1")
      w.WriteHeader(http.StatusTooManyRequests)

      resp := map[string]interface{}{
        "error":   "rate_limit_exceeded",
        "message": fmt.Sprintf("%s %s is rejected after waiting %v, QPS limit exceeded", r.Method, path, timeout),
        "code":    http.StatusTooManyRequests,
      }
      json.NewEncoder(w).Encode(resp)
      return
    }

    next.ServeHTTP(w, r)
  })
}

// Stats returns the QPS statistics for all configured paths.
func (l *QPSRateLimiter) Stats() []rate_.MethodQPSStats {
  return l.limiter.Stats()
}

// StatsHandler returns an HTTP handler that exposes rate limit stats.
func (l *QPSRateLimiter) StatsHandler() http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    stats := l.limiter.Stats()
    json.NewEncoder(w).Encode(map[string]interface{}{
      "rate_limits": stats,
    })
  })
}

// LimitAllQPS creates a simple QPS limiter that applies the same limit to all paths.
// This is a convenience function for simple use cases.
func LimitAllQPS(qps float64, burst int) *QPSRateLimiter {
  return NewQPSRateLimiter(qps, burst)
}

// ConcurrencyLimiter provides concurrency-based rate limiting for HTTP handlers.
// Unlike QPS limiting, concurrency control limits the number of requests being processed simultaneously.
// Tokens are returned when request processing completes.
type ConcurrencyLimiter struct {
  limiter *rate_.MethodLimiter
}

// NewConcurrencyLimiter creates a new concurrency-based rate limiter.
// defaultConcurrency: default max concurrency for paths without specific config, 0 means unlimited
func NewConcurrencyLimiter(defaultConcurrency int) *ConcurrencyLimiter {
  return &ConcurrencyLimiter{
    limiter: rate_.NewMethodLimiter(defaultConcurrency),
  }
}

// SetPathConcurrency sets the max concurrency for a specific path.
func (l *ConcurrencyLimiter) SetPathConcurrency(path string, maxConcurrency int) error {
  return l.limiter.AddLimiter(path, rate_.NewLimiter(maxConcurrency))
}

// Handler returns an HTTP middleware that applies concurrency rate limiting.
// The rate limit is based on the request URL path.
func (l *ConcurrencyLimiter) Handler(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path

    if !l.limiter.Allow(path) {
      logger := logs_.GetLogger(r.Context())
      logger.WithField("path", path).
        WithField("method", r.Method).
        Warn("request rejected by concurrency limiter")

      w.Header().Set("Content-Type", "application/json")
      w.Header().Set("Retry-After", "1")
      w.WriteHeader(http.StatusTooManyRequests)

      resp := map[string]interface{}{
        "error":   "concurrency_limit_exceeded",
        "message": fmt.Sprintf("%s %s is rejected by http_concurrency middleware, max concurrency exceeded", r.Method, path),
        "code":    http.StatusTooManyRequests,
      }
      json.NewEncoder(w).Encode(resp)
      return
    }
    defer l.limiter.Put(path)

    next.ServeHTTP(w, r)
  })
}

// HandlerFunc returns an HTTP middleware function.
func (l *ConcurrencyLimiter) HandlerFunc(next http.HandlerFunc) http.HandlerFunc {
  return l.Handler(next).ServeHTTP
}
