package time

import "time"

const (
	DayFormat           = "20060102"
	DefaultTimeFormat   = "2006-01-02 15:04:05"
	DefaultTimeMsFormat = "2006-01-02 15:04:05.000"
)

func Now() time.Time {
	now := time.Now()
	return now
}

func NowString(layout string) string {
	tm := Now()
	if layout == "" {
		layout = DefaultTimeFormat
	}

	return tm.Format(layout)
}

func BeginningOfDay(days int) time.Time {
	now := time.Now()
	theDay := now.AddDate(0, 0, days)
	y, m, d := theDay.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, theDay.Location())
}

func BeginningOfDayString(days int, layout string) string {
	now := time.Now()
	theDay := now.AddDate(0, 0, days)
	y, m, d := theDay.Date()

	if layout == "" {
		layout = DefaultTimeFormat
	}
	return time.Date(y, m, d, 0, 0, 0, 0, theDay.Location()).Format(layout)
}

func EndOfDay(days int) time.Time {
	now := time.Now()
	theDay := now.AddDate(0, 0, days)
	y, m, d := theDay.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), theDay.Location())
}

func EndOfDayString(days int, layout string) string {
	now := time.Now()
	theDay := now.AddDate(0, 0, days)
	y, m, d := theDay.Date()

	if layout == "" {
		layout = DefaultTimeFormat
	}
	return time.Date(
		y,
		m,
		d,
		23,
		59,
		59,
		int(time.Second-time.Nanosecond),
		theDay.Location(),
	).Format(layout)
}
