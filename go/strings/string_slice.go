package strings

import (
	set_ "github.com/kaydxh/golang/go/container/set"
)

func SliceUnique(s ...string) []string {
	ss := set_.NewString()
	for _, v := range s {
		ss.Insert(v)
	}

	return ss.List()
}

func SliceIntersection(s1 []string, s2 []string) []string {
	ss1 := set_.NewObject()
	for _, s := range s1 {
		ss1.Insert(s)
	}

	ss2 := set_.NewObject()
	for _, s := range s2 {
		ss2.Insert(s)
	}

	ss := []string{}
	for _, v := range ss1.Intersection(ss2).List() {
		s, ok := v.(string)
		if ok {
			ss = append(ss, s)
		}
	}
	return ss
}