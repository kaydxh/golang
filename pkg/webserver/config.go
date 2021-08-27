package webserver

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/sirupsen/logrus"

	gw_ "github.com/kaydxh/golang/pkg/grpc-gateway"
	"github.com/ory/viper"
)

const (
	defaultBindAddress           = ":80"
	defaultExternalAddress       = ":80"
	defaultShutdownDelayDuration = time.Duration(0)
)

type Config struct {
	Proto Web
	opts  struct {
		// If set, overrides params below
		viper       *viper.Viper
		bindAddress string
		// ExternalAddress is the address (hostname or IP and port) that should be used in
		// external (public internet) URLs for this GenericWebServer.
		externalAddress string
		// ShutdownDelayDuration allows to block shutdown for some time, e.g. until endpoints pointing to this API server
		// have converged on all node. During this time, the API server keeps serving, /healthz will return 200,
		// but /readyz will return failure.
		shutdownDelayDuration time.Duration
		gatewayOptions        []gw_.GRPCGatewayOption
	}
}

type completedConfig struct {
	*Config
	completeError error
}

type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

func (c *completedConfig) New() (*GenericWebServer, error) {
	opts := c.opts.gatewayOptions //[]gw_.GRPCGatewayOption{}
	opts = append(opts, gw_.WithServerUnaryInterceptorsTimerOptions())
	opts = append(
		opts,
		gw_.WithServerInterceptorsLogrusOptions(logrus.StandardLogger()),
	)
	opts = append(
		opts,
		gw_.WithServerUnaryInterceptorsRequestIdOptions(),
	)
	grpcBackend := gw_.NewGRPCGateWay(c.opts.bindAddress, opts...)
	//grpcBackend.ApplyOptions()
	ginBackend := gin.New()
	fmt.Printf(" - listen address[%s]\n", c.opts.bindAddress)

	return &GenericWebServer{
		ginBackend:       ginBackend,
		grpcBackend:      grpcBackend,
		postStartHooks:   map[string]postStartHookEntry{},
		preShutdownHooks: map[string]preShutdownHookEntry{},
		readinessStopCh:  make(chan struct{}),
	}, nil
}

// Complete set default ServerRunOptions.
func (c *Config) Complete() CompletedConfig {
	err := c.loadViper()
	if err != nil {
		return CompletedConfig{&completedConfig{
			Config:        c,
			completeError: err,
		}}
	}
	c.parseViper()
	return CompletedConfig{&completedConfig{Config: c}}
}

func (c *Config) parseViper() {
	c.opts.bindAddress = c.Proto.GetBindHostPort()
}

func (c *Config) loadViper() error {
	if c.opts.viper != nil {
		return viper_.UnmarshalProtoMessageWithJsonPb(c.opts.viper, &c.Proto)
	}

	return nil
}

// default bind port 80
func NewConfig(options ...ConfigOption) *Config {
	c := &Config{}
	c.ApplyOptions(options...)

	if c.opts.bindAddress == "" {
		c.opts.bindAddress = defaultBindAddress
	}

	if c.opts.externalAddress == "" {
		c.opts.externalAddress = defaultExternalAddress
	}

	return c
}
