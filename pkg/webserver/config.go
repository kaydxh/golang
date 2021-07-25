package webserver

import (
	"time"

	"github.com/gin-gonic/gin"
	viper_ "github.com/kaydxh/golang/pkg/viper"

	"github.com/ory/viper"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway/v2/grpc"
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
		getViper    func() *viper.Viper
		bindAddress string
		// ExternalAddress is the address (hostname or IP and port) that should be used in
		// external (public internet) URLs for this GenericWebServer.
		externalAddress string
		// ShutdownDelayDuration allows to block shutdown for some time, e.g. until endpoints pointing to this API server
		// have converged on all node. During this time, the API server keeps serving, /healthz will return 200,
		// but /readyz will return failure.
		shutdownDelayDuration time.Duration
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
	grpcBackend := grpc.NewGateway(c.opts.bindAddress)
	ginBackend := gin.New()

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
	if c.opts.getViper != nil {
		return viper_.UnmarshalProtoMessageWithJsonPb(c.opts.getViper(), &c.Proto)
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
