package opentelemetry

import (
	"context"
	"fmt"

	"github.com/kaydxh/golang/pkg/monitor/opentelemetry/metrics/prometheus"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/ory/viper"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Proto Monitor_OpenTelemetry
	opts  struct {
		// If set, overrides params below
		viper *viper.Viper
	}
}

type completedConfig struct {
	*Config
	completeError error
}

type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

func (c *completedConfig) New(ctx context.Context) error {

	logrus.Infof("Installing OpenTelemetry")

	if c.completeError != nil {
		return c.completeError
	}

	if !c.Proto.GetEnabled() {
		return nil
	}

	err := c.install(ctx)
	if err != nil {
		return err
	}
	logrus.Infof("Installed OpenTelemetry")

	return nil
}

func (c *completedConfig) install(ctx context.Context) error {

	var opts []OpenTelemetryOption
	metricType := c.Proto.OtelMetricExporterType
	switch metricType {
	case Monitor_OpenTelemetry_metric_prometheus:
		opts = append(opts, WithMeterPullExporter(&prometheus.PrometheusExporterBuiler{}))

	default:
		return fmt.Errorf("not support the metricType[%v]", metricType.String())

	}

	ot := NewOpenTelemetry(opts...)
	return ot.Install(ctx)
}

// Complete set default ServerRunOptions.
func (c *Config) Complete() CompletedConfig {
	err := c.loadViper()
	if err != nil {
		return CompletedConfig{&completedConfig{
			Config:        c,
			completeError: err,
		}}
	}

	return CompletedConfig{&completedConfig{Config: c}}
}

func (c *Config) loadViper() error {
	if c.opts.viper != nil {
		return viper_.UnmarshalProtoMessageWithJsonPb(c.opts.viper, &c.Proto)
	}

	return nil
}

func NewConfig(options ...ConfigOption) *Config {
	c := &Config{}
	c.ApplyOptions(options...)

	return c
}
