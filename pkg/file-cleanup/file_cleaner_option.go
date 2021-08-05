package filecleanup

// A FileCleanerOption sets options.
type FileCleanerOption interface {
	apply(*FileCleaner)
}

// EmptyFileCleanerOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyFileCleanerOption struct{}

func (EmptyFileCleanerOption) apply(*FileCleaner) {}

// FileCleanerOptionFunc wraps a function that modifies Client into an
// implementation of the FileCleanerOption interface.
type FileCleanerOptionFunc func(*FileCleaner)

func (f FileCleanerOptionFunc) apply(do *FileCleaner) {
	f(do)
}

// sample code for option, default for nothing to change
func _FileCleanerOptionWithDefault() FileCleanerOption {
	return FileCleanerOptionFunc(func(*FileCleaner) {
		// nothing to change
	})
}
func (o *FileCleaner) ApplyOptions(options ...FileCleanerOption) *FileCleaner {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
