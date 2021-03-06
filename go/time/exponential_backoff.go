package time

import (
	"math/rand"
	"time"
)

// Default values for ExponentialBackOff.
const (
	DefaultInitialInterval     = 500 * time.Millisecond
	DefaultRandomizationFactor = 0.5
	// The default multiplier value used for increment current interval
	DefaultMultiplier     = 1.5
	DefaultMaxInterval    = 60 * time.Second
	DefaultMaxElapsedTime = 15 * time.Minute
)

type ExponentialBackOff struct {
	currentInterval time.Duration
	startTime       time.Time

	opts struct {
		InitialInterval     time.Duration
		RandomizationFactor float64
		Multiplier          float64
		MaxInterval         time.Duration
		// After MaxElapsedTime the ExponentialBackOff returns Stop.
		// It never stops if MaxElapsedTime == 0.
		MaxElapsedTime time.Duration
	}
}

func NewExponentialBackOff(opts ...ExponentialBackOffOption) *ExponentialBackOff {
	bo := &ExponentialBackOff{}
	bo.opts.InitialInterval = DefaultInitialInterval
	bo.opts.RandomizationFactor = DefaultRandomizationFactor
	bo.opts.Multiplier = DefaultMultiplier

	bo.ApplyOptions(opts...)
	bo.Reset()
	return bo
}

func (b *ExponentialBackOff) Reset() {
	b.currentInterval = b.opts.InitialInterval
	b.startTime = time.Now()
}

func (b *ExponentialBackOff) GetCurrentInterval() time.Duration {
	return b.currentInterval
}

// false : have gone over the maximu elapsed time
// true : return remaining time
func (b *ExponentialBackOff) NextBackOff() (time.Duration, bool) {
	elapsed := b.GetElapsedTime()
	nextRandomizedInterval := getRandomValueFromInterval(b.opts.RandomizationFactor, b.currentInterval)

	if b.opts.MaxElapsedTime > 0 && elapsed > b.opts.MaxElapsedTime {
		return b.currentInterval, false
	}

	//update currentInterval
	b.incrementCurrentInterval()

	return nextRandomizedInterval, true
}

func (b *ExponentialBackOff) GetElapsedTime() time.Duration {
	return time.Now().Sub(b.startTime)
}

// Increments the current interval by multiplying it with the multiplier
func (b *ExponentialBackOff) incrementCurrentInterval() {
	if b.opts.MaxInterval > 0 && time.Duration(float64(b.currentInterval)*b.opts.Multiplier) > b.opts.MaxInterval {
		b.currentInterval = b.opts.MaxInterval
		return
	}

	b.currentInterval = time.Duration(float64(b.currentInterval) * b.opts.Multiplier)
}

func getRandomValueFromInterval(
	randomizationFactor float64,
	currentInterval time.Duration,
) time.Duration {
	var delta = randomizationFactor * float64(currentInterval)
	var minInterval = float64(currentInterval) - delta
	var maxInterval = float64(currentInterval) + delta

	// Get a random value from the range [minInterval, maxInterval].
	// The formula used below has a +1 because if the minInterval is 1 and the maxInterval is 3 then
	// we want a 33% chance for selecting either 1, 2 or 3.
	//Float64 returns, as a float64, a pseudo-random number in [0.0,1.0)
	//from the default Source.
	return time.Duration(minInterval + (rand.Float64() * (maxInterval - minInterval + 1)))

}
