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
package time_test

import (
	"testing"
	"time"

	time_ "github.com/kaydxh/golang/go/time"
)

func TestNowString(t *testing.T) {
	now := time_.NowString("")
	t.Logf(now)
}

func TestBeginningOfDayString(t *testing.T) {
	beginTime := time_.BeginningOfDayString(-1, "")
	t.Logf(beginTime)
}

func TestEndOfDayString(t *testing.T) {
	endTime := time_.EndOfDayString(-1, "")
	t.Logf(endTime)
}

func TestTruncateToUTCString(t *testing.T) {
	now := time.Now()
	tms := time_.TruncateToUTCString(now, time.Millisecond, time_.DefaultTimeMillFormat)
	tsc := time_.TruncateToUTCString(now, time.Second, time_.DefaultTimeMillFormat)
	tmt := time_.TruncateToUTCString(now, time.Minute, time_.DefaultTimeMillFormat)
	thr := time_.TruncateToUTCString(now, time.Hour, time_.DefaultTimeMillFormat)
	t.Logf("TruncateToUTC Millisecond: %v, Second: %v, Minute: %v, Hour: %v", tms, tsc, tmt, thr)
}

func TestNowFormat(t *testing.T) {
	now := time.Now()
	tm := now.Format(time_.ShortDashTimeFormat)
	t.Logf("Now: %v", tm)
}
