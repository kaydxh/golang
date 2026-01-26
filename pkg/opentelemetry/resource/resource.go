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

// ResourceOptions holds options for creating a Resource
type ResourceOptions struct {
	ServiceName    string
	ServiceVersion string
	Attrs          map[string]string
	EnableK8s      bool
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

	// Custom attributes
	for k, v := range options.Attrs {
		attrs = append(attrs, attribute.String(k, v))
	}

	return resource.NewWithAttributes(
		semconv.SchemaURL,
		attrs...,
	), nil
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
