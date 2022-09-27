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

	redis_ "github.com/kaydxh/golang/pkg/database/redis"
)

func TestHSetStruct(t *testing.T) {

	db := GetDBOrDie()

	// only get export Fields from testCases
	testCases := []struct {
		Key string `redis:"reids_key"`
		ID  int64  `redis:"redis_id"`
	}{
		{
			Key: "hset-test-1",
			ID:  1,
		},
		{
			Key: "hset-test-2",
			ID:  2,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		err := redis_.HSetStruct(ctx, db, testCase.Key, testCase)
		if err != nil {
			t.Fatalf("failed to hset, err: %v", err)
		}
		t.Logf("key: %v,  value: %v", testCase.Key, testCase)
	}
}

//get all fields by hash key
func TestHKeys(t *testing.T) {

	db := GetDBOrDie()

	testCases := []struct {
		Key string
	}{
		{
			Key: "hset-test-1",
		},
		{
			Key: "hset-test-2",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		result, err := db.HKeys(ctx, testCase.Key).Result()
		if err != nil {
			t.Fatalf("failed to HKeys, err: %v", err)
		}
		t.Logf("key: %v,  result: %v", testCase.Key, result)

	}
}

//get all field-value by hash key
func TestHGetAll(t *testing.T) {

	db := GetDBOrDie()

	testCases := []struct {
		Key string
	}{
		{
			Key: "hset-test-1",
		},
		{
			Key: "hset-test-2",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		result, err := db.HGetAll(ctx, testCase.Key).Result()
		if err != nil {
			t.Fatalf("failed to hGetAll, err: %v", err)
		}
		t.Logf("key: %v,  result: %v", testCase.Key, result)

	}
}

//get value by hash key and field
func TestHGet(t *testing.T) {

	db := GetDBOrDie()

	testCases := []struct {
		Key   string
		Field string
	}{
		{
			Key:   "hset-test-1",
			Field: "redis_id",
		},
		{
			Key:   "hset-test-2",
			Field: "redis_id",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		fieldValues, err := db.HGet(ctx, testCase.Key, testCase.Field).Result()
		if err != nil {
			t.Fatalf("failed to hGet, err: %v", err)
		}
		t.Logf("key: %v, field: %v,  result: %v", testCase.Key, testCase.Field, fieldValues)

	}
}

//exist hash key and field
func TestHExists(t *testing.T) {

	db := GetDBOrDie()

	testCases := []struct {
		Key   string
		Field string
	}{
		{
			Key:   "hset-test-1",
			Field: "redis_id",
		},
		{
			Key:   "hset-test-2",
			Field: "redis_id",
		},
		{
			Key:   "hset-test-no-exist",
			Field: "redis_id",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		ok, err := db.HExists(ctx, testCase.Key, testCase.Field).Result()
		if err != nil {
			t.Fatalf("failed to hGet, err: %v", err)
		}
		t.Logf("key: %v, field: %v,  result: %v", testCase.Key, testCase.Field, ok)

	}
}

// increment by hash key and field
func TestHIncrBy(t *testing.T) {

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

	for _, testCase := range testCases {
		// val is the result after incr
		val, err := db.HIncrBy(ctx, testCase.Key, testCase.Field, testCase.Incr).Result()
		if err != nil {
			t.Fatalf("failed to HIncrBy, err: %v", err)
		}
		t.Logf("key: %v, field: %v, val: %v", testCase.Key, testCase.Field, val)

	}
}

// get fields num by hash key
func TestHLen(t *testing.T) {

	db := GetDBOrDie()

	testCases := []struct {
		Key string
	}{
		{
			Key: "hset-test-1",
		},
		{
			Key: "hset-test-2",
		},
		{
			Key: "hset-test-no-exist",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		// count is the number of field by hash key
		count, err := db.HLen(ctx, testCase.Key).Result()
		if err != nil {
			t.Fatalf("failed to HLen, err: %v", err)
		}
		t.Logf("key: %v, field count: %v", testCase.Key, count)

	}
}

// get values by fields and  hash key
func TestHMGet(t *testing.T) {

	db := GetDBOrDie()

	testCases := []struct {
		Key   string
		Field string
	}{
		{
			Key:   "hset-test-1",
			Field: "redis_id",
		},
		{
			Key:   "hset-test-2",
			Field: "redis_id",
		},
		{
			Key:   "hset-test-no-exist",
			Field: "redis_id",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		// vals is value of field by hash key
		vals, err := db.HMGet(ctx, testCase.Key, testCase.Field).Result()
		if err != nil {
			t.Fatalf("failed to HMGet, err: %v", err)
		}
		t.Logf("key: %v, field: %v, vals: %v", testCase.Key, testCase.Field, vals)

	}
}

// set values by fields and  hash key
func TestHMSet(t *testing.T) {

	db := GetDBOrDie()

	testCases := []struct {
		Key   string
		Field string
		Value interface{}
	}{
		{
			Key:   "hset-test-1",
			Field: "redis_id",
			Value: 10,
		},
		{
			Key:   "hset-test-2",
			Field: "redis_id",
			Value: 20,
		},
		{
			Key:   "hset-test-2",
			Field: "redis_id_new",
			Value: "new field",
		},
		{
			Key:   "hset-test-no-exist",
			Field: "redis_id",
			Value: 30,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		// count is the number of field by hash key
		ok, err := db.HMSet(ctx, testCase.Key, testCase.Field, testCase.Value).Result()
		if err != nil {
			t.Fatalf("failed to HMSet, err: %v", err)
		}
		t.Logf("key: %v, field: %v, ok: %v", testCase.Key, testCase.Field, ok)

	}
}

// set values by fields and  hash key when field is not exist
func TestHSetNX(t *testing.T) {

	db := GetDBOrDie()

	testCases := []struct {
		Key   string
		Field string
		Value interface{}
	}{
		{
			Key:   "hset-test-1",
			Field: "redis_id",
			Value: 10,
		},
		{
			Key:   "hset-test-2",
			Field: "redis_id",
			Value: 20,
		},
		{
			Key:   "hset-test-2",
			Field: "redis_id_new_1",
			Value: "new field_1",
		},
		{
			Key:   "hset-test-no-exist-1",
			Field: "redis_id",
			Value: 30,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		// count is the number of field by hash key
		ok, err := db.HSetNX(ctx, testCase.Key, testCase.Field, testCase.Value).Result()
		if err != nil {
			t.Fatalf("failed to HSetNX, err: %v", err)
		}
		t.Logf("key: %v, field: %v, ok: %v", testCase.Key, testCase.Field, ok)

	}
}

// set values by fields and  hash key when field is not exist
/*
   redis_hset_test.go:397: key: hset-test-1, values: [10 hset-test-1]
   redis_hset_test.go:397: key: hset-test-2, values: [hset-test-2 20 new field new field_1]
   redis_hset_test.go:397: key: hset-test-no-exist-1, values: [30]
*/
func TestHVals(t *testing.T) {

	db := GetDBOrDie()

	testCases := []struct {
		Key string
	}{
		{
			Key: "hset-test-1",
		},
		{
			Key: "hset-test-2",
		},
		{
			Key: "hset-test-no-exist-1",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		// count is the number of field by hash key
		vals, err := db.HVals(ctx, testCase.Key).Result()
		if err != nil {
			t.Fatalf("failed to HVals, err: %v", err)
		}
		t.Logf("key: %v, values: %v", testCase.Key, vals)

	}
}
