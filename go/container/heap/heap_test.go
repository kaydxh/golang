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
package heap_test

import (
	"sync"
	"testing"

	heap_ "github.com/kaydxh/golang/go/container/heap"
)

func testHeapObjectKeyFunc(obj interface{}) string {
	return obj.(testHeapObject).name
}

type testHeapObject struct {
	name string
	val  interface{}
}

func mkHeapObj(name string, val interface{}) testHeapObject {
	return testHeapObject{name: name, val: val}
}

//minheap
func compareInts(val1 interface{}, val2 interface{}) bool {
	first := val1.(testHeapObject).val.(int)
	second := val2.(testHeapObject).val.(int)
	return first < second
}

// TestHeapBasic tests Heap invariant and synchronization.
func TestHeapBasic(t *testing.T) {
	h := heap_.NewHeap(testHeapObjectKeyFunc, compareInts)
	var wg sync.WaitGroup
	wg.Add(2)
	const amount = 10
	var i, u int
	// Insert items in the heap in opposite orders in two go routines.
	go func() {
		for i = amount; i > 0; i-- {
			h.Add(mkHeapObj(string([]rune{'a', rune(i)}), i))
		}
		wg.Done()
	}()
	go func() {
		for u = 0; u < amount; u++ {
			h.Add(mkHeapObj(string([]rune{'b', rune(u)}), u+1))
		}
		wg.Done()
	}()
	// Wait for the two go routines to finish.
	wg.Wait()

	t.Logf("heap: %+v", h.List())
	// Make sure that the numbers are popped in ascending order.
	prevNum := 0
	for i := 0; i < amount*2; i++ {
		obj, err := h.Pop()
		num := obj.(testHeapObject).val.(int)
		// All the items must be sorted.
		if err != nil || prevNum > num {
			t.Errorf("got %v out of order, last was %v", obj, prevNum)
		}
		t.Logf("get %v", num)
		prevNum = num
	}

}

// TestHeap_Get tests Heap.Get.
func TestHeap_Get(t *testing.T) {
	h := heap_.NewHeap(testHeapObjectKeyFunc, compareInts)
	h.Add(mkHeapObj("foo", 10))
	h.Add(mkHeapObj("bar", 1))
	h.Add(mkHeapObj("bal", 31))
	h.Add(mkHeapObj("baz", 11))

	// Get works with the key.
	obj, exists := h.Get(mkHeapObj("baz", 0))
	if !exists || obj.(testHeapObject).val != 11 {
		t.Fatalf("unexpected error in getting element")
	}
	// Get non-existing object.
	_, exists = h.Get(mkHeapObj("non-existing", 0))
	if exists {
		t.Fatalf("didn't expect to get any object")
	}
}

// TestHeap_GetByKey tests Heap.GetByKey and is very similar to TestHeap_Get.
func TestHeap_GetByKey(t *testing.T) {
	h := heap_.NewHeap(testHeapObjectKeyFunc, compareInts)
	h.Add(mkHeapObj("foo", 10))
	h.Add(mkHeapObj("bar", 1))
	h.Add(mkHeapObj("bal", 31))
	h.Add(mkHeapObj("baz", 11))

	obj, exists := h.GetByKey("baz")
	if exists == false || obj.(testHeapObject).val != 11 {
		t.Fatalf("unexpected error in getting element")
	}
	// Get non-existing object.
	_, exists = h.GetByKey("non-existing")
	if exists {
		t.Fatalf("didn't expect to get any object")
	}
}

type Student struct {
	Id    string
	Name  string
	Score float32
}

func testStudentObjectKeyFunc(obj interface{}) string {
	return obj.(*Student).Id
}

//maxheap
func compareStudentScore(val1 interface{}, val2 interface{}) bool {
	first := val1.(*Student).Score
	second := val2.(*Student).Score
	return first > second
}

func TestMaxHeap(t *testing.T) {
	h := heap_.NewHeap(testStudentObjectKeyFunc, compareStudentScore)
	h.AddIfHeapOrder(&Student{
		Id:    "id_1",
		Name:  "name_1",
		Score: 89.2,
	})
	h.AddIfHeapOrder(&Student{
		Id:    "id_2",
		Name:  "name_2",
		Score: 87.2,
	})
	h.AddIfHeapOrder(&Student{
		Id:    "id_3",
		Name:  "name_3",
		Score: 97.2,
	})
	h.AddIfHeapOrder(&Student{
		Id:    "id_1",
		Name:  "name_3",
		Score: 88.2,
	})
	t.Logf("list key: %v", h.ListKeys())

	sz := len(h.ListKeys())
	for i := 0; i < sz; i++ {
		obj, err := h.Pop()
		if err != nil {
			t.Fatalf("failed to pop, err: %v", err)
		}
		t.Logf("get obj: %v", obj.(*Student))
	}
}
