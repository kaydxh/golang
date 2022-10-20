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
	"testing"
	"time"

	time_ "github.com/kaydxh/golang/go/time"
	"gotest.tools/assert"
)

func TestExponentialBackOff(t *testing.T) {
	var (
		testInitialInterval     = 500 * time.Millisecond
		testRandomizationFactor = 0.1
		testMultiplier          = 2.0
		testMaxInterval         = 5 * time.Second
		testMaxElapsedTime      = 15 * time.Minute
		testMaxElasedCount      = 1
	)

	exp := time_.NewExponentialBackOff(
		time_.WithExponentialBackOffOptionInitialInterval(testInitialInterval),
		time_.WithExponentialBackOffOptionRandomizationFactor(testRandomizationFactor),
		time_.WithExponentialBackOffOptionMultiplier(testMultiplier),
		time_.WithExponentialBackOffOptionMaxInterval(testMaxInterval),
		time_.WithExponentialBackOffOptionMaxElapsedTime(testMaxElapsedTime),
		time_.WithExponentialBackOffOptionMaxElapsedCount(testMaxElasedCount),
	)

	expectedResults := []time.Duration{500, 1000, 2000, 4000, 5000, 5000, 5000, 5000, 5000, 5000}
	for i, d := range expectedResults {
		expectedResults[i] = d * time.Millisecond
	}
	for _, expected := range expectedResults {
		assert.Equal(t, expected, exp.GetCurrentInterval())
		// Assert that the next backoff falls in the expected range.
		var minInterval = expected - time.Duration(testRandomizationFactor*float64(expected))
		var maxInterval = expected + time.Duration(testRandomizationFactor*float64(expected))
		actualInterval, over := exp.NextBackOff()
		t.Logf("over: %v, actualInterval: %v", over, actualInterval)
		if !(minInterval <= actualInterval && actualInterval <= maxInterval) {
			t.Error("error")
		}
	}

}

func TestExponentialBackOffMaxElaspedTimeFailOver(t *testing.T) {
	var (
		testInitialInterval     = 500 * time.Millisecond
		testRandomizationFactor = 0.1
		testMultiplier          = 2.0
		testMaxInterval         = 5 * time.Second
		testMaxElapsedTime      = 10 * time.Second
	)

	exp := time_.NewExponentialBackOff(
		time_.WithExponentialBackOffOptionInitialInterval(testInitialInterval),
		time_.WithExponentialBackOffOptionRandomizationFactor(testRandomizationFactor),
		time_.WithExponentialBackOffOptionMultiplier(testMultiplier),
		time_.WithExponentialBackOffOptionMaxInterval(testMaxInterval),
		time_.WithExponentialBackOffOptionMaxElapsedTime(testMaxElapsedTime),
	)

	expectedResults := []time.Duration{500, 1000, 2000, 4000, 5000, 5000, 5000, 5000, 5000, 5000}
	for i, d := range expectedResults {
		expectedResults[i] = d * time.Millisecond
	}
	for _, expected := range expectedResults {
		assert.Equal(t, expected, exp.GetCurrentInterval())
		// Assert that the next backoff falls in the expected range.
		var minInterval = expected - time.Duration(testRandomizationFactor*float64(expected))
		var maxInterval = expected + time.Duration(testRandomizationFactor*float64(expected))
		actualInterval, over := exp.NextBackOff()
		t.Logf("over: %v, actualInterval: %v", over, actualInterval)
		if !(minInterval <= actualInterval && actualInterval <= maxInterval) {
			t.Error("error")
		}
		time.Sleep(actualInterval)
	}

}

func TestDescExponentialBackOff(t *testing.T) {
	var (
		testInitialInterval     = 5 * time.Second
		testRandomizationFactor = 0.1
		testMultiplier          = 0.5
		testMaxInterval         = testInitialInterval
		testMinInterval         = 100 * time.Millisecond
		testMaxElapsedTime      = time.Duration(0)
	)

	exp := time_.NewExponentialBackOff(
		time_.WithExponentialBackOffOptionInitialInterval(testInitialInterval),
		time_.WithExponentialBackOffOptionRandomizationFactor(testRandomizationFactor),
		time_.WithExponentialBackOffOptionMultiplier(testMultiplier),
		time_.WithExponentialBackOffOptionMaxInterval(testMaxInterval),
		time_.WithExponentialBackOffOptionMinInterval(testMinInterval),
		time_.WithExponentialBackOffOptionMaxElapsedTime(testMaxElapsedTime),
	)

	expectedResults := []time.Duration{500, 250, 125, 62, 31, 16, 8, 4, 2, 2, 2}
	for i, d := range expectedResults {
		expectedResults[i] = d * time.Millisecond
	}
	for _, expected := range expectedResults {
		//	assert.Equal(t, expected, exp.GetCurrentInterval())
		// Assert that the next backoff falls in the expected range.
		var minInterval = expected - time.Duration(testRandomizationFactor*float64(expected))
		var maxInterval = expected + time.Duration(testRandomizationFactor*float64(expected))
		actualInterval, over := exp.NextBackOff()
		t.Logf("over: %v, actualInterval: %v", over, actualInterval)
		if !(minInterval <= actualInterval && actualInterval <= maxInterval) {
		}
	}

	t.Logf("starting back...")
	for _, expected := range expectedResults {
		//	assert.Equal(t, expected, exp.GetCurrentInterval())
		// Assert that the next backoff falls in the expected range.
		var minInterval = expected - time.Duration(testRandomizationFactor*float64(expected))
		var maxInterval = expected + time.Duration(testRandomizationFactor*float64(expected))
		actualInterval, over := exp.PreBackOff()
		t.Logf("over: %v, actualInterval: %v", over, actualInterval)
		if !(minInterval <= actualInterval && actualInterval <= maxInterval) {
		}
	}

}
