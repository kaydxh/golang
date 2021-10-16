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
			Key: "set-test-1",
			Values: []string{
				"values-set-test-1-1",
				"values-set-test-1-2",
				"values-set-test-1-3",
				"values-set-test-1-4",
				"values-set-test-1-5",
			},
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

//delete values from set
func TestSRem(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key    string
		Values []string
	}{
		{
			Key:    "set-test-1",
			Values: []string{"values-set-test-1-2"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		//err: CROSSSLOT Keys in request don't hash to the same slot
		vals, err := db.SRem(ctx, testCase.Key, testCase.Values).Result()
		if err != nil {
			t.Fatalf("failed to SRem, err: %v", err)
		}
		t.Logf("sets: %v, SRem values: %v", testCase.Key, vals)
	}
}

// check value is in set
//need redis server version >= 6.2.0
func TestSMIsMember(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key    string
		Values []string
	}{
		{
			Key:    "set-test-1",
			Values: []string{"values-set-test-1-2"},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		//err: CROSSSLOT Keys in request don't hash to the same slot
		vals, err := db.SMIsMember(ctx, testCase.Key, testCase.Values).Result()
		if err != nil {
			t.Fatalf("failed to SMIsMember, err: %v", err)
		}
		t.Logf("sets: %v, SMIsMember values: %v", testCase.Key, vals)
	}
}

//random delete elemnt from set
func TestSPop(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key string
	}{
		{
			Key: "set-test-1",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		val, err := db.SPop(ctx, testCase.Key).Result()
		if err != nil {
			t.Fatalf("failed to SPop, err: %v", err)
		}
		t.Logf("sets: %v,  SPop value: %v", testCase.Key, val)
	}
}

// random get value from set
func TestSRandMember(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key string
	}{
		{
			Key: "set-test-1",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		val, err := db.SRandMember(ctx, testCase.Key).Result()
		if err != nil {
			t.Fatalf("failed to SRandMember, err: %v", err)
		}
		t.Logf("sets: %v, SRandMember value: %v", testCase.Key, val)
	}
}

// random get value from set
func TestSRandMemberN(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key   string
		Count int64
	}{
		{
			Key:   "set-test-1",
			Count: 2,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		vals, err := db.SRandMemberN(ctx, testCase.Key, testCase.Count).Result()
		if err != nil {
			t.Fatalf("failed to SRandMemberN, err: %v", err)
		}
		t.Logf("sets: %v, SRandMemberN values: %v", testCase.Key, vals)
	}
}

func TestSUnion(t *testing.T) {
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
		//err: CROSSSLOT Keys in request don't hash to the same slot
		vals, err := db.SUnion(ctx, testCase.Keys...).Result()
		if err != nil {
			t.Fatalf("failed to SUnion, err: %v", err)
		}
		t.Logf("sets: %v, SUnion values: %v", testCase.Keys, vals)
	}
}

func TestSMove(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Source      string
		Destination string
		Value       string
	}{
		{
			Source:      "set-test-1",
			Destination: "set-test-20",
			Value:       "values-set-test-1-5",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		//err: CROSSSLOT Keys in request don't hash to the same slot
		ok, err := db.SMove(ctx, testCase.Source, testCase.Destination, testCase.Value).Result()
		if err != nil {
			t.Fatalf("failed to SUnion, err: %v", err)
		}
		t.Logf("move sets: %v to %v, SUnion values: %v", testCase.Source, testCase.Destination, ok)
	}
}

func TestSUnionStore(t *testing.T) {
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
		val, err := db.SUnionStore(ctx, testCase.TargetSet, testCase.Keys...).Result()
		if err != nil {
			t.Fatalf("failed to SUnionStore, err: %v", err)
		}
		t.Logf("sets: %v, SUnionStore values: %v", testCase.Keys, val)
	}
}

func TestSScan(t *testing.T) {
	db := GetDBOrDie()

	testCases := []struct {
		Key    string
		Cursor uint64
		Match  string
		Count  int64
	}{
		{
			Key:    "set-test-1",
			Cursor: 0,
			Match:  "*set*",
			Count:  2,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, testCase := range testCases {
		//err: CROSSSLOT Keys in request don't hash to the same slot
		keys, cursor, err := db.SScan(ctx, testCase.Key, testCase.Cursor, testCase.Match, testCase.Count).Result()
		if err != nil {
			t.Fatalf("failed to SScan, err: %v", err)
		}
		t.Logf("sets: %v,SScan keys: %v, cursor: %v", testCase.Key, keys, cursor)
	}
}
