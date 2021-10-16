package redis_test

import (
	"context"
	"testing"
	"time"
)

func TestSAdd(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key    string
		Values []string
	}{
		{
			Key:    "set-test-1",
			Values: []string{"values-set-test-1-1", "values-set-test-1-2"},
		},
		{
			Key:    "set-test-20",
			Values: []string{"values-set-test-1-1", "values-set-test-20-1", "values-set-test-20-2"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		val, err := db.SAdd(ctx, testCase.Key, testCase.Values).Result()
		if err != nil {
			t.Fatalf("failed to SAdd, err: %v", err)
		}
		t.Logf("set: %v, val: %v", testCase.Key, val)
	}
}

func TestSCard(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key string
	}{
		{
			Key: "set-test-1",
		},
		{
			Key: "set-test-20",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		count, err := db.SCard(ctx, testCase.Key).Result()
		if err != nil {
			t.Fatalf("failed to SCard, err: %v", err)
		}
		t.Logf("set: %v, member count: %v", testCase.Key, count)
	}
}

//diff val from multi sets
// return the first set different elements from the other sets
func TestSDiff(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Keys []string
	}{
		{
			Keys: []string{"set-test-1", "set-test-20"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		vals, err := db.SDiff(ctx, testCase.Keys...).Result()
		if err != nil {
			t.Fatalf("failed to SDiff, err: %v", err)
		}
		t.Logf("sets: %v, diff values: %v", testCase.Keys, vals)
	}
}

func TestSDiffStore(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		TargetSet string
		Keys      []string
	}{
		{
			TargetSet: "target-set",
			Keys:      []string{"set-test-1", "set-test-20"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		//err: CROSSSLOT Keys in request don't hash to the same slot
		vals, err := db.SDiffStore(ctx, testCase.TargetSet, testCase.Keys...).Result()
		if err != nil {
			t.Fatalf("failed to SDiffStore, err: %v", err)
		}
		t.Logf("sets: %v,  SDiffStore values: %v", testCase.Keys, vals)
	}
}

func TestSInter(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Keys []string
	}{
		{
			Keys: []string{"set-test-1", "set-test-20"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		vals, err := db.SInter(ctx, testCase.Keys...).Result()
		if err != nil {
			t.Fatalf("failed to SDiffStore, err: %v", err)
		}
		t.Logf("sets: %v, SInter values: %v", testCase.Keys, vals)
	}
}

func TestSMembers(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key string
	}{
		{
			Key: "set-test-1",
		},
		{
			Key: "set-test-20",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		vals, err := db.SMembers(ctx, testCase.Key).Result()
		if err != nil {
			t.Fatalf("failed to SMembers, err: %v", err)
		}
		t.Logf("sets: %v,  SMembers values: %v", testCase.Key, vals)
	}
}

func TestSInterStore(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		TargetSet string
		Keys      []string
	}{
		{
			TargetSet: "target-inter-set",
			Keys:      []string{"set-test-1", "set-test-20"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		//err: CROSSSLOT Keys in request don't hash to the same slot
		vals, err := db.SInterStore(ctx, testCase.TargetSet, testCase.Keys...).Result()
		if err != nil {
			t.Fatalf("failed to SInterStore, err: %v", err)
		}
		t.Logf("sets: %v, SInterStore values: %v", testCase.Keys, vals)
	}
}
