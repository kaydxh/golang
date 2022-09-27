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
package strconv_test

import (
	"fmt"
	"testing"

	strconv_ "github.com/kaydxh/golang/go/strconv"
	"gotest.tools/v3/assert"
)

func TestParseInt64Batch(t *testing.T) {

	m := map[string]string{
		"fileId": "12345",
		"partId": "1",
	}
	nm, err := strconv_.ParseInt64Batch(m)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}

	assert.Equal(t, nm["fileId"], int64(12345))
	assert.Equal(t, nm["partId"], int64(1))

	fmt.Println(nm)
}
