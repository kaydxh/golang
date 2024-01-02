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
package opentelemetry

import (
	"context"
	"fmt"

	prometheus_ "github.com/kaydxh/golang/pkg/monitor/opentelemetry/metric/prometheus"
	stdoutmetric_ "github.com/kaydxh/golang/pkg/monitor/opentelemetry/metric/stdout"
	"github.com/kaydxh/golang/pkg/monitor/opentelemetry/resource"
	jaeger_ "github.com/kaydxh/golang/pkg/monitor/opentelemetry/tracer/jaeger"
	stdouttrace_ "github.com/kaydxh/golang/pkg/monitor/opentelemetry/tracer/stdout"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Proto Monitor_OpenTelemetry
	opts  struct {
		// If set, overrides params below
		viper                       *viper.Viper
		resourceStatsServiceOptions []resource.ResourceStatsServiceOption
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
	err = ot.Install(ctx)
	if err != nil {
		return err
	}

	_, err = c.installResource(ctx)
	return err
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

	case Monitor_OpenTelemetry_metric_stdout:
		builder := stdoutmetric_.NewStdoutExporterBuilder(
			stdoutmetric_.WithPrettyPrint(c.Proto.GetOtelMetricExporter().GetStdout().GetPrettyPrint()),
		)
		opts = append(opts, WithMeterPushExporter(builder))

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

	case Monitor_OpenTelemetry_trace_stdout:
		builder := stdouttrace_.NewStdoutExporterBuilder(
			stdouttrace_.WithPrettyPrint(c.Proto.GetOtelTraceExporter().GetStdout().GetPrettyPrint()),
		)
		opts = append(opts, WithTracerExporter(builder))

	case Monitor_OpenTelemetry_trace_none:
		// not enable tracer
		return nil, nil

	default:
		return nil, fmt.Errorf("not support the tracerType[%v]", tracerType.String())
	}

	return opts, nil
}

func (c *completedConfig) installResource(ctx context.Context) (*resource.ResourceStatsService, error) {

	var opts []resource.ResourceStatsServiceOption
	collectDuration := c.Proto.GetMetricCollectDuration().AsDuration()
	if collectDuration > 0 {
		opts = append(opts, resource.WithStatsCheckInterval(collectDuration))
	}
	opts = append(opts, c.opts.resourceStatsServiceOptions...)

	statsServer, err := resource.NewResourceStatsService(opts...)
	if err != nil {
		return nil, err
	}
	statsServer.Run(ctx)

	return statsServer, nil
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
