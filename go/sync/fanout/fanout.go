// Package fanout provides a small keyed event hub for in-process observers.
package fanout

import "sync"

// Hub fans events out to keyed subscribers. A zero buffer means the
// receive channel is unbuffered; a positive buffer bounds each subscriber's
// queue. The Hub carries no delivery guarantees: a full queue drops the
// event being published, not the subscriber.
type Hub[T any] struct {
	mu     sync.Mutex
	subs   map[string]map[chan T]struct{}
	buffer int
}

// New creates an empty hub. buffer is applied to every subscription and must
// not be negative.
func New[T any](buffer int) *Hub[T] {
	if buffer < 0 {
		panic("fanout: negative buffer")
	}
	return &Hub[T]{subs: make(map[string]map[chan T]struct{}), buffer: buffer}
}

// Subscribe returns subsequent events for key. The returned cancel function
// is idempotent and leaves already queued events readable until the channel
// is drained.
func (h *Hub[T]) Subscribe(key string) (<-chan T, func()) {
	if h == nil {
		closed := make(chan T)
		close(closed)
		return closed, func() {}
	}
	h.mu.Lock()
	if h.subs == nil {
		h.subs = make(map[string]map[chan T]struct{})
	}
	ch := make(chan T, h.bufferForLocked())
	if h.subs[key] == nil {
		h.subs[key] = make(map[chan T]struct{})
	}
	h.subs[key][ch] = struct{}{}
	h.mu.Unlock()

	var once sync.Once
	cancel := func() {
		once.Do(func() { h.unsubscribe(key, ch) })
	}
	return ch, cancel
}

// Publish delivers event to every current subscriber of key. If a subscriber's
// queue is full, this event is dropped for that subscriber and the
// subscription remains active. Publish never blocks and never races with a
// concurrent cancel: closing happens only while the same mutex is held.
func (h *Hub[T]) Publish(key string, event T) {
	if h == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.subs[key] {
		select {
		case ch <- event:
		default:
		}
	}
}

func (h *Hub[T]) unsubscribe(key string, ch chan T) {
	if h == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	set := h.subs[key]
	if _, ok := set[ch]; !ok {
		return
	}
	delete(set, ch)
	close(ch)
	if len(set) == 0 {
		delete(h.subs, key)
	}
}

func (h *Hub[T]) bufferForLocked() int {
	return h.buffer
}
