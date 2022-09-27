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
package jaeger

import (
	"fmt"
	url_ "net/url"
	"strings"

	"net"

	"go.opentelemetry.io/otel/exporters/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/net/context"
)

type JaegerExporterBuilderOptions struct {
}

type JaegerExporterBuilder struct {
	exporter *jaeger.Exporter
	opts     JaegerExporterBuilderOptions
}

func NewJaegerExporertBuilder(url string) (*JaegerExporterBuilder, error) {
	builder := &JaegerExporterBuilder{}
	u, err := url_.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("parse url: %v, err: %v", url, err)
	}

	switch strings.ToLower(u.Scheme) {
	case "http":
		builder.exporter, err = jaeger.New(jaeger.WithCollectorEndpoint([]jaeger.CollectorEndpointOption{
			jaeger.WithEndpoint(u.Host),
		}...))

	case "agent", "udp":
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return nil, fmt.Errorf("parse url: %v, err: %v", url, err)
		}
		builder.exporter, err = jaeger.New(jaeger.WithAgentEndpoint([]jaeger.AgentEndpointOption{
			jaeger.WithAgentHost(host),
			jaeger.WithAgentPort(port),
		}...))
	default:
		return nil, fmt.Errorf("unsupport scheme url: %v", url)
	}

	return builder, err
}

func (b *JaegerExporterBuilder) Build(ctx context.Context) (sdktrace.SpanExporter, error) {
	return b.exporter, nil
}
