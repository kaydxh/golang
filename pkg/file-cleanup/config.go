package filecleanup

import (
	"context"

	"github.com/go-playground/validator/v10"
	disk_ "github.com/kaydxh/golang/pkg/file-cleanup/disk"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/ory/viper"
	"github.com/sirupsen/logrus"
)

// Config ...
type Config struct {
	Proto disk_.DiskCleaner

	Validator *validator.Validate
	opts      struct {
		// If set, overrides params below
		viper       *viper.Viper
		diskOptions []disk_.DiskCleanerConfigOption
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

func (c *completedConfig) New(ctx context.Context) (*disk_.DiskCleanerSerivce, error) {

	logrus.Infof("Installing DiskCleanerSerivce")

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

func (c *completedConfig) install(ctx context.Context) (*disk_.DiskCleanerSerivce, error) {
	return disk_.NewDiskCleanerSerivce(
		c.Proto.GetDiskUsage(),
		c.Proto.GetPaths(),
		c.Proto.GetExts(),
		disk_.WithDiskCheckInterval(c.Proto.GetCheckInterval().AsDuration()),
		disk_.WithDiskBaseExpired(c.Proto.GetBaseExpired().AsDuration()),
		disk_.WithDiskMinExpired(c.Proto.GetMinExpired().AsDuration()),
	)
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
