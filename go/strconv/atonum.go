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
package strconv

import (
	"fmt"
	"strconv"
)

func ToFloat(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func ToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func ToInt64(str string) (int64, error) {
  return  strconv.ParseInt(str, 10, 64)
}

func ToUInt64(str string) (uint64, error) {
  return  strconv.ParseUint(str, 10, 64)
}

func ParseNumOrDefault[T any](str string, _default T, convert func(string) (T, error)) T {
	if str == "" {
		return _default
	}
	n, err := convert(str)
	if err != nil {
		//logurs.Infof("Invalid number value: %s", err)
		return _default
	}
	return n
}

func ParseNum[T any](str string, convert func(string) (T, error)) (T, error) {
    var t T
	if str == "" {
		return t,  fmt.Errorf("string is empty")
	}
	return convert(str)
}

func ParseNums[T any](strs []string, convert func(string) (T, error)) ([]T, error) {
    var ts []T
	if len(strs) == 0{
		return ts,  fmt.Errorf("string is empty")
	}

	for _, str := range strs {
	  t, err := convert(str)
	  if err != nil {
		return nil, err
	  }
	  ts = append(ts, t)
	}
	return  ts, nil
}
