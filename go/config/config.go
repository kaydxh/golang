package config

import (
	"context"
	"log"
)

type Config struct {
	// options
	opts struct {
		ErrorLog *log.Logger
	}
}

func New(ctx context.Context, options ...ConfigOption) (*Config, error) {
	c := &Config{}
	c.ApplyOptions(options...)

	return c, nil
}

func (c *Config) logf(format string, args ...interface{}) {
	if c.opts.ErrorLog != nil {
		c.opts.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}
