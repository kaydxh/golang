package logs

import (
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	time_ "github.com/kaydxh/golang/go/time"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/sirupsen/logrus"

	"github.com/ory/viper"
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
	if c.Proto.GetFormatter() == Log_json {
		logrus.SetFormatter(&logrus.JSONFormatter{})

	} else {
		//DisableColors set true, out format:
		//time="2021-08-07 20:21:46.468" level=info msg="Installing WebHandler" func="options.(*CompletedServerRunOptions).Run()" file="options.go:59"
		//DisableColors set false, out format:
		//INFO[2021-08-07T19:53:42+08:00]options.go:59 options.(*CompletedServerRunOptions).Run() Installing WebHandler
		logrus.SetFormatter(&logrus.TextFormatter{
			//ForceQuote:       true,
			DisableColors: true,
			//DisableQuote:     true,
			//FullTimestamp:    true,
			TimestampFormat:  time_.DefaultTimeMsFormat,
			CallerPrettyfier: GenShortCallPrettyfier(),
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
		"successed to install log",
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
			Level:     Log_info,
			Formatter: Log_text,
			Filepath:  "./log/" + filepath.Base(os.Args[0]),
		},
	}

	c.ApplyOptions(options...)

	return c
}
