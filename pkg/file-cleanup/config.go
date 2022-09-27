/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
package filecleanup

import (
	"context"

	"github.com/go-playground/validator/v10"
	disk_ "github.com/kaydxh/golang/pkg/file-cleanup/disk"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

	opts := []disk_.DiskCleanerConfigOption{
		disk_.WithDiskCheckInterval(c.Proto.GetCheckInterval().AsDuration()),
		disk_.WithDiskBaseExpired(c.Proto.GetBaseExpired().AsDuration()),
		disk_.WithDiskMinExpired(c.Proto.GetMinExpired().AsDuration()),
	}
	opts = append(opts, c.opts.diskOptions...)

	return disk_.NewDiskCleanerSerivce(
		c.Proto.GetDiskUsage(),
		c.Proto.GetPaths(),
		c.Proto.GetExts(),
		opts...,
	)
}

// Validate checks Config.
func (c *completedConfig) Validate() error {
	return c.Validator.Struct(c)
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to `ApplyOptions`, do that first. It's mutating the receiver.
func (c *Config) Complete(options ...ConfigOption) CompletedConfig {
	err := c.loadViper()
	if err != nil {
		return CompletedConfig{&completedConfig{
			Config:        c,
			completeError: err,
		}}
	}

	c.ApplyOptions(options...)

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
