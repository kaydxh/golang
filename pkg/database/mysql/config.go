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
package mysql

import (
	"context"

	"github.com/jmoiron/sqlx"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Proto Mysql
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

func (c *completedConfig) New(ctx context.Context) (*sqlx.DB, error) {

	logrus.Infof("Installing Mysql")

	if c.completeError != nil {
		return nil, c.completeError
	}

	if !c.Proto.GetEnabled() {
		logrus.Warnf("Mysql disenabled")
		return nil, nil
	}

	sqlxDB, err := c.install(ctx)
	if err != nil {
		return nil, err
	}
	logrus.Infof("Installed Mysql")

	return sqlxDB, nil
}

func (c *completedConfig) install(ctx context.Context) (*sqlx.DB, error) {
	db := NewDB(
		DBConfig{
			Address:  c.Proto.GetAddress(),
			DataName: c.Proto.GetDbName(),
			UserName: c.Proto.GetUsername(),
			Password: c.Proto.GetPassword(),
		},
		WithMaxConnections(int(c.Proto.GetMaxConnections())),
		WithMaxIdleConnections(int(c.Proto.GetMaxIdleConnections())),
		WithDialTimeout(c.Proto.GetDialTimeout().AsDuration()),
		WithReadTimeout(c.Proto.GetReadTimeout().AsDuration()),
		WithWriteTimeout(c.Proto.GetWriteTimeout().AsDuration()),
		WithConnMaxLifetime(c.Proto.GetMaxLifeTime().AsDuration()),
		WithInterpolateParams(c.Proto.GetInterpolateParams()),
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
