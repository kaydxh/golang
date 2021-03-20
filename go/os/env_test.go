package os_test

import (
	"os"
	"testing"

	os_ "github.com/kaydxh/golang/go/os"
	"github.com/stretchr/testify/assert"
)

func TestGetEnvAsStringOrFallback(t *testing.T) {
	const expected = "foo"

	assert := assert.New(t)
	key := "FLOCKER_SET_VAR"
	os.Setenv(key, expected)

	assert.Equal(expected, os_.GetEnvAsStringOrFallback(key, "~"+expected))

	key = "FLOCKER_UNSET_VAR"
	assert.Equal(expected, os_.GetEnvAsStringOrFallback(key, expected))
	assert.NotEqual(expected, os_.GetEnvAsStringOrFallback(key, "~"+expected))
}
