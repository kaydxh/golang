package filepath

import (
	"path"
	"path/filepath"
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
