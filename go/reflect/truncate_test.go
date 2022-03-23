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
					a     int
					Image []byte
				}
			}{
				RequestId: uuid.New().String(),
				Image:     []byte("12345678"),
				Item: struct {
					a     int
					Image []byte
				}{
					a:     1,
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
		{
			req: &struct {
				RequstId   string
				FrameImage []byte
				Jobs       []struct {
					JobType   int
					JobOutput struct {
						OccupyData struct {
							GroupCode  string
							FrameImage []byte
						}
					}
				}
			}{
				RequstId:   "RRRRRRID",
				FrameImage: []byte("frame data"),
				Jobs: []struct {
					JobType   int
					JobOutput struct {
						OccupyData struct {
							GroupCode  string
							FrameImage []byte
						}
					}
				}{
					{
						JobType: 1,
						JobOutput: struct {
							OccupyData struct {
								GroupCode  string
								FrameImage []byte
							}
						}{
							OccupyData: struct {
								GroupCode  string
								FrameImage []byte
							}{
								GroupCode:  "group code",
								FrameImage: []byte("frame data"),
							},
						},
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("case-%v", i), func(t *testing.T) {
			truncateReq := reflect_.TruncateBytes(testCase.req)
			t.Logf("req: %v\n, truncateReq: %s", testCase.req, truncateReq)
		})
	}
}
