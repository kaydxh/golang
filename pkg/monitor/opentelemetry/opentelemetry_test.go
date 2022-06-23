package opentelemetry_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	opentelemetry_ "github.com/kaydxh/golang/pkg/monitor/opentelemetry"
	viper_ "github.com/kaydxh/golang/pkg/viper"
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

func TestMetric(t *testing.T) {
	cfgFile := "./opentelemetry.yaml"
	config := opentelemetry_.NewConfig(opentelemetry_.WithViper(viper_.GetViper(cfgFile, "monitor.open_telemetry")))

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

//https://github.com/open-telemetry/opentelemetry-go/blob/main/example/jaeger/main.go
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

	for i := 0; i < 2; i++ {
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
