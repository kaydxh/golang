package redis

// A RedisOption sets options.
type RedisOption interface {
	apply(*RedisClient)
}

// EmptyRedisUrlOption does not alter the Redisuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyRedisOption struct{}

func (EmptyRedisOption) apply(*RedisClient) {}

// RedisOptionFunc wraps a function that modifies Redis into an
// implementation of the RedisOption interface.
type RedisOptionFunc func(*RedisClient)

func (f RedisOptionFunc) apply(do *RedisClient) {
	f(do)
}

// sample code for option, default for nothing to change
func _RedisOptionWithDefault() RedisOption {
	return RedisOptionFunc(func(*RedisClient) {
		// nothing to change
	})
}

func (o *RedisClient) ApplyOptions(options ...RedisOption) *RedisClient {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
