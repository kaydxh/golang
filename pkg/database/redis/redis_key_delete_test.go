/*
 *Copyright (c) 2023, kaydxh
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
	"fmt"
	"testing"
	"time"
)

func TestDeletePrefixKeys(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		KeyPrefix string
	}{
		{
			KeyPrefix: "zset",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		iter := db.Scan(ctx, 0, fmt.Sprintf("%s:*", testCase.KeyPrefix), 0).Iterator()
		for iter.Next(ctx) {
			key := iter.Val()
			value, err := db.Get(ctx, key).Result()
			if err != nil {
				t.Fatalf("failed to get key[%v], err: %v", key, err)
			}
			t.Logf("get key[%v] value[%v]", key, value)

			//todo delete
			d, err := db.TTL(ctx, key).Result()
			if err != nil {
				t.Fatalf("failed to get key[%v], err: %v", key, err)
			}
			if d == -1 { // -1 means no TTL
				t.Logf("no ttl key[%v]", key)
				/*
					if err := db.Del(ctx, key).Err(); err != nil {
						t.Fatalf("failed to delete key[%v], err: %v", key, err)
					}
				*/
			}
		}
		if err := iter.Err(); err != nil {
			t.Fatalf("failed to iter err: %v", err)
		}
	}
}
