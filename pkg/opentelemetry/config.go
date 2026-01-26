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

	"github.com/gin-gonic/gin"
	otlpmetric_ "github.com/kaydxh/golang/pkg/opentelemetry/metric/otlp"
	prometheus_ "github.com/kaydxh/golang/pkg/opentelemetry/metric/prometheus"
	stdoutmetric_ "github.com/kaydxh/golang/pkg/opentelemetry/metric/stdout"
	"github.com/kaydxh/golang/pkg/opentelemetry/resource"
	jaeger_ "github.com/kaydxh/golang/pkg/opentelemetry/tracer/jaeger"
	otlptrace_ "github.com/kaydxh/golang/pkg/opentelemetry/tracer/otlp"
	stdouttrace_ "github.com/kaydxh/golang/pkg/opentelemetry/tracer/stdout"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Proto OpenTelemetry
	opts  struct {
		// If set, overrides params below
		viper                       *viper.Viper
		resourceStatsServiceOptions []resource.ResourceStatsServiceOption
		// ginRouter is used to register /metrics endpoint
		ginRouter gin.IRouter
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

	var openTelemetryOpts []OpenTelemetryServiceOption

	// Install resource first (K8s attributes, service info, etc.)
	resOpts, err := c.installResourceAttributes(ctx)
	if err != nil {
		return err
	}
	openTelemetryOpts = append(openTelemetryOpts, resOpts...)

	opts, err := c.installMeter(ctx)
	if err != nil {
		return err
	}
	openTelemetryOpts = append(openTelemetryOpts, opts...)

	// Install App MeterProvider if configured
	appOpts, err := c.installAppMeter(ctx)
	if err != nil {
		return err
	}
	openTelemetryOpts = append(openTelemetryOpts, appOpts...)

	opts, err = c.installTracer(ctx)
	if err != nil {
		return err
	}
	openTelemetryOpts = append(openTelemetryOpts, opts...)

	ot := NewOpenTelemetryService(openTelemetryOpts...)
	err = ot.Install(ctx)
	if err != nil {
		return err
	}

	_, err = c.installResource(ctx)
	return err
}

func (c *completedConfig) installResourceAttributes(ctx context.Context) ([]OpenTelemetryServiceOption, error) {
	var opts []OpenTelemetryServiceOption

	resourceConfig := c.Proto.GetResource()

	// Build resource options
	var resourceOpts []resource.ResourceOption

	// Service name
	if resourceConfig.GetServiceName() != "" {
		resourceOpts = append(resourceOpts, resource.WithServiceName(resourceConfig.GetServiceName()))
	}

	// Custom attributes
	if len(resourceConfig.GetAttrs()) > 0 {
		resourceOpts = append(resourceOpts, resource.WithAttrs(resourceConfig.GetAttrs()))
	}

	// K8s attributes (enabled by default)
	// Note: After regenerating pb.go, use resourceConfig.GetK8S().GetEnabled()
	resourceOpts = append(resourceOpts, resource.WithK8s(true))

	// Create resource with K8s attributes
	res, err := resource.NewResource(resourceOpts...)
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	opts = append(opts, WithResource(res))
	return opts, nil
}

