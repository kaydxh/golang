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

// ServeHTTP, wrap g.gateMux httpServerOpts, and called by grpcHandlerFunc
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
