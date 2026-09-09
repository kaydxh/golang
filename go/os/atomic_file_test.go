package os

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAtomicWriteFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")
	if err := AtomicWriteFile(path, []byte(`{"a":1}`), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil || string(got) != `{"a":1}` {
		t.Fatalf("content = %q err=%v", got, err)
	}
	// Overwriting replaces the content and leaves no temp files behind.
	if err := AtomicWriteFile(path, []byte(`{"a":2}`), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err = os.ReadFile(path)
	if err != nil || string(got) != `{"a":2}` {
		t.Fatalf("rewrite = %q err=%v", got, err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "out.json" {
		t.Fatalf("stray temp files: %v", entries)
	}
}
