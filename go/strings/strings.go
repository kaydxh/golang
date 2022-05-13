package strings

import (
	"fmt"
	"strings"
	"unicode/utf8"
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
