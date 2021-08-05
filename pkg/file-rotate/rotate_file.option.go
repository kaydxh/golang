package rotatefile

import (
	"time"
)

func WithPrefixName(prefixName string) RotateFilerOption {
	return RotateFilerOptionFunc(func(c *RotateFiler) {
		c.opts.prefixName = prefixName
	})
}

func WithFileTimeLayout(fileTimeLayout string) RotateFilerOption {
	return RotateFilerOptionFunc(func(c *RotateFiler) {
		c.opts.fileTimeLayout = fileTimeLayout
	})
}

func WithSuffixName(subfixName string) RotateFilerOption {
	return RotateFilerOptionFunc(func(c *RotateFiler) {
		c.opts.subfixName = subfixName
	})
}

func WithMaxAge(maxAge time.Duration) RotateFilerOption {
	return RotateFilerOptionFunc(func(c *RotateFiler) {
		c.opts.maxAge = maxAge
	})
}

func WithMaxCount(maxCount int64) RotateFilerOption {
	return RotateFilerOptionFunc(func(c *RotateFiler) {
		c.opts.maxCount = maxCount
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
