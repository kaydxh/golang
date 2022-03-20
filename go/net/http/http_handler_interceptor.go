package http

import "net/http"

type HandlerInterceptors struct {
	Interceptors []func(h http.Handler) http.Handler
}

func NewHandlerInterceptors(opts ...HandlerInterceptorsOption) *HandlerInterceptors {
	handers := &HandlerInterceptors{}
	handers.ApplyOptions(opts...)
	return handers
}