func (c *completedConfig) installMeter(ctx context.Context) ([]OpenTelemetryServiceOption, error) {

	var opts []OpenTelemetryServiceOption
	collectDuration := c.Proto.GetMetricCollectDuration().AsDuration()
	if collectDuration > 0 {
		opts = append(opts, WithMetricCollectDuration(collectDuration))
	}

	metricType := c.Proto.OtelMetricExporterType
	switch metricType {
	case OtelMetricExporterType_metric_prometheus:
		urlPath := c.Proto.GetOtelMetricExporter().GetPrometheus().GetUrl()
		if urlPath == "" {
			urlPath = "/metrics"
		}
		builder := prometheus_.NewPrometheusExporterBuilder(
			prometheus_.WithMetricUrlPath(urlPath),
		)
		opts = append(opts, WithMeterPullExporter(builder))

		// Register /metrics route if gin router is provided
		c.registerMetricsRoute(urlPath)

	case OtelMetricExporterType_metric_stdout:
		builder := stdoutmetric_.NewStdoutExporterBuilder(
			stdoutmetric_.WithPrettyPrint(c.Proto.GetOtelMetricExporter().GetStdout().GetPrettyPrint()),
		)
		opts = append(opts, WithMeterPushExporter(builder))

	case OtelMetricExporterType_metric_otlp:
		otlpConfig := c.Proto.GetOtelMetricExporter().GetOtlp()
		var otlpOpts []otlpmetric_.OTLPExporterBuilderOption

		if otlpConfig.GetEndpoint() != "" {
			otlpOpts = append(otlpOpts, otlpmetric_.WithEndpoint(otlpConfig.GetEndpoint()))
		}

		// Set protocol
		protocol := otlpConfig.GetProtocol()
		if protocol == "grpc" {
			otlpOpts = append(otlpOpts, otlpmetric_.WithProtocol(otlpmetric_.ProtocolGRPC))
		} else {
			otlpOpts = append(otlpOpts, otlpmetric_.WithProtocol(otlpmetric_.ProtocolHTTP))
		}

		// Set insecure
		otlpOpts = append(otlpOpts, otlpmetric_.WithInsecure(otlpConfig.GetInsecure()))

		// Set URL path for HTTP
		if otlpConfig.GetUrlPath() != "" {
			otlpOpts = append(otlpOpts, otlpmetric_.WithURLPath(otlpConfig.GetUrlPath()))
		}

		// Set headers (including token)
		headers := make(map[string]string)
		for k, v := range otlpConfig.GetHeaders() {
			headers[k] = v
		}
		if otlpConfig.GetToken() != "" {
			headers["Authorization"] = "Bearer " + otlpConfig.GetToken()
		}
		if len(headers) > 0 {
			otlpOpts = append(otlpOpts, otlpmetric_.WithHeaders(headers))
		}

		builder := otlpmetric_.NewOTLPExporterBuilder(otlpOpts...)
		opts = append(opts, WithMeterPushExporter(builder))

	case OtelMetricExporterType_metric_none:
		// not enable metric
		return nil, nil

	default:
		return nil, fmt.Errorf("not support the metricType[%v]", metricType.String())

	}

	return opts, nil
}

// installAppMeter installs the App MeterProvider (separate from global)
func (c *completedConfig) installAppMeter(ctx context.Context) ([]OpenTelemetryServiceOption, error) {
	appConfig := c.Proto.GetAppMeterProvider()
	if appConfig == nil || !appConfig.GetEnabled() {
		return nil, nil
	}

	var opts []OpenTelemetryServiceOption

	// Set collect duration
	collectDuration := appConfig.GetCollectDuration().AsDuration()
	if collectDuration > 0 {
		opts = append(opts, WithAppMetricCollectDuration(collectDuration))
	} else if c.Proto.GetMetricCollectDuration().AsDuration() > 0 {
		// Fallback to global collect duration
		opts = append(opts, WithAppMetricCollectDuration(c.Proto.GetMetricCollectDuration().AsDuration()))
	}

	// Set resource for app meter
	appResource := appConfig.GetResource()
	if appResource != nil && appResource.GetServiceName() != "" {
		var resourceOpts []resource.ResourceOption
		resourceOpts = append(resourceOpts, resource.WithServiceName(appResource.GetServiceName()))
		if len(appResource.GetAttrs()) > 0 {
			resourceOpts = append(resourceOpts, resource.WithAttrs(appResource.GetAttrs()))
		}
		resourceOpts = append(resourceOpts, resource.WithK8s(true))
		res, err := resource.NewResource(resourceOpts...)
		if err != nil {
			return nil, fmt.Errorf("creating app resource: %w", err)
		}
		opts = append(opts, WithAppMeterResource(res))
	}

	// Determine exporter type (use app-specific or fallback to global)
	exporterType := appConfig.GetExporterType()
	if exporterType == OtelMetricExporterType_metric_none {
		exporterType = c.Proto.OtelMetricExporterType
	}

	// Get exporter config (use app-specific or fallback to global)
	exporterConfig := appConfig.GetExporter()
	if exporterConfig == nil {
		exporterConfig = c.Proto.GetOtelMetricExporter()
	}

	switch exporterType {
	case OtelMetricExporterType_metric_prometheus:
		builder := prometheus_.NewPrometheusExporterBuilder(
			prometheus_.WithMetricUrlPath(exporterConfig.GetPrometheus().GetUrl()),
		)
		opts = append(opts, WithAppMeterPullExporter(builder))

	case OtelMetricExporterType_metric_stdout:
		builder := stdoutmetric_.NewStdoutExporterBuilder(
			stdoutmetric_.WithPrettyPrint(exporterConfig.GetStdout().GetPrettyPrint()),
		)
		opts = append(opts, WithAppMeterPushExporter(builder))

	case OtelMetricExporterType_metric_otlp:
		otlpConfig := exporterConfig.GetOtlp()
		var otlpOpts []otlpmetric_.OTLPExporterBuilderOption

		if otlpConfig.GetEndpoint() != "" {
			otlpOpts = append(otlpOpts, otlpmetric_.WithEndpoint(otlpConfig.GetEndpoint()))
		}

		protocol := otlpConfig.GetProtocol()
		if protocol == "grpc" {
			otlpOpts = append(otlpOpts, otlpmetric_.WithProtocol(otlpmetric_.ProtocolGRPC))
		} else {
			otlpOpts = append(otlpOpts, otlpmetric_.WithProtocol(otlpmetric_.ProtocolHTTP))
		}

		otlpOpts = append(otlpOpts, otlpmetric_.WithInsecure(otlpConfig.GetInsecure()))

		if otlpConfig.GetUrlPath() != "" {
			otlpOpts = append(otlpOpts, otlpmetric_.WithURLPath(otlpConfig.GetUrlPath()))
		}

		headers := make(map[string]string)
		for k, v := range otlpConfig.GetHeaders() {
			headers[k] = v
		}
		if otlpConfig.GetToken() != "" {
			headers["Authorization"] = "Bearer " + otlpConfig.GetToken()
		}
		if len(headers) > 0 {
			otlpOpts = append(otlpOpts, otlpmetric_.WithHeaders(headers))
		}

		builder := otlpmetric_.NewOTLPExporterBuilder(otlpOpts...)
		opts = append(opts, WithAppMeterPushExporter(builder))

	case OtelMetricExporterType_metric_none:
		return nil, nil

	default:
		return nil, fmt.Errorf("app meter: not support the metricType[%v]", exporterType.String())
	}

	logrus.Infof("App MeterProvider configured with exporter type: %s", exporterType.String())
	return opts, nil
}

