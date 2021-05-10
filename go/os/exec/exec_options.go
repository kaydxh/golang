package exec

var _default_CommandBuilder_value = func() (val CommandBuilder) { return }()

type CommandBuilderOption interface {
	apply(*CommandBuilder)
}

// EmptyCommandBuilderOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyCommandBuilderOption struct{}

func (EmptyCommandBuilderOption) apply(*CommandBuilder) {}

// CopyDirOptionFunc wraps a function that modifies CopyDir into an
// implementation of the CopyDirOption interface.
type CommandBuilderOptionFunc func(*CommandBuilder)

func (f CommandBuilderOptionFunc) apply(do *CommandBuilder) {
	f(do)
}

// sample code for option, default for nothing to change
func _CommandBuilderOptionWithDefault() CommandBuilderOption {
	return CommandBuilderOptionFunc(func(*CommandBuilder) {
		// nothing to change
	})
}

func (o *CommandBuilder) ApplyOptions(
	options ...CommandBuilderOption,
) *CommandBuilder {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
