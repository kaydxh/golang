package cvtable_test

import (
	"fmt"
	"testing"

	cvtable_ "github.com/kaydxh/golang/pkg/gocv/cvtable"
)

func TestSim(t *testing.T) {
	cvtable, err := cvtable_.NewCVTable("./testdata/cvtable.conf")
	if err != nil {
		t.Fatalf("failed to new cvtable, err: %v", err)
	}

	testCases := []struct {
		score    float64
		expected string
	}{
		{
			score:    95.13,
			expected: "",
		},
		{
			score:    200.13,
			expected: "",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("case-%v", i), func(t *testing.T) {
			sim := cvtable.Sim(testCase.score)
			t.Logf("get sim %v for score %v", sim, testCase.score)
		})
	}
}

func TestScore(t *testing.T) {
	cvtable, err := cvtable_.NewCVTable("./testdata/cvtable.conf")
	if err != nil {
		t.Fatalf("failed to new cvtable, err: %v", err)
	}

	testCases := []struct {
		sim      float64
		expected string
	}{
		{
			sim:      0.9513,
			expected: "",
		},
		{
			sim:      1.20013,
			expected: "",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("case-%v", i), func(t *testing.T) {
			score := cvtable.Score(testCase.sim)
			t.Logf("get score %v for sim %v", score, testCase.sim)
		})
	}
}
