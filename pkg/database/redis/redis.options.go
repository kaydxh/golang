package redis

import "time"

func WithPoolSize(poolSize int) RedisOption {
	return RedisOptionFunc(func(c *RedisClient) {
		c.opts.poolSize = poolSize
	})
}

func WithMinIdleConnections(minIdleConns int) RedisOption {
	return RedisOptionFunc(func(c *RedisClient) {
		c.opts.minIdleConns = minIdleConns
	})
}

func WithDialTimeout(dialTimeout time.Duration) RedisOption {
	return RedisOptionFunc(func(c *RedisClient) {
		c.opts.dialTimeout = dialTimeout
	})
}

func WithReadTimeout(readTimeout time.Duration) RedisOption {
	return RedisOptionFunc(func(c *RedisClient) {
		c.opts.readTimeout = readTimeout
	})
}

func WithWriteTimeout(writeTimeout time.Duration) RedisOption {
	return RedisOptionFunc(func(c *RedisClient) {
		c.opts.writeTimeout = writeTimeout
	})
}

func WithMasterName(masterName string) RedisOption {
	return RedisOptionFunc(func(c *RedisClient) {
		c.opts.masterName = masterName
	})
}
