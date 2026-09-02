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
// Package jsonlines provides newline-delimited JSON framing helpers.
package jsonlines

import (
	"bufio"
	"encoding/json"
	"io"
)

// DefaultMaxTokenSize is the maximum line size used by NewScanner when max is
// not positive.
const DefaultMaxTokenSize = 4 * 1024 * 1024

// NewScanner returns a line scanner sized for JSON event streams.
func NewScanner(r io.Reader, maxTokenSize int) *bufio.Scanner {
	if maxTokenSize <= 0 {
		maxTokenSize = DefaultMaxTokenSize
	}
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), maxTokenSize)
	return scanner
}

// Write encodes one JSON value followed by a newline.
func Write(w io.Writer, value any) error {
	return json.NewEncoder(w).Encode(value)
}
