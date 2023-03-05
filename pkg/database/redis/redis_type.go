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
	"sync/atomic"

	"github.com/go-redis/redis/v8"
)

const (
	TypeString = "string"
	TypeHash   = "hash"
	TypeList   = "list"
	TypeSet    = "set"
	TypeZSet   = "zset"
	TypeOther  = "other"
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
