package exec

import "time"

//timeout
func WithTimeout(timeout time.Duration) CommandBuilderOption {
	return CommandBuilderOptionFunc(func(c *CommandBuilder) {
		c.opts.Timeout = timeout
	})
}
