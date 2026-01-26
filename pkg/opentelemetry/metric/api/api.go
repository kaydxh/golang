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
package api

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

// Default histogram bounds for common use cases
var (
	// DefaultDurationBounds for latency in milliseconds
	DefaultDurationBounds = []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}
	// DefaultSizeBounds for size in bytes
	DefaultSizeBounds = []float64{100, 1000, 10000, 100000, 1000000, 10000000}
)

var (
	appMeterProvider   otelmetric.MeterProvider
	appMeterProviderMu sync.RWMutex

	counterCache   = &instrumentCache[otelmetric.Int64Counter]{cache: make(map[string]otelmetric.Int64Counter)}
	histogramCache = &instrumentCache[otelmetric.Float64Histogram]{cache: make(map[string]otelmetric.Float64Histogram)}
)

// instrumentCache is a generic cache for metric instruments
type instrumentCache[T any] struct {
	cache map[string]T
	mu    sync.RWMutex
}

func (c *instrumentCache[T]) get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.cache[key]
	return v, ok
}

func (c *instrumentCache[T]) set(key string, v T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = v
}

// SetAppMeterProvider sets the application-level MeterProvider
// This allows separation between global infrastructure metrics and application metrics
func SetAppMeterProvider(mp otelmetric.MeterProvider) {
	appMeterProviderMu.Lock()
	defer appMeterProviderMu.Unlock()
	appMeterProvider = mp
}

// GetAppMeterProvider returns the application-level MeterProvider
// Falls back to global MeterProvider if not set
func GetAppMeterProvider() otelmetric.MeterProvider {
	appMeterProviderMu.RLock()
	defer appMeterProviderMu.RUnlock()
	if appMeterProvider != nil {
		return appMeterProvider
	}
	return otel.GetMeterProvider()
}

// attributeContextKey is used to store attributes in context
type attributeContextKey struct {
	meterName string
}

// WithAttribute adds a custom attribute to the context for a specific meter
// value supports string, int, float, bool types
func WithAttribute(ctx context.Context, meterName, key string, value any) context.Context {
	attrs := getContextAttributes(ctx, meterName)
	if attrs == nil {
		attrs = make(map[string]attribute.KeyValue)
	} else {
		// Copy to avoid mutation
		newAttrs := make(map[string]attribute.KeyValue, len(attrs)+1)
		for k, v := range attrs {
			newAttrs[k] = v
		}
		attrs = newAttrs
	}
	attrs[key] = ToAttribute(key, value)
	return context.WithValue(ctx, attributeContextKey{meterName: meterName}, attrs)
}

// WithAttributes adds multiple custom attributes to the context for a specific meter
func WithAttributes(ctx context.Context, meterName string, kvs map[string]any) context.Context {
	attrs := getContextAttributes(ctx, meterName)
	if attrs == nil {
		attrs = make(map[string]attribute.KeyValue, len(kvs))
	} else {
		// Copy to avoid mutation
		newAttrs := make(map[string]attribute.KeyValue, len(attrs)+len(kvs))
		for k, v := range attrs {
			newAttrs[k] = v
		}
		attrs = newAttrs
	}
	for k, v := range kvs {
		attrs[k] = ToAttribute(k, v)
	}
	return context.WithValue(ctx, attributeContextKey{meterName: meterName}, attrs)
}

// getContextAttributes retrieves attributes from context for a specific meter
func getContextAttributes(ctx context.Context, meterName string) map[string]attribute.KeyValue {
	if attrs, ok := ctx.Value(attributeContextKey{meterName: meterName}).(map[string]attribute.KeyValue); ok {
		return attrs
	}
	return nil
}

// getAttributeSlice converts context attributes to a slice
func getAttributeSlice(ctx context.Context, meterName string) []attribute.KeyValue {
	attrs := getContextAttributes(ctx, meterName)
	if len(attrs) == 0 {
		return nil
	}
	result := make([]attribute.KeyValue, 0, len(attrs))
	for _, kv := range attrs {
		result = append(result, kv)
	}
	return result
}

