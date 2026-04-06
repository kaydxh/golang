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
package resource

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// K8s environment variable names
const (
	EnvNodeIP            = "NODE_IP"
	EnvPodNamespace      = "POD_NAMESPACE"
	EnvPodName           = "POD_NAME"
	EnvPodIP             = "POD_IP"
	EnvContainerName     = "CONTAINER_NAME"
	EnvContainerPlatform = "CONTAINER_PLATFORM"
)

// Custom attribute keys for K8s
var (
	K8SNodeIPKey            = attribute.Key("k8s.node.ip")
	K8SPodIPKey             = attribute.Key("k8s.pod.ip")
	K8SContainerPlatformKey = attribute.Key("k8s.container.platform")
)

// APM Token key for Tencent Cloud APM
// https://console.cloud.tencent.com/apm/monitor/access
var (
	ApmTokenKey = attribute.Key("token")
)

// ZhiYan platform attribute keys
var (
	// ZhiYanAppMarkKey is the app mark for metric reporting (必填，上报应用标记)
	ZhiYanAppMarkKey = attribute.Key("__zhiyan_app_mark__")

	// ZhiYanInstanceMarkKey is the instance identifier for attribute reporting (选填，上报实例标识)
	ZhiYanInstanceMarkKey = attribute.Key("__zhiyan_instance_mark__")

	// ZhiYanEnvKey is the environment for attribute reporting
	ZhiYanEnvKey = attribute.Key("__zhiyan_env__")

	// ZhiYanExpandKey controls whether to expand resource attributes to metric dimensions
	ZhiYanExpandKey = attribute.Key("__zhiyan_expand_tag_enable__")

	// ZhiYanDataGrainKey is the data granularity (选填，数据粒度，int类型，默认60，可接受10,30,60)
	ZhiYanDataGrainKey = attribute.Key("__zhiyan_data_grain__")

	// ZhiYanDataTypeKey is the data type (选填，秒级粒度数据时填写"second")
	ZhiYanDataTypeKey = attribute.Key("__zhiyan_data_type__")

	// ZhiYanTpsTenantIDKey is the tenant ID for ZhiYan APM trace reporting
	// Format: "空间ID#日志租户#监控宝租户"
	ZhiYanTpsTenantIDKey = attribute.Key("tps.tenant.id")
)

// ResourceOptions holds options for creating a Resource
type ResourceOptions struct {
	ServiceName    string
	ServiceVersion string
	Attrs          map[string]string
	EnableK8s      bool
	ApmToken       string // APM Token for Tencent Cloud APM

	// ZhiYan platform options
	ZhiYanAppMark       string // App mark for business metric reporting
	ZhiYanGlobalAppMark string // Global app mark for infrastructure metrics
	ZhiYanEnv           string // Environment (prod/test/dev)
	ZhiYanInstanceMark  string // Instance identifier
	ZhiYanApmToken      string // APM Token for ZhiYan trace reporting
	ZhiYanExpandKey     string // Expand resource attrs to dimensions (yes/no)
	ZhiYanMetricGroup   string // Metric group (scope name) for ZhiYan
	ZhiYanDataGrain     int    // Data granularity: 10, 30, or 60 (default 60, minutes)
	ZhiYanDataType      string // Data type: "second" for sub-minute granularity

	// MeterType indicates whether this is for global or app metrics
	// Used to determine which ZhiYan app mark to use
	MeterType string
}

// ResourceOption is a function that configures ResourceOptions
type ResourceOption func(*ResourceOptions)

// WithServiceName sets the service name
func WithServiceName(name string) ResourceOption {
	return func(o *ResourceOptions) {
		o.ServiceName = name
	}
}

// WithServiceVersion sets the service version
func WithServiceVersion(version string) ResourceOption {
	return func(o *ResourceOptions) {
		o.ServiceVersion = version
	}
}

// WithAttrs sets additional attributes
func WithAttrs(attrs map[string]string) ResourceOption {
	return func(o *ResourceOptions) {
		o.Attrs = attrs
	}
}

// WithK8s enables K8s attribute detection
func WithK8s(enable bool) ResourceOption {
	return func(o *ResourceOptions) {
		o.EnableK8s = enable
	}
}

