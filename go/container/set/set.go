package set

import (
	"encoding/json"
	"reflect"
	"sort"
)

// sets.Object is a set of Objects, implemented via map[Object]struct{} for minimal memory consumption.
type Object map[interface{}]Empty

// NewObject creates a Object from a list of values.
func NewObject(items ...interface{}) Object {
	ss := Object{}
	ss.Insert(items...)
	return ss
}

// ObjectKeySet creates a Object from a keys of a map[Object](? extends interface{}).
// If the value passed in is not actually a map, this will panic.
func KeySet(theMap interface{}) Object {
	v := reflect.ValueOf(theMap)
	ret := Object{}

	for _, keyValue := range v.MapKeys() {
		ret.Insert(keyValue.Interface().(Object))
	}
	return ret
}

// Insert adds items to the set.
func (s Object) Insert(items ...interface{}) Object {
	for _, item := range items {
		s[item] = Empty{}
	}
	return s
}

// Delete removes all items from the set.
func (s Object) Delete(items ...interface{}) Object {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

// Has returns true if and only if item is contained in the set.
func (s Object) Has(item interface{}) bool {
	_, contained := s[item]
	return contained
}

// HasAll returns true if and only if all items are contained in the set.
func (s Object) HasAll(items ...interface{}) bool {
	for _, item := range items {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

// HasAny returns true if any items are contained in the set.
func (s Object) HasAny(items ...interface{}) bool {
	for _, item := range items {
		if s.Has(item) {
			return true
		}
	}
	return false
}

// Difference returns a set of objects that are not in s2
// For example:
// s1 = {a1, a2, a3}
// s2 = {a1, a2, a4, a5}
// s1.Difference(s2) = {a3}
// s2.Difference(s1) = {a4, a5}
func (s Object) Difference(s2 Object) Object {
	result := NewObject()
	for key := range s {
		if !s2.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// Union returns a new set which includes items in either s1 or s2.
// For example:
// s1 = {a1, a2}
// s2 = {a3, a4}
// s1.Union(s2) = {a1, a2, a3, a4}
// s2.Union(s1) = {a1, a2, a3, a4}
func (s1 Object) Union(s2 Object) Object {
	result := NewObject()
	for key := range s1 {
		result.Insert(key)
	}
	for key := range s2 {
		result.Insert(key)
	}
	return result
}

// Intersection returns a new set which includes the item in BOTH s1 and s2
// For example:
// s1 = {a1, a2}
// s2 = {a2, a3}
// s1.Intersection(s2) = {a2}
func (s1 Object) Intersection(s2 Object) Object {
	var walk, other Object
	result := NewObject()
	if s1.Len() < s2.Len() {
		walk = s1
		other = s2
	} else {
		walk = s2
		other = s1
	}
	for key := range walk {
		if other.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// IsSuperset returns true if and only if s1 is a superset of s2.
func (s1 Object) IsSuperset(s2 Object) bool {
	for item := range s2 {
		if !s1.Has(item) {
			return false
		}
	}
	return true
}

// Equal returns true if and only if s1 is equal (as a set) to s2.
// Two sets are equal if their membership is identical.
// (In practice, this means same elements, order doesn't matter)
func (s1 Object) Equal(s2 Object) bool {
	return len(s1) == len(s2) && s1.IsSuperset(s2)
}

type sortableSliceOfObject []interface{}

func (s sortableSliceOfObject) Len() int           { return len(s) }
func (s sortableSliceOfObject) Less(i, j int) bool { return lessObject(s[i], s[j]) }
func (s sortableSliceOfObject) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// List returns the contents as a sorted Object slice.
func (s Object) List() []interface{} {
	res := make(sortableSliceOfObject, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	sort.Sort(res)
	return []interface{}(res)
}

// UnsortedList returns the slice with contents in random order.
func (s Object) UnsortedList() []interface{} {
	res := make([]interface{}, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	return res
}

// Returns a single element from the set.
func (s Object) PopAny() (interface{}, bool) {
	for key := range s {
		s.Delete(key)
		return key, true
	}
	var zeroValue Object
	return zeroValue, false
}

// Len returns the size of the set.
func (s Object) Len() int {
	return len(s)
}

type Lesser interface {
	Less(a interface{}, b interface{}) bool
}

/*
func lessObject(lhs, rhs interface{}) bool {
	if lhs == nil {
		return false
	}

	if rhs == nil {
		return true
	}

	switch lhs := lhs.(type) {
	case nil:
		return false

	case *string:
		rhs, ok := rhs.(*string)
		if !ok {
			return false
		}

		return lhs != nil && rhs != nil && *lhs < *rhs

	case *[]byte:
		rhs, ok := rhs.(*[]byte)
		if !ok {
			return false
		}
		return lhs != nil && rhs != nil && string(*lhs) < string(*rhs)

	case *int:
		rhs, ok := rhs.(*int)
		if !ok {
			return false
		}
		return lhs != nil && rhs != nil && *lhs < *rhs

	case *int32:
		rhs, ok := rhs.(*int32)
		if !ok {
			return false
		}
		return lhs != nil && rhs != nil && *lhs < *rhs

	case *uint32:
		rhs, ok := rhs.(*uint32)
		if !ok {
			return false
		}
		return lhs != nil && rhs != nil && *lhs < *rhs

	case *int64:
		rhs, ok := rhs.(*int64)
		if !ok {
			return false
		}
		return lhs != nil && rhs != nil && *lhs < *rhs

	case *uint64:
		rhs, ok := rhs.(*uint64)
		if !ok {
			return false
		}
		return lhs != nil && rhs != nil && *lhs < *rhs

	case *float32:
		rhs, ok := rhs.(*float32)
		if !ok {
			return false
		}
		return lhs != nil && rhs != nil && *lhs < *rhs

	case *float64:
		rhs, ok := rhs.(*float64)
		if !ok {
			return false
		}
		return lhs != nil && rhs != nil && *lhs < *rhs


	default:
		var ll Lesser
		ll, ok = lhs.(Lesser)
		if !ok {

		}

			var err error
			jb, err = json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("jsonpb.Marshal: %v", err)
			}

	}
	return false
}
*/

func lessObject(lhs, rhs interface{}) bool {
	jbl, err := json.Marshal(lhs)
	if err != nil {
		return false
	}
	jbr, err := json.Marshal(rhs)
	if err != nil {
		return false
	}

	return string(jbl) < string(jbr)
}

/*
func LessObject(lhs, rhs interface{}) bool {
	jbl, err := json.Marshal(lhs)
	if err != nil {
		return false
	}
	jbr, err := json.Marshal(rhs)
	if err != nil {
		return false
	}

	return string(jbl) < string(jbr)
}
*/
