package time

import (
	"context"
	"time"

	runtime_ "github.com/kaydxh/golang/go/runtime"
)

// Until loops until context timout, running f every period.
// Until is syntactic sugar on top of JitterUntil with zero jitter factor and
// with sliding = true (which means the timer for period starts after the f
// completes).
func UntilWithContxt(
	ctx context.Context,
	f func(ctx context.Context), period time.Duration) {
	JitterUntilWithContext(ctx, f, period, nil)
}

func JitterUntilWithContext(
	ctx context.Context,
	f func(ctx context.Context),
	period time.Duration,
	stopCh <-chan struct{},
) {
	BackOffUntilWithContext(ctx, f,
		NewExponentialBackOff(
			// forever run
			WithExponentialBackOffOptionMaxElapsedTime(0),
			WithExponentialBackOffOptionInitialInterval(period),
			// ensure equal interval
			WithExponentialBackOffOptionMultiplier(1),
			WithExponentialBackOffOptionRandomizationFactor(0),
		), true, stopCh)

}

func BackOffUntilWithContext(
	ctx context.Context,
	f func(ctx context.Context),
	backoff Backoff,
	sliding bool,
	stopCh <-chan struct{},
) {
	var (
		t       time.Duration
		remain  time.Duration
		expired bool
	)

	for {
		select {
		case <-ctx.Done():
			return
		case <-stopCh:
			return
		default:
		}

		tc := New(true)
		if !sliding {
			// If it is false then period includes the runtime for f
			t, expired = backoff.NextBackOff()
		}

		func() {
			defer runtime_.Recover()
			f(ctx)
		}()

		if sliding {
			// If sliding is true, the period is computed after f runs
			tc.Reset()
			t, expired = backoff.NextBackOff()
		}
		if !expired {
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
			case <-ctx.Done():
				return
			case <-stopCh:
				return
			case <-timer.C:
			}
		}()
	}
}
