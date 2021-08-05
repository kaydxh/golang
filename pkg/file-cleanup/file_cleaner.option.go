package filecleanup

import (
	"time"
)

func WithMaxAge(maxAge time.Duration) FileCleanerOption {
	return FileCleanerOptionFunc(func(c *FileCleaner) {
		c.maxAge = maxAge
	})
}

func WithMaxCount(maxCount int64) FileCleanerOption {
	return FileCleanerOptionFunc(func(c *FileCleaner) {
		c.maxCount = maxCount
	})
}
