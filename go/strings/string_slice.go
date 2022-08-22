package strings

import (
	"strings"

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

func SliceDifference(s1 []string, s2 []string) []string {
	ss1 := set_.NewObject()
	for _, s := range s1 {
		ss1.Insert(s)
	}

	ss2 := set_.NewObject()
	for _, s := range s2 {
		ss2.Insert(s)
	}

	ss := []string{}
	for _, v := range ss1.Difference(ss2).List() {
		s, ok := v.(string)
		if ok {
			ss = append(ss, s)
		}
	}

	return ss
}

func RemoveEmpty(s []string) []string {
	var ss []string
	for _, v := range s {
		if v != "" {
			ss = append(ss, v)
		}
	}

	return ss
}

// sliceContains reports whether the provided string is present in the given slice of strings.
func SliceContainsCaseInSensitive(list []string, target string) bool {
	return SliceContains(list, target, false)
}

func SliceContains(list []string, target string, caseSensitive bool) bool {
	if !caseSensitive {
		target = strings.ToLower(target)
	}

	for _, s := range list {
		if !caseSensitive {
			s = strings.ToLower(s)
		}

		if s == target {
			return true
		}
	}
	return false
}

// Split2 returns the values from strings.SplitN(s, sep, 2).
// If sep is not found, it returns ("", "", false) instead.
func Split2(s, sep string) (string, string, bool) {
	spl := strings.SplitN(s, sep, 2)
	if len(spl) < 2 {
		return "", "", false
	}
	return spl[0], spl[1], true
}

/*
func SliceIntersectionInt(s1 []int, s2 []int) []int {
	ss1 := set_.NewObject(set_.GenerateArray([...]int(s1...)))
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
*/
