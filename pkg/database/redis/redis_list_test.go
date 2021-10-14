package redis_test

import (
	"context"
	"testing"
	"time"
)

//push value to list key
func TestRPush(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key   string
		Value string
	}{
		{
			Key:   "list-test-1",
			Value: "value-list-test-1",
		},
		{
			Key:   "list-test-2",
			Value: "value-list-test-2",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		val, err := db.RPush(ctx, testCase.Key, testCase.Value).Result()
		if err != nil {
			t.Fatalf("failed to RPush, err: %v", err)
		}
		t.Logf("key: %v, val: %v", testCase.Key, val)

	}

}

//push value to list key which is existed
func TestRPushX(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key   string
		Value string
	}{
		{
			Key:   "list-test-3",
			Value: "value-list-test-1",
		},
		{
			Key:   "list-test-4",
			Value: "value-list-test-2",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		val, err := db.RPushX(ctx, testCase.Key, testCase.Value).Result()
		if err != nil {
			t.Fatalf("failed to RPushX, err: %v", err)
		}
		t.Logf("key: %v, val: %v", testCase.Key, val)

	}
}

//get length for list key
func TestLLen(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key string
	}{
		{
			Key: "list-test-3",
		},
		{
			Key: "list-test-4",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		length, err := db.LLen(ctx, testCase.Key).Result()
		if err != nil {
			t.Fatalf("failed to LLen, err: %v", err)
		}
		t.Logf("list: %v, val: %v", testCase.Key, length)

	}
}

//push value to list key
//Start: 0 the first element ,End
//End: -1 the last element, -2, next to last element
func TestLRange(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key   string
		Start int64
		End   int64
	}{
		{
			Key:   "list-test-1",
			Start: 0,
			End:   -1,
		},
		{
			Key:   "list-test-2",
			Start: 1,
			End:   -1,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		vals, err := db.LRange(ctx, testCase.Key, testCase.Start, testCase.End).Result()
		if err != nil {
			t.Fatalf("failed to LRange, err: %v", err)
		}
		t.Logf("list: %v, vals: %v", testCase.Key, vals)

	}
}
