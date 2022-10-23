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
		f       func(context.Context) error
		loop    bool
	}{
		{
			name:    "test-sliding",
			period:  5 * time.Second,
			sliding: true,
			f: func(context.Context) error {
				time.Sleep(time.Second)
				fmt.Println("test-sliding")
				return nil
			},
			loop: true,
		},
		{
			name:    "test-nonsliding",
			sliding: false,
			period:  5 * time.Second,
			f: func(context.Context) error {
				time.Sleep(time.Second)
				fmt.Println("test-nonsliding")
				return nil
			},
			loop: true,
		},
		{
			name:    "test-nonsliding",
			sliding: false,
			period:  5 * time.Second,
			f: func(context.Context) error {
				time.Sleep(time.Second)
				fmt.Println("test-nonsliding")
				return nil
			},
			loop: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			//ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			//defer cancel()
			ctx := context.Background()

			time_.BackOffUntilWithContext(
				ctx,
				testCase.f,
				time_.NewExponentialBackOff(
					// forever run
					//time_.WithExponentialBackOffOptionMaxElapsedTime(0),
					time_.WithExponentialBackOffOptionMaxElapsedTime(20*time.Second),
					time_.WithExponentialBackOffOptionInitialInterval(testCase.period),
					time_.WithExponentialBackOffOptionMultiplier(1),
					time_.WithExponentialBackOffOptionRandomizationFactor(0),
				),
				testCase.sliding,
				true,
			)
			/*
				if err != nil {
					t.Fatalf("failed to write file: %v, got : %s", testCase.name, err)

				}
			*/

		})
	}
}

func TestRetryWithContext(t *testing.T) {
	testCases := []struct {
		name      string
		period    time.Duration
		retryTime int
		f         func(context.Context) error
	}{
		{
			name:      "test-retry-nil",
			period:    5 * time.Second,
			retryTime: 3,
			f: func(context.Context) error {
				time.Sleep(time.Second)
				fmt.Println("test-sliding")
				return nil
			},
		},
		{
			name:      "test-retry-error",
			period:    5 * time.Second,
			retryTime: 0,
			f: func(context.Context) error {
				time.Sleep(time.Second)
				fmt.Println("test-retry-error-sliding")
				return fmt.Errorf("error")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			//ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			//defer cancel()
			ctx := context.Background()

			err := time_.RetryWithContext(
				ctx,
				testCase.f,
				testCase.period,
				testCase.retryTime,
			)
			if err != nil {
				t.Fatalf("failed to call RetryWithContext: %v, got : %s", testCase.name, err)
			}

		})
	}
}
