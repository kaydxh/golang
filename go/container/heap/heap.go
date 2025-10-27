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
package heap

import (
	"container/heap"
	"fmt"
	"sync"
)

const (
	closedMsg = "heap is closed"
)

// LessFunc is used to compare two objects in the heap.
type LessFunc func(interface{}, interface{}) bool

// KeyFunc knows how to make a key from an object. Implementations
// should be deterministic.
type KeyFunc func(interface{}) string

type heapItem struct {
	obj   interface{} // The object which is stored in the heap.
	index int         // The index of the object's key in the Heap.queue.
}

type itemKeyValue struct {
	key string
	obj interface{}
}

// heapData is an internal struct that implements the standard heap interface
// and keeps the data stored in the heap.
type heapData struct {
	// items is a map from key of the objects to the objects and their index.
	// We depend on the property that items in the map are in the queue and vice versa.
	items map[string]*heapItem
	// queue implements a heap data structure and keeps the order of elements
	// according to the heap invariant. The queue keeps the keys of objects stored
	// in "items".
	queue []string

	// keyFunc is used to make the key used for queued item insertion and retrieval, and
	// should be deterministic.
	keyFunc KeyFunc
	// lessFunc is used to compare two objects in the heap.
	lessFunc LessFunc
}

var (
	_ = heap.Interface(&heapData{}) // heapData is a standard heap
)

// Less compares two objects and returns true if the first one should go
// in front of the second one in the heap.
func (h *heapData) Less(i, j int) bool {
	if i > len(h.queue) || j > len(h.queue) {
		return false
	}
	itemi, ok := h.items[h.queue[i]]
	if !ok {
		return false
	}
	itemj, ok := h.items[h.queue[j]]
	if !ok {
		return false
	}
	return h.lessFunc(itemi.obj, itemj.obj)
}

// Len returns the number of items in the Heap.
func (h *heapData) Len() int { return len(h.queue) }

// Swap implements swapping of two elements in the heap. This is a part of standard
// heap interface and should never be called directly.
func (h *heapData) Swap(i, j int) {
	h.queue[i], h.queue[j] = h.queue[j], h.queue[i]
	item := h.items[h.queue[i]]
	item.index = i
	item = h.items[h.queue[j]]
	item.index = j
}

// Push is supposed to be called by heap.Push only.
func (h *heapData) Push(kv interface{}) {
	keyValue := kv.(*itemKeyValue)
	n := len(h.queue)
	h.items[keyValue.key] = &heapItem{keyValue.obj, n}
	h.queue = append(h.queue, keyValue.key)
}

// Pop is supposed to be called by heap.Pop only.
func (h *heapData) Pop() interface{} {
	key := h.queue[len(h.queue)-1]
	h.queue = h.queue[0 : len(h.queue)-1]
	item, ok := h.items[key]
	if !ok {
		// This is an error
		return nil
	}
	delete(h.items, key)
	return item.obj
}

// Heap is a thread-safe producer/consumer queue that implements a heap data structure.
// It can be used to implement priority queues and similar data structures.
type Heap struct {
	lock sync.RWMutex
	cond sync.Cond

	// data stores objects and has a queue that keeps their ordering according
	// to the heap invariant.
	data *heapData

	// closed indicates that the queue is closed.
	// It is mainly used to let Pop() exit its control loop while waiting for an item.
	closed bool
}

// Close the Heap and signals condition variables that may be waiting to pop
// items from the heap.
func (h *Heap) Close() {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.closed = true
	h.cond.Broadcast()
}

// Add inserts an item, and puts it in the queue. The item is updated if it
// already exists.
func (h *Heap) Add(obj interface{}) error {
	key := h.data.keyFunc(obj)
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.closed {
		return fmt.Errorf(closedMsg)
	}
	if _, exists := h.data.items[key]; exists {
		h.data.items[key].obj = obj
		heap.Fix(h.data, h.data.items[key].index)
	} else {
		h.addIfNotPresentLocked(key, obj)
	}
	h.cond.Broadcast()
	return nil
}

// BulkAdd adds all the items in the list to the queue and then signals the condition
// variable. It is useful when the caller would like to add all of the items
// to the queue before consumer starts processing them.
func (h *Heap) BulkAdd(list []interface{}) error {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.closed {
		return fmt.Errorf(closedMsg)
	}
	for _, obj := range list {
		key := h.data.keyFunc(obj)
		if _, exists := h.data.items[key]; exists {
			h.data.items[key].obj = obj
			heap.Fix(h.data, h.data.items[key].index)
		} else {
			h.addIfNotPresentLocked(key, obj)
		}
	}
	h.cond.Broadcast()
	return nil
}

