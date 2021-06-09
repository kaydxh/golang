package time

import "time"

func WithExponentialBackOffOptionInitialInterval(initialInterval time.Duration) ExponentialBackOffOption {
	return ExponentialBackOffOptionFunc(func(opt *ExponentialBackOff) {
		opt.opts.InitialInterval = initialInterval
	})
}

func WithExponentialBackOffOptionRandomizationFactor(randomizationFactor float64) ExponentialBackOffOption {
	return ExponentialBackOffOptionFunc(func(opt *ExponentialBackOff) {
		opt.opts.RandomizationFactor = randomizationFactor
	})
}

func WithExponentialBackOffOptionMultiplier(multiplier float64) ExponentialBackOffOption {
	return ExponentialBackOffOptionFunc(func(opt *ExponentialBackOff) {
		opt.opts.Multiplier = multiplier
	})
}

func WithExponentialBackOffOptionMaxInterval(maxInterval time.Duration) ExponentialBackOffOption {
	return ExponentialBackOffOptionFunc(func(opt *ExponentialBackOff) {
		opt.opts.MaxInterval = maxInterval
	})
}

func WithExponentialBackOffOptionMaxElapsedTime(maxElapsedTime time.Duration) ExponentialBackOffOption {
	return ExponentialBackOffOptionFunc(func(opt *ExponentialBackOff) {
		opt.opts.MaxElapsedTime = maxElapsedTime
	})
}
