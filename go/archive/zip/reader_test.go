package zip

import (
	"archive/zip"
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestReadEntries(t *testing.T) {
	t.Parallel()
	var raw bytes.Buffer
	writer := zip.NewWriter(&raw)
	header := &zip.FileHeader{Name: "skill/scripts/run.sh", Method: zip.Deflate}
	header.SetMode(0o755)
	entry, err := writer.CreateHeader(header)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = entry.Write([]byte("#!/bin/sh\n"))
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	entries, err := ReadEntries(raw.Bytes(), ReadLimits{
		MaxFiles: 2, MaxFileBytes: 1024, MaxTotalBytes: 2048, RejectSymlink: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name != "skill/scripts/run.sh" ||
		entries[0].Mode&0o111 == 0 {
		t.Fatalf("entries = %+v", entries)
	}
}

func TestReadEntriesRejectsTraversalAndSymlink(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name string
		path string
		mode os.FileMode
		want error
	}{
		{name: "traversal", path: "../secret", mode: 0o644, want: ErrUnsafePath},
		{name: "nested traversal", path: "skill/../secret", mode: 0o644, want: ErrUnsafePath},
		{name: "windows absolute", path: `C:\secret`, mode: 0o644, want: ErrUnsafePath},
		{name: "symlink", path: "secret", mode: os.ModeSymlink | 0o777, want: ErrSymlink},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			var raw bytes.Buffer
			writer := zip.NewWriter(&raw)
			header := &zip.FileHeader{Name: test.path, Method: zip.Deflate}
			header.SetMode(test.mode)
			entry, err := writer.CreateHeader(header)
			if err != nil {
				t.Fatal(err)
			}
			_, _ = entry.Write([]byte("content"))
			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}
			if _, err := ReadEntries(raw.Bytes(), ReadLimits{RejectSymlink: true}); !errors.Is(err, test.want) {
				t.Fatalf("error = %v, want %v", err, test.want)
			}
		})
	}
}
