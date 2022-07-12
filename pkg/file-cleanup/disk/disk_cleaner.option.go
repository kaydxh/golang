package disk

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

func WithDiskUsageCallBack(f func(diskPath string, diskUsage float32)) DiskCleanerConfigOption {
	return DiskCleanerConfigOptionFunc(func(c *DiskCleanerConfig) {
		c.diskUsageCallBack = f
	})
}

func WithCleanPostCallBack(f func(file string, err error)) DiskCleanerConfigOption {
	return DiskCleanerConfigOptionFunc(func(c *DiskCleanerConfig) {
		c.cleanPostCallBack = f
	})
}
