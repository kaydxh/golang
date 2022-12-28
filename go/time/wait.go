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
package time

import (
	"context"
	"errors"
	"fmt"
	"time"

	errors_ "github.com/kaydxh/golang/go/errors"
	runtime_ "github.com/kaydxh/golang/go/runtime"
	"github.com/sirupsen/logrus"
)

var ErrTimeout = errors.New("timeout error")

// Until loops until context timout, running f every period.
// Until is syntactic sugar on top of JitterUntil with zero jitter factor and
// with sliding = true (which means the timer for period starts after the f
// completes).
func UntilWithContxt(
	ctx context.Context,
	f func(ctx context.Context) error, period time.Duration) {
	JitterUntilWithContext(ctx, f, period)
}

func JitterUntilWithContext(
	ctx context.Context,
	f func(ctx context.Context) error,
	period time.Duration,
) {
	BackOffUntilWithContext(ctx, f,
		NewExponentialBackOff(
			// forever run
			WithExponentialBackOffOptionMaxElapsedTime(0),
			WithExponentialBackOffOptionInitialInterval(period),
			// ensure equal interval
			WithExponentialBackOffOptionMultiplier(1),
			WithExponentialBackOffOptionRandomizationFactor(0),
		), true, true)

}

// RetryWithContext retryTime is not include the first call
func RetryWithContext(
	ctx context.Context,
	f func(ctx context.Context) error,
	period time.Duration,
	retryTimes int,
) error {
	return BackOffUntilWithContext(ctx, f,
		NewExponentialBackOff(
			// forever run
			WithExponentialBackOffOptionMaxElapsedTime(0),
			WithExponentialBackOffOptionInitialInterval(period),
			// ensure equal interval
			WithExponentialBackOffOptionMultiplier(1),
			WithExponentialBackOffOptionRandomizationFactor(0),
			WithExponentialBackOffOptionMaxElapsedCount(retryTimes),
		), true, false)
}

// loop true  -> BackOffUntilWithContext return  until time expired
// loop false ->  BackOffUntilWithContext return if f return nil,  or time expired
func BackOffUntilWithContext(
	ctx context.Context,
	f func(ctx context.Context) error,
	backoff Backoff,
	sliding bool,
	loop bool,
) (err error) {
	var (
		t       time.Duration
		remain  time.Duration
		expired bool
	)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %v", ctx.Err())
		default:
		}

		tc := New(true)
		if !sliding {
			// If it is false then period includes the runtime for f
			t, expired = backoff.NextBackOff()
		}

		func() {
			defer runtime_.Recover()
			err = f(ctx)
			logrus.Infof("finish call function, err[%v]", err)
		}()

		if !loop {
			if err == nil {
				return nil
			}
		}

		if sliding {
			// If sliding is true, the period is computed after f runs
			tc.Reset()
			t, expired = backoff.NextBackOff()
		}
		if !expired {
			return fmt.Errorf("got max wait time or max count")
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
			case <-timer.C:
			}
		}()
	}
}

func CallWithTimeout(ctx context.Context, timeout time.Duration, f func(ctx context.Context) error) error {

	tc := New(true)
	// nerver timeout
	if timeout <= 0 {
		err := f(ctx)
		tc.Tick("call func")
		logrus.WithField("modulel", "CallWithTimeout").WithField("timeout", timeout).Infof("finish call function %v, err: %v", tc.String(), err)
		return err
	}

	var errs []error
	done := make(chan struct{}, 1)
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	go func() {
		err := f(ctx)
		if err != nil {
			errs = append(errs, err)
		}
		tc.Tick("call func")

		done <- struct{}{}
		logrus.WithField("modulel", "CallWithTimeout").WithField("timeout", timeout).Infof("finish call function %v, err: %v", tc.String(), err)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	case <-timer.C:
		return ErrTimeout
	}

	return errors_.NewAggregate(errs)

}
