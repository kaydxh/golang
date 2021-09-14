package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"

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

func (r *RedisClient) GetRedis() (*redis.Client, error) {
	fmt.Println("reids: ", r.Conf)
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
	_, err := db.Ping().Result()
	if err != nil {
		return nil, err
	}

	r.db = db
	redisDB.Store(db)

	return r.db, nil
}

func (r *RedisClient) GetDatabaseUntil(maxWaitInterval time.Duration, failAfter time.Duration) (*redis.Client, error) {

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
			fmt.Println("actualInterval: ", actualInterval)
		}
	}
}

func (r *RedisClient) Close() error {
	if r.db == nil {
		return fmt.Errorf("no redis client")
	}
	return r.db.Close()
}
