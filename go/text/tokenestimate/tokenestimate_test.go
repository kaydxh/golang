package tokenestimate

import (
	"strings"
	"testing"
)

func TestEstimateTokensCJK(t *testing.T) {
	// One CJK rune ≈ one token: 保存配置 must not be treated as one token
	// for the whole run.
	if got := EstimateTokens("保存配置"); got < 4 {
		t.Fatalf("cjk estimate = %d, want >= 4", got)
	}
}

func TestEstimateTokensLatin(t *testing.T) {
	// ~4 bytes per token for latin text, rounded up.
	if got := EstimateTokens("abcdefgh"); got != 2 {
		t.Fatalf("latin estimate = %d, want 2", got)
	}
	if got := EstimateTokens("abcde"); got != 2 {
		t.Fatalf("latin rounding = %d, want 2", got)
	}
}

func TestEstimateTokensEmpty(t *testing.T) {
	if got := EstimateTokens(""); got != 0 {
		t.Fatalf("empty = %d", got)
	}
}

func TestEstimateTokensMixedConservative(t *testing.T) {
	mixed := "保存 save configuration 配置"
	// CJK runes: 保存配置 = 4; latin " save configuration " = 20 bytes → 5.
	if got := EstimateTokens(mixed); got != 9 {
		t.Fatalf("mixed = %d, want 9", got)
	}
}

func TestEstimateTokensLongDocument(t *testing.T) {
	long := strings.Repeat("word ", 1000)
	if got := EstimateTokens(long); got < 1000 {
		t.Fatalf("long = %d, want >= 1000", got)
	}
}
