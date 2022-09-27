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
package viper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gogo/protobuf/proto"
	jsonpb_ "github.com/kaydxh/golang/pkg/protobuf/jsonpb"
	"github.com/mitchellh/mapstructure"

	"github.com/spf13/viper"
)

func UnmarshalProtoMessageWithJsonPb(v *viper.Viper, msg proto.Message, options ...viper.DecoderConfigOption) error {
	if v == nil {
		return fmt.Errorf("viper is nil")
	}

	var opts []viper.DecoderConfigOption
	opts = append(opts, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.TagName = "json" // trick of protobuf, which generates json tag only
		decoderConfig.WeaklyTypedInput = true
		decoderConfig.DecodeHook = UnmarshalProtoMessageWithJsonpbHookFunc(msg)
	})
	opts = append(opts, options...)
	return v.Unmarshal(msg, opts...)
}

func UnmarshalProtoMessageWithJsonpbHookFunc(v proto.Message) mapstructure.DecodeHookFunc {
	return func(src reflect.Type, dst reflect.Type, data interface{}) (interface{}, error) {
		// Convert it by parsing
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		err = jsonpb_.UnmarshalWithAllowUnknownFields(bytes.NewReader(dataBytes), v)
		if err != nil {
			return data, err
		}

		return v, nil
	}
}
