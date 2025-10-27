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
	"flag"
	"fmt"
	"testing"
)

var (
	keyPrefix  = flag.String("keyPrefix", "_JOB_SERVICE", "key prefix")
	deleteKeys = flag.Bool("deleteKeys", false, "delete keys")
)

// GOOS=linux  GOARCH=amd64 go test  -c redis_key_delete_test.go redis_string_test.go
//go test -v --count=1 --test.timeout=0  -test.run  TestDeletePrefixKeys
//go test -v --count=1 --test.timeout=0  -test.run  TestDeletePrefixKeys -keyPrefix="test-" -deleteKeys=true
func TestDeletePrefixKeys(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		//	KeyPrefix string
		// 是否有过期时间
		ttl   bool
		batch int64
		//	deleteKeys bool
		dump bool
	}{
		{
			//		KeyPrefix:  "_JOB_SERVICE",
			ttl:   false,
			batch: 500,
			//		deleteKeys: false,
			dump: false,
		},
	}
	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	var count int64
	var deletedKeys []string
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("case-keyPrefix[%v]-deleteKeys[%v]", *keyPrefix, *deleteKeys), func(t *testing.T) {
			iter := db.Scan(ctx, 0, fmt.Sprintf("%s*", *keyPrefix), 0).Iterator()
			if err != nil {
				t.Fatalf("failed to scan key[%v], err: %v", *keyPrefix, err)
			}

			for iter.Next(ctx) {
				key := iter.Val()
				/*
					tp, err := db.Type(ctx, key).Result()
					if err != nil {
						t.Errorf("failed to type key[%v], err: %v", key, err)
						continue
					}

					switch tp {
					case redis_.TypeString:

					}
				*/
				/*
					value, err := db.Get(ctx, key).Result()
					if err != nil {
						t.Fatalf("failed to get key[%v], err: %v", key, err)
					}
				*/
				if testCase.dump {
					dumpVal, err := db.Dump(ctx, key).Result()
					if err != nil {
						t.Fatalf("failed to Dump, err: %v", err)
					}
					//t.Logf("key[%v], dumpValue[%v], value[%v]", key, dumpVal, value)
					t.Logf("key[%v], dumpValue[%v]", key, dumpVal)
				}
				//t.Logf("get key[%v] value[%v]", key, value)

				//todo delete
				if testCase.ttl {
					d, err := db.TTL(ctx, key).Result()
					if err != nil {
						t.Fatalf("failed to get key[%v], err: %v", key, err)
					}
					if d == -1 { // -1 means no TTL
						deletedKeys = append(deletedKeys, key)
						count++
					}

				} else {
					count++
					deletedKeys = append(deletedKeys, key)
				}

				if count%testCase.batch == 0 {
					if *deleteKeys {
						if err := db.Del(ctx, deletedKeys...).Err(); err != nil {
							t.Fatalf("failed to delete keys[%v], err: %v", key, err)
						}
					}
					deletedKeys = nil
					t.Logf("key number: %v ...", count)
				}
			}
			if err := iter.Err(); err != nil {
				t.Fatalf("failed to iter err: %v", err)
			}

			if len(deletedKeys) > 0 {
				if *deleteKeys {
					if err := db.Del(ctx, deletedKeys...).Err(); err != nil {
						t.Fatalf("failed to delete keys[%v], err: %v", deletedKeys, err)
					}
				}
				deletedKeys = nil
				t.Logf("key number: %v ...", count)
			}

		})
	}
	t.Logf("all key number: %v", count)

}
