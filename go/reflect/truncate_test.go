package reflect_test

import (
	"fmt"
	"testing"

	//	"github.com/google/uuid"
	"github.com/google/uuid"
	reflect_ "github.com/kaydxh/golang/go/reflect"
)

func TestTruncateBytes(t *testing.T) {

	tmp := []byte("12345678")
	tmp2 := [][]byte{[]byte("12345678"), []byte("12345678")}
	_ = tmp
	testCases := []struct {
		req interface{}
	}{
		{
			req: &struct {
				RequestId string
				Image     []byte
				Item      struct {
					Image []byte
				}
			}{
				RequestId: uuid.New().String(),
				Image:     []byte("12345678"),
				Item: struct {
					Image []byte
				}{
					Image: []byte("12345678"),
				},
			},
		},
		{
			req: []byte("12345678"),
		},
		{
			req: &tmp,
		},
		{
			req: [][]byte{[]byte("12345678"), []byte("12345678")},
		},
		{
			req: &tmp2,
		},
		{
			req: &struct {
				Images [][]byte
			}{
				Images: [][]byte{
					[]byte("12345678"),
					[]byte("12345678"),
				},
			},
		},
		{
			req: &struct {
				Images [][][]byte
			}{
				Images: [][][]byte{
					[][]byte{[]byte("12345678")},
					[][]byte{[]byte("12345678")},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("case-%v", i), func(t *testing.T) {
			truncateReq := reflect_.TruncateBytes(testCase.req)
			t.Logf("req: %v, truncateReq: %s", testCase.req, truncateReq)
		})
	}
}
