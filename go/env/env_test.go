package env_test

import (
	"os"
	"testing"

	"github.com/kaydxh/golang/go/env"
	"github.com/stretchr/testify/assert"
)

func TestGetEnvAsStringOrFallback(t *testing.T) {
	const expected = "foo"

	assert := assert.New(t)
	key := "FLOCKER_SET_VAR"
	os.Setenv(key, expected)

	assert.Equal(expected, env.GetEnvAsStringOrFallback(key, "~"+expected))

	key = "FLOCKER_UNSET_VAR"
	assert.Equal(expected, env.GetEnvAsStringOrFallback(key, expected))
	assert.NotEqual(expected, env.GetEnvAsStringOrFallback(key, "~"+expected))
}
