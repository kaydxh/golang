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
package taskqueue

import (
	"sync" // Used by sync.Map.
)

// Generate code that will fail if the constants change value.
func _() {
	// An "cannot convert TaskerMap literal (type TaskerMap) to type sync.Map" compiler error signifies that the base type have changed.
	// Re-run the go-syncmap command to generate them again.
	_ = (sync.Map)(TaskerMap{})
}

var _nil_TaskerMap_Tasker_value = func() (val Tasker) { return }()

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (m *TaskerMap) Load(key string) (Tasker, bool) {
	value, ok := (*sync.Map)(m).Load(key)
	if value == nil {
		return _nil_TaskerMap_Tasker_value, ok
	}
	return value.(Tasker), ok
}

// Store sets the value for a key.
func (m *TaskerMap) Store(key string, value Tasker) {
	(*sync.Map)(m).Store(key, value)
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *TaskerMap) LoadOrStore(key string, value Tasker) (Tasker, bool) {
	actual, loaded := (*sync.Map)(m).LoadOrStore(key, value)
	if actual == nil {
		return _nil_TaskerMap_Tasker_value, loaded
	}
	return actual.(Tasker), loaded
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m *TaskerMap) LoadAndDelete(key string) (value Tasker, loaded bool) {
	actual, loaded := (*sync.Map)(m).LoadAndDelete(key)
	if actual == nil {
		return _nil_TaskerMap_Tasker_value, loaded
	}
	return actual.(Tasker), loaded
}

// Delete deletes the value for a key.
func (m *TaskerMap) Delete(key string) {
	(*sync.Map)(m).Delete(key)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently, Range may reflect any mapping for that key
// from any point during the Range call.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (m *TaskerMap) Range(f func(key string, value Tasker) bool) {
	(*sync.Map)(m).Range(func(key, value interface{}) bool {
		return f(key.(string), value.(Tasker))
	})
}
