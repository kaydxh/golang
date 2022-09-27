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
)

//set key expired time s, if time > 0, time must > 1s
func TestExpire(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key        string
		Expiration time.Duration
	}{
		{
			Key:        "zset-test-1",
			Expiration: 30 * time.Millisecond,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		ok, err := db.Expire(ctx, testCase.Key, testCase.Expiration).Result()
		if err != nil {
			t.Fatalf("failed to Expire, err: %v", err)
		}
		t.Logf("set: %v, val: %v", testCase.Key, ok)
	}
}

//set key expired time with unix timestamp
func TestExpireAt(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key string
		Tm  time.Time
	}{
		{
			Key: "zset-test-1",
			Tm:  time.Now().Add(10 * time.Second),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		ok, err := db.ExpireAt(ctx, testCase.Key, testCase.Tm).Result()
		if err != nil {
			t.Fatalf("failed to ExpireAt, err: %v", err)
		}

		if !ok {
			t.Fatalf("failed to ExpireAt")
		}
		t.Logf("set: %v, val: %v", testCase.Key, ok)
	}
}

//set key expired time, if time > 0, time must > 1ms
func TestPExpire(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key        string
		Expiration time.Duration
	}{
		{
			Key:        "zset-test-1",
			Expiration: 30 * time.Second,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		ok, err := db.PExpire(ctx, testCase.Key, testCase.Expiration).Result()
		if err != nil {
			t.Fatalf("failed to PExpire, err: %v", err)
		}
		t.Logf("set: %v, val: %v", testCase.Key, ok)
	}
}

//serialize key
func TestDump(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key string
	}{
		{
			Key: "zset-test-1",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		val, err := db.Dump(ctx, testCase.Key).Result()
		if err != nil {
			t.Fatalf("failed to Dump, err: %v", err)
		}

		t.Logf("set: %v, val: %v", testCase.Key, val)
	}
}
