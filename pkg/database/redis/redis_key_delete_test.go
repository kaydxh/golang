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
)

func TestDeletePrefixKeys(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		KeyPrefix string
		// 是否有过期时间
		ttl           bool
		printInterval int64
		dump          bool
	}{
		{
			KeyPrefix:     "test-",
			ttl:           false,
			printInterval: 500,
			dump:          true,
		},
	}
	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	var count int64
	for _, testCase := range testCases {
		iter := db.Scan(ctx, 0, fmt.Sprintf("%s*", testCase.KeyPrefix), 0).Iterator()
		for iter.Next(ctx) {
			key := iter.Val()
			value, err := db.Get(ctx, key).Result()
			if err != nil {
				t.Fatalf("failed to get key[%v], err: %v", key, err)
			}
			if testCase.dump {
				dumpVal, err := db.Dump(ctx, key).Result()
				if err != nil {
					t.Fatalf("failed to Dump, err: %v", err)
				}
				t.Logf("key[%v], dumpValue[%v], value[%v]", key, dumpVal, value)
			}

			//t.Logf("get key[%v] value[%v]", key, value)

			//todo delete
			if testCase.ttl {
				d, err := db.TTL(ctx, key).Result()
				if err != nil {
					t.Fatalf("failed to get key[%v], err: %v", key, err)
				}
				if d == -1 { // -1 means no TTL
					count++
					//t.Logf("no ttl key[%v]", key)
					/*
						if err := db.Del(ctx, key).Err(); err != nil {
							t.Fatalf("failed to delete key[%v], err: %v", key, err)
						}
					*/
				}

			} else {
				count++
			}

			if count%testCase.printInterval == 0 {
				t.Logf("delete key number: %v ...", count)
			}
		}
		if err := iter.Err(); err != nil {
			t.Fatalf("failed to iter err: %v", err)
		}
	}
	t.Logf("delete all key number: %v", count)

}
