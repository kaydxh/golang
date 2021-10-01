package monitor

import (
	"context"

	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/ory/viper"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Proto Monitor
	opts  struct {
		// If set, overrides params below
		viper *viper.Viper
	}
}

type completedConfig struct {
	*Config
	completeError error
}

type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

func (c *completedConfig) Apply(ctx context.Context) error {

	logrus.Infof("Installing Monitor")

	if c.completeError != nil {
		return c.completeError
	}

	/*
		if !c.Proto.GetEnabled() {
			return nil
		}
	*/

	err := c.install(ctx)
	if err != nil {
		return err
	}
	logrus.Infof("Installed Monitor")

	return nil
}

func (c *completedConfig) install(ctx context.Context) error {
	return nil
}

// Complete set default ServerRunOptions.
func (c *Config) Complete() CompletedConfig {
	err := c.loadViper()
	if err != nil {
		return CompletedConfig{&completedConfig{
			Config:        c,
			completeError: err,
		}}
	}

	return CompletedConfig{&completedConfig{Config: c}}
}

func (c *Config) loadViper() error {
	if c.opts.viper != nil {
		return viper_.UnmarshalProtoMessageWithJsonPb(c.opts.viper, &c.Proto)
	}

	return nil
}

func NewConfig(options ...ConfigOption) *Config {
	c := &Config{}
	c.ApplyOptions(options...)

	return c
}
