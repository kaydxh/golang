package time

import (
	"fmt"
	"strings"
	"time"
)

type TimeCounter struct {
	starts  []time.Time
	message []string
}

func (t *TimeCounter) Tick(msg string) {
	t.starts = append(t.starts, time.Now())
	t.message = append(t.message, msg)
}

func (t *TimeCounter) String() string {
	var buf strings.Builder
	t.Summary(func(idx int, msg string, cost time.Duration, at time.Time) {
		buf.WriteString(fmt.Sprintf("#%d, msg: %s, cost: %s, at %s\n", idx, msg, cost, at))
	})

	return buf.String()
}

func (t *TimeCounter) Summary(f func(idx int, msg string, cost time.Duration, at time.Time)) {
	if f == nil || t == nil {
		return
	}

	if len(t.message) < len(t.starts) {
		return
	}

	for i := 1; i < len(t.starts); i++ {
		f(i, t.message[i], t.starts[i].Sub(t.starts[i-1]), t.starts[i])
	}
}
