package config

import (
	"log"
)

type Config struct {
	// options
	opts struct {
		Path     string
		ErrorLog *log.Logger
	}
}

func New(options ...ConfigOption) *Config {
	c := &Config{}
	c.ApplyOptions(options...)

	return c
}

func (c *Config) logf(format string, args ...interface{}) {
	if c.opts.ErrorLog != nil {
		c.opts.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}
