/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
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

func SliceWithCondition(s1 []string, cond func(s2 string) bool) []string {
	ss1 := set_.NewObject()
	for _, s := range s1 {
		if cond(s) {
			ss1.Insert(s)
		}
	}

	ss := []string{}
	for _, v := range ss1.List() {
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

func Filter(ss []string, cond func(string) bool) []string {
	var res []string
	for _, s := range ss {
		if cond(s) {
			res = append(res, s)
		}
	}

	return res
}
