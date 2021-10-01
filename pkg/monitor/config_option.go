package monitor

// A ConfigOption sets options.
type ConfigOption interface {
	apply(*Config)
}

// EmptyConfigOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyConfigOption struct{}

func (EmptyConfigOption) apply(*Config) {}

// ConfigOptionFunc wraps a function that modifies Client into an
// implementation of the ConfigOption interface.
type ConfigOptionFunc func(*Config)

func (f ConfigOptionFunc) apply(do *Config) {
	f(do)
}

// sample code for option, default for nothing to change
func _ConfigOptionWithDefault() ConfigOption {
	return ConfigOptionFunc(func(*Config) {
		// nothing to change
	})
}
func (o *Config) ApplyOptions(options ...ConfigOption) *Config {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
