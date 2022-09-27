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
package logs

import (
	"os"
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
	time_ "github.com/kaydxh/golang/go/time"
	logrus_ "github.com/kaydxh/golang/pkg/logs/logrus"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/spf13/viper"
)

const (
	DefaultMaxAge         = 72 * time.Hour
	DefaultMaxCount       = 72
	DefaultRotateInterval = time.Hour
	//100MB
	DefaultRotateSize = 104857600
)

type Config struct {
	Proto     Log
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

type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// Validate checks Config.
func (c *completedConfig) Validate() error {
	return c.Validator.Struct(c)
}

func (c *completedConfig) Apply() error {
	if c.completeError != nil {
		return c.completeError
	}
	if c.Validator == nil {
		c.Validator = validator.New()
	}

	return c.install()
}

func (c *completedConfig) install() error {
	logrus.Infof("Installing Logs")

	switch c.Proto.GetFormatter() {
	case Log_json:
		logrus.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: GenShortCallPrettyfier(),
		})

	case Log_text:
		//DisableColors set true, out format:
		//time="2021-08-07 20:21:46.468" level=info msg="Installing WebHandler" func="options.(*CompletedServerRunOptions).Run()" file="options.go:59"
		//DisableColors set false, out format:
		//INFO[2021-08-07T19:53:42+08:00]options.go:59 options.(*CompletedServerRunOptions).Run() Installing WebHandler

		//DisableQuote: ture, this config can format newline characters(\n) in the log message
		logrus.SetFormatter(&logrus.TextFormatter{
			//ForceQuote:       true,
			DisableColors:    true,
			DisableQuote:     true,
			FullTimestamp:    true,
			TimestampFormat:  time_.DefaultTimeMillFormat,
			CallerPrettyfier: GenShortCallPrettyfier(),
		})

	default:
		logrus.SetFormatter(&logrus_.GlogFormatter{
			//ForceQuote:       true,
			DisableColors:     true,
			DisableQuote:      true,
			FullTimestamp:     true,
			TimestampFormat:   time_.DefaultTimeMillFormat,
			CallerPrettyfier:  GenShortCallPrettyfier(),
			EnableGoroutineId: c.Proto.GetEnableGoroutineId(),
		})
	}

	level, err := logrus.ParseLevel(c.Proto.GetLevel().String())
	if err != nil {
		return err
	}
	logrus.SetLevel(level)
	logrus.SetReportCaller(c.Proto.GetReportCaller())

	err = WithRotate(
		logrus.StandardLogger(),
		c.Proto.GetFilepath(),
		c.Proto.GetRedirct(),
		WithMaxAge(c.Proto.GetMaxAge().AsDuration()),
		WithMaxCount(c.Proto.GetMaxCount()),
		WithRotateSize(c.Proto.GetRotateSize()),
		WithRotateInterval(c.Proto.GetRotateInterval().AsDuration()),
		WithPrefixName(filepath.Base(os.Args[0])),
		WithSuffixName(".log"),
	)
	if err != nil {
		return err
	}

	logrus.WithField(
		"path",
		c.Proto.GetFilepath(),
	).WithField(
		"rotate_interval", c.Proto.GetRotateInterval().AsDuration(),
	).WithField(
		"rotate_size", c.Proto.GetRotateSize(),
	).Infof(
		"Installed log",
	)

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
	c.parseViper()
	return CompletedConfig{&completedConfig{Config: c}}
}

func (c *Config) parseViper() {

}

func (c *Config) loadViper() error {
	if c.opts.viper != nil {
		return viper_.UnmarshalProtoMessageWithJsonPb(c.opts.viper, &c.Proto)
	}

	return nil
}

func NewConfig(options ...ConfigOption) *Config {
	c := &Config{
		Proto: Log{
			Level:          Log_info,
			Formatter:      Log_text,
			MaxAge:         durationpb.New(DefaultMaxAge),
			MaxCount:       DefaultMaxCount,
			RotateInterval: durationpb.New(DefaultRotateInterval),
			RotateSize:     DefaultRotateSize,

			Filepath: "./log/" + filepath.Base(os.Args[0]),
		},
	}

	c.ApplyOptions(options...)

	return c
}
