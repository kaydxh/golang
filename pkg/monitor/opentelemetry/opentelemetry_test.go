package opentelemetry_test

import (
	"fmt"
	"net/http"
	"testing"

	opentelemetry_ "github.com/kaydxh/golang/pkg/monitor/opentelemetry"
	viper_ "github.com/kaydxh/golang/pkg/viper"
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

func TestNew(t *testing.T) {
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

	metrics(ctx, 100)

	select {}
}

func metrics(ctx context.Context, n int) {
	funcNameKV := funcNameKey.String("metrics")
	for i := 0; i < n; i++ {
		funcLoopCounter.Add(ctx, 1, funcNameKV)
	}
}
