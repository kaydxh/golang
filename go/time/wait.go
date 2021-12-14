package time

import (
	"context"
	"time"
)

func BackOffUntilWithContext(ctx context.Context, f func(), backoff Backoff, sliding bool, stopCh <-chan struct{}) {
	var (
		t      time.Duration
		remain time.Duration
		ok     bool
	)

	for {
		select {
		case <-stopCh:
			return
		default:
		}

		tc := New(true)
		if !sliding {
			// If it is false then period includes the runtime for f
			t, ok = backoff.NextBackOff()
		}

		func() {
			f()
		}()

		if sliding {
			// If sliding is true, the period is computed after f runs
			tc.Reset()
			t, ok = backoff.NextBackOff()
		}
		if !ok {
			return
		}

		remain = t - tc.Elapse()
		//	fmt.Printf("remain: %v, data: %v\n", remain, time.Now().String())

		func() {
			if remain <= 0 {
				return
			}
			timer := time.NewTimer(remain)
			defer timer.Stop()

			// NOTE: b/c there is no priority selection in golang
			// it is possible for this to race, meaning we could
			// trigger t.C and stopCh, and t.C select falls through.
			// In order to mitigate we re-check stopCh at the beginning
			// of every loop to prevent extra executions of f().

			select {
			case <-stopCh:
				return
			case <-timer.C:
			}
		}()
	}
}
