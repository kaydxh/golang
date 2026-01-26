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
package api_test

import (
	"context"
	"time"

	"github.com/kaydxh/golang/pkg/opentelemetry/metric/api"
	"go.opentelemetry.io/otel/attribute"
)

func Example_simpleCounter() {
	ctx := context.Background()

	// Method 1: Using context attributes
	ctx = api.WithAttribute(ctx, "my_app", "user_id", "12345")
	ctx = api.WithAttribute(ctx, "my_app", "action", "login")
	api.IncrCounter(ctx, "my_app", "user_actions")

	// Method 2: Using explicit attributes
	api.IncrCounterWithAttrs(ctx, "my_app", "requests",
		attribute.String("endpoint", "/api/v1/users"),
		attribute.Int("status", 200),
	)

	// Method 3: Using Instrument object
	counter := api.NewCounter("my_app", "api_calls")
	counter.SetAttribute("service", "user-service")
	counter.Incr(ctx)

	// Method 4: Chaining with With
	counter.With("method", "GET").With("path", "/users").Incr(ctx)
}

func Example_histogram() {
	ctx := context.Background()

	// Record latency with context attributes
	ctx = api.WithAttribute(ctx, "latency", "endpoint", "/api/v1/query")
	api.RecordDuration(ctx, "latency", "request_duration", 125.5)

	// Record with explicit attributes
	api.RecordDurationWithAttrs(ctx, "latency", "db_query_duration", 45.2,
		attribute.String("table", "users"),
		attribute.String("operation", "select"),
	)

	// Using Histogram object
	histogram := api.NewHistogram("latency", "processing_time",
		api.WithBounds([]float64{10, 50, 100, 500, 1000}),
	)
	histogram.SetAttribute("processor", "image")
	histogram.Record(ctx, 234.5)
}

func Example_timer() {
	ctx := context.Background()

	// Create a timer for measuring operation duration
	timer := api.NewTimer("performance", "operation_duration")
	timer.SetAttribute("operation", "data_processing")

	start := time.Now()
	// ... do some work ...
	_ = start // simulate work

	// Record the duration
	elapsed := float64(time.Since(start).Milliseconds())
	timer.Record(ctx, elapsed)

	// With chaining for different operations
	timer.With("operation", "file_upload").Record(ctx, 150.0)
	timer.With("operation", "database_write").Record(ctx, 25.5)
}

func Example_globalVsApp() {
	ctx := context.Background()

	// Global metrics (infrastructure level)
	api.GlobalIncrCounter(ctx, "infrastructure", "gc_cycles")
	api.GlobalRecordDuration(ctx, "infrastructure", "gc_pause", 5.2)

	// App metrics (application level)
	api.IncrCounter(ctx, "business", "orders_created")
	api.RecordDuration(ctx, "business", "checkout_time", 1250.0)

	// Using Instrument objects
	globalCounter := api.NewGlobalCounter("system", "memory_allocations")
	globalCounter.Incr(ctx)

	appCounter := api.NewCounter("app", "user_signups")
	appCounter.Incr(ctx)
}

func Example_batchAttributes() {
	ctx := context.Background()

	// Add multiple attributes at once
	ctx = api.WithAttributes(ctx, "my_service", map[string]any{
		"user_id":    "12345",
		"session_id": "abc-def-ghi",
		"region":     "us-west-2",
		"version":    "1.2.3",
	})

	api.IncrCounter(ctx, "my_service", "api_calls")
	api.RecordDuration(ctx, "my_service", "latency", 45.2)
}
