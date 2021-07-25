package webserver

import (
	"time"

	"github.com/ory/viper"
)

func WithViper(v *viper.Viper) ConfigOption {
	return ConfigOptionFunc(func(c *Config) {
		c.opts.viper = v
	})
}

func WithShutdownDelayDuration(shutdownDelayDuration time.Duration) ConfigOption {
	return ConfigOptionFunc(func(c *Config) {
		c.opts.shutdownDelayDuration = shutdownDelayDuration
	})
}
