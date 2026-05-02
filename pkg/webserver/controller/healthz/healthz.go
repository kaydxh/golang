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
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	gw_ "github.com/kaydxh/golang/pkg/grpc-gateway"
)

const (
	defaultCheckTimeout = 10 * time.Second
)

// Controller handles health check endpoints.
// Supports Kubernetes-style /healthz, /livez, /readyz endpoints.
type Controller struct {
	// livezCheckers are used for liveness checks.
	// Liveness probes determine if the application is running.
	livezCheckers *CompositeHealthChecker

	// readyzCheckers are used for readiness checks.
	// Readiness probes determine if the application is ready to receive traffic.
	readyzCheckers *CompositeHealthChecker

	// ready indicates if the server is ready to receive traffic.
	// This is set to false during shutdown.
	ready atomic.Bool

	// checkTimeout is the timeout for health checks.
	checkTimeout time.Duration

	// disableRootRoute 是否禁用根路径 "/" 的健康检查路由注册。
	// 当项目使用 SPA 前端（static 控制器的 SPAMode=true）时，应设为 true，
	// 避免与 static 控制器的 GET "/" 路由冲突。
	// 禁用后，负载均衡器应使用 /healthz 或 /livez 进行健康探测。
	disableRootRoute bool
}

// ControllerOption is a functional option for Controller.
type ControllerOption func(*Controller)

// WithCheckTimeout sets the timeout for health checks.
func WithCheckTimeout(timeout time.Duration) ControllerOption {
	return func(c *Controller) {
		c.checkTimeout = timeout
	}
}

// WithLivezCheckers sets the liveness checkers.
func WithLivezCheckers(checkers ...HealthChecker) ControllerOption {
	return func(c *Controller) {
		for _, checker := range checkers {
			c.livezCheckers.AddChecker(checker)
		}
	}
}

// WithReadyzCheckers sets the readiness checkers.
func WithReadyzCheckers(checkers ...HealthChecker) ControllerOption {
	return func(c *Controller) {
		for _, checker := range checkers {
			c.readyzCheckers.AddChecker(checker)
		}
	}
}

// WithDisableRootRoute 禁用根路径 "/" 的健康检查路由。
// 当项目使用 SPA 前端时，根路径应由 static 控制器处理（返回 index.html），
// 而非返回健康检查结果。负载均衡器应改用 /healthz 或 /livez 进行探测。
func WithDisableRootRoute() ControllerOption {
	return func(c *Controller) {
		c.disableRootRoute = true
	}
}

// NewController creates a new health check controller.
func NewController(opts ...ControllerOption) *Controller {
	c := &Controller{
		livezCheckers:  NewCompositeHealthChecker(PingHealthChecker{}),
		readyzCheckers: NewCompositeHealthChecker(PingHealthChecker{}),
		checkTimeout:   defaultCheckTimeout,
	}
	c.ready.Store(true)

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// SetRoutes registers health check endpoints.
func (c *Controller) SetRoutes(ginRouter gin.IRouter, grpcRouter *gw_.GRPCGateway) {
	// / - root path health check for load balancer probes (supports GET and HEAD)
	// 当 disableRootRoute=true 时跳过，避免与 SPA 前端的 static 控制器冲突。
	if !c.disableRootRoute {
		ginRouter.GET("/", c.RootHealthz())
		ginRouter.HEAD("/", c.RootHealthz())
	}

	// /healthz - general health check (combines livez and readyz)
	ginRouter.GET("/healthz", c.Healthz())

	// /livez - liveness probe
	// Returns 200 if the application is alive
	ginRouter.GET("/livez", c.Livez())

	// /readyz - readiness probe
	// Returns 200 if the application is ready to receive traffic
	ginRouter.GET("/readyz", c.Readyz())

	// /healthz/verbose - detailed health check with all checker results
	ginRouter.GET("/healthz/verbose", c.HealthzVerbose())

	// /livez/verbose - detailed liveness check
	ginRouter.GET("/livez/verbose", c.LivezVerbose())

	// /readyz/verbose - detailed readiness check
	ginRouter.GET("/readyz/verbose", c.ReadyzVerbose())
}

// RootHealthz returns a handler for the root path "/" endpoint.
// This is a lightweight health check designed for load balancer probes
// that send HEAD or GET requests to the root path.
func (c *Controller) RootHealthz() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !c.ready.Load() {
			ctx.Status(http.StatusServiceUnavailable)
			return
		}
		ctx.String(http.StatusOK, "ok")
	}
}

