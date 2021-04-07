package strconv_test

import (
	"fmt"
	"testing"

	strconv_ "github.com/kaydxh/golang/go/strconv"
	"gotest.tools/v3/assert"
)

func TestParseInt64Batch(t *testing.T) {

	m := map[string]string{
		"fileId": "12345",
		"partId": "1",
	}
	nm, err := strconv_.ParseInt64Batch(m)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	assert.Equal(t, nm["fileId"], int64(12345))
	assert.Equal(t, nm["partId"], int64(1))

	fmt.Println(nm)
}
