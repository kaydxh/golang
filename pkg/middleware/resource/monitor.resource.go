package resource

import (
	"fmt"

	net_ "github.com/kaydxh/golang/go/net"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		var errCode codes.Code
		s, _ := status.FromError(dim.Error)
		if s != nil {
			attrs = append(attrs, ErrorCodeKey.String(fmt.Sprintf("%d:%s", int64(errCode), s.Message())))
		}
	} else {
		attrs = append(attrs, ErrorCodeKey.String(fmt.Sprintf("%d:%s", 0, "OK")))
	}

	return attrs
}
