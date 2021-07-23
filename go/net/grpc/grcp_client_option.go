package grpc

// A GrpcClientOption sets options.
type GrpcClientOption interface {
	apply(*GrpcClient)
}

// EmptyGrpcClientOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyGrpcClientOption struct{}

func (EmptyGrpcClientOption) apply(*GrpcClient) {}

// GrpcClientOptionFunc wraps a function that modifies Client into an
// implementation of the GrpcClientOption interface.
type GrpcClientOptionFunc func(*GrpcClient)

func (f GrpcClientOptionFunc) apply(do *GrpcClient) {
	f(do)
}

// sample code for option, default for nothing to change
func _GrpcClientOptionWithDefault() GrpcClientOption {
	return GrpcClientOptionFunc(func(*GrpcClient) {
		// nothing to change
	})
}
func (o *GrpcClient) ApplyOptions(options ...GrpcClientOption) *GrpcClient {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
