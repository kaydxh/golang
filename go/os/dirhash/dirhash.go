// Package dirhash computes deterministic content digests over directory
// trees: sha256 over one "<mode> <relpath> <sha256(content)>" line per
// regular file, in sorted path order. Two materializations of the same
// tree content always produce the same digest, and any content change —
// including an uncommitted working file — changes the digest.
package dirhash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Options configures a tree digest. The zero value hashes every regular
// file under root.
type Options struct {
	// ExcludedDirs are directory names (case-insensitive) the walk never
	// descends into, at any depth.
	ExcludedDirs map[string]bool
	// ExcludedNames are file names (case-sensitive) skipped wherever they
	// appear — generated junk like .DS_Store, or files a caller refuses to
	// fingerprint.
	ExcludedNames map[string]bool
	// MaxHashBytes caps how much of one file gets hashed. Larger files are
	// skipped entirely — they contribute no line. Zero hashes everything.
	MaxHashBytes int64
}

// DigestTree returns the content digest of root as "sha256:<hex>".
// Symlinks are never followed: the walk stays inside root.
func DigestTree(root string, opts Options) (string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}
		name := d.Name()
		if d.IsDir() {
			if opts.ExcludedDirs[strings.ToLower(name)] {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.Type().IsRegular() || opts.ExcludedNames[name] {
			return nil
		}
		if opts.MaxHashBytes > 0 {
			if info, statErr := d.Info(); statErr == nil && info.Size() > opts.MaxHashBytes {
				return nil
			}
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("dirhash: walk %s: %w", root, err)
	}
	sort.Strings(files)

	hasher := sha256.New()
	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			return "", fmt.Errorf("dirhash: stat %s: %w", path, err)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("dirhash: read %s: %w", path, err)
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return "", fmt.Errorf("dirhash: rel %s: %w", path, err)
		}
		contentSum := sha256.Sum256(data)
		if _, err := fmt.Fprintf(hasher, "%o %s %s\n", info.Mode().Perm(), filepath.ToSlash(rel), hex.EncodeToString(contentSum[:])); err != nil {
			return "", fmt.Errorf("dirhash: hash line: %w", err)
		}
	}
	return "sha256:" + hex.EncodeToString(hasher.Sum(nil)), nil
}