// AddIfNotPresent inserts an item, and puts it in the queue. If an item with
// the key is present in the map, no changes is made to the item.
//
// This is useful in a single producer/consumer scenario so that the consumer can
// safely retry items without contending with the producer and potentially enqueueing
// stale items.
func (h *Heap) AddIfNotPresent(obj interface{}) error {
	id := h.data.keyFunc(obj)
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.closed {
		return fmt.Errorf(closedMsg)
	}
	h.addIfNotPresentLocked(id, obj)
	h.cond.Broadcast()
	return nil
}

// AddIfHeapOrder inserts an item, and puts it in the queue. If an item with
// the key is present in the map, and new obj meet the sort of heap,
// then update the item, or  no changes is made to the item.
func (h *Heap) AddIfHeapOrder(obj interface{}) error {
	key := h.data.keyFunc(obj)
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.closed {
		return fmt.Errorf(closedMsg)
	}
	if existedObj, exists := h.data.items[key]; exists {
		if !h.data.lessFunc(existedObj.obj, obj) {
			h.data.items[key].obj = obj
			heap.Fix(h.data, h.data.items[key].index)
		}
	} else {
		h.addIfNotPresentLocked(key, obj)
	}
	h.cond.Broadcast()
	return nil
}

// addIfNotPresentLocked assumes the lock is already held and adds the provided
// item to the queue if it does not already exist.
func (h *Heap) addIfNotPresentLocked(key string, obj interface{}) {
	if _, exists := h.data.items[key]; exists {
		return
	}
	heap.Push(h.data, &itemKeyValue{key, obj})
}

// Update is the same as Add in this implementation. When the item does not
// exist, it is added.
func (h *Heap) Update(obj interface{}) error {
	return h.Add(obj)
}

// Delete removes an item.
func (h *Heap) Delete(obj interface{}) error {
	key := h.data.keyFunc(obj)
	h.lock.Lock()
	defer h.lock.Unlock()
	if item, ok := h.data.items[key]; ok {
		heap.Remove(h.data, item.index)
		return nil
	}
	return fmt.Errorf("object not found")
}

// Pop waits until an item is ready. If multiple items are
// ready, they are returned in the order given by Heap.data.lessFunc.
func (h *Heap) Pop() (interface{}, error) {
	h.lock.Lock()
	defer h.lock.Unlock()
	for len(h.data.queue) == 0 {
		// When the queue is empty, invocation of Pop() is blocked until new item is enqueued.
		// When Close() is called, the h.closed is set and the condition is broadcast,
		// which causes this loop to continue and return from the Pop().
		if h.closed {
			return nil, fmt.Errorf("heap is closed")
		}
		h.cond.Wait()
	}
	obj := heap.Pop(h.data)
	if obj == nil {
		return nil, fmt.Errorf("object was removed from heap data")
	}

	return obj, nil
}

// List returns a list of all the items.
func (h *Heap) List() []interface{} {
	h.lock.RLock()
	defer h.lock.RUnlock()
	list := make([]interface{}, 0, len(h.data.items))
	for _, item := range h.data.items {
		list = append(list, item.obj)
	}
	return list
}

// ListKeys returns a list of all the keys of the objects currently in the Heap.
// Note: the key order is random, because it's data structure is map
func (h *Heap) ListKeys() []string {
	h.lock.RLock()
	defer h.lock.RUnlock()
	list := make([]string, 0, len(h.data.items))
	for key := range h.data.items {
		list = append(list, key)
	}
	return list
}

// Get returns the requested item, or sets exists=false.
// return item.obj, exists
func (h *Heap) Get(obj interface{}) (interface{}, bool) {
	key := h.data.keyFunc(obj)
	return h.GetByKey(key)
}

// GetByKey returns the requested item, or sets exists=false.
func (h *Heap) GetByKey(key string) (interface{}, bool) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	item, exists := h.data.items[key]
	if !exists {
		return nil, false
	}
	return item.obj, true
}

// IsClosed returns true if the queue is closed.
func (h *Heap) IsClosed() bool {
	h.lock.RLock()
	defer h.lock.RUnlock()
	if h.closed {
		return true
	}
	return false
}

// NewHeap returns a Heap which can be used to queue up items to process.
func NewHeap(keyFn KeyFunc, lessFn LessFunc) *Heap {
	h := &Heap{
		data: &heapData{
			items:    map[string]*heapItem{},
			queue:    []string{},
			keyFunc:  keyFn,
			lessFunc: lessFn,
		},
	}
	h.cond.L = &h.lock
	return h
}
