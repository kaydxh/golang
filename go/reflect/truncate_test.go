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
	"strings"
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

func TestTruncateBytesWithThreshold(t *testing.T) {
	// 生成超过 1024 字节的数据
	longBytes := make([]byte, 2048)
	for i := range longBytes {
		longBytes[i] = byte('A' + (i % 26))
	}
	shortBytes := []byte("short data")

	testCases := []struct {
		name        string
		req         interface{}
		expectTrunc bool // 是否期望被截断
	}{
		{
			name: "bytes超过阈值应被截断",
			req: &struct {
				RequestId string
				Image     []byte
			}{
				RequestId: "test-id",
				Image:     longBytes,
			},
			expectTrunc: true,
		},
		{
			name: "bytes未超过阈值不应被截断",
			req: &struct {
				RequestId string
				Image     []byte
			}{
				RequestId: "test-id",
				Image:     shortBytes,
			},
			expectTrunc: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := reflect_.TruncateBytes(tc.req)
			t.Logf("truncateReq: %+v", result)
		})
	}
}

func TestTruncateBytesAndStrings(t *testing.T) {
	// 生成超过 1024 字节的数据
	longString := strings.Repeat("ABCDEFGHIJ", 200) // 2000 字节
	shortString := "short string"
	longBytes := make([]byte, 2048)
	for i := range longBytes {
		longBytes[i] = byte('0' + (i % 10))
	}
	shortBytes := []byte("short bytes")

	testCases := []struct {
		name string
		req  interface{}
	}{
		{
			name: "string和bytes都超过阈值",
			req: &struct {
				RequestId string
				Data      string
				Image     []byte
			}{
				RequestId: "test-id",
				Data:      longString,
				Image:     longBytes,
			},
		},
		{
			name: "string和bytes都未超过阈值",
			req: &struct {
				RequestId string
				Data      string
				Image     []byte
			}{
				RequestId: "test-id",
				Data:      shortString,
				Image:     shortBytes,
			},
		},
		{
			name: "嵌套结构体中的长string和bytes",
			req: &struct {
				RequestId string
				Item      struct {
					Name  string
					Image []byte
				}
			}{
				RequestId: "test-id",
				Item: struct {
					Name  string
					Image []byte
				}{
					Name:  longString,
					Image: longBytes,
				},
			},
		},
		{
			name: "自定义阈值测试",
			req: &struct {
				Data string
			}{
				Data: "hello world, this is a test string",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("before: %+v", tc.req)
			result := reflect_.TruncateBytesAndStrings(tc.req)
			t.Logf("after:  %+v", result)
		})
	}

	// 单独测试自定义阈值
	t.Run("自定义阈值-threshold=10-prefix=5", func(t *testing.T) {
		req := &struct {
			Data  string
			Image []byte
		}{
			Data:  "hello world, this is a long string",
			Image: []byte("hello world, this is long bytes"),
		}
		t.Logf("before: %+v", req)
		result := reflect_.TruncateBytesAndStringsWithThreshold(req, 10, 5)
		t.Logf("after:  %+v", result)
	})
}

// TestTruncateWithMap 测试包含 map 类型的结构体截断（模拟 google.protobuf.Struct 场景）
func TestTruncateWithMap(t *testing.T) {
	longBytes := make([]byte, 2048)
	for i := range longBytes {
		longBytes[i] = byte('X')
	}

	testCases := []struct {
		name string
		req  interface{}
	}{
		{
			name: "包含map的结构体",
			req: &struct {
				Fields map[string]*struct {
					Data  string
					Image []byte
				}
			}{
				Fields: map[string]*struct {
					Data  string
					Image []byte
				}{
					"key1": {
						Data:  strings.Repeat("A", 2000),
						Image: longBytes,
					},
					"key2": {
						Data:  "short",
						Image: []byte("short"),
					},
				},
			},
		},
		{
			name: "包含map[string]interface{}的结构体",
			req: &struct {
				Meta map[string]interface{}
			}{
				Meta: map[string]interface{}{
					"name": "test",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("before: %+v", tc.req)
			result := reflect_.TruncateBytesAndStrings(tc.req)
			t.Logf("after:  %+v", result)
		})
	}
}

// TestTruncateWithCircularLikeStruct 测试模拟循环引用结构（类似 google.protobuf.Struct）不会堆栈溢出
func TestTruncateWithCircularLikeStruct(t *testing.T) {
	// 模拟 google.protobuf.Struct 的循环引用结构
	type Value struct {
		StringValue string
		BytesValue  []byte
		StructValue *struct {
			Fields map[string]*Value
		}
	}

	type Struct struct {
		Fields map[string]*Value
	}

	req := &Struct{
		Fields: map[string]*Value{
			"image_data": {
				BytesValue: make([]byte, 2048),
			},
			"long_text": {
				StringValue: strings.Repeat("Z", 2000),
			},
			"nested": {
				StructValue: &struct {
					Fields map[string]*Value
				}{
					Fields: map[string]*Value{
						"inner_data": {
							BytesValue: make([]byte, 1500),
						},
					},
				},
			},
		},
	}

	// 这个测试的关键是：不会因为类似循环引用的结构导致堆栈溢出
	t.Logf("before: %+v", req)
	result := reflect_.TruncateBytesAndStrings(req)
	t.Logf("after:  %+v", result)
}