// ToAttribute converts a value to an attribute.KeyValue
func ToAttribute(key string, value any) attribute.KeyValue {
	switch v := value.(type) {
	case string:
		return attribute.String(key, v)
	case int:
		return attribute.Int(key, v)
	case int64:
		return attribute.Int64(key, v)
	case float64:
		return attribute.Float64(key, v)
	case bool:
		return attribute.Bool(key, v)
	case int8:
		return attribute.Int64(key, int64(v))
	case int16:
		return attribute.Int64(key, int64(v))
	case int32:
		return attribute.Int64(key, int64(v))
	case uint:
		return attribute.Int64(key, int64(v))
	case uint8:
		return attribute.Int64(key, int64(v))
	case uint16:
		return attribute.Int64(key, int64(v))
	case uint32:
		return attribute.Int64(key, int64(v))
	case uint64:
		return attribute.Int64(key, int64(v))
	case float32:
		return attribute.Float64(key, float64(v))
	default:
		return attribute.String(key, fmt.Sprintf("%v", v))
	}
}

// ToAttributeString converts a value to string for attribute
func ToAttributeString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%g", v)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// getCounter retrieves or creates a counter instrument
func getCounter(meterName, instrumentName string, provider otelmetric.MeterProvider) otelmetric.Int64Counter {
	key := fmt.Sprintf("%s_%s", meterName, instrumentName)

	if counter, ok := counterCache.get(key); ok {
		return counter
	}

	meter := provider.Meter(meterName)
	counter, err := meter.Int64Counter(instrumentName)
	if err != nil {
		otel.Handle(err)
		return noop.Int64Counter{}
	}

	counterCache.set(key, counter)
	return counter
}

// getHistogram retrieves or creates a histogram instrument
func getHistogram(meterName, instrumentName string, bounds []float64, provider otelmetric.MeterProvider) otelmetric.Float64Histogram {
	key := fmt.Sprintf("%s_%s", meterName, instrumentName)

	if histogram, ok := histogramCache.get(key); ok {
		return histogram
	}

	meter := provider.Meter(meterName)
	var opts []otelmetric.Float64HistogramOption
	if len(bounds) > 0 {
		opts = append(opts, otelmetric.WithExplicitBucketBoundaries(bounds...))
	}
	histogram, err := meter.Float64Histogram(instrumentName, opts...)
	if err != nil {
		otel.Handle(err)
		return noop.Float64Histogram{}
	}

	histogramCache.set(key, histogram)
	return histogram
}

// ========================================
// Application-level Metric APIs (使用 App MeterProvider)
// ========================================

// AddCounter reports a counter value to the application MeterProvider
// meterName: metric group name
// instrumentName: metric name
// value: counter value to add
func AddCounter(ctx context.Context, meterName, instrumentName string, value int64) {
	counter := getCounter(meterName, instrumentName, GetAppMeterProvider())
	attrs := getAttributeSlice(ctx, meterName)
	counter.Add(ctx, value, otelmetric.WithAttributes(attrs...))
}

// IncrCounter increments a counter by 1 to the application MeterProvider
// meterName: metric group name
// instrumentName: metric name
func IncrCounter(ctx context.Context, meterName, instrumentName string) {
	AddCounter(ctx, meterName, instrumentName, 1)
}

// RecordHistogram records a histogram value to the application MeterProvider
// meterName: metric group name
// instrumentName: metric name
// value: histogram value to record
// bounds: histogram bucket boundaries (use nil for default)
func RecordHistogram(ctx context.Context, meterName, instrumentName string, value float64, bounds []float64) {
	if bounds == nil {
		bounds = DefaultDurationBounds
	}
	histogram := getHistogram(meterName, instrumentName, bounds, GetAppMeterProvider())
	attrs := getAttributeSlice(ctx, meterName)
	histogram.Record(ctx, value, otelmetric.WithAttributes(attrs...))
}

// RecordDuration records a duration value (in milliseconds) with default duration bounds
func RecordDuration(ctx context.Context, meterName, instrumentName string, durationMs float64) {
	RecordHistogram(ctx, meterName, instrumentName, durationMs, DefaultDurationBounds)
}

// ========================================
// Global-level Metric APIs (使用 Global MeterProvider)
// ========================================

// GlobalAddCounter reports a counter value to the global MeterProvider
// meterName: metric group name
// instrumentName: metric name
// value: counter value to add
func GlobalAddCounter(ctx context.Context, meterName, instrumentName string, value int64) {
	counter := getCounter(meterName, instrumentName, otel.GetMeterProvider())
	attrs := getAttributeSlice(ctx, meterName)
	counter.Add(ctx, value, otelmetric.WithAttributes(attrs...))
}

