package filecleanup

import (
	"time"
)

func WithDiskCheckInterval(interval time.Duration) DiskCleanerConfigOption {
	return DiskCleanerConfigOptionFunc(func(c *DiskCleanerConfig) {
		c.checkInterval = interval
	})
}

func WithDiskBaseExpired(expired time.Duration) DiskCleanerConfigOption {
	return DiskCleanerConfigOptionFunc(func(c *DiskCleanerConfig) {
		c.baseExpired = expired
	})
}

func WithDiskMinExpired(expired time.Duration) DiskCleanerConfigOption {
	return DiskCleanerConfigOptionFunc(func(c *DiskCleanerConfig) {
		c.minExpired = expired
	})
}
