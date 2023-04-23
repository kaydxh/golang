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
package datastore

import (
	"reflect"
	"sort"
)

type MessageKey struct {
	Key     string
	MsgType MsgType

	Fields []string
	// To db is table name
	Path string
}

func (m MessageKey) Equual(s MessageKey) bool {
	if m.Key != s.Key {
		return false
	}

	if m.MsgType != s.MsgType {
		return false
	}

	if m.Path != s.Path {
		return false
	}

	sort.Strings(m.Fields)
	sort.Strings(s.Fields)
	return reflect.DeepEqual(m.Fields, s.Fields)
}

type Message struct {
	Key   []byte
	Value []byte
}

type MsgType int32

const (
	MsgType_Insert MsgType = 0
	MsgType_Delete MsgType = 1
	MsgType_Update MsgType = 2
	MsgType_Get    MsgType = 3
)