// WithApmToken sets the APM Token for Tencent Cloud APM
func WithApmToken(token string) ResourceOption {
	return func(o *ResourceOptions) {
		o.ApmToken = token
	}
}

// ZhiYan platform options

// WithZhiYanAppMark sets the app mark for business metric reporting
func WithZhiYanAppMark(appMark string) ResourceOption {
	return func(o *ResourceOptions) {
		o.ZhiYanAppMark = appMark
	}
}

// WithZhiYanGlobalAppMark sets the global app mark for infrastructure metrics
func WithZhiYanGlobalAppMark(globalAppMark string) ResourceOption {
	return func(o *ResourceOptions) {
		o.ZhiYanGlobalAppMark = globalAppMark
	}
}

// WithZhiYanEnv sets the ZhiYan environment
func WithZhiYanEnv(env string) ResourceOption {
	return func(o *ResourceOptions) {
		o.ZhiYanEnv = env
	}
}

// WithZhiYanInstanceMark sets the instance identifier
func WithZhiYanInstanceMark(instanceMark string) ResourceOption {
	return func(o *ResourceOptions) {
		o.ZhiYanInstanceMark = instanceMark
	}
}

// WithZhiYanApmToken sets the APM token for ZhiYan trace reporting
func WithZhiYanApmToken(token string) ResourceOption {
	return func(o *ResourceOptions) {
		o.ZhiYanApmToken = token
	}
}

// WithZhiYanExpandKey sets whether to expand resource attrs to dimensions
func WithZhiYanExpandKey(expand string) ResourceOption {
	return func(o *ResourceOptions) {
		o.ZhiYanExpandKey = expand
	}
}

// WithZhiYanMetricGroup sets the metric group (scope name) for ZhiYan
// Common values: "default", "client_report", "server_report"
func WithZhiYanMetricGroup(group string) ResourceOption {
	return func(o *ResourceOptions) {
		o.ZhiYanMetricGroup = group
	}
}

// WithZhiYanDataGrain sets the data granularity for ZhiYan
// Valid values: 10, 30, 60 (default 60, minutes)
func WithZhiYanDataGrain(grain int) ResourceOption {
	return func(o *ResourceOptions) {
		o.ZhiYanDataGrain = grain
	}
}

// WithZhiYanDataType sets the data type for ZhiYan
// Set to "second" for sub-minute granularity data
func WithZhiYanDataType(dataType string) ResourceOption {
	return func(o *ResourceOptions) {
		o.ZhiYanDataType = dataType
	}
}

// WithMeterType sets the meter type (global or app)
func WithMeterType(meterType string) ResourceOption {
	return func(o *ResourceOptions) {
		o.MeterType = meterType
	}
}

// MeterType constants
const (
	MeterTypeGlobal = "global"
	MeterTypeApp    = "app"
)

// NewResource creates a new OpenTelemetry Resource with the given options
func NewResource(opts ...ResourceOption) (*resource.Resource, error) {
	options := &ResourceOptions{
		EnableK8s: true, // Enable K8s by default
	}
	for _, opt := range opts {
		opt(options)
	}

	// Base attributes
	attrs := []attribute.KeyValue{}

	// Service name (default to process name)
	serviceName := options.ServiceName
	if serviceName == "" {
		serviceName = filepath.Base(os.Args[0])
	}
	attrs = append(attrs, semconv.ServiceName(serviceName))

	// Service version
	if options.ServiceVersion != "" {
		attrs = append(attrs, semconv.ServiceVersion(options.ServiceVersion))
	}

	// K8s attributes from environment variables
	if options.EnableK8s {
		k8sAttrs := GetK8sAttributes()
		attrs = append(attrs, k8sAttrs...)
	}

	// APM Token (Tencent Cloud APM)
	if options.ApmToken != "" {
		attrs = append(attrs, ApmTokenKey.String(options.ApmToken))
	}

	// ZhiYan platform attributes
	attrs = append(attrs, getZhiYanAttributes(options)...)

	// Custom attributes
	for k, v := range options.Attrs {
		attrs = append(attrs, attribute.String(k, v))
	}

	return resource.NewWithAttributes(
		semconv.SchemaURL,
		attrs...,
	), nil
}

