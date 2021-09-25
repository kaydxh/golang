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
			Key: "hset-test1",
			ID:  1,
		},
		{
			Key: "hset-test2",
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
