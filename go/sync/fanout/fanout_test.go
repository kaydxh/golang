package fanout

import (
	"sync"
	"testing"
)

func TestHubBroadcastsByKey(t *testing.T) {
	t.Parallel()

	hub := New[int](8)
	first, cancelFirst := hub.Subscribe("issue")
	second, cancelSecond := hub.Subscribe("issue")
	other, cancelOther := hub.Subscribe("run")
	defer cancelFirst()
	defer cancelSecond()
	defer cancelOther()

	hub.Publish("issue", 7)
	if got := <-first; got != 7 {
		t.Fatalf("first = %d", got)
	}
	if got := <-second; got != 7 {
		t.Fatalf("second = %d", got)
	}
	select {
	case got := <-other:
		t.Fatalf("other key received %d", got)
	default:
	}
}

func TestFullBufferDropsNewEventAndKeepsSubscription(t *testing.T) {
	t.Parallel()

	hub := New[string](1)
	events, cancel := hub.Subscribe("issue")
	defer cancel()

	hub.Publish("issue", "first")
	hub.Publish("issue", "dropped")

	if got := <-events; got != "first" {
		t.Fatalf("first event = %q", got)
	}
	hub.Publish("issue", "later")
	if got := <-events; got != "later" {
		t.Fatalf("later event = %q", got)
	}
	cancel()
	if got, ok := <-events; ok {
		t.Fatalf("unexpected event %q", got)
	}
}

func TestCancelIsIdempotentAndDrainable(t *testing.T) {
	t.Parallel()

	hub := New[int](2)
	events, cancel := hub.Subscribe("issue")
	hub.Publish("issue", 1)
	cancel()
	cancel()
	hub.Publish("issue", 2)

	if got := <-events; got != 1 {
		t.Fatalf("queued event = %d", got)
	}
	if got, ok := <-events; ok {
		t.Fatalf("subscription received after cancel: %d", got)
	}
}

func TestConcurrentPublishAndCancel(t *testing.T) {
	t.Parallel()

	hub := New[int](32)
	var wg sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				events, cancel := hub.Subscribe("issue")
				if worker%2 == 0 {
					hub.Publish("issue", i)
				}
				cancel()
				for range events {
				}
			}
		}(worker)
	}
	wg.Wait()
}
