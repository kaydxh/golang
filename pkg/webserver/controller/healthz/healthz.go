package healthz

import (
	"net/http"

	"github.com/gin-gonic/gin"
	gw_ "github.com/kaydxh/golang/pkg/grpc-gateway"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) SetRoutes(ginRouter gin.IRouter, grpcRouter *gw_.GRPCGateway) {
	ginRouter.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "Healthz Ok")
	})
}
