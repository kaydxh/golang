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
package marshaler

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type JSONPb struct {
	runtime.JSONPb
	opts struct {
		useProtoNames   bool
		useEnumNumbers  bool
		emitUnpopulated bool
		discardUnknown  bool
		indent          string
	}
}

func NewDefaultJSONPb() *JSONPb {
	return NewJSONPb(
		// use json name
		WithUseProtoNames(false),
		WithUseEnumNumbers(false),
		WithEmitUnpopulated(true),
		WithDiscardUnknown(true),
		WithIndent("\t"),
	)
}

func NewJSONPb(options ...JSONPbOption) *JSONPb {
	j := &JSONPb{}
	j.ApplyOptions(options...)
	return j
}

/*
// if implemet the function, can parse some field from req data
func (j *JSONPb) NewDecoder(r io.Reader) runtime.Decoder {
	return runtime.DecoderFunc(func(v interface{}) error {
		rawData, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		err = j.JSONPb.Unmarshal(rawData, v)
		if err != nil {
			return err
		}
		id := reflect_.RetrieveId(v, reflect_.FieldNameRequestId)
		_ = id
		return nil
	})
}
*/