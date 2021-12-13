package protojson_test

import (
	"testing"

	"github.com/google/uuid"
	protojson_ "github.com/kaydxh/golang/go/encoding/protojson"
	testdata_ "github.com/kaydxh/golang/go/encoding/protojson/testdata"
)

func TestMarshal(t *testing.T) {
	request := &testdata_.DateRequest{
		RequestId: uuid.New().String(),
	}

	data, err := protojson_.Marshal(request)
	if err != nil {
		t.Fatalf("failed t marshal, err: %v", err)
	}
	t.Logf("marshal data: %v", string(data))

}
