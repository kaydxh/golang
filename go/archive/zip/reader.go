package zip

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

// ReadEntries sentinel errors.
var (
	ErrInvalidArchive = errors.New("zip: invalid archive")
	ErrUnsafePath     = errors.New("zip: unsafe entry path")
	ErrSymlink        = errors.New("zip: symbolic links are not allowed")
	ErrTooManyFiles   = errors.New("zip: too many files")
	ErrFileTooLarge   = errors.New("zip: file too large")
	ErrTotalTooLarge  = errors.New("zip: uncompressed content too large")
)

// ReadLimits bound in-memory ZIP processing.
type ReadLimits struct {
	MaxFiles      int
	MaxFileBytes  int64
	MaxTotalBytes int64
	RejectSymlink bool
}

// Entry is one regular file read from a ZIP archive.
type Entry struct {
	Name string
	Data []byte
	Mode os.FileMode
}

// ReadEntries reads regular files from a ZIP archive without extracting them.
// Names are normalized to slash-separated relative paths and traversal paths
// are rejected.
func ReadEntries(raw []byte, limits ReadLimits) ([]Entry, error) {
	reader, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidArchive, err)
	}
	entries := make([]Entry, 0, len(reader.File))
	var total int64
	for _, file := range reader.File {
		name, err := cleanRelativeName(file.Name)
		if err != nil {
			return nil, err
		}
		if name == "" || file.FileInfo().IsDir() {
			continue
		}
		mode := file.Mode()
		if limits.RejectSymlink && mode&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("%w: %q", ErrSymlink, file.Name)
		}
		if !mode.IsRegular() {
			return nil, fmt.Errorf("%w: entry %q is not a regular file", ErrInvalidArchive, file.Name)
		}
		if limits.MaxFiles > 0 && len(entries) >= limits.MaxFiles {
			return nil, fmt.Errorf("%w: maximum is %d", ErrTooManyFiles, limits.MaxFiles)
		}
		content, err := readEntry(file, limits.MaxFileBytes)
		if err != nil {
			return nil, err
		}
		total += int64(len(content))
		if limits.MaxTotalBytes > 0 && total > limits.MaxTotalBytes {
			return nil, fmt.Errorf("%w: maximum is %d bytes", ErrTotalTooLarge, limits.MaxTotalBytes)
		}
		entries = append(entries, Entry{Name: name, Data: content, Mode: mode})
	}
	return entries, nil
}

func cleanRelativeName(name string) (string, error) {
	name = strings.ReplaceAll(name, "\\", "/")
	name = strings.TrimPrefix(name, "./")
	for _, part := range strings.Split(name, "/") {
		if part == ".." {
			return "", fmt.Errorf("%w: %q", ErrUnsafePath, name)
		}
	}
	cleaned := path.Clean(name)
	if cleaned == "." || cleaned == "" {
		return "", nil
	}
	if path.IsAbs(cleaned) || cleaned == ".." || strings.HasPrefix(cleaned, "../") ||
		(len(cleaned) >= 2 && cleaned[1] == ':') {
		return "", fmt.Errorf("%w: %q", ErrUnsafePath, name)
	}
	return cleaned, nil
}

func readEntry(file *zip.File, maxBytes int64) ([]byte, error) {
	reader, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("%w: open entry %q: %v", ErrInvalidArchive, file.Name, err)
	}
	defer reader.Close()
	source := io.Reader(reader)
	if maxBytes > 0 {
		source = io.LimitReader(reader, maxBytes+1)
	}
	content, err := io.ReadAll(source)
	if err != nil {
		return nil, fmt.Errorf("%w: read entry %q: %v", ErrInvalidArchive, file.Name, err)
	}
	if maxBytes > 0 && int64(len(content)) > maxBytes {
		return nil, fmt.Errorf("%w: entry %q exceeds %d bytes", ErrFileTooLarge, file.Name, maxBytes)
	}
	return content, nil
}
