package tcloud

import (
	reflect_ "github.com/kaydxh/golang/go/reflect"
)

// fieldNameRequestId is variable name of API 3.0
var fieldNameRequestId = "RequestId"

type RequestIdRetriever interface {
	GetRequestId() string
}

func retrieveRequestId(req interface{}) string {
	if req == nil {
		return ""
	}
	// if defined GetRequestId()string method, use it get requestId
	if req, ok := req.(RequestIdRetriever); ok {
		if req != nil {
			return req.GetRequestId()
		}
	}
	return reflect_.RetrieveStructField(req, fieldNameRequestId)
}

func trySetRequestId(req interface{}, id string) {
	reflect_.TrySetStructFiled(req, fieldNameRequestId, id)
}
