package set_test

import (
	"testing"

	set_ "github.com/kaydxh/golang/go/container/set"
)

func TestSetInsert(t *testing.T) {
	s := set_.NewObject("10", "2", "5")
	s.Insert("a", "d", "e")
	if len(s) != 5 {
		t.Errorf("Expected  len=5: %d", len(s))
	}

	if !s.Has("a") || !s.Has("b") || !s.Has("c") || !s.Has("d") || !s.Has("e") {
		t.Errorf("UnExpected  contents: %#v", s)
	}

	//%v output value
	//map[a:{} b:{} c:{} d:{} e:{}]
	t.Logf("s: %v", s)

	//%+v output field name + value
	//map[a:{} b:{} c:{} d:{} e:{}]
	t.Logf("s: %+v", s)

	//%#v output struct name + field name + value
	//set.Object{"a":set.Empty{}, "b":set.Empty{}, "c":set.Empty{}, "d":set.Empty{}, "e":set.Empty{}}
	t.Logf("s: %#v", s)

	//[a b c d e]
	t.Logf("s: %v", s.List())

}

func TestSetEquals(t *testing.T) {

	a := set_.NewObject("1", "2")
	b := set_.NewObject("2", "1")

	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	//It is a set; duplicates are ignored
	b = set_.NewObject("2", "1", "1")
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}
}

func TestSetUnion(t *testing.T) {
	tests := []struct {
		s1       set_.Object
		s2       set_.Object
		expected set_.Object
	}{
		{
			set_.NewObject("1", "2", "3", "4"),
			set_.NewObject("3", "4", "5", "6"),
			set_.NewObject("1", "2", "3", "4", "5", "6"),
		},
		{
			set_.NewObject("1", "2", "3", "4"),
			set_.NewObject(),
			set_.NewObject("1", "2", "3", "4"),
		},
		{
			set_.NewObject(),
			set_.NewObject("1", "2", "3", "4"),
			set_.NewObject("1", "2", "3", "4"),
		},
		{
			set_.NewObject(),
			set_.NewObject(),
			set_.NewObject(),
		},
	}

	for _, test := range tests {
		union := test.s1.Union(test.s2)
		if union.Len() != test.expected.Len() {
			t.Errorf("Expected union.Len()=%d but got %d", test.expected.Len(), union.Len())
		}

		if !union.Equal(test.expected) {
			t.Errorf(
				"Expected union.Equal(expected) but not true.  union:%v expected:%v",
				union.List(),
				test.expected.List(),
			)
		}
	}

}

func TestSetIntersection(t *testing.T) {
	tests := []struct {
		s1       set_.Object
		s2       set_.Object
		expected set_.Object
	}{
		{
			set_.NewObject("1", "2", "3", "4"),
			set_.NewObject("3", "4", "5", "6"),
			set_.NewObject("3", "4"),
		},
		{
			set_.NewObject("1", "2", "3", "4"),
			set_.NewObject("1", "2", "3", "4"),
			set_.NewObject("1", "2", "3", "4"),
		},
		{
			set_.NewObject("1", "2", "3", "4"),
			set_.NewObject(),
			set_.NewObject(),
		},
		{
			set_.NewObject(),
			set_.NewObject("1", "2", "3", "4"),
			set_.NewObject(),
		},
		{
			set_.NewObject(),
			set_.NewObject(),
			set_.NewObject(),
		},
	}

	for _, test := range tests {
		intersection := test.s1.Intersection(test.s2)
		if intersection.Len() != test.expected.Len() {
			t.Errorf("Expected intersection.Len()=%d but got %d", test.expected.Len(), intersection.Len())
		}

		if !intersection.Equal(test.expected) {
			t.Errorf(
				"Expected intersection.Equal(expected) but not true.  intersection:%v expected:%v",
				intersection.List(),
				test.expected.List(),
			)
		}
	}

}