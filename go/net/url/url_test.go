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
package url_test

import (
	"fmt"
	"testing"

	url_ "github.com/kaydxh/golang/go/net/url"
	"golang.org/x/net/context"
)

func TestUrlEncode(t *testing.T) {
	type SimpleChild struct {
		Status bool
		Name   string
	}

	type SimpleData struct {
		Id           int
		Name         string
		Child        SimpleChild
		ParamsInt8   map[string]int8
		ParamsString map[string]string
		Array        [3]uint16
	}

	data := SimpleData{
		Id:   2,
		Name: "http://localhost/test.php?id=2",
		Child: SimpleChild{
			Status: true,
		},
		ParamsInt8: map[string]int8{
			"one": 1,
		},
		ParamsString: map[string]string{
			"two": "你好",
		},

		Array: [3]uint16{2, 3, 300},
	}

	c, _ := url_.New(context.Background())
	bytes, err := c.Encode(data)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
	fmt.Println(string(bytes))
}
