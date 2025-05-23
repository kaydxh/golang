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
package url

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

type Client struct {
	mutex  sync.Mutex
	buffer *bytes.Buffer

	opts struct {
		needEmptyValue bool
		urlCodec       UrlCodec
	}
}

func New(ctx context.Context, options ...ClientOption) (*Client, error) {
	c := &Client{}
	c.opts.urlCodec = DefaultUrlCodec{}
	c.ApplyOptions(options...)

	return c, nil
}

//can convert struct to url encode for url paratment
func (c *Client) Encode(data interface{}) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.buffer = new(bytes.Buffer)
	rv := reflect.ValueOf(data)

	err := c.build(rv, "", reflect.Interface)
	if err != nil {
		return "", err
	}

	buf := c.buffer.Bytes()
	c.buffer = nil

	return string(buf[0 : len(buf)-1]), nil

}

func (c *Client) encode(rv reflect.Value) (string, error) {
	encoder := getEncoder(rv.Kind())
	if encoder == nil {
		return "", fmt.Errorf("unsupport type: %v", rv.Type().String())
	}

	return encoder.Encode(rv), nil
}

func (c *Client) build(
	rv reflect.Value,
	parentKey string,
	parentKind reflect.Kind,
) error {
	fmt.Println("-- parentKey: ", parentKey)
	switch rv.Kind() {
	case reflect.Map:
		for _, key := range rv.MapKeys() {
			checkKey := key
			if key.Kind() == reflect.Interface || key.Kind() == reflect.Ptr {
				checkKey = checkKey.Elem()
			}

			keyStr, err := c.encode(checkKey)
			if err != nil {
				return err
			}

			c.build(rv.MapIndex(key), keyStr, rv.Kind())

		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			c.build(rv.Index(i), parentKey+"["+strconv.Itoa(i)+"]", rv.Kind())
		}

	case reflect.Struct:
		rt := rv.Type()
		for i := 0; i < rt.NumField(); i++ {
			ft := rt.Field(i)
			//unexported
			if ft.PkgPath != "" && !ft.Anonymous {
				continue
			}

			//specially handle anonymous fields
			if ft.Anonymous && rv.Field(i).Kind() == reflect.Struct {
				c.build(rv.Field(i), parentKey, rv.Kind())
				continue
			}

			/*
				tag := ft.Tag.Get("query")
				//all ignore
				if tag == "-" {
					continue
				}

				t := newTag(tag)
				//get the related name
				name := t.getName()
				if name == "" {
					name = ft.Name
				}
			*/

			name := ft.Name
			c.build(rv.Field(i), name, rv.Kind())
		}

	case reflect.Ptr, reflect.Interface:
		if !rv.IsNil() {
			c.build(rv.Elem(), parentKey, parentKind)
		}

	default:
		c.appendKeyValue(parentKey, rv, parentKind)

	}

	return nil

}

//basic structure can be translated directly
func (c *Client) appendKeyValue(key string, rv reflect.Value, parentKind reflect.Kind) error {
	//If parent type is struct and empty value will be ignored by default. unless needEmptyValue is true.
	if parentKind == reflect.Struct && !c.opts.needEmptyValue && isEmptyValue(rv) {
		return nil
	}

	//If parent type is slice or array, then repack key. eg. students[0] -> students[]
	if parentKind == reflect.Slice || parentKind == reflect.Array {
		key = repackArrayQueryKey(key)
	}

	s, err := c.encode(rv)
	if err != nil {
		return err
	}

	_, err = c.buffer.WriteString(
		c.opts.urlCodec.Escape(key) + "=" + c.opts.urlCodec.Escape(s) + "&",
	)

	return err
}

//Is Zero-value
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

//if key like `students[0]` , repack it to `students[]`
func repackArrayQueryKey(key string) string {
	l := len(key)
	if l > 0 && key[l-1] == ']' {
		for l--; l >= 0; l-- {
			if key[l] == '[' {
				return key[:l+1] + "]"
			}
		}
	}
	return key
}
