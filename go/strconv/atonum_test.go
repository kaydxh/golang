package strconv_test

import (
	"fmt"
	"testing"

	strconv_ "github.com/kaydxh/golang/go/strconv"
)

func TestParseNumOrDefault(t *testing.T) {

	testCases := []struct {
		str      string
		expected string
	}{
		{
			str:      "12345",
			expected: "",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			num := strconv_.ParseNumOrDefault(testCase.str, 0, strconv_.ToInt)
			t.Logf("get num: %v", num)

		})
	}

}

func TestParseNum(t *testing.T) {

	testCases := []struct {
		str      string
		expected string
	}{
		{
			str:      "12345",
			expected: "",
		},
		{
			str:      "badcase",
			expected: "",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			num, err := strconv_.ParseNum(testCase.str, strconv_.ToInt)
			if err != nil {
				t.Errorf("expecet nil, got %v", err)
			}
			t.Logf("get num: %v", num)
		})
	}

}