// getZhiYanAttributes returns ZhiYan-specific resource attributes
func getZhiYanAttributes(options *ResourceOptions) []attribute.KeyValue {
	var attrs []attribute.KeyValue

	// Determine which app mark to use based on meter type
	appMark := ""
	if options.MeterType == MeterTypeApp && options.ZhiYanAppMark != "" {
		appMark = options.ZhiYanAppMark
	} else if options.MeterType == MeterTypeGlobal && options.ZhiYanGlobalAppMark != "" {
		appMark = options.ZhiYanGlobalAppMark
	}

	// Debug log for ZhiYan attributes selection
	logrus.Debugf("getZhiYanAttributes: MeterType=%s, ZhiYanAppMark=%s, ZhiYanGlobalAppMark=%s, selected appMark=%s",
		options.MeterType, options.ZhiYanAppMark, options.ZhiYanGlobalAppMark, appMark)

	// Only add ZhiYan attributes if app mark is set
	if appMark != "" {
		attrs = append(attrs, ZhiYanAppMarkKey.String(appMark))
		logrus.Infof("ZhiYan resource attribute added: %s=%s", ZhiYanAppMarkKey, appMark)

		// Environment (default to "prod" if not set)
		env := options.ZhiYanEnv
		if env == "" {
			env = "prod"
		}
		attrs = append(attrs, ZhiYanEnvKey.String(env))

		// Instance mark
		if options.ZhiYanInstanceMark != "" {
			attrs = append(attrs, ZhiYanInstanceMarkKey.String(options.ZhiYanInstanceMark))
		}

		// Expand key (default to "no")
		expandKey := options.ZhiYanExpandKey
		if expandKey != "yes" {
			expandKey = "no"
		}
		attrs = append(attrs, ZhiYanExpandKey.String(expandKey))

		// Data grain (选填，数据粒度，默认60)
		if options.ZhiYanDataGrain > 0 {
			attrs = append(attrs, ZhiYanDataGrainKey.Int(options.ZhiYanDataGrain))
		}

		// Data type (选填，秒级粒度数据时填写"second")
		if options.ZhiYanDataType != "" {
			attrs = append(attrs, ZhiYanDataTypeKey.String(options.ZhiYanDataType))
		}
	}

	// ZhiYan APM Token for trace reporting
	if options.ZhiYanApmToken != "" {
		attrs = append(attrs, ZhiYanTpsTenantIDKey.String(options.ZhiYanApmToken))
		// Also add service namespace for ZhiYan APM
		env := options.ZhiYanEnv
		if env == "" {
			env = "prod"
		}
		attrs = append(attrs, semconv.ServiceNamespace(env))
	}

	return attrs
}

// GetK8sAttributes returns K8s-related attributes from environment variables
func GetK8sAttributes() []attribute.KeyValue {
	var attrs []attribute.KeyValue

	// Node IP
	if nodeIP := os.Getenv(EnvNodeIP); nodeIP != "" {
		attrs = append(attrs, K8SNodeIPKey.String(nodeIP))
	}

	// Pod Namespace
	if podNamespace := os.Getenv(EnvPodNamespace); podNamespace != "" {
		attrs = append(attrs, semconv.K8SNamespaceName(podNamespace))
	}

	// Pod Name
	if podName := os.Getenv(EnvPodName); podName != "" {
		attrs = append(attrs, semconv.K8SPodName(podName))
	}

	// Pod IP
	if podIP := os.Getenv(EnvPodIP); podIP != "" {
		attrs = append(attrs, K8SPodIPKey.String(podIP))
	}

	// Container Name
	if containerName := os.Getenv(EnvContainerName); containerName != "" {
		attrs = append(attrs, semconv.K8SContainerName(containerName))
	}

	// Container Platform (e.g., STKE, TKE)
	if containerPlatform := os.Getenv(EnvContainerPlatform); containerPlatform != "" {
		attrs = append(attrs, K8SContainerPlatformKey.String(containerPlatform))
	}

	return attrs
}

// IsInK8s checks if the application is running in a K8s environment
func IsInK8s() bool {
	// Check common K8s environment variables
	return os.Getenv(EnvPodName) != "" || os.Getenv(EnvPodNamespace) != ""
}
