package crontab

import (
	"context"

	"github.com/go-playground/validator/v10"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	Proto Crontab

	Validator *validator.Validate
	opts      struct {
		// If set, overrides params below
		viper *viper.Viper
	}
}

type completedConfig struct {
	*Config
	completeError error
}

// CompletedConfig ...
type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

func (c *completedConfig) New(ctx context.Context) (*CrontabSerivce, error) {

	logrus.Infof("Installing Crontab")

	if c.completeError != nil {
		return nil, c.completeError
	}

	if !c.Proto.GetEnabled() {
		return nil, nil
	}

	rs, err := c.install(ctx)
	if err != nil {
		return nil, err
	}
	logrus.Infof("Installed Crontab")

	return rs, nil
}

func (c *completedConfig) install(ctx context.Context) (*CrontabSerivce, error) {
	checkInterval := c.Proto.GetCheckInterval().AsDuration()
	rs := NewCrontabSerivce(checkInterval)
	return rs, nil
}

// Validate checks Config.
func (c *completedConfig) Validate() error {
	return c.Validator.Struct(c)
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to `ApplyOptions`, do that first. It's mutating the receiver.
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

// NewConfig returns a Config struct with the default values
func NewConfig(options ...ConfigOption) *Config {
	c := &Config{}
	c.ApplyOptions(options...)

	return c
}
