package grpcgateway

import (
	"context"
	"net"
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

func (g *GRPCGateway) ListenAndServe() error {
	g.initOnce()
	//	g.Server.ListenAndServe()
	/*
		if g.Server.shuttingDown() {
			return http.ErrServerClosed
		}
	*/
	addr := g.Server.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", g.gatewayMux)
	srv := &http.Server{
		Addr:    addr,
		Handler: grpcHandlerFunc(g.grpcServer, mux),
	}
	return srv.Serve(ln)

}

func (g *GRPCGateway) registerGRPCFunc(h GRPCHandler) {
	h.Register(g.grpcServer)
}

func (g *GRPCGateway) RegisterGRPCHandler(h func(srv *grpc.Server)) {
	g.initOnce()
	g.registerGRPCFunc(GRPCHandlerFunc(h))
}

func (g *GRPCGateway) registerHTTPFunc(ctx context.Context, h HTTPHandler) error {
	return h.Register(ctx, g.gatewayMux, g.Server.Addr, g.opts.clientDialOptions)
}

func (g *GRPCGateway) RegisterHTTPHandler(ctx context.Context,
	h func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error,
) error {
	g.initOnce()
	return g.registerHTTPFunc(ctx, HTTPHandlerFunc(h))
}
