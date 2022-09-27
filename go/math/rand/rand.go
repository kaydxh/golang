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
package rand

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	globalRand  = rand.New(rand.NewSource(time.Now().UnixNano()))
	mu          sync.Mutex
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
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

//  RangeString generate string length [0, n].
func RangeString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
