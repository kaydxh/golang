package filecleanup

import (
	disk_ "github.com/kaydxh/golang/pkg/file-cleanup/disk"
	"github.com/ory/viper"
)

func WithViper(v *viper.Viper) ConfigOption {
	return ConfigOptionFunc(func(c *Config) {
		c.opts.viper = v
	})
}

func WithDiskUsageCallBack(f func(diskUsage float32)) ConfigOption {
	return ConfigOptionFunc(func(c *Config) {
		c.opts.diskOptions = append(c.opts.diskOptions, disk_.WithDiskUsageCallBack(f))
	})
}

func WithCleanPostCallBack(f func(file string)) ConfigOption {
	return ConfigOptionFunc(func(c *Config) {
		c.opts.diskOptions = append(c.opts.diskOptions, disk_.WithCleanPostCallBack(f))
	})
}
