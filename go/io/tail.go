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
package io

import (
	"io"
	"os"
	"strings"
	"sync"
)

// TailWriter is an io.Writer that keeps only the last limit bytes written to
// it — the end of a child process's stderr, say, where the reason it failed
// usually is, without holding a chatty process's whole output in memory.
// It is safe for concurrent use.
type TailWriter struct {
	mu        sync.Mutex
	limit     int
	buf       []byte
	truncated bool
}

// NewTailWriter returns a TailWriter keeping at most limit bytes.
func NewTailWriter(limit int) *TailWriter {
	if limit < 0 {
		limit = 0
	}
	return &TailWriter{limit: limit}
}

// Write keeps the tail of p; it never fails.
func (w *TailWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	n := len(p)
	if n >= w.limit {
		w.truncated = w.truncated || n > w.limit || len(w.buf) > 0
		w.buf = append(w.buf[:0], p[n-w.limit:]...)
		return n, nil
	}
	if over := len(w.buf) + n - w.limit; over > 0 {
		w.buf = append(w.buf[:0], w.buf[over:]...)
		w.truncated = true
	}
	w.buf = append(w.buf, p...)
	return n, nil
}

// Bytes returns a copy of the kept tail.
func (w *TailWriter) Bytes() []byte {
	w.mu.Lock()
	defer w.mu.Unlock()
	return append([]byte(nil), w.buf...)
}

// String returns the kept tail.
func (w *TailWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return string(w.buf)
}

// Truncated reports whether earlier output was dropped.
func (w *TailWriter) Truncated() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.truncated
}

// ReadFileTail reads at most the last limit bytes of the file at path, and
// reports whether anything before them was left out.
func ReadFileTail(path string, limit int64) ([]byte, bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, false, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, false, err
	}
	truncated := info.Size() > limit
	if truncated {
		if _, err := f.Seek(info.Size()-limit, io.SeekStart); err != nil {
			return nil, false, err
		}
	}
	data, err := io.ReadAll(io.LimitReader(f, limit))
	return data, truncated, err
}

// TailText renders a byte-bounded tail — from TailWriter or ReadFileTail — as
// text: a rune the byte cut split is dropped, surrounding blank space is
// trimmed, and only the last maxLines lines are kept (all when maxLines <= 0).
func TailText(data []byte, maxLines int) string {
	text := strings.TrimSpace(strings.ToValidUTF8(string(data), ""))
	if maxLines <= 0 || text == "" {
		return text
	}
	lines := strings.Split(text, "\n")
	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}
	return strings.Join(lines, "\n")
}
