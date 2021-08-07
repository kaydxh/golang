package logs

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

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

type CallerPrettyfierFunc func(f *runtime.Frame) (function string, file string)

func GenShortCallPrettyfier() CallerPrettyfierFunc {
	return func(f *runtime.Frame) (function string, file string) {
		funcname := path.Base(f.Function)
		filename := path.Base(f.File)
		return fmt.Sprintf("%s()", funcname), fmt.Sprintf("%s:%d", filename, f.Line)
	}
}
