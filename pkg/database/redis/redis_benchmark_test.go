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
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

// -count the benchmark times, -benchtime the test execute times(用例执行次数) or execute time(用例执行时间)
// -v detail log info , -benchmem summary of memory
//go test -bench="Set" -benchtime=5s -count=3 .
//go test -bench="Set" -benchtime=50x -count=3 .
func BenchmarkSet(t *testing.B) {
	db := GetDBOrDie()
	//	defer db.Close()

	keyPrefix := "test"

	for n := 0; n < t.N; n++ {
		fmt.Println("n: ", n)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		key := fmt.Sprintf("%s-%d", keyPrefix, n)
		result, err := db.Set(ctx, key, key, 0).Result()
		if err != nil {
			t.Fatalf("failed to set string, err: %v", err)
		}

		t.Logf("result of %v: %v", key, result)
	}
}

//go test -v -run=redis_benchmark_test.go -test.bench="ParallelSet" -benchtime=5s -count=3 .
//go test -v -run=redis_benchmark_test.go -test.bench="ParallelSet" -benchtime=5s -count=3 .
func BenchmarkParallelSet(t *testing.B) {
	db := GetDBOrDie()
	//	defer db.Close()

	keyPrefix := "test"

	t.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			key := fmt.Sprintf("%s-%s", keyPrefix, uuid.New().String())

			result, err := db.Set(ctx, key, key, 0).Result()
			if err != nil {
				t.Fatalf("failed to set string, err: %v", err)
			}

			t.Logf("result of %v: %v", key, result)
		}
	})

}
