package time

import "time"

const (
	DayFormat           = "20060102"
	ShortTimeFormat     = "20060102150405"
	ShortDashTimeFormat = "2006-01-02-15:04:05"
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

// Truncate only happens in UTC semantics, apparently.
// observed values for truncating given time with 86400 secs:
//
// before truncation: 2018/06/01 03:54:54 2018-06-01T03:18:00+09:00
// after  truncation: 2018/06/01 03:54:54 2018-05-31T09:00:00+09:00
//
// This is really annoying when we want to truncate in local time
// so we hack: we take the apparent local time in the local zone,
// and pretend that it's in UTC. do our math, and put it back to
// the local zone
func TruncateToUTC(t time.Time, d time.Duration) time.Time {
	if t.Location() == time.UTC {
		return t.Truncate(d)
	}

	base := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
	base = base.Truncate(d)

	return time.Date(
		base.Year(),
		base.Month(),
		base.Day(),
		base.Hour(),
		base.Minute(),
		base.Second(),
		base.Nanosecond(),
		base.Location(),
	)
}

func TruncateToUTCString(t time.Time, d time.Duration, layout string) string {
	utc := TruncateToUTC(t, d)

	if layout == "" {
		layout = DefaultTimeFormat
	}

	return utc.Format(layout)
}
