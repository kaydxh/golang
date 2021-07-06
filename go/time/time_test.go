package time_test

import (
	"testing"

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
