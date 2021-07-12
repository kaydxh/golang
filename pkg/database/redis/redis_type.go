package redis

import (
	"sync/atomic"

	"github.com/go-redis/redis"
)

type RedisDB atomic.Value

//check type RedisDB* whether to implement interface of atomic.Value
//var _ atomic.Value = (*RedisDB)(nil)

//check type  RedisDB whether to implement interface of atomic.Value
//var _ atomic.Value = RedisDB{}

func (m *RedisDB) Store(value *redis.Client) {
	(*atomic.Value)(m).Store(value)
}

func (m *RedisDB) Load() *redis.Client {
	value := (*atomic.Value)(m).Load()
	if value == nil {
		return nil
	}
	return value.(*redis.Client)
}
