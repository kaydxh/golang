/*
 *Copyright (c) 2023, kaydxh
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

package vector_test

import (
	"math/rand"
	"testing"

	"github.com/kaydxh/golang/pkg/gocv/vector"
)

func TestVector(t *testing.T) {
	var ve vector.Vector[float32]
	for i := 0; i < 1024; i++ {
		ve.Append(rand.Float32())
	}
	t.Logf("len: %v", ve.Len())
	s, err := ve.String()
	if err != nil {
		t.Fatalf("failed to get string, err: %v", err)
	}
	t.Logf("String: %v", s)

	ve.Norm()
	s, err = ve.String()
	if err != nil {
		t.Fatalf("failed to get norm string, err: %v", err)
	}
	t.Logf("Norm string: %v", s)
}

func TestCreateNormalizedVector(t *testing.T) {
	dim := 1024
	data := vector.CreateNormalizedVector[float32](dim)
	t.Logf("data: %v", data)
}

func TestCosineDistance(t *testing.T) {
	dim := 1024
	data1 := vector.CreateNormalizedVector[float32](dim)
	data2 := vector.CreateNormalizedVector[float32](dim)
	v1 := vector.NewVector(data1...)
	v2 := vector.NewVector(data2...)
	distance := v1.CosineDistance(v2)
	t.Logf("distance: %v", distance)

	distance = vector.CosineDistance(data1, data2)
	t.Logf("distance: %v", distance)
}
