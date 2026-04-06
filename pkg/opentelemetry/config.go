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
	middleware_resource "github.com/kaydxh/golang/pkg/middleware/resource"
	otlpmetric_ "github.com/kaydxh/golang/pkg/opentelemetry/metric/otlp"
	prometheus_ "github.com/kaydxh/golang/pkg/opentelemetry/metric/prometheus"
	stdoutmetric_ "github.com/kaydxh/golang/pkg/opentelemetry/metric/stdout"
	"github.com/kaydxh/golang/pkg/opentelemetry/resource"
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

// New installs all OpenTelemetry components (Tracer + Meter) together
// Use this when you don't need to control the initialization order
func (c *completedConfig) New(ctx context.Context) error {
	logrus.Infof("Installing OpenTelemetry")

	if c.completeError != nil {
		return c.completeError
	}

	if !c.Proto.GetEnabled() {
		return nil
	}

	// Install tracer first
	if err := c.InstallTracer(ctx); err != nil {
		return err
	}

	// Then install meter
	if err := c.InstallMeter(ctx); err != nil {
		return err
	}

	logrus.Infof("Installed OpenTelemetry")
	return nil
}

// InstallTracer installs only the TracerProvider
// This should be called BEFORE webserver creation so that trace interceptors can use the correct TracerProvider
func (c *completedConfig) InstallTracer(ctx context.Context) error {
	if c.completeError != nil {
		return c.completeError
	}

	if !c.Proto.GetEnabled() {
		return nil
	}

	logrus.Infof("======== OpenTelemetry Tracer install starting ========")

	var openTelemetryOpts []OpenTelemetryServiceOption

	// Install resource for tracer
	resOpts, err := c.installResourceAttributesForTracer(ctx)
	if err != nil {
		return err
	}
	openTelemetryOpts = append(openTelemetryOpts, resOpts...)

	// Install tracer exporter
	opts, err := c.installTracerExporter(ctx)
	if err != nil {
		return err
	}
	openTelemetryOpts = append(openTelemetryOpts, opts...)

	if len(opts) > 0 {
		ot := NewOpenTelemetryService(openTelemetryOpts...)
		err = ot.Install(ctx)
		if err != nil {
			return err
		}
		logrus.Infof("Installed OpenTelemetry Tracer")
	}

	return nil
}

// InstallMeter installs MeterProvider (Global + App if configured)
// This can be called after InstallTracer, or independently if tracer is not needed
func (c *completedConfig) InstallMeter(ctx context.Context) error {
	if c.completeError != nil {
		return c.completeError
	}

	if !c.Proto.GetEnabled() {
		return nil
	}

	logrus.Infof("======== OpenTelemetry Meter install starting ========")

	var openTelemetryOpts []OpenTelemetryServiceOption

	// Install resource for meter only
	resOpts, err := c.installResourceAttributesForMeter(ctx)
	if err != nil {
		return err
	}
	openTelemetryOpts = append(openTelemetryOpts, resOpts...)

	// Install global meter exporter
	opts, err := c.installMeterExporter(ctx)
	if err != nil {
		return err
	}
	openTelemetryOpts = append(openTelemetryOpts, opts...)

	// Install App MeterProvider if configured
	appOpts, err := c.installAppMeterExporter(ctx)
	if err != nil {
		return err
	}
	openTelemetryOpts = append(openTelemetryOpts, appOpts...)

	ot := NewOpenTelemetryService(openTelemetryOpts...)
	err = ot.Install(ctx)
	if err != nil {
		return err
	}

	// Install resource stats service
	_, err = c.installResourceStats(ctx)
	if err != nil {
		return err
	}

	logrus.Infof("Installed OpenTelemetry Meter")
	return nil
}

