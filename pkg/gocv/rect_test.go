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
package gocv_test

import (
	"fmt"
	"testing"

	gocv_ "github.com/kaydxh/golang/pkg/gocv"
)

func TestScale(t *testing.T) {
	testCases := []struct {
		X      int32
		Y      int32
		Width  int32
		Height int32
		factor float32
	}{
		{
			X:      100,
			Y:      100,
			Width:  100,
			Height: 100,
			factor: 1.1,
		},
		{
			X:      100,
			Y:      100,
			Width:  100,
			Height: 100,
			factor: 0.9,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			r := gocv_.Rect{
				X:      testCase.X,
				Y:      testCase.Y,
				Width:  testCase.Width,
				Height: testCase.Height,
			}
			r2 := r.Scale(testCase.factor).Intersect(r)
			t.Logf("r2: %v", r2)
		})
	}
}

func TestClosest(t *testing.T) {
	testCases := []struct {
		r  gocv_.Rect
		rs []gocv_.Rect
	}{
		{
			r: gocv_.Rect{
				X:      1042,
				Y:      518,
				Width:  389,
				Height: 467,
			},
			rs: []gocv_.Rect{
				// out
				{
					X:      0,
					Y:      10,
					Width:  10,
					Height: 10,
				},
				{
					X:      400,
					Y:      600,
					Width:  700,
					Height: 700,
				},

				// small in
				{
					X:      1100,
					Y:      600,
					Width:  300,
					Height: 300,
				},

				//big in
				{
					X:      938,
					Y:      245,
					Width:  807,
					Height: 819,
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			_, cr := testCase.r.Closest(testCase.rs...)
			t.Logf("cr: %v", cr)
		})
	}
}
