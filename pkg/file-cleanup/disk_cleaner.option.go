package filecleanup

import (
	"time"
)

func WithDiskCheckInterval(interval time.Duration) DiskCleanerConfigOption {
	return DiskCleanerConfigOptionFunc(func(c *DiskCleanerConfig) {
		c.checkInterval = interval
	})
}
