/*
 *Copyright (c) 2026, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
package exec_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	exec_ "github.com/kaydxh/golang/go/os/exec"
)

func writeExecutable(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLookPathInUsesOnlyTheGivenPATH(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX executable bits")
	}
	listed, ambient := t.TempDir(), t.TempDir()
	want := writeExecutable(t, listed, "tool")
	writeExecutable(t, ambient, "ambient-only")
	t.Setenv("PATH", ambient)

	environ := []string{"HOME=/nowhere", "PATH=" + listed}
	got, err := exec_.LookPathIn("tool", environ)
	if err != nil || got != want {
		t.Fatalf("LookPathIn(tool) = %q, %v; want %q", got, err, want)
	}
	if _, err := exec_.LookPathIn("ambient-only", environ); !errors.Is(err, exec.ErrNotFound) {
		t.Fatalf("a name only on the process PATH must not resolve: %v", err)
	}
}

func TestLookPathInWithoutPATHFindsNothing(t *testing.T) {
	if _, err := exec_.LookPathIn("sh", []string{"HOME=/nowhere"}); !errors.Is(err, exec.ErrNotFound) {
		t.Fatalf("no PATH entry must resolve nothing: %v", err)
	}
}

func TestLookPathInSkipsNonExecutableAndRelativeEntries(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX executable bits")
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "data"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := exec_.LookPathIn("data", []string{"PATH=" + dir}); !errors.Is(err, exec.ErrNotFound) {
		t.Fatalf("a non-executable file must not resolve: %v", err)
	}
	if _, err := exec_.LookPathIn("data", []string{"PATH=relative/bin"}); !errors.Is(err, exec.ErrNotFound) {
		t.Fatalf("a relative PATH entry must not resolve: %v", err)
	}
}

func TestLookPathInRejectsPaths(t *testing.T) {
	if _, err := exec_.LookPathIn("./tool", []string{"PATH=/bin"}); err == nil {
		t.Fatal("a name with a separator is not a PATH lookup")
	}
}
