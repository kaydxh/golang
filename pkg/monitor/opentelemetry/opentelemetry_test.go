package opentelemetry_test

import (
	"fmt"
	"testing"

	opentelemetry_ "github.com/kaydxh/golang/pkg/monitor/opentelemetry"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"golang.org/x/net/context"
)

func TestNew(t *testing.T) {
	cfgFile := "./opentelemetry.yaml"
	config := opentelemetry_.NewConfig(opentelemetry_.WithViper(viper_.GetViper(cfgFile, "monitor.open_telemetry")))

	err := config.Complete().New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}
	fmt.Printf("config: %+v", config.Proto.String())

}
