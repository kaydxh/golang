package redis

import (
	"context"

	"github.com/go-redis/redis"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/ory/viper"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Proto Redis
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

func (c *completedConfig) New(ctx context.Context) (*redis.Client, error) {

	logrus.Infof("Installing Redis")

	if c.completeError != nil {
		return nil, c.completeError
	}

	if !c.Proto.GetEnabled() {
		return nil, nil
	}

	redisDB, err := c.install(ctx)
	if err != nil {
		return nil, err
	}
	logrus.Infof("Installed Redis")

	return redisDB, nil
}

func (c *completedConfig) install(ctx context.Context) (*redis.Client, error) {
	db := NewRedisClient(

		DBConfig{
			Addresses: c.Proto.GetAddresses(),
			UserName:  c.Proto.GetUsername(),
			Password:  c.Proto.GetPassword(),
			DB:        int(c.Proto.GetDb()),
		},
		WithPoolSize(int(c.Proto.GetPoolSize())),
		WithMinIdleConnections(int(c.Proto.GetMinIdleConns())),
		WithDialTimeout(c.Proto.GetDialTimeout().AsDuration()),
		WithReadTimeout(c.Proto.GetReadTimeout().AsDuration()),
		WithWriteTimeout(c.Proto.GetWriteTimeout().AsDuration()),
	)

	return db.GetDatabaseUntil(
		ctx,
		c.Proto.GetMaxWaitDuration().AsDuration(),
		c.Proto.GetFailAfterDuration().AsDuration(),
	)

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
