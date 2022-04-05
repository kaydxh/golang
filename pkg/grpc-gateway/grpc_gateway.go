package grpcgateway

import (
	"context"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpc_ "github.com/kaydxh/golang/go/net/grpc"
	http_ "github.com/kaydxh/golang/go/net/http"

	"google.golang.org/grpc"
)

type InterceptorOption struct {
	grpcServerOpts struct {
		unaryInterceptors  []grpc.UnaryServerInterceptor
		streamInterceptors []grpc.StreamServerInterceptor
	}
	httpServerOpts struct {
		handlerChain http_.HandlerChain
		/*
			//invoke before http handler
			PreHttpInterceptors []http_.HandlerInterceptor
			//invoke after http handler
			PostHttpInterceptors []http_.HandlerInterceptor
		*/
	}
}

type GRPCGateway struct {
	grpcServer *grpc.Server
	http.Server
	//assigned by ginRouter in PrepareRun
	Handler http.Handler
	//gatewayMux for http handler
	gatewayMux *runtime.ServeMux
	once       sync.Once

	opts struct {
		interceptionOptions InterceptorOption
		serverOptions       []grpc.ServerOption
		gatewayMuxOptions   []runtime.ServeMuxOption
		clientDialOptions   []grpc.DialOption
	}
}

func NewGRPCGateWay(addr string, options ...GRPCGatewayOption) *GRPCGateway {
	server := &GRPCGateway{
		Server: http.Server{
			Addr: addr,
		},
	}
	server.ApplyOptions(options...)

	return server
}

func (g *GRPCGateway) initOnce() {
	g.once.Do(func() {
		//now not support tls
		g.opts.clientDialOptions = append(g.opts.clientDialOptions, grpc_.ClientDialOptions()...)

		serverOptions := []grpc.ServerOption{}
		serverOptions = append(
			g.opts.serverOptions,
			grpc.ChainUnaryInterceptor(g.opts.interceptionOptions.grpcServerOpts.unaryInterceptors...),
			grpc.ChainStreamInterceptor(g.opts.interceptionOptions.grpcServerOpts.streamInterceptors...),
		)

		g.opts.gatewayMuxOptions = append(g.opts.gatewayMuxOptions,
			runtime.WithRoutingErrorHandler(
				func(ctx context.Context, mux *runtime.ServeMux,
					marshaler runtime.Marshaler,
					w http.ResponseWriter, r *http.Request, code int) {

					//g.Handler is gin handler
					httpHandler := g.Handler
					if httpHandler == nil {
						httpHandler = http.DefaultServeMux
					}

					// NotFound and NotAllowed, use gin handler
					if code == http.StatusNotFound || code == http.StatusMethodNotAllowed {
						httpHandler.ServeHTTP(w, r)
						return
					}
					runtime.DefaultRoutingErrorHandler(ctx, mux, marshaler, w, r, code)
				}))

		g.grpcServer = grpc.NewServer(serverOptions...)
		g.gatewayMux = runtime.NewServeMux(g.opts.gatewayMuxOptions...)
		g.Server.Handler = grpcHandlerFunc(g.grpcServer, g)
	})
}

//called by grpcHandlerFunc
func (g *GRPCGateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.opts.interceptionOptions.httpServerOpts.handlerChain.WrapH(g.gatewayMux).ServeHTTP(w, r)
}

func (g *GRPCGateway) ListenAndServe() error {
	g.initOnce()
	return g.Server.ListenAndServe()
}

func (g *GRPCGateway) registerGRPCFunc(h GRPCHandler) {
	g.initOnce()
	h.Register(g.grpcServer)
}

func (g *GRPCGateway) RegisterGRPCHandler(h func(srv *grpc.Server)) {
	g.initOnce()
	g.registerGRPCFunc(GRPCHandlerFunc(h))
}

func (g *GRPCGateway) registerHTTPFunc(ctx context.Context, h HTTPHandler) error {
	g.initOnce()
	return h.Register(ctx, g.gatewayMux, g.Server.Addr, g.opts.clientDialOptions)
}

func (g *GRPCGateway) RegisterHTTPHandler(ctx context.Context,
	h func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error,
) error {
	g.initOnce()
	return g.registerHTTPFunc(ctx, HTTPHandlerFunc(h))
}
