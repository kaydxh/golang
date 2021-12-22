package reflect_test

import (
	"testing"

	"github.com/google/uuid"
	reflect_ "github.com/kaydxh/golang/go/reflect"
)

func TestTruncateBytes(t *testing.T) {
	type HttpRequest struct {
		RequestId string
		Username  string
		Image     []byte
		Item      struct {
			Image []byte
		}
	}

	id := uuid.NewString()
	req := &HttpRequest{
		RequestId: id,
		Image:     []byte("12345678"),
		Item: struct {
			Image []byte
		}{
			Image: []byte("12345678"),
		},
	}

	//		req := []byte("12345678")
	truncateReq := reflect_.TruncateBytes(req)
	t.Logf("truncateReq: %s", truncateReq)
}
