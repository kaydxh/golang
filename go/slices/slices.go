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
package slices

import (
	set_ "github.com/kaydxh/golang/go/container/set"
)

func Unique[S ~[]E, E comparable](s S) S {
	ss := set_.New[E]()
	for _, v := range s {
		ss.Insert(v)
	}

	return ss.List()
}

func SliceIntersection[S ~[]E, E comparable](s1, s2 S) S {
	ss1 := set_.New[E]()
	for _, s := range s1 {
		ss1.Insert(s)
	}

	ss2 := set_.New[E]()
	for _, s := range s2 {
		ss2.Insert(s)
	}

	var ss S
	for _, v := range ss1.Intersection(ss2).List() {
		ss = append(ss, v)
	}
	return ss
}

func SliceDifference[S ~[]E, E comparable](s1, s2 S) S {
	ss1 := set_.New[E]()
	for _, s := range s1 {
		ss1.Insert(s)
	}

	ss2 := set_.New[E]()
	for _, s := range s2 {
		ss2.Insert(s)
	}

	var ss S
	for _, v := range ss1.Difference(ss2).List() {
		ss = append(ss, v)
	}

	return ss
}

func SliceWithCondition[S ~[]E, E comparable](s1 S, cond func(e E) bool) S {
	ss1 := set_.New[E]()
	for _, s := range s1 {
		if cond(s) {
			ss1.Insert(s)
		}
	}

	var ss S
	for _, v := range ss1.List() {
		ss = append(ss, v)
	}

	return ss
}

func RemoveEmpty[S ~[]E, E comparable](s S) S {
	var ss S
	var zero E
	for _, v := range s {
		if v != zero {
			ss = append(ss, v)
		}
	}

	return ss
}
