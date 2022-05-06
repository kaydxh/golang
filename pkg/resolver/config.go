package reslover

import (
	"context"

	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/ory/viper"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Proto Reslover
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

func (c *completedConfig) New(ctx context.Context) (*ResloverService, error) {

	logrus.Infof("Installing Reslover")

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
	logrus.Infof("Installed Reslover")

	return rs, nil
}

func (c *completedConfig) install(ctx context.Context) (*ResloverService, error) {
	resolverInterval := c.Proto.GetResolveInterval().AsDuration()
	rs := NewResloverService(resolverInterval)
	for _, domain := range c.Proto.GetDomains() {
		rq := ResloverQuery{
			Domain: domain,
			Opts: ResloverOptions{
				ResloverType:    c.Proto.GetResloverType(),
				LoadBalanceMode: c.Proto.GetLoadBalanceMode(),
			},
		}
		rs.AddService(rq)
	}

	return rs, nil
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
