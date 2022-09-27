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
	reflect_ "github.com/kaydxh/golang/go/reflect"
)

// only get export Fields from arg
// set tag-values for field-value to hash
func HSetStruct(ctx context.Context, db *redis.Client, key string, arg interface{}) error {

	if db == nil {
		return fmt.Errorf("found unexpected nil redis client")
	}

	tagsValues := reflect_.AllTagsValues(arg, "redis")
	_, err := db.HSet(ctx, key, tagsValues).Result()
	if err != nil {
		return fmt.Errorf("failed to HSet key: %v, err: %v", key, err)
	}

	return nil
}
