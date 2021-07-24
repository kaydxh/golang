package webserver

import (
	"time"

	"github.com/ory/viper"
)

func WithGetViper(f func() *viper.Viper) ConfigOption {
	return ConfigOptionFunc(func(c *Config) {
		c.opts.getViper = f
	})
}

func WithShutdownDelayDuration(shutdownDelayDuration time.Duration) ConfigOption {
	return ConfigOptionFunc(func(c *Config) {
		c.opts.shutdownDelayDuration = shutdownDelayDuration
	})
}
