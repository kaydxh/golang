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
package strconv

import (
	"strconv"
)

func ParseUintOrFallback(s string, base int, bitSize int, defaultValue uint64) (uint64, error) {
	ns, err := strconv.ParseUint(s, base, bitSize)
	if err != nil {
		return defaultValue, err
	}

	return ns, nil
}

func ParseInt64Batch(m map[string]string) (map[string]int64, error) {
	nm := make(map[string]int64, 0)

	for k, v := range m {
		ns, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		nm[k] = ns
	}

	return nm, nil
}
