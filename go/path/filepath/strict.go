package filepath

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrUnsafePath means the path or one of its components cannot be used
	// without following a symlink or accepting an unexpected file type.
	ErrUnsafePath = errors.New("unsafe path")
	// ErrPathNotFound means a requested path component does not exist.
	ErrPathNotFound = errors.New("path not found")
)

// StrictRegularFile resolves a slash-separated relative path under root and
// returns the inspected regular file. It refuses absolute paths, Windows
// volume names, backslashes, NUL, non-canonical paths, symlinks, non-directory
// intermediate components, and non-regular final components. The returned path
// is still subject to replacement between inspection and use; callers that
// need stronger guarantees must open through a root-restricted file handle.
func StrictRegularFile(root, relative string) (string, fs.FileInfo, error) {
	if err := validateStrictRelativePath(relative); err != nil {
		return "", nil, err
	}

	rootInfo, err := os.Lstat(root)
	if err != nil || rootInfo.Mode()&os.ModeSymlink != 0 || !rootInfo.IsDir() {
		return "", nil, fmt.Errorf("%w: root is not a real directory: %v", ErrUnsafePath, err)
	}

	current := root
	parts := strings.Split(relative, "/")
	for index, part := range parts {
		current = filepath.Join(current, filepath.FromSlash(part))
		info, statErr := os.Lstat(current)
		if errors.Is(statErr, os.ErrNotExist) {
			return "", nil, fmt.Errorf("%w: %s", ErrPathNotFound, relative)
		}
		if statErr != nil {
			return "", nil, fmt.Errorf("inspect %s: %w", relative, statErr)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", nil, fmt.Errorf("%w: symlink component %q", ErrUnsafePath, relative)
		}
		if index < len(parts)-1 {
			if !info.IsDir() {
				return "", nil, fmt.Errorf("%w: non-directory component %q", ErrUnsafePath, relative)
			}
			continue
		}
		if !info.Mode().IsRegular() {
			return "", nil, fmt.Errorf("%w: non-regular file %q", ErrUnsafePath, relative)
		}
		return current, info, nil
	}
	return "", nil, fmt.Errorf("%w: empty path", ErrUnsafePath)
}

func validateStrictRelativePath(relative string) error {
	if relative == "" || relative != strings.TrimSpace(relative) ||
		strings.Contains(relative, `\`) || strings.ContainsRune(relative, 0) ||
		filepath.IsAbs(relative) || filepath.VolumeName(relative) != "" ||
		filepath.ToSlash(filepath.Clean(relative)) != relative ||
		relative == "." || relative == ".." || strings.HasPrefix(relative, "../") {
		return fmt.Errorf("%w: invalid relative path %q", ErrUnsafePath, relative)
	}
	return nil
}
