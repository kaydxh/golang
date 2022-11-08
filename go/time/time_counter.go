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
package time

import (
	"fmt"
	"strings"
	"time"
)

type TimeCounter struct {
	starts  []time.Time
	message []string
	effect  bool
}

func New(effect bool) *TimeCounter {
	t := &TimeCounter{
		effect: effect,
	}
	if effect {
		t.starts = append(t.starts, time.Now())
		t.message = append(t.message, "start")
	}

	return t
}

func (t *TimeCounter) Tick(msg string) {
	if t.effect {
		t.starts = append(t.starts, time.Now())
		t.message = append(t.message, msg)
	}
}

func (t *TimeCounter) Elapse() time.Duration {
	if !t.effect {
		return time.Duration(0)
	}

	if len(t.starts) == 0 {
		return time.Duration(0)
	}

	return time.Now().Sub(t.starts[0])
}

func (t *TimeCounter) String() string {
	if !t.effect {
		return ""
	}

	var buf strings.Builder
	t.Summary(func(idx int, msg string, cost time.Duration, at time.Time) {
		buf.WriteString(fmt.Sprintf("#%d, msg: %s, cost: %s, at %s ", idx, msg, cost, at.Format(time.RFC3339)))
	})

	return buf.String()
}

func (t *TimeCounter) Summary(f func(idx int, msg string, cost time.Duration, at time.Time)) {
	if !t.effect {
		return
	}
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

func (t *TimeCounter) Reset() {
	t.starts = nil
	t.message = nil
}
