package filepath

import (
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

func GetParentRelPath(filePath string) string {
	fileName := path.Base(filePath)
	parentDir := GetParentRelDir(filePath)
	return filepath.Join(parentDir, fileName)
}

func GetParentRelDir(filePath string) string {
	dir := path.Dir(filePath)
	return path.Base(dir)
}

func LastChar(str string) uint8 {
	if str == "" {
		// char '0' means dec 48, 0 means null
		return 0
	}

	return str[len(str)-1]
}

// join paths, keep relativePath suffix
func JoinPaths(rootPath, relativePath string) (string, error) {
	absolutePath, err := filepath.Abs(rootPath)
	if err != nil {
		return "", err
	}

	if relativePath == "" {
		return absolutePath, nil
	}

	finalPath := path.Join(absolutePath, relativePath)
	if LastChar(relativePath) == '/' && LastChar(finalPath) != '/' {
		return finalPath + "/", nil
	}

	return finalPath, nil
}

// CanonicalizePath turns path into an absolute path without symlinks.
func CanonicalizePath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	path, err = filepath.EvalSymlinks(path)

	// Get a better error if we have an invalid path
	if pathErr, ok := err.(*os.PathError); ok {
		err = errors.Wrap(pathErr.Err, pathErr.Path)
	}

	return path, err
}
