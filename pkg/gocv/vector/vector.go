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
package vector

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

type Vector[T float32 | float64 | int8] struct {
	data []T
}

func NewVector[T float32 | float64 | int8]() *Vector[T] {
	return &Vector[T]{
		data: make([]T, 0),
	}
}

func (v *Vector[T]) Append(t T) {
	v.data = append(v.data, t)
}

func (v *Vector[T]) Norm() {
	var sum T
	for _, data := range v.data {
		sum += data * data
	}
	magnitude := math.Sqrt(float64(sum))

	for i, value := range v.data {
		v.data[i] = T(float64(value) / magnitude)
	}
}
func (v *Vector[T]) Assign(data []byte) error {
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, v.data)
	return err
}

func (v *Vector[T]) Dims() int {
	return len(v.data)
}

func (v *Vector[T]) Len() int {
	return len(v.data)
}

func (v Vector[T]) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, v.data)
	return buf.Bytes(), err
}

func (v Vector[T]) String() (string, error) {
	var err error
	buf := bytes.NewBuffer([]byte{'['})
	for i, value := range v.data {
		if i < len(v.data)-1 {
			_, err = buf.WriteString(fmt.Sprintf("%v ", value))
		} else {
			_, err = buf.WriteString(fmt.Sprintf("%v", value))
		}
		if err != nil {
			return "", err
		}
	}
	_, err = buf.Write([]byte{']'})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
