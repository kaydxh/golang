package grpcgateway

import "google.golang.org/grpc"

type GRPCHandler interface {
	Register(srv *grpc.Server)
}

type GRPCHandlerFunc func(srv *grpc.Server)

func (h GRPCHandlerFunc) Register(srv *grpc.Server) {
	h(srv)
}