// installResourceAttributesForTracer creates resource attributes for tracer
// Uses WithResource which sets both tracer and meter resource options
func (c *completedConfig) installResourceAttributesForTracer(ctx context.Context) ([]OpenTelemetryServiceOption, error) {
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

	// Set meter type as Global for ZhiYan global_app_mark selection
	resourceOpts = append(resourceOpts, resource.WithMeterType(resource.MeterTypeGlobal))

	// APM Token (Tencent Cloud APM)
	apmConfig := resourceConfig.GetApm()
	if apmConfig != nil && apmConfig.GetToken() != "" {
		resourceOpts = append(resourceOpts, resource.WithApmToken(apmConfig.GetToken()))
		logrus.Infof("APM Token configured for resource attributes")
	}

	// ZhiYan platform configuration
	zhiyanConfig := resourceConfig.GetZhiyan()
	if zhiyanConfig != nil {
		if zhiyanConfig.GetAppMark() != "" {
			resourceOpts = append(resourceOpts, resource.WithZhiYanAppMark(zhiyanConfig.GetAppMark()))
		}
		if zhiyanConfig.GetGlobalAppMark() != "" {
			resourceOpts = append(resourceOpts, resource.WithZhiYanGlobalAppMark(zhiyanConfig.GetGlobalAppMark()))
		}
		if zhiyanConfig.GetEnv() != "" {
			resourceOpts = append(resourceOpts, resource.WithZhiYanEnv(zhiyanConfig.GetEnv()))
		}
		if zhiyanConfig.GetInstanceMark() != "" {
			resourceOpts = append(resourceOpts, resource.WithZhiYanInstanceMark(zhiyanConfig.GetInstanceMark()))
		}
		if zhiyanConfig.GetZhiyanApmToken() != "" {
			resourceOpts = append(resourceOpts, resource.WithZhiYanApmToken(zhiyanConfig.GetZhiyanApmToken()))
			logrus.Infof("ZhiYan APM Token configured for trace reporting")
		}
		if zhiyanConfig.GetExpandKey() != "" {
			resourceOpts = append(resourceOpts, resource.WithZhiYanExpandKey(zhiyanConfig.GetExpandKey()))
		}
		if zhiyanConfig.GetMetricGroup() != "" {
			resourceOpts = append(resourceOpts, resource.WithZhiYanMetricGroup(zhiyanConfig.GetMetricGroup()))
			// Set metric group for middleware metrics (must be called before metrics are created)
			middleware_resource.SetMetricGroup(zhiyanConfig.GetMetricGroup())
		}
		// TODO: Uncomment after regenerating pb.go with new proto fields (data_grain, data_type)
		// if zhiyanConfig.GetDataGrain() > 0 {
		// 	resourceOpts = append(resourceOpts, resource.WithZhiYanDataGrain(int(zhiyanConfig.GetDataGrain())))
		// }
		// if zhiyanConfig.GetDataType() != "" {
		// 	resourceOpts = append(resourceOpts, resource.WithZhiYanDataType(zhiyanConfig.GetDataType()))
		// }

		// Log ZhiYan configuration
		if zhiyanConfig.GetAppMark() != "" || zhiyanConfig.GetGlobalAppMark() != "" {
			logrus.Infof("ZhiYan platform configured: app_mark=%s, global_app_mark=%s, env=%s, metric_group=%s",
				zhiyanConfig.GetAppMark(), zhiyanConfig.GetGlobalAppMark(), zhiyanConfig.GetEnv(), zhiyanConfig.GetMetricGroup())
		}
	}

	// Create resource with K8s attributes
	res, err := resource.NewResource(resourceOpts...)
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	opts = append(opts, WithResource(res))
	return opts, nil
}

