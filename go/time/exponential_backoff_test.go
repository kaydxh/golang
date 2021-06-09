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
		actualInterval, maxElapsed := exp.NextBackOff()
		t.Logf("maxElapsed: %v, actualInterval: %v", maxElapsed, actualInterval)
		if !(minInterval <= actualInterval && actualInterval <= maxInterval) {
			t.Error("error")
		}

	}

}
