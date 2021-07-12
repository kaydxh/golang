package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"

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
)

type Config struct {
	Addresses []string
	DataName  string
	UserName  string
	Password  string
	DB        int
}

type RedisClient struct {
	Conf Config
	db   *redis.Client

	opts struct {
		PoolSize     int
		MinIdleConns int
		DialTimeout  time.Duration
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
}

func NewRedisClient(conf Config, opts ...RedisOption) *RedisClient {
	c := &RedisClient{
		Conf: conf,
	}
	c.opts.PoolSize = DefaultPoolSize
	c.opts.MinIdleConns = DefaultMinIdleConns
	c.opts.DialTimeout = DefaultDialTimeout
	c.opts.ReadTimeout = DefaultReadTimeout
	c.opts.WriteTimeout = DefaultWriteTimeout

	c.ApplyOptions(opts...)

	return c
}

func GetDB() *redis.Client {
	return redisDB.Load()
}

func (r *RedisClient) GetRedis() (*redis.Client, error) {
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
			PoolSize:     r.opts.PoolSize,
			MinIdleConns: r.opts.MinIdleConns,
			//	MaxRetries:   config.MaxRetries,
			DialTimeout:  r.opts.DialTimeout,
			ReadTimeout:  r.opts.ReadTimeout,
			WriteTimeout: r.opts.WriteTimeout,
		})
	}

	if len(r.Conf.Addresses) > 1 {
		db = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    "mymaster",
			SentinelAddrs: r.Conf.Addresses,
			Password:      r.Conf.Password,
			DB:            r.Conf.DB,
			PoolSize:      r.opts.PoolSize,
			MinIdleConns:  r.opts.MinIdleConns,
			//MaxRetries:    config.MaxRetries,
			DialTimeout:  r.opts.DialTimeout,
			ReadTimeout:  r.opts.ReadTimeout,
			WriteTimeout: r.opts.WriteTimeout,
		})
	}
	_, err := db.Ping().Result()
	if err != nil {
		return nil, err
	}

	r.db = db
	redisDB.Store(db)

	return r.db, nil
}

func (r *RedisClient) GetRedisUntil(maxWaitInterval time.Duration, failAfter time.Duration) (*redis.Client, error) {

	exp := time_.NewExponentialBackOff(
		time_.WithExponentialBackOffOptionMaxInterval(maxWaitInterval),
		time_.WithExponentialBackOffOptionMaxElapsedTime(failAfter),
	)
	for {
		rc, err := r.GetRedis()
		if err == nil {
			return rc, nil

		} else {
			actualInterval, ok := exp.NextBackOff()
			if !ok {
				return nil, fmt.Errorf("get database fail after: %v", failAfter)
			}

			time.Sleep(actualInterval)
		}
	}
}

/*
func (r *RedisClient) Close() error {
	if r.rc == nil {
		return fmt.Errorf("no redis client")
	}
	return d.db.Close()
}
*/
