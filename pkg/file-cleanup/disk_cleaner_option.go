package filecleanup

// A DiskCleanerConfigOption sets options.
type DiskCleanerConfigOption interface {
	apply(*DiskCleanerConfig)
}

// EmptyDiskCleanerConfigOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyDiskCleanerConfigOption struct{}

func (EmptyDiskCleanerConfigOption) apply(*DiskCleanerConfig) {}

// DiskCleanerConfigOptionFunc wraps a function that modifies Client into an
// implementation of the DiskCleanerConfigOption interface.
type DiskCleanerConfigOptionFunc func(*DiskCleanerConfig)

func (f DiskCleanerConfigOptionFunc) apply(do *DiskCleanerConfig) {
	f(do)
}

// sample code for option, default for nothing to change
func _DiskCleanerConfigOptionWithDefault() DiskCleanerConfigOption {
	return DiskCleanerConfigOptionFunc(func(*DiskCleanerConfig) {
		// nothing to change
	})
}
func (o *DiskCleanerConfig) ApplyOptions(options ...DiskCleanerConfigOption) *DiskCleanerConfig {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
