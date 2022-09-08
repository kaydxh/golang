package marshaler

// A JSONPbOption sets options.
type JSONPbOption interface {
	apply(*JSONPb)
}

// EmptyJSONPbOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyJSONPbOption struct{}

func (EmptyJSONPbOption) apply(*JSONPb) {}

// JSONPbOptionFunc wraps a function that modifies Client into an
// implementation of the JSONPbOption interface.
type JSONPbOptionFunc func(*JSONPb)

func (f JSONPbOptionFunc) apply(do *JSONPb) {
	f(do)
}

// sample code for option, default for nothing to change
func _JSONPbOptionWithDefault() JSONPbOption {
	return JSONPbOptionFunc(func(*JSONPb) {
		// nothing to change
	})
}
func (o *JSONPb) ApplyOptions(options ...JSONPbOption) *JSONPb {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
