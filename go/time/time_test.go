package time_test

import (
	"testing"

	time_ "github.com/kaydxh/golang/go/time"
)

func TestBeginningOfDayString(t *testing.T) {
	beginTime := time_.BeginningOfDayString(-1, "")
	t.Logf(beginTime)
}

func TestEndOfDayString(t *testing.T) {
	endTime := time_.EndOfDayString(-1, "")
	t.Logf(endTime)
}
