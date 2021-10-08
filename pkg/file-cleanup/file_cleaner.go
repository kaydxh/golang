package filecleanup

import (
	"os"
	"sort"
	"time"

	errors_ "github.com/kaydxh/golang/go/errors"
	filepath_ "github.com/kaydxh/golang/go/path/filepath"
)

type FileCleaner struct {
	//filedir string
	//pattern string
	//maxAge is the maximum number of time to retain old files, 0 is unlimited
	maxAge time.Duration
	//maxCount is the maximum number to retain old files, 0 is unlimited
	maxCount int64
}

func FileCleanup(pattern string, options ...FileCleanerOption) error {
	var cleaner FileCleaner
	cleaner.ApplyOptions(options...)

	matches, err := filepath_.Glob(pattern)
	if err != nil {
		return err
	}

	now := time.Now()

	removeMatches := make([]string, 0, len(matches))
	for _, path := range matches {
		fi, err := os.Stat(path)
		if err != nil {
			continue
		}

		if cleaner.maxAge > 0 && now.Sub(fi.ModTime()) < cleaner.maxAge {
			continue
		}

		removeMatches = append(removeMatches, path)
	}

	if cleaner.maxCount > 0 {
		if cleaner.maxCount < int64(len(matches)) {
			sort.Sort(RotatedFiles(matches))
			removeMatches = append(
				removeMatches,
				matches[len(removeMatches):len(matches)-int(cleaner.maxCount)-len(removeMatches)]...,
			)
		}

	}

	//clean
	var errs []error
	for _, file := range removeMatches {
		err = os.Remove(file)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors_.NewAggregate(errs)
}

type RotatedFiles []string

func (f RotatedFiles) Len() int {
	return len(f)
}

func (f RotatedFiles) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f RotatedFiles) Less(i, j int) bool {
	fi, err := os.Stat(f[i])
	if err != nil {
		return false
	}

	fj, err := os.Stat(f[j])
	if err != nil {
		return false
	}

	if fi.ModTime().Equal(fj.ModTime()) {
		if len(f[i]) == len(f[j]) {
			return f[i] < f[j]
		}

		return len(f[i]) < len(f[j]) //  foo.9  < foo.10
	}

	return fi.ModTime().Before(fj.ModTime())

}