func (c *completedConfig) installTracer(ctx context.Context) ([]OpenTelemetryServiceOption, error) {

	var opts []OpenTelemetryServiceOption
	tracerType := c.Proto.OtelTraceExporterType
	switch tracerType {
	case OtelTraceExporterType_trace_jaeger:
		builder, err := jaeger_.NewJaegerExporertBuilder(c.Proto.GetOtelTraceExporter().GetJaeger().GetUrl())
		if err != nil {
			return nil, fmt.Errorf("new jaeger exporter builder err: %v", err)
		}
		opts = append(opts, WithTracerExporter(builder))

	case OtelTraceExporterType_trace_stdout:
		builder := stdouttrace_.NewStdoutExporterBuilder(
			stdouttrace_.WithPrettyPrint(c.Proto.GetOtelTraceExporter().GetStdout().GetPrettyPrint()),
		)
		opts = append(opts, WithTracerExporter(builder))

	case OtelTraceExporterType_trace_otlp:
		otlpConfig := c.Proto.GetOtelTraceExporter().GetOtlp()
		var otlpOpts []otlptrace_.OTLPTraceExporterBuilderOption

		if otlpConfig.GetEndpoint() != "" {
			otlpOpts = append(otlpOpts, otlptrace_.WithEndpoint(otlpConfig.GetEndpoint()))
		}

		// Set protocol
		protocol := otlpConfig.GetProtocol()
		if protocol == "grpc" {
			otlpOpts = append(otlpOpts, otlptrace_.WithProtocol(otlptrace_.ProtocolGRPC))
		} else {
			otlpOpts = append(otlpOpts, otlptrace_.WithProtocol(otlptrace_.ProtocolHTTP))
		}

		// Set insecure
		otlpOpts = append(otlpOpts, otlptrace_.WithInsecure(otlpConfig.GetInsecure()))

		// Set URL path for HTTP (default: /v1/traces)
		if otlpConfig.GetUrlPath() != "" {
			otlpOpts = append(otlpOpts, otlptrace_.WithURLPath(otlpConfig.GetUrlPath()))
		}

		// Set headers (including token)
		headers := make(map[string]string)
		for k, v := range otlpConfig.GetHeaders() {
			headers[k] = v
		}
		if otlpConfig.GetToken() != "" {
			headers["Authorization"] = "Bearer " + otlpConfig.GetToken()
		}
		if len(headers) > 0 {
			otlpOpts = append(otlpOpts, otlptrace_.WithHeaders(headers))
		}

		builder := otlptrace_.NewOTLPTraceExporterBuilder(otlpOpts...)
		opts = append(opts, WithTracerExporter(builder))

	case OtelTraceExporterType_trace_none:
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

// registerMetricsRoute registers /metrics endpoint to gin router if provided.
func (c *completedConfig) registerMetricsRoute(urlPath string) {
	if c.opts.ginRouter == nil {
		logrus.Warnf("Prometheus exporter enabled but no gin router provided, /metrics endpoint will not be registered. Use WithGinRouter() option to enable automatic registration.")
		return
	}

	c.opts.ginRouter.GET(urlPath, gin.WrapH(prometheus_.GetMetricsHandler()))
	logrus.Infof("Registered Prometheus metrics handler at %s", urlPath)
}
