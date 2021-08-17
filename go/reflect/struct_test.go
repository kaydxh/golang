package reflect_test

import (
	"testing"

	"github.com/google/uuid"
	reflect_ "github.com/kaydxh/golang/go/reflect"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveStructField(t *testing.T) {
	type HttpRequest struct {
		RequestId string
		Username  string
	}

	id := uuid.NewString()
	req := &HttpRequest{
		RequestId: id,
	}

	requestId := reflect_.RetrieveStructField(req, "RequestId")
	t.Logf("requestId: %v", requestId)
	assert.Equal(t, id, requestId)
}

func TestTrySetStructField(t *testing.T) {
	type HttpRequest struct {
		RequestId string
		Username  string
	}

	id := uuid.NewString()
	req := &HttpRequest{
		//	RequestId: id,
	}

	reflect_.TrySetStructFiled(req, "RequestId", id)
	t.Logf("requestId: %v", req.RequestId)
	assert.Equal(t, id, req.RequestId)
}
