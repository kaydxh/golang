package mysql

// A DBOption sets options.
type DBOption interface {
	apply(*DB)
}

// EmptyDBUrlOption does not alter the DBuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyDBOption struct{}

func (EmptyDBOption) apply(*DB) {}

// DBOptionFunc wraps a function that modifies DB into an
// implementation of the DBOption interface.
type DBOptionFunc func(*DB)

func (f DBOptionFunc) apply(do *DB) {
	f(do)
}

// sample code for option, default for nothing to change
func _DBOptionWithDefault() DBOption {
	return DBOptionFunc(func(*DB) {
		// nothing to change
	})
}
func (o *DB) ApplyOptions(options ...DBOption) *DB {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
