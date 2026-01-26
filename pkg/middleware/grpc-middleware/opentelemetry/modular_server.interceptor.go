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
	"strings"
	"time"

	net_ "github.com/kaydxh/golang/go/net"
	grpc_ "github.com/kaydxh/golang/go/net/grpc"
	"github.com/kaydxh/golang/pkg/opentelemetry/metric/report"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ModularServerConfig holds configuration for the modular server interceptor
type ModularServerConfig struct {
	// AppName is the application name for this service
	AppName string
	// ServerName is the server name for this service
	ServerName string
}

// UnaryServerModularInterceptor returns a unary server interceptor for modular reporting (被调上报)
func UnaryServerModularInterceptor(cfg ModularServerConfig) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		// Execute handler
		resp, err := handler(ctx, req)

		// Calculate cost
		costMs := float64(time.Since(startTime).Milliseconds())

		// Build dimension
		dim := buildServerDimension(ctx, info.FullMethod, err, cfg)

		// Report metric
		report.ReportServerMetric(ctx, dim, costMs)

		return resp, err
	}
}

// StreamServerModularInterceptor returns a stream server interceptor for modular reporting (被调上报)
func StreamServerModularInterceptor(cfg ModularServerConfig) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()

		// Execute handler
		err := handler(srv, ss)

		// Calculate cost
		costMs := float64(time.Since(startTime).Milliseconds())

		// Build dimension
		dim := buildServerDimension(ss.Context(), info.FullMethod, err, cfg)

		// Report metric
		report.ReportServerMetric(ss.Context(), dim, costMs)

		return err
	}
}

func buildServerDimension(ctx context.Context, fullMethod string, err error, cfg ModularServerConfig) *report.ServerDimension {
	// Parse service and method from fullMethod
	// fullMethod format: /package.Service/Method
	service, method := parseFullMethod(fullMethod)

	// Get local IP
	localIP := ""
	if hostIP, ipErr := net_.GetHostIP(); ipErr == nil {
		localIP = hostIP.String()
	}

	// Get peer IP
	peerIP := ""
	if peerAddr, peerErr := grpc_.GetIPFromContext(ctx); peerErr == nil {
		peerIP = peerAddr.String()
	}

	// Get caller app from metadata
	callerApp := ""
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if apps := md.Get("x-caller-app"); len(apps) > 0 {
			callerApp = apps[0]
		}
	}

	return &report.ServerDimension{
		RetCode:    report.ErrorToRetCode(err),
		PIp:        localIP,
		PApp:       cfg.AppName,
		PServer:    cfg.ServerName,
		PService:   service,
		PInterface: method,
		AIp:        peerIP,
		AApp:       callerApp,
	}
}

// parseFullMethod parses the gRPC full method into service and method
// fullMethod format: /package.Service/Method
func parseFullMethod(fullMethod string) (service, method string) {
	// Remove leading slash
	fullMethod = strings.TrimPrefix(fullMethod, "/")

	// Split by last slash
	idx := strings.LastIndex(fullMethod, "/")
	if idx < 0 {
		return fullMethod, ""
	}

	return fullMethod[:idx], fullMethod[idx+1:]
}
