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
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	time_ "github.com/kaydxh/golang/go/time"
)

var redisDB RedisDB

// Default values for Redis.
const (
	DefaultMinIdleConns = 10
	DefaultPoolSize     = 10
	DefaultDialTimeout  = 5 * time.Second
	DefaultReadTimeout  = 5 * time.Second
	DefaultWriteTimeout = 5 * time.Second
	DefaultMasterName   = "mymaster"
)

type DBConfig struct {
	Addresses []string
	UserName  string
	Password  string
	DB        int
}

type RedisClient struct {
	Conf DBConfig
	db   *redis.Client

	opts struct {
		poolSize     int
		minIdleConns int
		dialTimeout  time.Duration
		readTimeout  time.Duration
		writeTimeout time.Duration
		masterName   string
	}
}

func NewRedisClient(conf DBConfig, opts ...RedisOption) *RedisClient {
	c := &RedisClient{
		Conf: conf,
	}
	c.opts.poolSize = DefaultPoolSize
	c.opts.minIdleConns = DefaultMinIdleConns
	c.opts.dialTimeout = DefaultDialTimeout
	c.opts.readTimeout = DefaultReadTimeout
	c.opts.writeTimeout = DefaultWriteTimeout
	c.opts.masterName = DefaultMasterName

	c.ApplyOptions(opts...)

	return c
}

func GetDB() *redis.Client {
	return redisDB.Load()
}

func (r *RedisClient) GetRedis(ctx context.Context) (*redis.Client, error) {
	if r.db != nil {
		return r.db, nil
	}

	if len(r.Conf.Addresses) == 0 {
		return nil, fmt.Errorf("invalid redis address")
	}

	var db *redis.Client
	if len(r.Conf.Addresses) == 1 {
		db = redis.NewClient(&redis.Options{
			Addr:         r.Conf.Addresses[0],
			Password:     r.Conf.Password,
			DB:           r.Conf.DB,
			PoolSize:     r.opts.poolSize,
			MinIdleConns: r.opts.minIdleConns,
			DialTimeout:  r.opts.dialTimeout,
			ReadTimeout:  r.opts.readTimeout,
			WriteTimeout: r.opts.writeTimeout,
		})
	}

	if len(r.Conf.Addresses) > 1 {
		db = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    r.opts.masterName,
			SentinelAddrs: r.Conf.Addresses,
			Password:      r.Conf.Password,
			DB:            r.Conf.DB,
			PoolSize:      r.opts.poolSize,
			MinIdleConns:  r.opts.minIdleConns,
			DialTimeout:   r.opts.dialTimeout,
			ReadTimeout:   r.opts.readTimeout,
			WriteTimeout:  r.opts.writeTimeout,
		})
	}
	_, err := db.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	r.db = db
	redisDB.Store(db)

	return r.db, nil
}

func (r *RedisClient) GetDatabaseUntil(
	ctx context.Context,
	maxWaitInterval time.Duration, failAfter time.Duration) (*redis.Client, error) {

	exp := time_.NewExponentialBackOff(
		time_.WithExponentialBackOffOptionMaxInterval(maxWaitInterval),
		time_.WithExponentialBackOffOptionMaxElapsedTime(failAfter),
	)
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("cancel get database: %v", ctx.Err())

		default:
			rc, err := r.GetRedis(ctx)
			if err == nil {
				return rc, nil

			} else {
				actualInterval, ok := exp.NextBackOff()
				if !ok {
					return nil, fmt.Errorf("get database fail after: %v", failAfter)
				}

				logrus.Infof("delay %v, try again", actualInterval)
				time.Sleep(actualInterval)
			}
		}
	}
}

func (r *RedisClient) Close() error {
	if r.db == nil {
		return fmt.Errorf("no redis client")
	}
	return r.db.Close()
}