// GlobalIncrCounter increments a counter by 1 to the global MeterProvider
// meterName: metric group name
// instrumentName: metric name
func GlobalIncrCounter(ctx context.Context, meterName, instrumentName string) {
	GlobalAddCounter(ctx, meterName, instrumentName, 1)
}

// GlobalRecordHistogram records a histogram value to the global MeterProvider
// meterName: metric group name
// instrumentName: metric name
// value: histogram value to record
// bounds: histogram bucket boundaries (use nil for default)
func GlobalRecordHistogram(ctx context.Context, meterName, instrumentName string, value float64, bounds []float64) {
	if bounds == nil {
		bounds = DefaultDurationBounds
	}
	histogram := getHistogram(meterName, instrumentName, bounds, otel.GetMeterProvider())
	attrs := getAttributeSlice(ctx, meterName)
	histogram.Record(ctx, value, otelmetric.WithAttributes(attrs...))
}

// GlobalRecordDuration records a duration value (in milliseconds) with default duration bounds
func GlobalRecordDuration(ctx context.Context, meterName, instrumentName string, durationMs float64) {
	GlobalRecordHistogram(ctx, meterName, instrumentName, durationMs, DefaultDurationBounds)
}

// ========================================
// Convenient APIs with attributes (直接传入属性，不依赖 context)
// ========================================

// AddCounterWithAttrs reports a counter value with explicit attributes
func AddCounterWithAttrs(ctx context.Context, meterName, instrumentName string, value int64, attrs ...attribute.KeyValue) {
	counter := getCounter(meterName, instrumentName, GetAppMeterProvider())
	counter.Add(ctx, value, otelmetric.WithAttributes(attrs...))
}

// IncrCounterWithAttrs increments a counter by 1 with explicit attributes
func IncrCounterWithAttrs(ctx context.Context, meterName, instrumentName string, attrs ...attribute.KeyValue) {
	AddCounterWithAttrs(ctx, meterName, instrumentName, 1, attrs...)
}

// RecordHistogramWithAttrs records a histogram value with explicit attributes
func RecordHistogramWithAttrs(ctx context.Context, meterName, instrumentName string, value float64, bounds []float64, attrs ...attribute.KeyValue) {
	if bounds == nil {
		bounds = DefaultDurationBounds
	}
	histogram := getHistogram(meterName, instrumentName, bounds, GetAppMeterProvider())
	histogram.Record(ctx, value, otelmetric.WithAttributes(attrs...))
}

// RecordDurationWithAttrs records a duration value with explicit attributes
func RecordDurationWithAttrs(ctx context.Context, meterName, instrumentName string, durationMs float64, attrs ...attribute.KeyValue) {
	RecordHistogramWithAttrs(ctx, meterName, instrumentName, durationMs, DefaultDurationBounds, attrs...)
}

// GlobalAddCounterWithAttrs reports a counter value with explicit attributes to global provider
func GlobalAddCounterWithAttrs(ctx context.Context, meterName, instrumentName string, value int64, attrs ...attribute.KeyValue) {
	counter := getCounter(meterName, instrumentName, otel.GetMeterProvider())
	counter.Add(ctx, value, otelmetric.WithAttributes(attrs...))
}

// GlobalIncrCounterWithAttrs increments a counter by 1 with explicit attributes to global provider
func GlobalIncrCounterWithAttrs(ctx context.Context, meterName, instrumentName string, attrs ...attribute.KeyValue) {
	GlobalAddCounterWithAttrs(ctx, meterName, instrumentName, 1, attrs...)
}

// GlobalRecordHistogramWithAttrs records a histogram value with explicit attributes to global provider
func GlobalRecordHistogramWithAttrs(ctx context.Context, meterName, instrumentName string, value float64, bounds []float64, attrs ...attribute.KeyValue) {
	if bounds == nil {
		bounds = DefaultDurationBounds
	}
	histogram := getHistogram(meterName, instrumentName, bounds, otel.GetMeterProvider())
	histogram.Record(ctx, value, otelmetric.WithAttributes(attrs...))
}

// GlobalRecordDurationWithAttrs records a duration value with explicit attributes to global provider
func GlobalRecordDurationWithAttrs(ctx context.Context, meterName, instrumentName string, durationMs float64, attrs ...attribute.KeyValue) {
	GlobalRecordHistogramWithAttrs(ctx, meterName, instrumentName, durationMs, DefaultDurationBounds, attrs...)
}
