package resource

import (
	"fmt"

	errors_ "github.com/kaydxh/golang/go/errors"
	net_ "github.com/kaydxh/golang/go/net"
	app_ "github.com/kaydxh/golang/pkg/webserver/app"
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

type Dimension struct {
	CalleeMethod string
	Error        error
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
