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
package taskqueue

import (
	"context"
	"fmt"

	redis_ "github.com/kaydxh/golang/pkg/database/redis"
	queue_ "github.com/kaydxh/golang/pkg/pool/taskqueue/queue"
	redisq_ "github.com/kaydxh/golang/pkg/pool/taskqueue/queue/redis"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Proto TaskQueue
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

func (c *completedConfig) New(ctx context.Context, opts ...PoolOption) (*Pool, error) {

	logrus.Infof("Installing TaskQueue")

	if c.completeError != nil {
		return nil, c.completeError
	}

	if !c.Proto.GetEnabled() {
		logrus.Warnf("TaskQueue disenabled")
		return nil, nil
	}

	rs, err := c.install(ctx, opts...)
	if err != nil {
		return nil, err
	}
	logrus.Infof("Installed TaskQueue")

	return rs, nil
}

func (c *completedConfig) install(ctx context.Context, opts ...PoolOption) (*Pool, error) {

	config := &c.Proto

	var queue queue_.Queue
	switch config.GetQueueType() {
	case TaskQueue_queue_type_redis:
		db := redis_.GetDB()
		if db == nil {
			return nil, fmt.Errorf("db is not first installed")
		}
		queue = redisq_.NewQueue(db, queue_.QueueOptions{Name: "redis"})

	default:
		return nil, fmt.Errorf("queue type %v is not support", config.GetQueueType())
	}

	pool := NewPool(queue,
		WithFetcherBurst(config.GetFetcherBurst()),
		WithWorkerBurst(config.GetWorkerBurst()),
		WithWorkTimeout(config.GetWorkTimeout().AsDuration()),
		WithFetchTimeout(config.GetFetchTimeout().AsDuration()),
		WithResultExpired(config.GetResultExpired().AsDuration()),
	)
	pool.ApplyOptions(opts...)

	err := pool.Consume(ctx)
	if err != nil {
		logrus.Errorf("failed to start consume task, err: %v", err)
		return nil, err
	}

	return pool, nil
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
