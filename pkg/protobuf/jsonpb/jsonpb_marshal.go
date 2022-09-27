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
package jsonpb

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	structpb "github.com/golang/protobuf/ptypes/struct"
)

func MarshaToStructpb(v interface{}) (*structpb.Struct, error) {
	var jb []byte
	switch v := v.(type) {
	case nil:
		return &structpb.Struct{}, nil

	case *structpb.Struct:
		return v, nil

	case proto.Message:
		m := jsonpb.Marshaler{EmitDefaults: true}
		data, err := m.MarshalToString(v)
		if err != nil {
			return nil, fmt.Errorf("failed to Marshal json: %v", err)
		}
		jb = []byte(data)

	case []byte:
		jb = v
	case *[]byte:
		jb = *v
	case string:
		jb = []byte(v)
	case *string:
		jb = []byte(*v)
	default:
		var err error
		jb, err = json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to Marshal json: %v", err)
		}
	}

	var dataStructpb structpb.Struct
	if err := UnmarshalWithAllowUnknownFields(bytes.NewReader(jb), &dataStructpb); err != nil {
		return nil, err
	}
	return &dataStructpb, nil
}
