/*
 *Copyright (c) 2023, kaydxh
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
package binlog

import (
	"context"
	"fmt"

	ds_ "github.com/kaydxh/golang/pkg/binlog/datastore"
	dbstore_ "github.com/kaydxh/golang/pkg/binlog/datastore/dbstore"
	filestore_ "github.com/kaydxh/golang/pkg/binlog/datastore/filestore"
	mysql_ "github.com/kaydxh/golang/pkg/database/mysql"
	rotate_ "github.com/kaydxh/golang/pkg/file-rotate"
	mq_ "github.com/kaydxh/golang/pkg/mq"
	taskq_ "github.com/kaydxh/golang/pkg/pool/taskqueue"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Proto Binlog
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

func (c *completedConfig) New(ctx context.Context, taskq *taskq_.Pool, consumers []mq_.Consumer, opts ...BinlogServiceOption) (*BinlogService, error) {

	logrus.Infof("Installing BinlogService")

	if c.completeError != nil {
		return nil, c.completeError
	}

	if !c.Proto.GetEnabled() {
		logrus.Warnf("BinlogService disenabled")
		return nil, nil
	}

	sqlxDB, err := c.install(ctx, taskq, consumers, opts...)
	if err != nil {
		return nil, err
	}
	logrus.Infof("Installed BinlogService")

	return sqlxDB, nil
}

func (c *completedConfig) install(ctx context.Context, taskq *taskq_.Pool, consumers []mq_.Consumer, opts ...BinlogServiceOption) (*BinlogService, error) {
	config := &c.Proto
	var (
		dataStore ds_.DataStore
		err       error
	)
	switch config.GetBinlogType() {
	case BinlogType_BinlogType_DB:
		db := mysql_.GetDB()
		if db == nil {
			return nil, fmt.Errorf("db is not installed")
		}
		dataStore, err = dbstore_.NewDBDataStore(db)
		if err != nil {
			return nil, fmt.Errorf("db is not installed")
		}
	case BinlogType_BinlogType_File:
		filelogConfig := config.GetFileLog()

		dataStore, err = filestore_.NewFileDataStore(
			filelogConfig.GetFilepath(),
			rotate_.WithRotateSize(filelogConfig.GetRotateSize()),
			rotate_.WithRotateInterval(filelogConfig.GetRotateInterval().AsDuration()),
			//rotate_.WithRotateCallback(bs.rotateCallback),
		)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("binlog type %v is not support", config.GetBinlogType())
	}
	bs, err := NewBinlogService(dataStore, taskq, consumers, opts...)
	if err != nil {
		return nil, err
	}
	err = bs.Run(ctx)
	return bs, err
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
