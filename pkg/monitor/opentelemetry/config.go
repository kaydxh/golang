package opentelemetry

import (
	"context"
	"fmt"

	prometheus_ "github.com/kaydxh/golang/pkg/monitor/opentelemetry/metric/prometheus"
	jaeger_ "github.com/kaydxh/golang/pkg/monitor/opentelemetry/tracer/jaeger"
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

	var openTelemetryOpts []OpenTelemetryOption
	opts, err := c.installMeter(ctx)
	if err != nil {
		return err
	}
	openTelemetryOpts = append(openTelemetryOpts, opts...)

	opts, err = c.installTracer(ctx)
	if err != nil {
		return err
	}
	openTelemetryOpts = append(openTelemetryOpts, opts...)

	ot := NewOpenTelemetry(openTelemetryOpts...)
	return ot.Install(ctx)
}

func (c *completedConfig) installMeter(ctx context.Context) ([]OpenTelemetryOption, error) {

	var opts []OpenTelemetryOption
	collectDuration := c.Proto.GetMetricCollectDuration().AsDuration()
	if collectDuration > 0 {
		opts = append(opts, WithMetricCollectDuration(collectDuration))
	}

	metricType := c.Proto.OtelMetricExporterType
	switch metricType {
	case Monitor_OpenTelemetry_metric_prometheus:
		builder := prometheus_.NewPrometheusExporterBuilder(
			prometheus_.WithMetricUrlPath(c.Proto.GetOtelMetricExporter().GetPrometheus().GetUrl()),
		)
		opts = append(opts, WithMeterPullExporter(builder))

	case Monitor_OpenTelemetry_metric_none:
		// not enable metric
		return nil, nil

	default:
		return nil, fmt.Errorf("not support the metricType[%v]", metricType.String())

	}

	return opts, nil
}

func (c *completedConfig) installTracer(ctx context.Context) ([]OpenTelemetryOption, error) {

	var opts []OpenTelemetryOption
	tracerType := c.Proto.OtelTraceExporterType
	switch tracerType {
	case Monitor_OpenTelemetry_trace_jaeger:
		builder, err := jaeger_.NewJaegerExporertBuilder(c.Proto.GetOtelTraceExporter().GetJaeger().GetUrl())
		if err != nil {
			return nil, fmt.Errorf("new jaeger exporter builder err: %v", err)
		}
		opts = append(opts, WithTracerExporter(builder))

	case Monitor_OpenTelemetry_trace_none:
		// not enable tracer
		return nil, nil

	default:
		return nil, fmt.Errorf("not support the tracerType[%v]", tracerType.String())
	}

	return opts, nil
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