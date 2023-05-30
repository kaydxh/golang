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
package strings

import (
	"fmt"
	"strings"
	"unicode/utf8"

	strconv_ "github.com/kaydxh/golang/go/strconv"
)

func GetStringOrFallback(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}

	return ""
}

/*
func Replace(s string, old string, news []string, n int) string {
	if len(news) == 0 || n == 0 {
		return s
	}

	if m := strings.Count(s, old); m == 0 {
		return s // avoid allocation
	} else if n < 0 || m < n {
		n = m
	}
	// if len(news) < n , padding news use last element in news
	for i := 0; i < n-len(news); i++ {
		news = append(news, news[len(news)-1])
	}

	incLen := 0
	for i := 0; i < n; i++ {
		incLen += len(news[i]) - len(old)
	}

	// Apply replacements to buffer.
	var b strings.Builder
	b.Grow(len(s) + incLen)
	start := 0
	for i := 0; i < n; i++ {
		j := start
		if len(old) == 0 {
			if i > 0 {
				_, wid := utf8.DecodeRuneInString(s[start:])
				j += wid
			}
		} else {
			j += strings.Index(s[start:], old)
		}
		b.WriteString(s[start:j])
		b.WriteString(news[i])
		start = j + len(old)
	}
	b.WriteString(s[start:])
	return b.String()

}
*/

/*
func ReplaceAll(s, old string, news []string) string {
	return Replace(s, old, news, -1)
}
*/

func Replace(s string, old string, news []interface{}, useQuote bool, n int) string {
	if len(news) == 0 || n == 0 {
		return s
	}

	if m := strings.Count(s, old); m == 0 {
		return s // avoid allocation
	} else if n < 0 || m < n {
		n = m
	}
	// if len(news) < n , padding news use last element in news
	for i := 0; i < n-len(news); i++ {
		news = append(news, news[len(news)-1])
	}

	var b strings.Builder
	start := 0
	for i := 0; i < n; i++ {
		j := start
		if len(old) == 0 {
			if i > 0 {
				_, wid := utf8.DecodeRuneInString(s[start:])
				j += wid
			}
		} else {
			j += strings.Index(s[start:], old)
		}
		b.WriteString(s[start:j])
		if useQuote {
			b.WriteString(fmt.Sprintf(`"%v"`, news[i]))
		} else {
			b.WriteString(fmt.Sprintf("%v", news[i]))
		}
		start = j + len(old)
	}
	b.WriteString(s[start:])
	return b.String()

}

func ReplaceAll(s, old string, news []interface{}, useQuote bool) string {
	return Replace(s, old, news, useQuote, -1)
}

// Split slices s into all substrings separated by sep and returns a slice of
// the substrings between those separators.
//
// If s does not contain sep and sep is not empty, Split returns a
// slice of length 0.
//
// If sep is empty, Split splits after each UTF-8 sequence. If both s
// and sep are empty, Split returns an empty slice.
//
// It is equivalent to SplitN with a count of -1.
func SplitOmitEmpty(s, sep string) []string {
	var res []string
	a := strings.Split(s, sep)
	for _, v := range a {
		if v != "" {
			res = append(res, v)
		}
	}

	return res
}

// Split2 returns the values from strings.SplitN(s, sep, 2).
// If sep is not found, it returns ("", "", false) instead.
func Split2(s, sep string) (string, string, bool) {
	spl := strings.SplitN(s, sep, 2)
	if len(spl) < 2 {
		return "", "", false
	}
	return spl[0], spl[1], true
}

func SplitToNums[T any](s, sep string, convert func(string) (T, error)) ([]T, error) {
	ss := SplitOmitEmpty(s, sep)
	return strconv_.ParseNums(ss, convert)
}

func EqualCaseInsensitive(src, dst string) bool {
	return strings.ToLower(src) == strings.ToLower(dst)
}

func EmptyString(str string) bool {
	return strings.TrimSpace(str) == ""
}
