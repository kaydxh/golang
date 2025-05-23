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

func TestZAdd(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key   string
		Value string
		Score float64
	}{
		{
			Key:   "zset-test-1",
			Value: "values-zset-test-1-1",
			Score: 0.5,
		},
		{
			Key:   "zset-test-1",
			Value: "values-zset-test-1-2",
			Score: 50,
		},
		{
			Key:   "zset-test-1",
			Value: "values-zset-test-1-3",
			Score: 500,
		},
		{
			Key:   "zset-test-1",
			Value: "values-zset-test-1-4",
			Score: 100,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		val, err := db.ZAdd(ctx, testCase.Key, &redis.Z{
			Score:  testCase.Score,
			Member: testCase.Value,
		}).Result()
		if err != nil {
			t.Fatalf("failed to ZAdd, err: %v", err)
		}
		t.Logf("set: %v, val: %v", testCase.Key, val)
	}
}

func TestZCard(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key string
	}{
		{
			Key: "zset-test-1",
		},
		{
			Key: "zset-test-20",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		count, err := db.ZCard(ctx, testCase.Key).Result()
		if err != nil {
			t.Fatalf("failed to ZCard, err: %v", err)
		}
		t.Logf("set: %v, member count: %v", testCase.Key, count)
	}
}

// return values sorted by score
func TestZRange(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key   string
		Start int64
		Stop  int64
	}{
		{
			Key:   "zset-test-1",
			Start: 0,
			Stop:  -1,
		},
		{
			Key:   "zset-test-20",
			Start: 0,
			Stop:  -1,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		vals, err := db.ZRange(ctx, testCase.Key, testCase.Start, testCase.Stop).Result()
		if err != nil {
			t.Fatalf("failed to ZRange, err: %v", err)
		}
		t.Logf("set: %v, members: %v", testCase.Key, vals)
	}
}

func TestZRangeByScore(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key string
		Min string
		Max string
	}{
		{
			Key: "zset-test-1",
			Min: "50",
			Max: "100",
		},
		{
			Key: "zset-test-1",
			Min: "-inf",
			Max: "+inf",
		},
		{
			Key: "zset-test-1",
			Min: "-1",
			Max: "100",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		vals, err := db.ZRangeByScore(ctx, testCase.Key, &redis.ZRangeBy{
			Min: testCase.Min,
			Max: testCase.Max,
		}).Result()
		if err != nil {
			t.Fatalf("failed to ZRangeByScore, err: %v", err)
		}
		t.Logf("set: %v, ZRangeByScore: %v", testCase.Key, vals)
	}
}
