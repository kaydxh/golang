package config

func WithConfigOptionPath(path string) ConfigOption {
	return ConfigOptionFunc(func(opt *Config) {
		opt.opts.Path = path
	})
}
