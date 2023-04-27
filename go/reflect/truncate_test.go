/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
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
			t.Logf("req: %+v\n, ", testCase.req)
			truncateReq := reflect_.TruncateBytes(testCase.req)
			//t.Logf("req: %+v\n, truncateReq: %+v", testCase.req, truncateReq)
			t.Logf("truncateReq: %+v", truncateReq)
		})
	}
}

func TestTruncateBytesWithMaxArraySize(t *testing.T) {

	testCases := []struct {
		req interface{}
	}{
		{
			req: &struct {
				RequestId string
				Image     []byte
				Item      []struct {
					a     int
					Image []byte
				}
			}{
				RequestId: uuid.New().String(),
				Image:     []byte("12345678"),
				Item: []struct {
					a     int
					Image []byte
				}{
					{
						a:     1,
						Image: []byte("12345678"),
					},
					{
						a:     2,
						Image: []byte("12345678"),
					},
					{
						a:     3,
						Image: []byte("12345678"),
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("case-%v", i), func(t *testing.T) {
			t.Logf("req: %+v\n", testCase.req)
			//	truncateReq := reflect_.TruncateBytesWithMaxArraySize(testCase.req, 1)
			//	t.Logf("truncateReq: %+v", truncateReq)
		})
	}

}
