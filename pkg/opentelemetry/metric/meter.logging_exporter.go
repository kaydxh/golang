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
package metric

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// LoggingExporter wraps a metric.Exporter and logs export results
type LoggingExporter struct {
	inner    metric.Exporter
	name     string // exporter name for logging (e.g., "OTLP", "stdout")
	endpoint string // endpoint for logging
}

// NewLoggingExporter creates a new LoggingExporter that wraps the given exporter
func NewLoggingExporter(inner metric.Exporter, name, endpoint string) *LoggingExporter {
	return &LoggingExporter{
		inner:    inner,
		name:     name,
		endpoint: endpoint,
	}
}

// Temporality returns the Temporality to use for an instrument kind.
func (e *LoggingExporter) Temporality(kind metric.InstrumentKind) metricdata.Temporality {
	return e.inner.Temporality(kind)
}

// Aggregation returns the Aggregation to use for an instrument kind.
func (e *LoggingExporter) Aggregation(kind metric.InstrumentKind) metric.Aggregation {
	return e.inner.Aggregation(kind)
}

// Export exports metrics and logs the result
func (e *LoggingExporter) Export(ctx context.Context, rm *metricdata.ResourceMetrics) error {
	start := time.Now()

	// Count metrics and check temporality
	metricCount := 0
	dataPointCount := 0
	temporality := "unknown"
	for _, sm := range rm.ScopeMetrics {
		metricCount += len(sm.Metrics)
		for _, m := range sm.Metrics {
			dataPointCount += countDataPoints(m)
			if temporality == "unknown" {
				temporality = getTemporality(m)
			}
		}
	}

	// Log resource attributes for debugging (first export only)
	//if rm.Resource != nil {
	//	attrs := rm.Resource.Attributes()
	//	for _, attr := range attrs {
	//		if attr.Key == "__zhiyan_app_mark__" {
	//			logrus.Debugf("[%s] resource attribute: %s=%s", e.name, attr.Key, attr.Value.AsString())
	//		}
	//	}
	//}

	err := e.inner.Export(ctx, rm)
	duration := time.Since(start)

	if err != nil {
		logrus.Errorf("[%s] metric export failed: endpoint=%s, metrics=%d, datapoints=%d, temporality=%s, duration=%v, error=%v",
			e.name, e.endpoint, metricCount, dataPointCount, temporality, duration, err)
	} else {
		logrus.Infof("[%s] metric export success: endpoint=%s, metrics=%d, datapoints=%d, temporality=%s, duration=%v",
			e.name, e.endpoint, metricCount, dataPointCount, temporality, duration)
	}

	return err
}

// ForceFlush flushes the exporter
func (e *LoggingExporter) ForceFlush(ctx context.Context) error {
	return e.inner.ForceFlush(ctx)
}

// Shutdown shuts down the exporter
func (e *LoggingExporter) Shutdown(ctx context.Context) error {
	logrus.Infof("[%s] metric exporter shutting down: endpoint=%s", e.name, e.endpoint)
	return e.inner.Shutdown(ctx)
}

// countDataPoints counts the number of data points in a metric
func countDataPoints(m metricdata.Metrics) int {
	switch data := m.Data.(type) {
	case metricdata.Sum[int64]:
		return len(data.DataPoints)
	case metricdata.Sum[float64]:
		return len(data.DataPoints)
	case metricdata.Gauge[int64]:
		return len(data.DataPoints)
	case metricdata.Gauge[float64]:
		return len(data.DataPoints)
	case metricdata.Histogram[int64]:
		return len(data.DataPoints)
	case metricdata.Histogram[float64]:
		return len(data.DataPoints)
	case metricdata.ExponentialHistogram[int64]:
		return len(data.DataPoints)
	case metricdata.ExponentialHistogram[float64]:
		return len(data.DataPoints)
	default:
		return 0
	}
}

// getTemporality returns the temporality string of a metric
func getTemporality(m metricdata.Metrics) string {
	switch data := m.Data.(type) {
	case metricdata.Sum[int64]:
		return data.Temporality.String()
	case metricdata.Sum[float64]:
		return data.Temporality.String()
	case metricdata.Histogram[int64]:
		return data.Temporality.String()
	case metricdata.Histogram[float64]:
		return data.Temporality.String()
	case metricdata.ExponentialHistogram[int64]:
		return data.Temporality.String()
	case metricdata.ExponentialHistogram[float64]:
		return data.Temporality.String()
	default:
		return "gauge"
	}
}
