package time

import "time"

func BeginningOfDay(days int) time.Time {
	now := time.Now()
	theDay := now.AddDate(0, 0, days)
	y, m, d := theDay.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, theDay.Location())
}

func EndOfDay(days int) time.Time {
	now := time.Now()
	theDay := now.AddDate(0, 0, days)
	y, m, d := theDay.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), theDay.Location())
}

func BeginningOfDayString(days int, layout string) string {
	now := time.Now()
	theDay := now.AddDate(0, 0, days)
	y, m, d := theDay.Date()

	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	return time.Date(y, m, d, 0, 0, 0, 0, theDay.Location()).Format(layout)
}

func EndOfDayString(days int, layout string) string {
	now := time.Now()
	theDay := now.AddDate(0, 0, days)
	y, m, d := theDay.Date()

	if layout == "" {
		layout = "2006-01-02 15:04:05"
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
