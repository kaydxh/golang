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
