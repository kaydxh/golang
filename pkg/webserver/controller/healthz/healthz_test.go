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
package healthz_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	healthz_ "github.com/kaydxh/golang/pkg/webserver/controller/healthz"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestPingHealthChecker(t *testing.T) {
	checker := healthz_.PingHealthChecker{}

	if checker.Name() != "ping" {
		t.Errorf("expected name 'ping', got '%s'", checker.Name())
	}

	if err := checker.Check(context.Background()); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestFuncHealthChecker(t *testing.T) {
	t.Run("healthy check", func(t *testing.T) {
		checker := healthz_.NewFuncHealthChecker("test", func(ctx context.Context) error {
			return nil
		})

		if checker.Name() != "test" {
			t.Errorf("expected name 'test', got '%s'", checker.Name())
		}

		if err := checker.Check(context.Background()); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("unhealthy check", func(t *testing.T) {
		expectedErr := errors.New("service unavailable")
		checker := healthz_.NewFuncHealthChecker("test", func(ctx context.Context) error {
			return expectedErr
		})

		if err := checker.Check(context.Background()); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("nil check func", func(t *testing.T) {
		checker := healthz_.NewFuncHealthChecker("test", nil)

		if err := checker.Check(context.Background()); err != nil {
			t.Errorf("expected no error for nil func, got: %v", err)
		}
	})
}

func TestCompositeHealthChecker(t *testing.T) {
	t.Run("all healthy", func(t *testing.T) {
		composite := healthz_.NewCompositeHealthChecker(
			healthz_.PingHealthChecker{},
			healthz_.NewFuncHealthChecker("check1", func(ctx context.Context) error {
				return nil
			}),
		)

		if err := composite.Check(context.Background()); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		results, healthy := composite.CheckAll(context.Background())
		if !healthy {
			t.Error("expected all healthy")
		}
		if len(results) != 2 {
			t.Errorf("expected 2 results, got %d", len(results))
		}
	})

	t.Run("one unhealthy", func(t *testing.T) {
		composite := healthz_.NewCompositeHealthChecker(
			healthz_.PingHealthChecker{},
			healthz_.NewFuncHealthChecker("failing", func(ctx context.Context) error {
				return errors.New("failed")
			}),
		)

		if err := composite.Check(context.Background()); err == nil {
			t.Error("expected error, got nil")
		}

		results, healthy := composite.CheckAll(context.Background())
		if healthy {
			t.Error("expected not all healthy")
		}
		if len(results) != 2 {
			t.Errorf("expected 2 results, got %d", len(results))
		}
	})

	t.Run("add and remove checker", func(t *testing.T) {
		composite := healthz_.NewCompositeHealthChecker()

		checker := healthz_.NewFuncHealthChecker("dynamic", func(ctx context.Context) error {
			return nil
		})

		composite.AddChecker(checker)
		if len(composite.Checkers()) != 1 {
			t.Errorf("expected 1 checker, got %d", len(composite.Checkers()))
		}

		composite.RemoveChecker("dynamic")
		if len(composite.Checkers()) != 0 {
			t.Errorf("expected 0 checkers, got %d", len(composite.Checkers()))
		}
	})
}

func TestController_Healthz(t *testing.T) {
	controller := healthz_.NewController()

	router := gin.New()
	controller.SetRoutes(router, nil)

	t.Run("healthz returns ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
		if w.Body.String() != "ok" {
			t.Errorf("expected body 'ok', got '%s'", w.Body.String())
		}
	})

	t.Run("livez returns ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/livez", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("readyz returns ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestController_NotReady(t *testing.T) {
	controller := healthz_.NewController()
	controller.SetReady(false)

	router := gin.New()
	controller.SetRoutes(router, nil)

	t.Run("readyz returns service unavailable when not ready", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
		}
	})

	t.Run("healthz returns service unavailable when not ready", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
		}
	})

	t.Run("livez still returns ok when not ready", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/livez", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestController_WithFailingChecker(t *testing.T) {
	failingChecker := healthz_.NewFuncHealthChecker("failing", func(ctx context.Context) error {
		return errors.New("service down")
	})

	controller := healthz_.NewController(
		healthz_.WithLivezCheckers(failingChecker),
	)

	router := gin.New()
	controller.SetRoutes(router, nil)

	t.Run("livez returns service unavailable with failing checker", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/livez", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
		}
	})
}

func TestController_VerboseEndpoints(t *testing.T) {
	controller := healthz_.NewController()

	router := gin.New()
	controller.SetRoutes(router, nil)

	endpoints := []string{"/healthz/verbose", "/livez/verbose", "/readyz/verbose"}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, endpoint, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json; charset=utf-8" {
				t.Errorf("expected JSON content type, got '%s'", contentType)
			}
		})
	}
}

func TestController_WithOptions(t *testing.T) {
	controller := healthz_.NewController(
		healthz_.WithCheckTimeout(5*time.Second),
		healthz_.WithLivezCheckers(healthz_.PingHealthChecker{}),
		healthz_.WithReadyzCheckers(healthz_.PingHealthChecker{}),
	)

	if !controller.IsReady() {
		t.Error("expected controller to be ready by default")
	}
}

func TestController_AddRemoveCheckers(t *testing.T) {
	controller := healthz_.NewController()

	checker := healthz_.NewFuncHealthChecker("dynamic", func(ctx context.Context) error {
		return nil
	})

	controller.AddLivezChecker(checker)
	controller.AddReadyzChecker(checker)

	controller.RemoveLivezChecker("dynamic")
	controller.RemoveReadyzChecker("dynamic")
}
