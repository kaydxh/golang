package os

import (
	"os"

	errors_ "github.com/kaydxh/golang/go/errors"
	filepath_ "github.com/kaydxh/golang/go/path/filepath"
)

func RemoveBatch(filenames []string) error {
	var errs []error
	for _, path := range filenames {
		err := os.Remove(path)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors_.NewAggregate(errs)
}

func RemoveWithGlob(pattern string) error {
	matches, err := filepath_.Glob(pattern)
	if err != nil {
		return err
	}

	return RemoveBatch(matches)

}
