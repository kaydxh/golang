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
package opentelemetry_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	opentelemetry_ "github.com/kaydxh/golang/pkg/monitor/opentelemetry"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"golang.org/x/net/context"
)

const (
	instrumentationName    = "github/kaydxh/instrument"
	instrumentationVersion = "v0.0.1"
)

var (
	meter = global.MeterProvider().Meter(
		"",
		metric.WithInstrumentationVersion(instrumentationVersion),
	)

	funcLoopCounter, _ = meter.SyncInt64().Counter("function_loops")
	funcNameKey        = attribute.Key("function_name")
)

func memoryUsageCallBack(total, free uint64, usage float64) {
	logrus.Infof("memory total: %v, free: %v, usage: %v", total, free, usage)
}

func TestMetric(t *testing.T) {
	cfgFile := "./opentelemetry.yaml"
	config := opentelemetry_.NewConfig(
		opentelemetry_.WithViper(viper_.GetViper(cfgFile, "monitor.open_telemetry")),
		opentelemetry_.WithMemoryCallBack(memoryUsageCallBack),
	)

	ctx := context.Background()
	err := config.Complete().New(ctx)
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}
	fmt.Printf("config: %+v", config.Proto.String())

	go func() {
		_ = http.ListenAndServe(":2222", nil)
	}()

	for {
		metrics(ctx, 100)
		time.Sleep(time.Second)
	}

}

func metrics(ctx context.Context, n int) {
	funcNameKV := funcNameKey.String("metrics")
	for i := 0; i < n; i++ {
		funcLoopCounter.Add(ctx, 1, funcNameKV)
	}
}

// https://github.com/open-telemetry/opentelemetry-go/blob/main/example/jaeger/main.go
func TestTrace(t *testing.T) {
	cfgFile := "./opentelemetry.yaml"
	config := opentelemetry_.NewConfig(opentelemetry_.WithViper(viper_.GetViper(cfgFile, "monitor.open_telemetry")))

	ctx := context.Background()
	err := config.Complete().New(ctx)
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}
	fmt.Printf("config: %+v", config.Proto.String())

	trace := otel.GetTracerProvider()
	tr := trace.Tracer(instrumentationName)

	ctx, span := tr.Start(ctx, "traceFunc")
	defer span.End()

	for {
		doTrace(ctx)
		time.Sleep(time.Second)
	}

}

func doTrace(ctx context.Context) {
	// Use the global TracerProvider.
	tr := otel.Tracer("component-doTrace")
	_, span := tr.Start(ctx, "doTrace")
	span.SetAttributes(attribute.Key("testset").String("value"))
	defer span.End()

	time.Sleep(200 * time.Millisecond)
}
