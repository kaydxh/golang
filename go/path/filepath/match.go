package filepath

import (
	"os"
	"path/filepath"
)

func Glob(pattern string) (matches []string, err error) {
	matches, err = filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	matchFilters := make([]string, 0, len(matches))

	for _, path := range matches {
		_, err := os.Stat(path)
		if err != nil {
			continue
		}

		_, err = os.Lstat(path)
		if err != nil {
			continue
		}
		matchFilters = append(matchFilters, path)
	}

	return matchFilters, nil
}
