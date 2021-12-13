package protojson

// A MarshalerOption sets options.
type MarshalerOption interface {
	apply(*Marshaler)
}

// EmptyMarshalerOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyMarshalerOption struct{}

func (EmptyMarshalerOption) apply(*Marshaler) {}

// MarshalerOptionFunc wraps a function that modifies Client into an
// implementation of the MarshalerOption interface.
type MarshalerOptionFunc func(*Marshaler)

func (f MarshalerOptionFunc) apply(do *Marshaler) {
	f(do)
}

// sample code for option, default for nothing to change
func _MarshalerOptionWithDefault() MarshalerOption {
	return MarshalerOptionFunc(func(*Marshaler) {
		// nothing to change
	})
}
func (o *Marshaler) ApplyOptions(options ...MarshalerOption) *Marshaler {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
