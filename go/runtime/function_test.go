package runtime_test

import (
	"fmt"
	"testing"

	runtime_ "github.com/kaydxh/golang/go/runtime"
	"github.com/stretchr/testify/assert"
)

func testA(a int) int {
	return a
}

func TestNameOfFunction(t *testing.T) {
	testCases := []struct {
		f        interface{}
		expected string
	}{
		{
			f: func(a int) int {
				return a
			},
			expected: "github.com/kaydxh/golang/go/runtime_test.TestNameOfFunction.func1",
		},
		{
			f: func(a int) int {
				return a
			},
			expected: "github.com/kaydxh/golang/go/runtime_test.TestNameOfFunction.func2",
		},
		{
			f:        testA,
			expected: "github.com/kaydxh/golang/go/runtime_test.testA",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test-case-%d", i), func(t *testing.T) {
			funcName := runtime_.NameOfFunction(testCase.f)
			assert.Equal(t, testCase.expected, funcName)
			t.Logf("funcName: %v", funcName)
		})
	}
}
