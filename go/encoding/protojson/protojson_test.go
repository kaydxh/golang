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
		t.Fatalf("failed to marshal, err: %v", err)
	}
	t.Logf("marshal data: %v", string(data))

	var req testdata_.DateRequest
	err = protojson_.Unmarshal(data, &req)
	if err != nil {
		t.Fatalf("failed to unmarshal, err: %v", err)
	}
	t.Logf("unmarshal object: %v", req.String())
}
