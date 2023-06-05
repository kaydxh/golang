package resource

import (
	"context"
	"fmt"

	context_ "github.com/kaydxh/golang/go/context"
	errors_ "github.com/kaydxh/golang/go/errors"
	net_ "github.com/kaydxh/golang/go/net"
	app_ "github.com/kaydxh/golang/pkg/webserver/app"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

//dims
var (
	CallerMethodKey = attribute.Key("caller_method") // caller method
	CalleeMethodKey = attribute.Key("callee_method") // callee method
	PodIpKey        = attribute.Key("pod_ip")        // pod ip
	ServerNameKey   = attribute.Key("server_name")   // server name
	ErrorCodeKey    = attribute.Key("error_code")    // error code
)

var OpentelemetryDimsKey = "opentelemetry_dims_key"
var OpentelemetryMetricsKey = "opentelemetry_metrics_key"

type Dimension struct {
	CalleeMethod string
	Error        error
}

type AddKeyContexFunc func(ctx context.Context) context.Context

func AddKeysContext(ctx context.Context, ops ...AddKeyContexFunc) context.Context {
	for _, op := range ops {
		ctx = op(ctx)
	}
	return ctx
}

func AddAttrKeysContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, OpentelemetryDimsKey, map[string]interface{}{})
}

func UpdateAttrsContext(ctx context.Context, values map[string]interface{}) error {
	return context_.UpdateContext(ctx, OpentelemetryDimsKey, values)
}

func AddMetricKeysContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, OpentelemetryMetricsKey, map[string]interface{}{})
}

func UpdateMetricContext(ctx context.Context, values map[string]interface{}) error {
	return context_.UpdateContext(ctx, OpentelemetryMetricsKey, values)
}

func AppendAttrsContext(ctx context.Context, values ...string) context.Context {
	return context_.AppendContext(ctx, OpentelemetryDimsKey, values...)
}

func AppendMetricCountContext(ctx context.Context, values ...string) context.Context {
	return context_.AppendContext(ctx, OpentelemetryMetricsKey, values...)
}

func ExtractAttrsWithContext(ctx context.Context) []attribute.KeyValue {
	var (
		attrs []attribute.KeyValue
		dims  []string
	)

	switch value := ctx.Value(OpentelemetryDimsKey).(type) {
	case string:
		dims = append(dims, value)

	case []string:
		dims = append(dims, value...)
	}

	for _, dim := range dims {
		var dimKey = attribute.Key(dim)
		sv := context_.ExtractStringFromContext(ctx, dim)
		if sv != "" {
			attrs = append(attrs, dimKey.String(sv))
			continue
		}

		iv, err := context_.ExtractIntegerFromContext(ctx, dim)
		if err == nil {
			attrs = append(attrs, dimKey.Int64(iv))
			continue
		}
	}

	return attrs
}

func ReportBusinessMetric(ctx context.Context, attrs []attribute.KeyValue) {
	var metrics []string

	switch value := ctx.Value(OpentelemetryMetricsKey).(type) {
	case string:
		metrics = append(metrics, value)

	case []string:
		metrics = append(metrics, value...)
	}

	for _, metric := range metrics {
		counter, err := DefaultMetricMonitor.GetOrNewBusinessCounter(metric)
		if err != nil {
			otel.Handle(err)
			continue
		}

		iv, err := context_.ExtractIntegerFromContext(ctx, metric)
		if err == nil {
			counter.Add(ctx, iv, attrs...)
			continue
		}
	}

	return
}

func Attrs(dim Dimension) []attribute.KeyValue {
	var attrs []attribute.KeyValue
	hostIP, err := net_.GetHostIP()
	if err == nil && hostIP.String() != "" {
		attrs = append(attrs, PodIpKey.String(hostIP.String()))
	}
	if dim.CalleeMethod != "" {
		attrs = append(attrs, CalleeMethodKey.String(dim.CalleeMethod))
	}

	if dim.Error != nil {
		errorCode := int64(errors_.ErrorToCode(dim.Error))
		message := errors_.ErrorToString(dim.Error)
		attrs = append(attrs, ErrorCodeKey.String(fmt.Sprintf("%d:%s", errorCode, message)))
	} else {
		attrs = append(attrs, ErrorCodeKey.String("0:OK"))
	}

	appName := app_.GetVersion().AppName
	if appName != "" {
		attrs = append(attrs, ServerNameKey.String(appName))
	}

	return attrs
}
