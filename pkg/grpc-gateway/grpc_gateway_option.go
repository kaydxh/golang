package grpcgateway

// A GRPCGatewayOption sets options.
type GRPCGatewayOption interface {
	apply(*GRPCGateway)
}

// EmptyGRPCGatewayOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyGRPCGatewayOption struct{}

func (EmptyGRPCGatewayOption) apply(*GRPCGateway) {}

// GRPCGatewayOptionFunc wraps a function that modifies Client into an
// implementation of the GRPCGatewayOption interface.
type GRPCGatewayOptionFunc func(*GRPCGateway)

func (f GRPCGatewayOptionFunc) apply(do *GRPCGateway) {
	f(do)
}

// sample code for option, default for nothing to change
func _GRPCGatewayOptionWithDefault() GRPCGatewayOption {
	return GRPCGatewayOptionFunc(func(*GRPCGateway) {
		// nothing to change
	})
}
func (o *GRPCGateway) ApplyOptions(options ...GRPCGatewayOption) *GRPCGateway {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
