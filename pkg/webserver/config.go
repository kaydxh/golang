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
package webserver

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"google.golang.org/grpc"

	errors_ "github.com/kaydxh/golang/go/errors"
	gw_ "github.com/kaydxh/golang/pkg/grpc-gateway"
	healthz_ "github.com/kaydxh/golang/pkg/webserver/controller/healthz"
	profiler_ "github.com/kaydxh/golang/pkg/webserver/controller/profiler"
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
	gin.SetMode(gin.ReleaseMode)
	ginBackend := gin.New()
	fmt.Printf(" - listen address[%s]\n", c.opts.bindAddress)

	ws := &GenericWebServer{
		ginBackend:       ginBackend,
		grpcBackend:      grpcBackend,
		postStartHooks:   map[string]postStartHookEntry{},
		preShutdownHooks: map[string]preShutdownHookEntry{},
		readinessStopCh:  make(chan struct{}),
	}

	var errs []error
	if c.Proto.GetDebug().GetEnableProfiling() {
		fmt.Printf("- install debug handler")
		err := ws.AddPostStartHook("debug_hanlder", func(ctx context.Context) error {
			ws.InstallWebHandlers(profiler_.NewController())
			return nil
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	err := c.installDefaultHander(ws)
	if err != nil {
		errs = append(errs, err)
	}

	return ws, errors_.NewAggregate(errs)
}

func (c *completedConfig) installDefaultHander(ws *GenericWebServer) error {
	return ws.AddPostStartHook("default_hanlder", func(ctx context.Context) error {
		ws.InstallWebHandlers(healthz_.NewController())
		return nil
	})
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

// use LocalMiddlewareWrap instead of request id, metric, cost time, print request and response middleware
func (c *Config) installHttpMiddlewareChain() []gw_.GRPCGatewayOption {

	httpConfig := c.Proto.GetHttp()
	var opts []gw_.GRPCGatewayOption
	opts = append(
		opts,
		// request id
		//	gw_.WithHttpHandlerInterceptorRequestIDOptions(),

		// http recoverer
		gw_.WithHttpHandlerInterceptorRecoveryOptions(),

		// clean path
		gw_.WithHttpHandlerInterceptorCleanPathOptions(),

		// http body proto
		gw_.WithServerInterceptorsHttpBodyProtoOptions(),

		//gw_.WithHttpHandlerInterceptorsMetricOptions(),

		// limit rate
		gw_.WithHttpHandlerInterceptorsLimitAllOptions(
			int(httpConfig.GetMaxConcurrency()),
		),
	)

	//options
	switch httpConfig.GetApiFormatter() {
	case Web_Http_tcloud_api_v30:
		opts = append(opts,
			gw_.WithServerInterceptorsTCloud30HTTPResponseOptions(),
			gw_.WithServerInterceptorsTCloud30HttpErrorOptions(),
		)
	case Web_Http_trivial_api_v10:
		opts = append(opts,
			gw_.WithServerInterceptorsTrivialV1HTTPResponseOptions(),
			gw_.WithServerInterceptorsTrivialV1HttpErrorOptions(),
		)
	case Web_Http_trivial_api_v20:
		opts = append(opts,
			gw_.WithServerInterceptorsTrivialV2HTTPResponseOptions(),
			gw_.WithServerInterceptorsTrivialV2HttpErrorOptions(),
		)
	default:
		opts = append(opts,
			gw_.WithServerInterceptorsNoopHttpErrorOptions(),
		)
	}

	/*
		enableInoutPrinter := httpConfig.GetEnableInoutputPrinter()
		opts = append(opts,
			// inout header printer
			gw_.WithHttpHandlerInterceptorInOutputHeaderPrinterOptions(enableInoutPrinter),
			// print inoutput body
			gw_.WithHttpHandlerInterceptorInOutputPrinterOptions(enableInoutPrinter),
		)
	*/

	return opts
}

func (c *Config) installGrpcMiddlewareChain() []gw_.GRPCGatewayOption {

	grpcConfig := c.Proto.GetGrpc()
	var opts []gw_.GRPCGatewayOption
	opts = append(
		opts,

		// requestId
		//gw_.WithServerUnaryInterceptorsRequestIdOptions(),

		// recovery
		gw_.WithServerInterceptorsRecoveryOptions(),

		// limit rate
		gw_.WithServerInterceptorsLimitRateOptions(
			int(grpcConfig.GetMaxConcurrencyUnary()),
			int(grpcConfig.GetMaxConcurrencyStream()),
		),

		// total req, fail req, cost time metrics, errorcode ip dims
		gw_.WithServerUnaryMetricInterceptorOptions(),

		// print input and output body
		//gw_.WithServerUnaryInterceptorsInOutPacketOptions(),
		//gw_.WithServerInterceptorTimeoutOptions(grpcConfig.GetTimeout().AsDuration()),
	)

	return opts
}

func (c *Config) WithWebConfigOptions(opts ...ConfigOption) {
	c.ApplyOptions(opts...)
}

func (c *Config) AppendGRPCGatewayOptions(opts ...gw_.GRPCGatewayOption) {
	c.opts.gatewayOptions = append(c.opts.gatewayOptions, opts...)
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
