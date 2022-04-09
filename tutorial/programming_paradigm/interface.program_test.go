package tutorial

import (
	"fmt"
	"testing"
)

//Program to an interface not an implementation
type Country struct {
	Name string
}

type City struct {
	Name string
}

type Stringable interface {
	ToString() string
}

func (c Country) ToString() string {
	return "Country = " + c.Name
}
func (c City) ToString() string {
	return "City = " + c.Name
}

// Program to an interface not an implementation
func PrintStr(p Stringable) {
	fmt.Println(p.ToString())
}

func TestPrintStr(t *testing.T) {
	testCases := []struct {
		obj      Stringable
		expected string
	}{
		{
			obj:      Country{"USA"},
			expected: "",
		},
		{
			obj:      City{"Los Angeles"},
			expected: "",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			PrintStr(testCase.obj)
		})
	}
}
