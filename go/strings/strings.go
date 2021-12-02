package strings

import (
	"strings"
	"unicode/utf8"
)

func GetStringOrFallback(values []string, defaultValue string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}

	return defaultValue
}

func Replace(s string, old string, news []string, n int) string {
	if len(news) == 0 || n == 0 {
		return s
	}

	// if len(news) < n , padding news use last element in news
	for i := 0; i < n-len(news); i++ {
		news = append(news, news[len(news)-1])
	}

	if m := strings.Count(s, old); m == 0 {
		return s // avoid allocation
	} else if n < 0 || m < n {
		n = m
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
