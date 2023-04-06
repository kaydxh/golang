/*
 *Copyright (c) 2023, kaydxh
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
	"testing"

	reflect_ "github.com/kaydxh/golang/go/reflect"
)

func TestArrayAllTagsValues(t *testing.T) {
	type HttpRequest struct {
		RequestId string `db:"request_id"`
		Username  string `db:"username"`
	}

	req := []HttpRequest{
		HttpRequest{
			RequestId: "123",
			Username:  "123-username",
		},
		HttpRequest{
			RequestId: "456",
			Username:  "456-username",
		},
	}

	tagsValues := reflect_.ArrayAllTagsVaules(req, "db")
	t.Logf("tagsValues: %v", tagsValues)
	//assert.Equal(t, []string{"request_id"}, fields)
}
