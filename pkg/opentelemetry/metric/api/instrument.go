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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
)

// Instrument represents a metric instrument with pre-configured meter and attributes
type Instrument struct {
	meterProvider  otelmetric.MeterProvider
	meterName      string
	instrumentName string
	attributes     map[string]attribute.KeyValue
	bounds         []float64 // for histogram
}

// InstrumentOption configures an Instrument
type InstrumentOption func(*Instrument)

// WithBounds sets histogram bucket boundaries
func WithBounds(bounds []float64) InstrumentOption {
	return func(i *Instrument) {
		i.bounds = bounds
	}
}

// WithInstrumentAttrs sets initial attributes
func WithInstrumentAttrs(attrs map[string]any) InstrumentOption {
	return func(i *Instrument) {
		for k, v := range attrs {
			i.attributes[k] = ToAttribute(k, v)
		}
	}
}

// NewInstrument creates an instrument using the application MeterProvider
func NewInstrument(meterName, instrumentName string, opts ...InstrumentOption) *Instrument {
	i := &Instrument{
		meterProvider:  GetAppMeterProvider(),
		meterName:      meterName,
		instrumentName: instrumentName,
		attributes:     make(map[string]attribute.KeyValue),
		bounds:         DefaultDurationBounds,
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

// NewGlobalInstrument creates an instrument using the global MeterProvider
func NewGlobalInstrument(meterName, instrumentName string, opts ...InstrumentOption) *Instrument {
	i := &Instrument{
		meterProvider:  otel.GetMeterProvider(),
		meterName:      meterName,
		instrumentName: instrumentName,
		attributes:     make(map[string]attribute.KeyValue),
		bounds:         DefaultDurationBounds,
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

// Clone creates a copy of the instrument with the same configuration
func (i *Instrument) Clone() *Instrument {
	newAttrs := make(map[string]attribute.KeyValue, len(i.attributes))
	for k, v := range i.attributes {
		newAttrs[k] = v
	}
	return &Instrument{
		meterProvider:  i.meterProvider,
		meterName:      i.meterName,
		instrumentName: i.instrumentName,
		attributes:     newAttrs,
		bounds:         i.bounds,
	}
}

// SetAttribute sets an attribute on the instrument
// Returns the instrument for chaining
func (i *Instrument) SetAttribute(key string, value any) *Instrument {
	i.attributes[key] = ToAttribute(key, value)
	return i
}

// SetAttributes sets multiple attributes on the instrument
// Returns the instrument for chaining
func (i *Instrument) SetAttributes(attrs map[string]any) *Instrument {
	for k, v := range attrs {
		i.attributes[k] = ToAttribute(k, v)
	}
	return i
}

// ClearAttributes clears all attributes
func (i *Instrument) ClearAttributes() *Instrument {
	i.attributes = make(map[string]attribute.KeyValue)
	return i
}

// GetAttributes returns all configured attributes as a slice
func (i *Instrument) GetAttributes() []attribute.KeyValue {
	result := make([]attribute.KeyValue, 0, len(i.attributes))
	for _, kv := range i.attributes {
		result = append(result, kv)
	}
	return result
}

// AddCounter adds a value to the counter
func (i *Instrument) AddCounter(ctx context.Context, value int64) {
	counter := getCounter(i.meterName, i.instrumentName, i.meterProvider)
	counter.Add(ctx, value, otelmetric.WithAttributes(i.GetAttributes()...))
}

// IncrCounter increments the counter by 1
func (i *Instrument) IncrCounter(ctx context.Context) {
	i.AddCounter(ctx, 1)
}

// RecordHistogram records a value to the histogram
func (i *Instrument) RecordHistogram(ctx context.Context, value float64) {
	histogram := getHistogram(i.meterName, i.instrumentName, i.bounds, i.meterProvider)
	histogram.Record(ctx, value, otelmetric.WithAttributes(i.GetAttributes()...))
}

// RecordDuration records a duration value in milliseconds
func (i *Instrument) RecordDuration(ctx context.Context, durationMs float64) {
	i.RecordHistogram(ctx, durationMs)
}

// ========================================
// Counter is a specialized instrument for counter metrics
// ========================================

// Counter represents a counter metric instrument
type Counter struct {
	*Instrument
}

// NewCounter creates a counter using the application MeterProvider
func NewCounter(meterName, instrumentName string, opts ...InstrumentOption) *Counter {
	return &Counter{Instrument: NewInstrument(meterName, instrumentName, opts...)}
}

// NewGlobalCounter creates a counter using the global MeterProvider
func NewGlobalCounter(meterName, instrumentName string, opts ...InstrumentOption) *Counter {
	return &Counter{Instrument: NewGlobalInstrument(meterName, instrumentName, opts...)}
}

// Add adds a value to the counter
func (c *Counter) Add(ctx context.Context, value int64) {
	c.AddCounter(ctx, value)
}

// Incr increments the counter by 1
func (c *Counter) Incr(ctx context.Context) {
	c.IncrCounter(ctx)
}

// With returns a new Counter with additional attributes (for chaining)
func (c *Counter) With(key string, value any) *Counter {
	newCounter := &Counter{Instrument: c.Clone()}
	newCounter.SetAttribute(key, value)
	return newCounter
}

// ========================================
// Histogram is a specialized instrument for histogram metrics
// ========================================

// Histogram represents a histogram metric instrument
type Histogram struct {
	*Instrument
}

// NewHistogram creates a histogram using the application MeterProvider
func NewHistogram(meterName, instrumentName string, opts ...InstrumentOption) *Histogram {
	return &Histogram{Instrument: NewInstrument(meterName, instrumentName, opts...)}
}

// NewGlobalHistogram creates a histogram using the global MeterProvider
func NewGlobalHistogram(meterName, instrumentName string, opts ...InstrumentOption) *Histogram {
	return &Histogram{Instrument: NewGlobalInstrument(meterName, instrumentName, opts...)}
}

// Record records a value to the histogram
func (h *Histogram) Record(ctx context.Context, value float64) {
	h.RecordHistogram(ctx, value)
}

// With returns a new Histogram with additional attributes (for chaining)
func (h *Histogram) With(key string, value any) *Histogram {
	newHistogram := &Histogram{Instrument: h.Clone()}
	newHistogram.SetAttribute(key, value)
	return newHistogram
}

// ========================================
// Timer is a helper for measuring duration
// ========================================

// Timer is a helper for measuring and recording duration
type Timer struct {
	histogram *Histogram
}

// NewTimer creates a timer using the application MeterProvider
func NewTimer(meterName, instrumentName string, opts ...InstrumentOption) *Timer {
	return &Timer{histogram: NewHistogram(meterName, instrumentName, opts...)}
}

// NewGlobalTimer creates a timer using the global MeterProvider
func NewGlobalTimer(meterName, instrumentName string, opts ...InstrumentOption) *Timer {
	return &Timer{histogram: NewGlobalHistogram(meterName, instrumentName, opts...)}
}

// SetAttribute sets an attribute on the timer
func (t *Timer) SetAttribute(key string, value any) *Timer {
	t.histogram.SetAttribute(key, value)
	return t
}

// Record records a duration in milliseconds
func (t *Timer) Record(ctx context.Context, durationMs float64) {
	t.histogram.Record(ctx, durationMs)
}

// With returns a new Timer with additional attributes
func (t *Timer) With(key string, value any) *Timer {
	return &Timer{histogram: t.histogram.With(key, value)}
}
