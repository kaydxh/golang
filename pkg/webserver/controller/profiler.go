package controller

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	gw_ "github.com/kaydxh/golang/pkg/grpc-gateway"
)

type Controller struct {
	EnableProfiling bool
}

func NewController(enableProfiling bool) *Controller {
	return &Controller{
		EnableProfiling: enableProfiling,
	}
}

func (c *Controller) SetRoutes(ginRouter gin.IRouter, grpcRouter *gw_.GRPCGateway) {
	if c.EnableProfiling {
		ginRouter.GET("/debug/pprof/*path", c.PProf())
	}
}

func (c *Controller) PProf() gin.HandlerFunc {
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	return gin.WrapH(http.DefaultServeMux)
}
