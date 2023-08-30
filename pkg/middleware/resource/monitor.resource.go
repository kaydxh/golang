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
	"go.opentelemetry.io/otel/metric"
)

// dims
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
	var attrs []attribute.KeyValue

	values, ok := ctx.Value(OpentelemetryDimsKey).(map[string]interface{})
	if !ok {
		return nil
	}
	for key, value := range values {
		var dimKey = attribute.Key(key)
		switch value := value.(type) {
		case string:
			attrs = append(attrs, dimKey.String(value))

		case int:
			attrs = append(attrs, dimKey.Int(value))
		case int32:
			attrs = append(attrs, dimKey.Int(int(value)))
		case int64:
			attrs = append(attrs, dimKey.Int64(int64(value)))

		case uint:
			attrs = append(attrs, dimKey.Int(int(value)))
		case uint32:
			attrs = append(attrs, dimKey.Int(int(value)))
		case uint64:
			attrs = append(attrs, dimKey.Int64(int64(value)))

		case float32:
			attrs = append(attrs, dimKey.Float64(float64(value)))
		case float64:
			attrs = append(attrs, dimKey.Float64(value))

		case bool:
			attrs = append(attrs, dimKey.Bool(value))
		}
	}

	return attrs
}

func ReportBusinessMetric(ctx context.Context, attrs []attribute.KeyValue) {
	values, ok := ctx.Value(OpentelemetryMetricsKey).(map[string]interface{})
	if !ok {
		return
	}

	for key, value := range values {

		var (
			n           int64
			f           float64
			counterType bool
		)
		switch value := value.(type) {
		case int:
			n = int64(value)
			counterType = true

		case int32:
			n = int64(value)
			counterType = true
		case int64:
			n = int64(value)
			counterType = true

		case uint:
			n = int64(value)
			counterType = true
		case uint32:
			n = int64(value)
			counterType = true
		case uint64:
			n = int64(value)
			counterType = true

		case float32:
			f = float64(value)
		case float64:
			f = float64(value)
		}

		if counterType {
			counter, err := DefaultMetricMonitor.GetOrNewBusinessCounter(key)
			if err != nil {
				otel.Handle(err)
				continue
			}
			counter.Add(ctx, n, metric.WithAttributes(attrs...))

		} else {

			histogram, err := DefaultMetricMonitor.GetOrNewBusinessHistogram(key)
			if err != nil {
				otel.Handle(err)
				continue
			}
			histogram.Record(ctx, f, metric.WithAttributes(attrs...))
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
