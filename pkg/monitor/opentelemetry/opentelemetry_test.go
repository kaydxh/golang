package opentelemetry_test

import (
	"fmt"
	"testing"

	opentelemetry_ "github.com/kaydxh/golang/pkg/monitor/opentelemetry"
	viper_ "github.com/kaydxh/golang/pkg/viper"
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

	metrics(ctx, 100)

}

func metrics(ctx context.Context, n int) {
	for i := 0; i < n; i++ {
		funcLoopCounter.Add(ctx, 1)
	}
}
