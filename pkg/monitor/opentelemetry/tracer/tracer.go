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
package tracer

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type TracerOptions struct {
	builer           TracerExporterBuilder
	serviceName      string
	serviceVersion   string
	serviceNamespace string
	tracerProvider   *sdktrace.TracerProvider
}

type Tracer struct {
	opts TracerOptions
}

func NewTracer(opts ...TracerOption) *Tracer {
	t := &Tracer{}
	t.ApplyOptions(opts...)

	return t
}

//https://github.com/open-telemetry/opentelemetry-go/blob/main/example/jaeger/main.go
func (t *Tracer) Install(ctx context.Context) (err error) {
	exp, err := t.createExporter(ctx)
	if err != nil {
		return err
	}

	// Build resource attributes
	resourceAttrs := []resource.Option{
		resource.WithSchemaURL(semconv.SchemaURL),
	}

	// Add service information if provided
	if t.opts.serviceName != "" {
		attrs := []attribute.KeyValue{
			semconv.ServiceName(t.opts.serviceName),
		}
		if t.opts.serviceVersion != "" {
			attrs = append(attrs, semconv.ServiceVersion(t.opts.serviceVersion))
		}
		if t.opts.serviceNamespace != "" {
			attrs = append(attrs, semconv.ServiceNamespace(t.opts.serviceNamespace))
		}
		resourceAttrs = append(resourceAttrs, resource.WithAttributes(attrs...))
	}

	res, err := resource.New(ctx, resourceAttrs...)
	if err != nil {
		return fmt.Errorf("creating resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		// Record information about this application in a Resource.
		sdktrace.WithResource(res),
	)

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	// Store the tracer provider for later use
	t.opts.tracerProvider = tp

	return nil
}

func (t *Tracer) createExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	if t.opts.builer == nil {
		return nil, fmt.Errorf("trace exporter builder is nil")
	}

	return t.opts.builer.Build(ctx)
}

// TracerProvider returns the configured TracerProvider
func (t *Tracer) TracerProvider() trace.TracerProvider {
	if t.opts.tracerProvider != nil {
		return t.opts.tracerProvider
	}
	return otel.GetTracerProvider()
}

// Shutdown shuts down the tracer provider
func (t *Tracer) Shutdown(ctx context.Context) error {
	if t.opts.tracerProvider != nil {
		return t.opts.tracerProvider.Shutdown(ctx)
	}
	return nil
}
