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
			r2 := r.Scale(testCase.factor)
			t.Logf("r2: %v", r2)
			r3 := r.Intersect(r2)
			t.Logf("r3: %v", r3)
		})
	}
}
