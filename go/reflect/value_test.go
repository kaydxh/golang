package reflect_test

import (
	"testing"

	reflect_ "github.com/kaydxh/golang/go/reflect"
)

func TestReflectValue(t *testing.T) {
	testCases := []struct {
		valueType string
		value     interface{}
		expected  string
	}{
		{
			valueType: "bool",
			value:     true,
			expected:  "",
		},
		{
			valueType: "int",
			value:     123456789,
			expected:  "",
		},
		{
			valueType: "uint",
			value:     uint(123456789),
			expected:  "",
		},
		{
			valueType: "float32",
			value:     0.123456789,
			expected:  "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.valueType, func(t *testing.T) {
			rvalue, err := reflect_.ReflectValue(testCase.valueType, testCase.value)
			if err != nil {
				t.Fatalf("failed to reflect value: %v, got : %s", testCase.value, err)

			}

			t.Logf("reflect value: %v", rvalue)

		})
	}

}
