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

	"github.com/go-redis/redis/v8"
)

func GetValue(ctx context.Context, db *redis.Client, key string) ([]string, error) {

	if db == nil {
		return nil, fmt.Errorf("found unexpected nil redis client")
	}

	var values []string

	typ, err := db.Type(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get type of key: %v, err: %v", key, err)
	}

	switch typ {
	case "string":
		data, err := db.Get(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get value of key: %v, err: %v", key, err)
		}
		values = append(values, data)
	case "list":
		llen, err := db.LLen(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get len of key: %v, err: %v", key, err)
		}

		var i int64
		for i = 0; i < llen; i++ {
			data, err := db.LIndex(ctx, key, i).Result()
			if err != nil {
				return nil, fmt.Errorf("failed to get list data of key: %v, err: %v", key, err)
			}
			values = append(values, data)
		}

	case "set":
		members, err := db.SMembers(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get members of key: %v, err: %v", key, err)
		}
		values = append(values, members...)

	default:

	}

	return values, nil

}

func GetValues(ctx context.Context, db *redis.Client, keys ...string) ([][]string, error) {
	if db == nil {
		return nil, fmt.Errorf("found unexpected nil redis client")
	}

	var values [][]string

	for _, key := range keys {

		value, err := GetValue(ctx, db, key)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil

}
