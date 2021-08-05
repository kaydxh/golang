package logs

import "time"

func WithPrefixName(prefixName string) RotateOption {
	return RotateOptionFunc(func(c *Rotate) {
		c.prefixName = prefixName
	})
}

func WithSuffixName(suffixName string) RotateOption {
	return RotateOptionFunc(func(c *Rotate) {
		c.suffixName = suffixName
	})
}

func WithMaxAge(maxAge time.Duration) RotateOption {
	return RotateOptionFunc(func(c *Rotate) {
		c.maxAge = maxAge
	})
}

func WithMaxCount(maxCount int64) RotateOption {
	return RotateOptionFunc(func(c *Rotate) {
		c.maxCount = maxCount
	})
}

func WithRotateSize(rotateSize int64) RotateOption {
	return RotateOptionFunc(func(c *Rotate) {
		c.rotateSize = rotateSize
	})
}

func WithRotateInterval(rotateInterval time.Duration) RotateOption {
	return RotateOptionFunc(func(c *Rotate) {
		c.rotateInterval = rotateInterval
	})
}
