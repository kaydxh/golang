package resolver

import (
	"context"

	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Proto Resolver
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

func (c *completedConfig) New(ctx context.Context) (*ResolverService, error) {

	logrus.Infof("Installing Resolver")

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
	logrus.Infof("Installed Resolver")

	return rs, nil
}

func (c *completedConfig) install(ctx context.Context) (*ResolverService, error) {
	resolverInterval := c.Proto.GetResolveInterval().AsDuration()
	rs := NewResolverService(resolverInterval)

	if c.Proto.ResolverType == Resolver_resolver_type_k8s {
		for _, svc := range c.Proto.GetK8S().GetServiceNames() {
			rq, err := NewResolverQuery(
				svc,
				WithResolverType(c.Proto.GetResolverType()),
				WithLoadBalanceMode(c.Proto.GetLoadBalanceMode()),
				WithNodeGroup(c.Proto.GetK8S().GetNodeGroup()),
				WithNodeUnit(c.Proto.GetK8S().GetNodeUnit()),
			)
			if err != nil {
				return nil, err
			}
			rs.AddService(rq)
		}

	} else {

		for _, domain := range c.Proto.GetDomains() {
			rq, err := NewResolverQuery(
				domain,
				WithResolverType(c.Proto.GetResolverType()),
				WithLoadBalanceMode(c.Proto.GetLoadBalanceMode()),
			)
			if err != nil {
				return nil, err
			}
			rs.AddService(rq)
		}
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
