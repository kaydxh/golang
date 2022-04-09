package tutorial

import (
	"fmt"
	"testing"
)

type Shape interface {
	Sides() int
	Area() int
}
type Square struct {
	length int
}

func NewSquare(length int) *Square {
	return &Square{
		length: length,
	}
}

func (s *Square) Sides() int {
	return 4
}

func (s *Square) Area() int {
	return 16
}

// check interface complete
// if not implete Area method, it will report InvalidIfaceAssign: cannot use (*Square)(nil) (value of type *Square)
//as Shape value in variable declaration: missing method Area
var _ Shape = (*Square)(nil)

func TestNewSquare(t *testing.T) {
	testCases := []struct {
		length   int
		expected string
	}{
		{
			length:   5,
			expected: "",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			s := NewSquare(testCase.length)
			t.Logf("s: %v", s.Sides())
		})
	}
}