// installResourceAttributesForMeter creates resource attributes for meter only
// Uses WithMeterResource to avoid setting tracer options (tracer already initialized)
func (c *completedConfig) installResourceAttributesForMeter(ctx context.Context) ([]OpenTelemetryServiceOption, error) {
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
	resourceOpts = append(resourceOpts, resource.WithK8s(true))

	// Set meter type as Global for ZhiYan global_app_mark selection
	resourceOpts = append(resourceOpts, resource.WithMeterType(resource.MeterTypeGlobal))

	// ZhiYan platform configuration
	zhiyanConfig := resourceConfig.GetZhiyan()
	if zhiyanConfig != nil {
		if zhiyanConfig.GetAppMark() != "" {
			resourceOpts = append(resourceOpts, resource.WithZhiYanAppMark(zhiyanConfig.GetAppMark()))
		}
		if zhiyanConfig.GetGlobalAppMark() != "" {
			resourceOpts = append(resourceOpts, resource.WithZhiYanGlobalAppMark(zhiyanConfig.GetGlobalAppMark()))
		}
		if zhiyanConfig.GetEnv() != "" {
			resourceOpts = append(resourceOpts, resource.WithZhiYanEnv(zhiyanConfig.GetEnv()))
		}
		if zhiyanConfig.GetInstanceMark() != "" {
			resourceOpts = append(resourceOpts, resource.WithZhiYanInstanceMark(zhiyanConfig.GetInstanceMark()))
		}
		if zhiyanConfig.GetMetricGroup() != "" {
			resourceOpts = append(resourceOpts, resource.WithZhiYanMetricGroup(zhiyanConfig.GetMetricGroup()))
		}
	}

	// Create resource with K8s attributes
	res, err := resource.NewResource(resourceOpts...)
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	// Use WithMeterResource instead of WithResource to avoid setting tracer options
	opts = append(opts, WithMeterResource(res))
	return opts, nil
}

// installMeterExporter configures the global meter exporter
func (c *completedConfig) installMeterExporter(ctx context.Context) ([]OpenTelemetryServiceOption, error) {

	var opts []OpenTelemetryServiceOption
	collectDuration := c.Proto.GetMetricCollectDuration().AsDuration()
	if collectDuration > 0 {
		opts = append(opts, WithMetricCollectDuration(collectDuration))
	}

	metricType := c.Proto.OtelMetricExporterType
	logrus.Infof("installMeter: metricType=%v (%s)", metricType, metricType.String())
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

		// Set compression (gzip)
		compression := otlpConfig.GetCompression()
		otlpOpts = append(otlpOpts, otlpmetric_.WithCompression(compression))
		logrus.Debugf("OTLP config: compression=%v (from config)", compression)

		// Set temporality (Delta for ZhiYan)
		temporalityDelta := otlpConfig.GetTemporalityDelta()
		otlpOpts = append(otlpOpts, otlpmetric_.WithTemporalityDelta(temporalityDelta))
		logrus.Debugf("OTLP config: temporality_delta=%v (from config)", temporalityDelta)

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
		// Enable export logging for OTLP exporter
		opts = append(opts, WithMeterExporterLogging("OTLP", otlpConfig.GetEndpoint()))
		logrus.Infof("OTLP metric exporter configured: endpoint=%s, protocol=%s, insecure=%v, compression=%v, temporality_delta=%v",
			otlpConfig.GetEndpoint(), otlpConfig.GetProtocol(), otlpConfig.GetInsecure(),
			otlpConfig.GetCompression(), otlpConfig.GetTemporalityDelta())

	case OtelMetricExporterType_metric_none:
		// not enable metric
		return nil, nil

	default:
		return nil, fmt.Errorf("not support the metricType[%v]", metricType.String())

	}

	return opts, nil
}

