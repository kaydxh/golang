/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
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

		// maxAge 0 is unlimited
		if cleaner.maxAge <= 0 {
			continue
		}

		if now.Sub(fi.ModTime()) < cleaner.maxAge {
			continue
		}

		removeMatches = append(removeMatches, path)
	}

	if cleaner.maxCount > 0 {
		if cleaner.maxCount < int64(len(matches)) {

			removeCount := len(matches) - int(cleaner.maxCount) - len(removeMatches)
			if removeCount > 0 {
				sort.Sort(RotatedFiles(matches))
				removeMatches = append(
					removeMatches,
					matches[len(removeMatches):len(removeMatches)+removeCount]...,
				)
			}
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
