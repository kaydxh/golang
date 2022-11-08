package s3

import (
	"context"

	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gocloud.dev/blob"
)

type Config struct {
	Proto S3
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

func (c *completedConfig) New(ctx context.Context) (*blob.Bucket, error) {

	logrus.Infof("Installing S3")

	if c.completeError != nil {
		return nil, c.completeError
	}

	/*
		if !c.Proto.GetEnabled() {
			return nil, nil
		}
	*/

	bucket, err := c.install(ctx)
	if err != nil {
		return nil, err
	}
	logrus.Infof("Installed S3")

	return bucket, nil
}

func (c *completedConfig) install(ctx context.Context) (*blob.Bucket, error) {
	s3Config := &c.Proto
	storageConfig, err := ParseUrl(s3Config.Url)
	if err != nil {
		return nil, err
	}
	storageConfig.SecretId = s3Config.GetSecretId()
	storageConfig.SecretKey = s3Config.GetSecretKey()
	s, err := NewStorage(
		ctx,
		storageConfig,
	)

	return s.bucket, err
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