// installAppMeter installs the App MeterProvider (separate from global)
// installAppMeterExporter configures the App MeterProvider exporter (separate from global)
func (c *completedConfig) installAppMeterExporter(ctx context.Context) ([]OpenTelemetryServiceOption, error) {
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
	globalResource := c.Proto.GetResource()

	var resourceOpts []resource.ResourceOption

	// Service name
	if appResource != nil && appResource.GetServiceName() != "" {
		resourceOpts = append(resourceOpts, resource.WithServiceName(appResource.GetServiceName()))
	} else if globalResource.GetServiceName() != "" {
		resourceOpts = append(resourceOpts, resource.WithServiceName(globalResource.GetServiceName()))
	}

	// Custom attributes
	if appResource != nil && len(appResource.GetAttrs()) > 0 {
		resourceOpts = append(resourceOpts, resource.WithAttrs(appResource.GetAttrs()))
	}

	// K8s attributes
	resourceOpts = append(resourceOpts, resource.WithK8s(true))

	// Set meter type as App for ZhiYan app mark selection
	resourceOpts = append(resourceOpts, resource.WithMeterType(resource.MeterTypeApp))

	// ZhiYan configuration from global resource (for App MeterProvider)
	zhiyanConfig := globalResource.GetZhiyan()
	if zhiyanConfig != nil {
		if zhiyanConfig.GetAppMark() != "" {
			resourceOpts = append(resourceOpts, resource.WithZhiYanAppMark(zhiyanConfig.GetAppMark()))
		}
		if zhiyanConfig.GetEnv() != "" {
			resourceOpts = append(resourceOpts, resource.WithZhiYanEnv(zhiyanConfig.GetEnv()))
		}
		if zhiyanConfig.GetInstanceMark() != "" {
			resourceOpts = append(resourceOpts, resource.WithZhiYanInstanceMark(zhiyanConfig.GetInstanceMark()))
		}
		if zhiyanConfig.GetExpandKey() != "" {
			resourceOpts = append(resourceOpts, resource.WithZhiYanExpandKey(zhiyanConfig.GetExpandKey()))
		}
	}

	res, err := resource.NewResource(resourceOpts...)
	if err != nil {
		return nil, fmt.Errorf("creating app resource: %w", err)
	}
	opts = append(opts, WithAppMeterResource(res))

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
		// Enable export logging for App OTLP exporter
		opts = append(opts, WithAppMeterExporterLogging("AppOTLP", otlpConfig.GetEndpoint()))

	case OtelMetricExporterType_metric_none:
		return nil, nil

	default:
		return nil, fmt.Errorf("app meter: not support the metricType[%v]", exporterType.String())
	}

	logrus.Infof("App MeterProvider configured with exporter type: %s", exporterType.String())
	return opts, nil
}

// installTracerExporter configures the tracer exporter
func (c *completedConfig) installTracerExporter(ctx context.Context) ([]OpenTelemetryServiceOption, error) {

	var opts []OpenTelemetryServiceOption
	tracerType := c.Proto.OtelTraceExporterType
	logrus.Infof("installTracerExporter: tracerType=%s (%s)", tracerType.String(), tracerType)

	switch tracerType {
	case OtelTraceExporterType_trace_stdout:
		prettyPrint := c.Proto.GetOtelTraceExporter().GetStdout().GetPrettyPrint()
		builder := stdouttrace_.NewStdoutExporterBuilder(
			stdouttrace_.WithPrettyPrint(prettyPrint),
		)
		opts = append(opts, WithTracerExporter(builder))
		logrus.Infof("Stdout trace exporter configured: pretty_print=%v", prettyPrint)

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
		// Enable export logging for OTLP trace exporter
		opts = append(opts, WithTracerExporterLogging("OTLP", otlpConfig.GetEndpoint()))
		logrus.Infof("OTLP trace exporter configured: endpoint=%s, protocol=%s, insecure=%v, url_path=%s",
			otlpConfig.GetEndpoint(), otlpConfig.GetProtocol(), otlpConfig.GetInsecure(), otlpConfig.GetUrlPath())

	case OtelTraceExporterType_trace_none:
		logrus.Infof("Trace exporter disabled (trace_none)")
		return nil, nil

	default:
		return nil, fmt.Errorf("not support the tracerType[%v]", tracerType.String())
	}

	return opts, nil
}

// installResourceStats installs the resource stats service for monitoring
func (c *completedConfig) installResourceStats(ctx context.Context) (*resource.ResourceStatsService, error) {

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
