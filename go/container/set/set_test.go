package set_test

import (
	"testing"

	set_ "github.com/kaydxh/golang/go/container/set"
)

func TestGenericSetNew(t *testing.T) {
	s := set_.New[int]()
	s.Insert(2)
	t.Logf("s: %v", s)
}

func TestGenericSetInsert(t *testing.T) {
	s := set_.New("10", "2", "5")
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

func TestGenericSetEquals(t *testing.T) {

	a := set_.New("1", "2")
	b := set_.New("2", "1")

	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	//It is a set; duplicates are ignored
	b = set_.New("2", "1", "1")
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}
}

func TestGenericSetUnion(t *testing.T) {
	tests := []struct {
		s1       set_.Set[string]
		s2       set_.Set[string]
		expected set_.Set[string]
	}{
		{
			set_.New("1", "2", "3", "4"),
			set_.New("3", "4", "5", "6"),
			set_.New("1", "2", "3", "4", "5", "6"),
		},
		{
			set_.New("1", "2", "3", "4"),
			set_.New[string](),
			set_.New("1", "2", "3", "4"),
		},
		{
			set_.New[string](),
			set_.New("1", "2", "3", "4"),
			set_.New("1", "2", "3", "4"),
		},
		{
			set_.New[string](),
			set_.New[string](),
			set_.New[string](),
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
