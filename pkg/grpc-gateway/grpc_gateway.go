package grpcgateway

import (
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpc_ "github.com/kaydxh/golang/go/net/grpc"
	"google.golang.org/grpc"
)

type GRPCGateway struct {
	grpcServer *grpc.Server
	http.Server
	handler    http.Handler
	gatewayMux *runtime.ServeMux
	once       sync.Once

	opts struct {
		/*
			ServerOptions struct {
				opts []grpc.ServerOption
			}
		*/
		serverOptions     []grpc.ServerOption
		gatewayMuxOptions []runtime.ServeMuxOption
		clientDialOptions []grpc.DialOption
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
		g.grpcServer = grpc.NewServer(g.opts.serverOptions...)
		g.gatewayMux = runtime.NewServeMux(g.opts.gatewayMuxOptions...)
		g.opts.clientDialOptions = append(g.opts.clientDialOptions, grpc_.ClientDialOptions()...)
	})
}

func (g *GRPCGateway) ListenAndServe() {
	g.initOnce()
	g.Server.ListenAndServe()
}

func (g *GRPCGateway) register(h GRPCHandler) {
	h.Register(g.grpcServer)
}

func (g *GRPCGateway) RegisterGrpcHandler(h func(srv *grpc.Server)) {
	g.initOnce()
	g.register(GRPCHandlerFunc(h))
}
