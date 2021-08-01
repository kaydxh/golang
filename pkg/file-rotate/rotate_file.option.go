package rotatefile

import (
	"time"
)

func WithMaxAge(maxAge time.Duration) RotateFilerOption {
	return RotateFilerOptionFunc(func(c *RotateFiler) {
		c.opts.maxAge = maxAge
	})
}

func WithRotateSize(rotateSize int64) RotateFilerOption {
	return RotateFilerOptionFunc(func(c *RotateFiler) {
		c.opts.rotateSize = rotateSize
	})
}

func WithRotateInterval(rotateInterval time.Duration) RotateFilerOption {
	return RotateFilerOptionFunc(func(c *RotateFiler) {
		c.opts.rotateInterval = rotateInterval
	})
}

func WithPrefixName(prefixName string) RotateFilerOption {
	return RotateFilerOptionFunc(func(c *RotateFiler) {
		c.opts.prefixName = prefixName
	})
}

func WithSuffixName(subfixName string) RotateFilerOption {
	return RotateFilerOptionFunc(func(c *RotateFiler) {
		c.opts.subfixName = subfixName
	})
}
