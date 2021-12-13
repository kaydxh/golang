package protojson

// A UnmarshalerOption sets options.
type UnmarshalerOption interface {
	apply(*Unmarshaler)
}

// EmptyUnmarshalerOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyUnmarshalerOption struct{}

func (EmptyUnmarshalerOption) apply(*Unmarshaler) {}

// UnmarshalerOptionFunc wraps a function that modifies Client into an
// implementation of the UnmarshalerOption interface.
type UnmarshalerOptionFunc func(*Unmarshaler)

func (f UnmarshalerOptionFunc) apply(do *Unmarshaler) {
	f(do)
}

// sample code for option, default for nothing to change
func _UnmarshalerOptionWithDefault() UnmarshalerOption {
	return UnmarshalerOptionFunc(func(*Unmarshaler) {
		// nothing to change
	})
}
func (o *Unmarshaler) ApplyOptions(options ...UnmarshalerOption) *Unmarshaler {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
