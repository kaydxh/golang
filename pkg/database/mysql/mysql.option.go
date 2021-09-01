package mysql

import (
	"time"
)

func WithMaxConnections(maxConns int) DBOption {
	return DBOptionFunc(func(c *DB) {
		c.opts.maxConns = maxConns
	})
}

func WithMaxIdleConnections(maxIdleConns int) DBOption {
	return DBOptionFunc(func(c *DB) {
		c.opts.maxIdleConns = maxIdleConns
	})
}

func WithDialTimeout(dialTimeout time.Duration) DBOption {
	return DBOptionFunc(func(c *DB) {
		c.opts.dialTimeout = dialTimeout
	})
}

func WithReadTimeout(readTimeout time.Duration) DBOption {
	return DBOptionFunc(func(c *DB) {
		c.opts.readTimeout = readTimeout
	})
}

func WithWriteTimeout(writeTimeout time.Duration) DBOption {
	return DBOptionFunc(func(c *DB) {
		c.opts.writeTimeout = writeTimeout
	})
}

func WithConnMaxLifetime(connMaxLifetime time.Duration) DBOption {
	return DBOptionFunc(func(c *DB) {
		c.opts.connMaxLifetime = connMaxLifetime
	})
}
