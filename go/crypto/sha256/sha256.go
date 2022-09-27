/*
 *Copyright (c) 2022, kaydxh
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
package sha256

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	os_ "github.com/kaydxh/golang/go/os"
)

func SumBytes(b []byte) string {
	h := sha256.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

func SumString(s string) string {
	return SumBytes([]byte(s))
}

func SumReader(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func SumReaderN(r io.Reader, n int64) (string, error) {
	h := sha256.New()
	if _, err := io.CopyN(h, r, n); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func SumReaderAt(r io.ReaderAt, offset, length int64) (string, error) {
	h := sha256.New()
	buf := make([]byte, 1024)
	var total int64

	for total < length {
		n, err := r.ReadAt(buf, offset)
		if err == nil || err == io.EOF {
			if n > 0 {
				_, tmpErr := io.CopyN(h, bytes.NewReader(buf), int64(n))
				if tmpErr != nil {
					return "", tmpErr
				}
				offset += int64(n)
				total += int64(n)

			}

			if err == io.EOF {
				break
			}
		} else {
			return "", err
		}

	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func SumFile(fileName string) (string, error) {
	file, err := os_.OpenFile(fileName, true)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return SumReader(file)
}

func SumFileAt(fileName string, offset int64, length int64) (string, error) {
	file, err := os_.OpenAll(fileName, os.O_RDWR, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return SumReaderAt(file, offset, length)
}
