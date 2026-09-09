// Package tokenestimate approximates token counts without a tokenizer.
// It exists for budget accounting where a real tokenizer is unavailable:
// CJK-heavy text counts roughly one token per CJK rune, everything else
// roughly one token per four bytes. Results are estimates by definition —
// callers must mark them as estimated, never as exact.
package tokenestimate

import "unicode/utf8"

// EstimateTokens returns a conservative token estimate for s.
// CJK runes count one token each; the remaining bytes count one token per
// four bytes, rounded up, so the estimate errs high rather than under
// budget.
func EstimateTokens(s string) int {
	cjk := 0
	for _, r := range s {
		if isCJK(r) {
			cjk++
		}
	}
	// Non-CJK bytes count precisely, four bytes per token.
	non := 0
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if !isCJK(r) {
			non += size
		}
		i += size
	}
	tokens := cjk + (non+3)/4
	return tokens
}

func isCJK(r rune) bool {
	switch {
	case r >= 0x4E00 && r <= 0x9FFF: // CJK Unified Ideographs
		return true
	case r >= 0x3400 && r <= 0x4DBF: // Extension A
		return true
	case r >= 0x3000 && r <= 0x303F: // CJK punctuation
		return true
	case r >= 0xFF00 && r <= 0xFFEF: // fullwidth forms
		return true
	}
	return false
}
