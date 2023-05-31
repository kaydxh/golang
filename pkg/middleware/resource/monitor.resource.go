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

var OpentelemetryDimKeys = "opentelemetry_dim_keys"
var OpentelemetryMetricCountKeys = "opentelemetry_metric_count_keys"

type Dimension struct {
	CalleeMethod string
	Error        error
}

func ExtractAttrsWithContext(ctx context.Context) []attribute.KeyValue {
	var (
		attrs []attribute.KeyValue
		dims  []string
	)

	switch value := ctx.Value(OpentelemetryDimKeys).(type) {
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

	switch value := ctx.Value(OpentelemetryMetricCountKeys).(type) {
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
