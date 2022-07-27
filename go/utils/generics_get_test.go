package utils_test

import (
	"testing"

	utils_ "github.com/kaydxh/golang/go/utils"
)

func TestGetValueOrFallback(t *testing.T) {
	result := utils_.GetValueOrFallback(0, 20)
	t.Logf("result: %v", result)
}
