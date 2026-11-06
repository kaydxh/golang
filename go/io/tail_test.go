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
package io_test

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	io_ "github.com/kaydxh/golang/go/io"
)

func TestTailWriterKeepsEverythingUnderTheLimit(t *testing.T) {
	w := io_.NewTailWriter(16)
	if n, err := w.Write([]byte("hello ")); n != 6 || err != nil {
		t.Fatalf("Write = %d, %v", n, err)
	}
	_, _ = w.Write([]byte("world"))
	if got := w.String(); got != "hello world" || w.Truncated() {
		t.Fatalf("String = %q, truncated = %v", got, w.Truncated())
	}
}

func TestTailWriterKeepsOnlyTheLastBytes(t *testing.T) {
	w := io_.NewTailWriter(8)
	_, _ = w.Write([]byte("0123456789"))
	_, _ = w.Write([]byte("ab"))
	if got := w.String(); got != "456789ab" || !w.Truncated() {
		t.Fatalf("String = %q, truncated = %v", got, w.Truncated())
	}
	_, _ = w.Write([]byte(strings.Repeat("z", 20)))
	if got := w.String(); got != strings.Repeat("z", 8) {
		t.Fatalf("one large write keeps its own tail: %q", got)
	}
}

func TestTailWriterBytesIsACopy(t *testing.T) {
	w := io_.NewTailWriter(4)
	_, _ = w.Write([]byte("abcd"))
	got := w.Bytes()
	got[0] = 'x'
	if w.String() != "abcd" {
		t.Fatalf("Bytes must not alias the buffer: %q", w.String())
	}
}

func TestTailWriterIsSafeForConcurrentUse(t *testing.T) {
	w := io_.NewTailWriter(64)
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				_, _ = w.Write([]byte("line\n"))
				_ = w.String()
			}
		}()
	}
	wg.Wait()
	if got := len(w.String()); got != 64 {
		t.Fatalf("len = %d, want 64", got)
	}
}

func TestReadFileTail(t *testing.T) {
	path := filepath.Join(t.TempDir(), "service.log")
	if err := os.WriteFile(path, []byte("0123456789"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, truncated, err := io_.ReadFileTail(path, 4)
	if err != nil || string(got) != "6789" || !truncated {
		t.Fatalf("ReadFileTail(4) = %q, %v, %v", got, truncated, err)
	}
	got, truncated, err = io_.ReadFileTail(path, 64)
	if err != nil || string(got) != "0123456789" || truncated {
		t.Fatalf("ReadFileTail(64) = %q, %v, %v", got, truncated, err)
	}
	if _, _, err := io_.ReadFileTail(filepath.Join(t.TempDir(), "missing"), 4); !os.IsNotExist(err) {
		t.Fatalf("missing file: %v", err)
	}
}

func TestTailText(t *testing.T) {
	cut := []byte("中文")[1:] // a byte-bounded tail starting mid-rune
	for _, tc := range []struct {
		name     string
		data     []byte
		maxLines int
		want     string
	}{
		{"empty", nil, 3, ""},
		{"blank", []byte("\n  \n"), 3, ""},
		{"trimmed", []byte("\n  one\ntwo  \n"), 3, "one\ntwo"},
		{"last lines", []byte("a\nb\nc\nd\n"), 2, "c\nd"},
		{"all lines", []byte("a\nb\nc"), 0, "a\nb\nc"},
		{"cut rune", append(cut, "\nok"...), 3, "文\nok"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := io_.TailText(tc.data, tc.maxLines); got != tc.want {
				t.Fatalf("TailText = %q, want %q", got, tc.want)
			}
		})
	}
}
