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
package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

// transaction, but not support rollback, ensure only the client command execed
func TestTxPipelined(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key   string
		Field string
		Incr  int64
	}{
		{
			Key:   "hset-test-1",
			Field: "redis_id",
			Incr:  1,
		},
		{
			Key:   "hset-test-2",
			Field: "redis_id",
			Incr:  2,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.TxPipelined(ctx, func(redis.Pipeliner) error {
		for _, testCase := range testCases {
			val, err := db.HIncrBy(ctx, testCase.Key, testCase.Field, testCase.Incr).Result()
			if err != nil {
				return err
			}

			t.Logf("set: %v, val: %v", testCase.Key, val)
		}

		return nil
	})

	if err != nil {
		t.Fatalf("failed to TxPipelined, err: %v", err)
	}

}
