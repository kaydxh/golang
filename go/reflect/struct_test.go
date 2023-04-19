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
	"testing"

	"github.com/google/uuid"
	reflect_ "github.com/kaydxh/golang/go/reflect"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveStructField(t *testing.T) {
	type HttpRequest struct {
		RequestId string
		Username  string
	}

	id := uuid.NewString()
	req := &HttpRequest{
		RequestId: id,
	}

	requestId := reflect_.RetrieveStructField(req, "RequestId")
	t.Logf("requestId: %v", requestId)
	assert.Equal(t, id, requestId)
}

func TestTrySetStructField(t *testing.T) {
	type HttpRequest struct {
		RequestId string
		Username  string
	}

	id := uuid.NewString()
	req := &HttpRequest{
		//	RequestId: id,
	}

	reflect_.TrySetStructFiled(req, "RequestId", id)
	t.Logf("requestId: %v", req.RequestId)
	assert.Equal(t, id, req.RequestId)
}

func TestNonzeroFieldTags(t *testing.T) {
	type HttpRequest struct {
		RequestId string `db:"request_id"`
		Username  string `db:"username"`
	}

	id := uuid.NewString()
	req := &HttpRequest{
		RequestId: id,
		//	Username:  "username 1",
	}
	fields := reflect_.NonzeroFieldTags(req, "db")
	t.Logf("fields: %v", fields)
	assert.Equal(t, []string{"request_id"}, fields)
}

func TestAllFieldTags(t *testing.T) {
	type HttpRequest struct {
		RequestId string `db:"request_id"`
		Username  string `db:"username"`
	}

	id := uuid.NewString()
	req := HttpRequest{
		RequestId: id,
		//	Username:  "username 1",
	}
	fields := reflect_.AllFieldTags(req, "db")
	t.Logf("fields: %v", fields)
	//assert.Equal(t, []string{"request_id"}, fields)
}

func TestAllTagsValues(t *testing.T) {
	type HttpRequest struct {
		RequestId string `db:"request_id"`
		Username  string `db:"username"`
	}

	id := uuid.NewString()
	req := &HttpRequest{
		RequestId: id,
		Username:  "admin",
	}
	tagsValues := reflect_.AllTagsValues(req, "db")
	t.Logf("tagsValues: %v", tagsValues)
	//assert.Equal(t, []string{"request_id"}, fields)
}
