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
	tms := time_.TruncateToUTCString(now, time.Millisecond, time_.DefaultTimeMsFormat)
	tsc := time_.TruncateToUTCString(now, time.Second, time_.DefaultTimeMsFormat)
	tmt := time_.TruncateToUTCString(now, time.Minute, time_.DefaultTimeMsFormat)
	thr := time_.TruncateToUTCString(now, time.Hour, time_.DefaultTimeMsFormat)
	t.Logf("TruncateToUTC Millisecond: %v, Second: %v, Minute: %v, Hour: %v", tms, tsc, tmt, thr)
}
