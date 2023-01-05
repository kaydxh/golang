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
package interceptorprometheus

import (
	"context"
	"fmt"

	grpc_ "github.com/kaydxh/golang/go/net/grpc"
	time_ "github.com/kaydxh/golang/go/time"
	logs_ "github.com/kaydxh/golang/pkg/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptorOfTimer returns a new unary server interceptors that timing request
func UnaryServerInterceptorOfTimer(enabledMetric bool) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		tc := time_.New(true)
		logger := logs_.GetLogger(ctx)
		summary := func() {
			tc.Tick(info.FullMethod)
			if enabledMetric {
				M.durationCost.WithLabelValues(info.FullMethod).Observe(float64(tc.Elapse().Milliseconds()))
			}

			logger.WithField("method", info.FullMethod).Infof(tc.String())
		}
		defer summary()

		return handler(ctx, req)
	}
}

// UnaryServerInterceptorOfCodeMessage returns a new unary server interceptors that timing request
func UnaryServerInterceptorOfCodeMessage(enabledMetric bool) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		peerAddr, _ := grpc_.GetIPFromContext(ctx)
		var (
			resp    interface{}
			code    uint32
			message string
			err     error
		)

		logger := logs_.GetLogger(ctx)
		summary := func() {
			codeMessage := fmt.Sprintf("%d:%s", code, message)
			if enabledMetric {
				metircLabels := map[string]string{
					MetircLabelMethod:      info.FullMethod,
					MetircLabelClientIP:    peerAddr.String(),
					MetircLabelCodeMessage: codeMessage,
				}
				M.calledTotal.With(metircLabels).Inc()
			}

			logger.WithField(
				"method",
				info.FullMethod,
			).WithField(
				"code_messge",
				codeMessage,
			).Infof(
				"called by peer addr: %v",
				peerAddr.String(),
			)
		}
		defer summary()

		resp, err = handler(ctx, req)
		grpcErr, ok := status.FromError(err)
		if ok {
			code = uint32(grpcErr.Code())
			message = grpcErr.Message()

		}
		return resp, err
	}
}
