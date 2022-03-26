package webserver

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	viper_ "github.com/kaydxh/golang/pkg/viper"

	gw_ "github.com/kaydxh/golang/pkg/grpc-gateway"
	"github.com/ory/viper"
)

const (
	defaultBindAddress             = ":80"
	defaultExternalAddress         = ":80"
	defaultShutdownDelayDuration   = time.Duration(0)
	defaultShutdownTimeoutDuration = 5 * time.Second
)

type Config struct {
	Proto     Web
	Validator *validator.Validate
	opts      struct {
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
		//shutdownTimeoutDuration force shutdonw server after some time
		shutdownTimeoutDuration time.Duration
		gatewayOptions          []gw_.GRPCGatewayOption
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

// Validate checks Config.
func (c *completedConfig) Validate() error {
	return c.Validator.Struct(c)
}

func (c *completedConfig) New() (*GenericWebServer, error) {
	if c.completeError != nil {
		return nil, c.completeError
	}
	err := c.Validate()
	if err != nil {
		return nil, err
	}

	return c.install()
}

func (c *completedConfig) install() (*GenericWebServer, error) {
	opts := c.installMiddleware()
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

	if c.Validator == nil {
		c.Validator = validator.New()
	}

	return CompletedConfig{&completedConfig{Config: c}}
}

func (c *Config) installMiddleware() []gw_.GRPCGatewayOption {
	// recovery
	c.opts.gatewayOptions = append(
		c.opts.gatewayOptions,
		gw_.WithServerInterceptorsRecoveryOptions(),
	)

	//auto generate requestId
	c.opts.gatewayOptions = append(
		c.opts.gatewayOptions,
		gw_.WithServerUnaryInterceptorsRequestIdOptions(),
	)

	// trace id
	c.opts.gatewayOptions = append(
		c.opts.gatewayOptions,
		gw_.WithHttpHandlerInterceptorTraceIDOptions(),
	)

	// limit rate
	c.opts.gatewayOptions = append(
		c.opts.gatewayOptions,
		gw_.WithServerInterceptorsLimitRateOptions(
			int(c.Proto.GetMaxConcurrencyUnary()),
			int(c.Proto.GetMaxConcurrencyStream()),
		),
	)

	// time cost
	c.opts.gatewayOptions = append(
		c.opts.gatewayOptions,
		gw_.WithServerUnaryInterceptorsTimerOptions(c.Proto.GetMonitor().GetPrometheus().GetEnabledMetricTimerCost()),
	)

	// code,message and client ip
	c.opts.gatewayOptions = append(
		c.opts.gatewayOptions,
		gw_.WithServerUnaryInterceptorsCodeMessageOptions(
			c.Proto.GetMonitor().GetPrometheus().GetEnabledMetricCodeMessage(),
		),
	)

	// log input and output
	c.opts.gatewayOptions = append(
		c.opts.gatewayOptions,
		gw_.WithServerUnaryInterceptorsInOutPacketOptions(),
	)

	return c.opts.gatewayOptions
}

func (c *Config) WithWebConfigOptions(opts ...ConfigOption) {
	c.ApplyOptions(opts...)
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
