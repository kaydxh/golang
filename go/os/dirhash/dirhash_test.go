package dirhash

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestDigestTreeDeterministic(t *testing.T) {
	tree := map[string]string{"a/b.go": "package b\n", "a.txt": "x", "z/c.go": "package c\n"}
	first, err := DigestTree(writeTree(t, tree), Options{})
	if err != nil {
		t.Fatal(err)
	}
	second, err := DigestTree(writeTree(t, tree), Options{})
	if err != nil {
		t.Fatal(err)
	}
	if first == "" || first != second {
		t.Fatalf("digests not deterministic: %q vs %q", first, second)
	}
}

func TestDigestTreeTracksContent(t *testing.T) {
	root := writeTree(t, map[string]string{"a.go": "package a\n"})
	before, err := DigestTree(root, Options{})
	if err != nil {
		t.Fatal(err)
	}
	// Dirty a file: the digest must move.
	if err := os.WriteFile(filepath.Join(root, "a.go"), []byte("package a // changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	after, err := DigestTree(root, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if before == after {
		t.Fatal("content change must change the digest")
	}
	// Adding a file moves it too.
	if err := os.WriteFile(filepath.Join(root, "b.go"), []byte("package a\n\nvar B = 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	grown, err := DigestTree(root, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if grown == after {
		t.Fatal("file addition must change the digest")
	}
}

func TestDigestTreeExclusions(t *testing.T) {
	root := writeTree(t, map[string]string{
		"src/main.go":     "package main\n",
		"node_modules/x":  "dep",
		"dist/bundle.js":  "built",
		".git/config":     "git",
		".DS_Store":       "junk",
		"src/.DS_Store":   "junk",
	})
	digest, err := DigestTree(root, Options{
		ExcludedDirs:  map[string]bool{"node_modules": true, "dist": true, ".git": true},
		ExcludedNames: map[string]bool{".DS_Store": true},
	})
	if err != nil {
		t.Fatal(err)
	}
	// The same tree with the excluded files deleted must produce the same
	// digest: exclusions are not part of the identity.
	clean := writeTree(t, map[string]string{"src/main.go": "package main\n"})
	cleanDigest, err := DigestTree(clean, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if digest != cleanDigest {
		t.Fatalf("excluded files leaked into digest: %q vs %q", digest, cleanDigest)
	}
}

func TestDigestTreeMaxHashBytes(t *testing.T) {
	root := t.TempDir()
	big := make([]byte, 1024)
	if err := os.WriteFile(filepath.Join(root, "big.bin"), big, 0o644); err != nil {
		t.Fatal(err)
	}
	capped, err := DigestTree(root, Options{MaxHashBytes: 512})
	if err != nil {
		t.Fatal(err)
	}
	empty, err := DigestTree(t.TempDir(), Options{})
	if err != nil {
		t.Fatal(err)
	}
	if capped != empty {
		t.Fatal("oversized files must be skipped, not hashed")
	}
}

func TestDigestTreeEmptyRoot(t *testing.T) {
	digest, err := DigestTree(t.TempDir(), Options{})
	if err != nil || digest == "" {
		t.Fatalf("empty root digest = %q err=%v", digest, err)
	}
}
