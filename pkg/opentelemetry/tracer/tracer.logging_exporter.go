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
	"time"

	"github.com/sirupsen/logrus"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// LoggingExporter wraps a trace.SpanExporter and logs export results
type LoggingExporter struct {
	inner    sdktrace.SpanExporter
	name     string // exporter name for logging (e.g., "OTLP", "Jaeger", "stdout")
	endpoint string // endpoint for logging
}

// NewLoggingExporter creates a new LoggingExporter that wraps the given exporter
func NewLoggingExporter(inner sdktrace.SpanExporter, name, endpoint string) *LoggingExporter {
	return &LoggingExporter{
		inner:    inner,
		name:     name,
		endpoint: endpoint,
	}
}

// ExportSpans exports spans and logs the result
func (e *LoggingExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	start := time.Now()

	spanCount := len(spans)

	// Log span details for debugging
	if spanCount > 0 {
		logrus.Debugf("[%s] exporting %d spans to %s", e.name, spanCount, e.endpoint)
	}

	err := e.inner.ExportSpans(ctx, spans)
	duration := time.Since(start)

	if err != nil {
		logrus.Errorf("[%s] trace export failed: endpoint=%s, spans=%d, duration=%v, error=%v",
			e.name, e.endpoint, spanCount, duration, err)
	} else {
		if spanCount > 0 {
			logrus.Infof("[%s] trace export success: endpoint=%s, spans=%d, duration=%v",
				e.name, e.endpoint, spanCount, duration)
		}
	}

	return err
}

// Shutdown shuts down the exporter
func (e *LoggingExporter) Shutdown(ctx context.Context) error {
	logrus.Infof("[%s] trace exporter shutting down: endpoint=%s", e.name, e.endpoint)
	return e.inner.Shutdown(ctx)
}
