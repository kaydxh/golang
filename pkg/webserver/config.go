package webserver

import (
	"fmt"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"google.golang.org/grpc"

	gw_ "github.com/kaydxh/golang/pkg/grpc-gateway"
	"github.com/spf13/viper"
)

const (
	defaultBindAddress             = ":80"
	defaultExternalAddress         = ":80"
	defaultShutdownDelayDuration   = time.Duration(0)
	defaultShutdownTimeoutDuration = 5 * time.Second

	defaultMaxReceiveMessageSize = math.MaxInt32
	defaultMaxSendMessageSize    = math.MaxInt32
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
	c.Config.opts.gatewayOptions = append(c.Config.opts.gatewayOptions, c.installGrpcMessageSizeOptions()...)
	c.Config.opts.gatewayOptions = append(c.Config.opts.gatewayOptions, c.installHttpMiddlewareChain()...)
	c.Config.opts.gatewayOptions = append(c.Config.opts.gatewayOptions, c.installGrpcMiddlewareChain()...)
	grpcBackend := gw_.NewGRPCGateWay(c.opts.bindAddress, c.Config.opts.gatewayOptions...)
	//grpcBackend.ApplyOptions()
	gin.SetMode(c.Proto.Mode.String())
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

func (c *Config) installGrpcMessageSizeOptions() []gw_.GRPCGatewayOption {
	maxRecvMsgSize := defaultMaxReceiveMessageSize
	maxSendMsgSize := defaultMaxSendMessageSize

	var opts []gw_.GRPCGatewayOption
	// request
	// http -> grpc client -> grpc server
	// ---------------------------------  -> gin server
	if c.Proto.GetGrpc().GetMaxReceiveMessageSize() > 0 {

		maxRecvMsgSize = int(c.Proto.GetGrpc().GetMaxReceiveMessageSize())
	}
	opts = append(
		opts,
		gw_.WithServerOptions(grpc.MaxRecvMsgSize(maxRecvMsgSize)),
	)
	opts = append(
		opts,
		gw_.WithClientDialOptions(grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(maxRecvMsgSize))),
	)

	// response
	// http <- grpc client <- grpc server
	if c.Proto.GetGrpc().GetMaxSendMessageSize() > 0 {
		maxSendMsgSize = int(c.Proto.GetGrpc().GetMaxSendMessageSize())
	}
	opts = append(
		opts,
		gw_.WithServerOptions(grpc.MaxSendMsgSize(maxSendMsgSize)),
	)
	opts = append(
		opts,
		gw_.WithClientDialOptions(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxSendMsgSize))),
	)

	return opts
}

func (c *Config) installHttpMiddlewareChain() []gw_.GRPCGatewayOption {

	var opts []gw_.GRPCGatewayOption

	// recoverer
	opts = append(
		opts,
		gw_.WithHttpHandlerInterceptorRecoveryOptions(),
	)

	// trace id
	opts = append(
		opts,
		gw_.WithHttpHandlerInterceptorTraceIDOptions(),
	)

	// time cost
	opts = append(
		opts,
		gw_.WithHttpHandlerInterceptorsTimerOptions(c.Proto.GetMonitor().GetPrometheus().GetEnabledMetricTimerCost()),
	)

	//inout
	opts = append(
		opts,
		gw_.WithHttpHandlerInterceptorInOutPacketOptions(),
	)

	return opts
}

func (c *Config) installGrpcMiddlewareChain() []gw_.GRPCGatewayOption {
	var opts []gw_.GRPCGatewayOption

	// recovery
	opts = append(
		opts,
		gw_.WithServerInterceptorsRecoveryOptions(),
	)

	//auto generate requestId
	opts = append(
		opts,
		gw_.WithServerUnaryInterceptorsRequestIdOptions(),
	)

	// limit rate
	opts = append(
		opts,
		gw_.WithServerInterceptorsLimitRateOptions(
			int(c.Proto.GetGrpc().GetMaxConcurrencyUnary()),
			int(c.Proto.GetGrpc().GetMaxConcurrencyStream()),
		),
	)

	/*
		// time cost
		opts = append(
			opts,
			gw_.WithServerUnaryInterceptorsTimerOptions(c.Proto.GetMonitor().GetPrometheus().GetEnabledMetricTimerCost()),
		)
	*/

	// code,message and client ip
	opts = append(
		opts,
		gw_.WithServerUnaryInterceptorsCodeMessageOptions(
			c.Proto.GetMonitor().GetPrometheus().GetEnabledMetricCodeMessage(),
		),
	)

	// log input and output
	opts = append(
		opts,
		gw_.WithServerUnaryInterceptorsInOutPacketOptions(),
	)

	return opts
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
