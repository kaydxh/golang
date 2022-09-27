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
package gocv_test

import (
	"testing"

	io_ "github.com/kaydxh/golang/go/io"
	gocv_ "github.com/kaydxh/golang/pkg/gocv/cgo/gocv"
)

func TestMagickInitializeMagick(t *testing.T) {
}

func TestMagickImageDecode(t *testing.T) {
	filename := "testdata/test.jpg"
	data, err := io_.ReadFile(filename)
	if err != nil {
		t.Error("Invalid ReadFile in TestMagickImageDecode")
		return
	}
	t.Logf("data size: %v", len(data))
	req := gocv_.NewMagickImageDecodeRequest()
	req.Image = data
	resp, err := gocv_.MagickImageDecode(req)
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Logf("resp: %v", resp)
}
