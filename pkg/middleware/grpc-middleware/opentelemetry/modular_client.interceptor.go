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
package interceptoropentelemetry

import (
	"context"
	"time"

	net_ "github.com/kaydxh/golang/go/net"
	"github.com/kaydxh/golang/pkg/opentelemetry/metric/report"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ModularClientConfig holds configuration for the modular client interceptor
type ModularClientConfig struct {
	// AppName is the application name for this client
	AppName string
	// ServerName is the server name for this client
	ServerName string
	// ServiceName is the service name for this client
	ServiceName string
}

// UnaryClientModularInterceptor returns a unary client interceptor for modular reporting (主调上报)
func UnaryClientModularInterceptor(cfg ModularClientConfig) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		startTime := time.Now()

		// Add caller info to metadata
		ctx = injectCallerMetadata(ctx, cfg)

		// Execute call
		err := invoker(ctx, method, req, reply, cc, opts...)

		// Calculate cost
		costMs := float64(time.Since(startTime).Milliseconds())

		// Build dimension
		dim := buildClientDimension(ctx, method, cc.Target(), err, cfg)

		// Report metric
		report.ReportClientMetric(ctx, dim, costMs)

		return err
	}
}

// StreamClientModularInterceptor returns a stream client interceptor for modular reporting (主调上报)
func StreamClientModularInterceptor(cfg ModularClientConfig) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		startTime := time.Now()

		// Add caller info to metadata
		ctx = injectCallerMetadata(ctx, cfg)

		// Execute call
		stream, err := streamer(ctx, desc, cc, method, opts...)

		// Calculate cost
		costMs := float64(time.Since(startTime).Milliseconds())

		// Build dimension
		dim := buildClientDimension(ctx, method, cc.Target(), err, cfg)

		// Report metric
		report.ReportClientMetric(ctx, dim, costMs)

		return stream, err
	}
}

func injectCallerMetadata(ctx context.Context, cfg ModularClientConfig) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}

	// Add caller app name
	if cfg.AppName != "" {
		md.Set("x-caller-app", cfg.AppName)
	}

	return metadata.NewOutgoingContext(ctx, md)
}

func buildClientDimension(_ context.Context, fullMethod, target string, err error, cfg ModularClientConfig) *report.ClientDimension {
	// Parse service and method from fullMethod
	service, method := parseFullMethod(fullMethod)

	// Get local IP
	localIP := ""
	if hostIP, ipErr := net_.GetHostIP(); ipErr == nil {
		localIP = hostIP.String()
	}

	return &report.ClientDimension{
		RetCode:    report.ErrorToRetCode(err),
		AIp:        localIP,
		AApp:       cfg.AppName,
		AServer:    cfg.ServerName,
		AService:   cfg.ServiceName,
		AInterface: method,
		PIp:        target,
		PService:   service,
		PInterface: method,
	}
}
