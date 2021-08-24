package set_test

import (
	"testing"

	set_ "github.com/kaydxh/golang/go/container/set"
)

func TestStringSetInsert(t *testing.T) {
	s := set_.NewString("a", "b", "c")
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
	//set.String{"a":set.Empty{}, "b":set.Empty{}, "c":set.Empty{}, "d":set.Empty{}, "e":set.Empty{}}
	t.Logf("s: %#v", s)

	//[a b c d e]
	t.Logf("s: %v", s.List())

}

func TestStringSetEquals(t *testing.T) {

	a := set_.NewString("1", "2")
	b := set_.NewString("2", "1")

	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}

	//It is a set; duplicates are ignored
	b = set_.NewString("2", "1", "1")
	if !a.Equal(b) {
		t.Errorf("Expected to be equal: %v vs %v", a, b)
	}
}
