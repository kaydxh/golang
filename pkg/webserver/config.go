package webserver

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway/v2/grpc"
)

type Config struct {
	InternalAddress string
	// ExternalAddress is the address (hostname or IP and port) that should be used in
	// external (public internet) URLs for this GenericAPIServer.
	ExternalAddress string
	// ShutdownDelayDuration allows to block shutdown for some time, e.g. until endpoints pointing to this API server
	// have converged on all node. During this time, the API server keeps serving, /healthz will return 200,
	// but /readyz will return failure.
	ShutdownDelayDuration time.Duration

	// The limit on the request body size that would be accepted and decoded in a write request.
	// 0 means no limit.
	maxRequestBodyBytes int64
}

type completedConfig struct {
	*Config
}

type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// Complete set default ServerRunOptions.
func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{&completedConfig{c}}
}

func (c *completedConfig) New() (*GenericWebServer, error) {
	grpcBackend := grpc.NewGateway(c.InternalAddress)
	ginBackend := gin.New()

	return &GenericWebServer{
		ginBackend:      ginBackend,
		grpcBackend:     grpcBackend,
		readinessStopCh: make(chan struct{}),
	}, nil
}

func NewConfig() *Config {
	return &Config{
		ShutdownDelayDuration: time.Duration(0),
	}
}
