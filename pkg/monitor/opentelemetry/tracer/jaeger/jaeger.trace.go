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
