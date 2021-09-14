package redis_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

// -count the benchmark times, -benchtime the test execute times(用例执行次数) or execute time(用例执行时间)
//go test -bench="Set" -benchtime=5s -count=3 .
//go test -bench="Set" -benchtime=50x -count=3 .
func BenchmarkSet(t *testing.B) {
	db := GetDBOrDie()
	defer db.Close()

	keyPrefix := "test"

	for n := 0; n < t.N; n++ {
		fmt.Println("n: ", n)

		key := fmt.Sprintf("%s-%d", keyPrefix, n)

		result, err := db.Set(key, key, 0).Result()
		if err != nil {
			t.Fatalf("failed to set string, err: %v", err)
		}

		t.Logf("result of %v: %v", key, result)
	}
}

//go test -run=redis_benchmark_test.go -test.bench="ParallelSet" -benchtime=5s -count=3 .
func BenchmarkParallelSet(t *testing.B) {
	db := GetDBOrDie()
	defer db.Close()

	keyPrefix := "test"

	t.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := fmt.Sprintf("%s-%s", keyPrefix, uuid.New().String())

			result, err := db.Set(key, key, 0).Result()
			if err != nil {
				t.Fatalf("failed to set string, err: %v", err)
			}

			t.Logf("result of %v: %v", key, result)
		}
	})

}
