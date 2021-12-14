package time

import "time"

type Backoff interface {
	NextBackOff() (time.Duration, bool)
	Reset()
}
