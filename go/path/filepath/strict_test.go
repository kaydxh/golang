package filepath_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	filepath_ "github.com/kaydxh/golang/go/path/filepath"
)

func TestStrictRegularFileAcceptsRealFile(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "dir"), 0o700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "dir", "file.txt")
	if err := os.WriteFile(path, []byte("body"), 0o600); err != nil {
		t.Fatal(err)
	}

	got, info, err := filepath_.StrictRegularFile(root, "dir/file.txt")
	if err != nil {
		t.Fatal(err)
	}
	if got != path || info.Size() != 4 || !info.Mode().IsRegular() {
		t.Fatalf("got %q %#v", got, info)
	}
}

func TestStrictRegularFileRejectsUnsafePaths(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "file"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, relative := range []string{
		"", ".", "..", "../file", " /file", "file ", "/file", `dir\file`,
		"dir//file", "dir/../file", "dir/./file",
	} {
		if _, _, err := filepath_.StrictRegularFile(root, relative); !errors.Is(err, filepath_.ErrUnsafePath) {
			t.Fatalf("path %q: err = %v", relative, err)
		}
	}
}

func TestStrictRegularFileRejectsSymlinksAndUnexpectedTypes(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "dir"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "file"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "dir"), filepath.Join(root, "link-dir")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "file"), filepath.Join(root, "dir", "link-file")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "missing"), filepath.Join(root, "dangling")); err != nil {
		t.Fatal(err)
	}

	tests := map[string]error{
		"link-dir/file": filepath_.ErrUnsafePath,
		"dir/link-file": filepath_.ErrUnsafePath,
		"dangling":      filepath_.ErrUnsafePath,
		"dir":           filepath_.ErrUnsafePath,
		"file/child":    filepath_.ErrUnsafePath,
		"missing":       filepath_.ErrPathNotFound,
	}
	for relative, want := range tests {
		if _, _, err := filepath_.StrictRegularFile(root, relative); !errors.Is(err, want) {
			t.Fatalf("path %q: err = %v, want %v", relative, err, want)
		}
	}
}
