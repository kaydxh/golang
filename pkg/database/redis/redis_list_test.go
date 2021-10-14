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
			Start: 0,
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

//can use negative indexes
//-1 the last element, -2, next to last element
func TestLIndex(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key   string
		Index int64
	}{
		{
			Key:   "list-test-1",
			Index: 0,
		},
		{
			Key:   "list-test-2",
			Index: -1,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		val, err := db.LIndex(ctx, testCase.Key, testCase.Index).Result()
		if err != nil {
			t.Fatalf("failed to LRange, err: %v", err)
		}
		t.Logf("list: %v, index: %v, val: %v", testCase.Key, testCase.Index, val)

	}
}

//insert value before the pivot
func TestLInsertBefore(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key   string
		Pivot string
		Value string
	}{
		{
			Key:   "list-test-1",
			Pivot: "value-list-test-1",
			Value: "value-list-test-2",
		},
		{
			Key:   "list-test-2",
			Pivot: "value-list-test-2",
			Value: "value-list-test-3",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		val, err := db.LInsertBefore(ctx, testCase.Key, testCase.Pivot, testCase.Value).Result()
		if err != nil {
			t.Fatalf("failed to LInsertBefore, err: %v", err)
		}
		t.Logf("list: %v, val: %v", testCase.Key, val)

	}
}

//insert value to list head
func TestLPush(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key    string
		Values []string
	}{
		{
			Key:    "list-test-1",
			Values: []string{"value-list-1-test-10", "value-list-1-test-20"},
		},
		{
			Key:    "list-test-2",
			Values: []string{"value-list-2-test-10", "value-list-2-test-20"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		val, err := db.LPush(ctx, testCase.Key, testCase.Values).Result()
		if err != nil {
			t.Fatalf("failed to LInsertBefore, err: %v", err)
		}
		t.Logf("list: %v, val: %v", testCase.Key, val)

	}
}

//keep elements in range, the others will deleted
func TestLTrim(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key   string
		Start int64
		End   int64
	}{
		{
			Key:   "list-test-1",
			Start: 1,
			End:   3,
		},
		{
			Key:   "list-test-2",
			Start: 0,
			End:   -1,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		val, err := db.LTrim(ctx, testCase.Key, testCase.Start, testCase.End).Result()
		if err != nil {
			t.Fatalf("failed to LInsertBefore, err: %v", err)
		}
		t.Logf("list: %v, val: %v", testCase.Key, val)

	}
}

//delete element from list head, if the list is empty, will block times to find elements that can be ejected or wait timout
func TestBLPop(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key string
	}{
		{
			Key: "list-test-1",
		},
		{
			Key: "list-test-20",
		},
	}

	timout := 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		val, err := db.BLPop(ctx, timout, testCase.Key).Result()
		if err != nil {
			t.Fatalf("failed to BLPop, err: %v", err)
		}
		t.Logf("list: %v, val: %v", testCase.Key, val)

	}
}
