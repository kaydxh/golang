package rand

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	mu         sync.Mutex
)

// Int implements rand.Int on the global source.
func Int() int {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Int()
}

// Int31n implements rand.Int31n on the global source.
func Int31n(n int32) int32 {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Int31n(n)
}

// Uint32 implements rand.Uint32 on the global source.
func Uint32() uint32 {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Uint32()
}

// Int63n implements rand.Int63n on the global source.
func Int63n(n int64) int64 {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Int63n(n)
}

// Intn implements rand.Intn on the global source.
func Intn(n int) int {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Intn(n)
}

// Float32 implements rand.Float32 on the global source.
func Float32() float32 {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Float32()
}

// Float64 implements rand.Float64 on the global source.
func Float64() float64 {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Float64()
}

// Uint64 implements rand.Uint64 on the global source.
func Uint64() uint64 {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Uint64()
}

// Read implements rand.Read on the global source.
func Read(p []byte) (n int, err error) {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Read(p)
}

// RandInt generate number [min, max).
func RangeInt(min, max int) (int, error) {
	if min < 0 || max < 0 || max <= min {
		return 0, fmt.Errorf("min or max must > 0 and max > min")
	}
	return Intn(max-min) + min, nil
}
