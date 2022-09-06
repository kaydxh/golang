package controller

import (
	"net/http"
	"runtime"

	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
	gw_ "github.com/kaydxh/golang/pkg/grpc-gateway"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) SetRoutes(ginRouter gin.IRouter, grpcRouter *gw_.GRPCGateway) {
	ginRouter.GET("/debug/pprof/*path", c.Profile())
}

func (c *Controller) Profile() gin.HandlerFunc {
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	return gin.WrapH(http.DefaultServeMux)
}