// Healthz returns a handler for the /healthz endpoint.
// It combines liveness and readiness checks.
func (c *Controller) Healthz() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkCtx, cancel := context.WithTimeout(ctx.Request.Context(), c.checkTimeout)
		defer cancel()

		// Check liveness
		if err := c.livezCheckers.Check(checkCtx); err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"type":   "livez",
				"error":  err.Error(),
			})
			return
		}

		// Check readiness
		if !c.ready.Load() {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"type":   "readyz",
				"error":  "server is shutting down",
			})
			return
		}

		if err := c.readyzCheckers.Check(checkCtx); err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"type":   "readyz",
				"error":  err.Error(),
			})
			return
		}

		ctx.String(http.StatusOK, "ok")
	}
}

// Livez returns a handler for the /livez endpoint.
func (c *Controller) Livez() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkCtx, cancel := context.WithTimeout(ctx.Request.Context(), c.checkTimeout)
		defer cancel()

		if err := c.livezCheckers.Check(checkCtx); err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  err.Error(),
			})
			return
		}

		ctx.String(http.StatusOK, "ok")
	}
}

// Readyz returns a handler for the /readyz endpoint.
func (c *Controller) Readyz() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check if server is marked as not ready (e.g., during shutdown)
		if !c.ready.Load() {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  "server is shutting down",
			})
			return
		}

		checkCtx, cancel := context.WithTimeout(ctx.Request.Context(), c.checkTimeout)
		defer cancel()

		if err := c.readyzCheckers.Check(checkCtx); err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  err.Error(),
			})
			return
		}

		ctx.String(http.StatusOK, "ok")
	}
}

// HealthzVerbose returns a handler for detailed health check information.
func (c *Controller) HealthzVerbose() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkCtx, cancel := context.WithTimeout(ctx.Request.Context(), c.checkTimeout)
		defer cancel()

		response := HealthCheckResponse{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}

		// Check liveness
		livezResults, livezHealthy := c.livezCheckers.CheckAll(checkCtx)

		// Check readiness
		readyzResults, readyzHealthy := c.readyzCheckers.CheckAll(checkCtx)

		// Combine results
		response.Checks = append(response.Checks, livezResults...)
		response.Checks = append(response.Checks, readyzResults...)

		allHealthy := livezHealthy && readyzHealthy && c.ready.Load()
		if allHealthy {
			response.Status = "healthy"
			ctx.JSON(http.StatusOK, response)
		} else {
			response.Status = "unhealthy"
			ctx.JSON(http.StatusServiceUnavailable, response)
		}
	}
}

// LivezVerbose returns a handler for detailed liveness check information.
func (c *Controller) LivezVerbose() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkCtx, cancel := context.WithTimeout(ctx.Request.Context(), c.checkTimeout)
		defer cancel()

		results, allHealthy := c.livezCheckers.CheckAll(checkCtx)

		response := HealthCheckResponse{
			Checks:    results,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}

		if allHealthy {
			response.Status = "healthy"
			ctx.JSON(http.StatusOK, response)
		} else {
			response.Status = "unhealthy"
			ctx.JSON(http.StatusServiceUnavailable, response)
		}
	}
}

// ReadyzVerbose returns a handler for detailed readiness check information.
func (c *Controller) ReadyzVerbose() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response := HealthCheckResponse{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}

		// Check if server is marked as not ready
		if !c.ready.Load() {
			response.Status = "not ready"
			response.Checks = []HealthCheckResult{
				{
					Name:    "shutdown",
					Healthy: false,
					Error:   "server is shutting down",
				},
			}
			ctx.JSON(http.StatusServiceUnavailable, response)
			return
		}

		checkCtx, cancel := context.WithTimeout(ctx.Request.Context(), c.checkTimeout)
		defer cancel()

		results, allHealthy := c.readyzCheckers.CheckAll(checkCtx)
		response.Checks = results

		if allHealthy {
			response.Status = "ready"
			ctx.JSON(http.StatusOK, response)
		} else {
			response.Status = "not ready"
			ctx.JSON(http.StatusServiceUnavailable, response)
		}
	}
}

// SetReady sets the readiness state of the server.
func (c *Controller) SetReady(ready bool) {
	c.ready.Store(ready)
}

// IsReady returns the current readiness state.
func (c *Controller) IsReady() bool {
	return c.ready.Load()
}

// AddLivezChecker adds a liveness checker.
func (c *Controller) AddLivezChecker(checker HealthChecker) {
	c.livezCheckers.AddChecker(checker)
}

// AddReadyzChecker adds a readiness checker.
func (c *Controller) AddReadyzChecker(checker HealthChecker) {
	c.readyzCheckers.AddChecker(checker)
}

// RemoveLivezChecker removes a liveness checker by name.
func (c *Controller) RemoveLivezChecker(name string) {
	c.livezCheckers.RemoveChecker(name)
}

// RemoveReadyzChecker removes a readiness checker by name.
func (c *Controller) RemoveReadyzChecker(name string) {
	c.readyzCheckers.RemoveChecker(name)
}
