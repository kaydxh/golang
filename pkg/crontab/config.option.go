package crontab

import (
	"github.com/spf13/viper"
)

func WithViper(v *viper.Viper) ConfigOption {
	return ConfigOptionFunc(func(c *Config) {
		c.opts.viper = v
	})
}
