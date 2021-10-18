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
