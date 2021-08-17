package reflect_test

import (
	"testing"

	reflect_ "github.com/kaydxh/golang/go/reflect"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveStructField(t *testing.T) {
	type HttpRequest struct {
		RequestId string
		Username  string
	}

	id := "123"
	req := &HttpRequest{
		RequestId: id,
	}

	requestId := reflect_.RetrieveStructField(req, "RequestId")
	t.Logf("requestId: %v", requestId)
	assert.Equal(t, id, requestId)
}
