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
package report

import (
	"context"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RetCode constants
const (
	RetCodeSuccess = 0
	RetCodeTimeout = -1 // 超时错误码
)

// ServerDimension represents server-side (被调) dimension attributes
type ServerDimension struct {
	// RetCode - 返回码
	RetCode int

	// Passive side (被调方 P)
	PIp        string // 被调 IP
	PApp       string // 被调应用名
	PServer    string // 被调服务名
	PService   string // 被调 Service (gRPC service name)
	PInterface string // 被调接口 (gRPC method)

	// Active side (主调方 A)
	AIp  string // 主调 IP (from peer)
	AApp string // 主调应用名 (from metadata)
}

// IsSuccess returns true if the request was successful
func (d *ServerDimension) IsSuccess() bool {
	return d.RetCode == RetCodeSuccess
}

// IsTimeout returns true if the request timed out
func (d *ServerDimension) IsTimeout() bool {
	return d.RetCode == RetCodeTimeout
}

// ToAttributes converts dimension to OpenTelemetry attributes
func (d *ServerDimension) ToAttributes() []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		RetCodeKey.String(strconv.Itoa(d.RetCode)),
	}

	if d.PIp != "" {
		attrs = append(attrs, PIpKey.String(d.PIp))
	}
	if d.PApp != "" {
		attrs = append(attrs, PAppKey.String(d.PApp))
	}
	if d.PServer != "" {
		attrs = append(attrs, PServerKey.String(d.PServer))
	}
	if d.PService != "" {
		attrs = append(attrs, PServiceKey.String(d.PService))
	}
	if d.PInterface != "" {
		attrs = append(attrs, PInterfaceKey.String(d.PInterface))
	}
	if d.AIp != "" {
		attrs = append(attrs, AIpKey.String(d.AIp))
	}
	if d.AApp != "" {
		attrs = append(attrs, AAppKey.String(d.AApp))
	}

	return attrs
}

// ClientDimension represents client-side (主调) dimension attributes
type ClientDimension struct {
	// RetCode - 返回码
	RetCode int

	// Active side (主调方 A)
	AIp        string // 主调 IP
	AApp       string // 主调应用名
	AServer    string // 主调服务名
	AService   string // 主调 Service
	AInterface string // 主调接口

	// Passive side (被调方 P)
	PIp        string // 被调 IP (target address)
	PApp       string // 被调应用名
	PServer    string // 被调服务名
	PService   string // 被调 Service (gRPC service name)
	PInterface string // 被调接口 (gRPC method)
}

// IsSuccess returns true if the request was successful
func (d *ClientDimension) IsSuccess() bool {
	return d.RetCode == RetCodeSuccess
}

// IsTimeout returns true if the request timed out
func (d *ClientDimension) IsTimeout() bool {
	return d.RetCode == RetCodeTimeout
}

// ToAttributes converts dimension to OpenTelemetry attributes
func (d *ClientDimension) ToAttributes() []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		RetCodeKey.String(strconv.Itoa(d.RetCode)),
	}

	// Active side (主调方)
	if d.AIp != "" {
		attrs = append(attrs, AIpKey.String(d.AIp))
	}
	if d.AApp != "" {
		attrs = append(attrs, AAppKey.String(d.AApp))
	}
	if d.AServer != "" {
		attrs = append(attrs, AServerKey.String(d.AServer))
	}
	if d.AService != "" {
		attrs = append(attrs, AServiceKey.String(d.AService))
	}
	if d.AInterface != "" {
		attrs = append(attrs, AInterfaceKey.String(d.AInterface))
	}

	// Passive side (被调方)
	if d.PIp != "" {
		attrs = append(attrs, PIpKey.String(d.PIp))
	}
	if d.PApp != "" {
		attrs = append(attrs, PAppKey.String(d.PApp))
	}
	if d.PServer != "" {
		attrs = append(attrs, PServerKey.String(d.PServer))
	}
	if d.PService != "" {
		attrs = append(attrs, PServiceKey.String(d.PService))
	}
	if d.PInterface != "" {
		attrs = append(attrs, PInterfaceKey.String(d.PInterface))
	}

	return attrs
}

// ErrorToRetCode converts an error to a return code
func ErrorToRetCode(err error) int {
	if err == nil {
		return RetCodeSuccess
	}

	// Check for gRPC status
	st, ok := status.FromError(err)
	if ok {
		switch st.Code() {
		case codes.OK:
			return RetCodeSuccess
		case codes.DeadlineExceeded:
			return RetCodeTimeout
		case codes.Canceled:
			return -2 // Canceled
		default:
			return int(st.Code())
		}
	}

	// Check for context errors
	if err == context.DeadlineExceeded {
		return RetCodeTimeout
	}
	if err == context.Canceled {
		return -2
	}

	// Default error code
	return -999
}
