package time_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	time_ "github.com/kaydxh/golang/go/time"
)

func TestBackOffUntilWithContext(t *testing.T) {
	testCases := []struct {
		name    string
		period  time.Duration
		sliding bool
		f       func(context.Context)
	}{
		{
			name:    "test-sliding",
			period:  5 * time.Second,
			sliding: true,
			f: func(context.Context) {
				time.Sleep(time.Second)
				fmt.Println("test-sliding")
			},
		},
		{
			name:    "test-nonsliding",
			sliding: false,
			period:  5 * time.Second,
			f: func(context.Context) {
				fmt.Println("test-nonsliding")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var stopCh chan struct{}
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			time_.BackOffUntilWithContext(
				ctx,
				testCase.f,
				time_.NewExponentialBackOff(
					// forever run
					//time_.WithExponentialBackOffOptionMaxElapsedTime(0),
					time_.WithExponentialBackOffOptionMaxElapsedTime(time.Minute),
					time_.WithExponentialBackOffOptionInitialInterval(testCase.period),
					time_.WithExponentialBackOffOptionMultiplier(1),
					time_.WithExponentialBackOffOptionRandomizationFactor(0),
				),
				testCase.sliding,
				stopCh,
			)
			/*
				if err != nil {
					t.Fatalf("failed to write file: %v, got : %s", testCase.name, err)

				}
			*/

		})
	}
}
