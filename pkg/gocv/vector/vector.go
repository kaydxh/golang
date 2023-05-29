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

	rand_ "github.com/kaydxh/golang/go/math/rand"
	"golang.org/x/exp/constraints"
)

type Vector[T constraints.Integer | constraints.Float] struct {
	data []T
}

func NewVector[T constraints.Integer | constraints.Float](data ...T) *Vector[T] {
	v := &Vector[T]{
		data: make([]T, 0),
	}
	if len(data) > 0 {
		v.data = data
	}
	return v
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

func (v *Vector[T]) Data() []T {
	return v.data
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
			_, err = buf.WriteString(fmt.Sprintf("%v,", value))
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

func (v Vector[T]) CosineDistance(r *Vector[T]) float64 {
	if v.Len() != r.Len() {
		return 0
	}

	var (
		sum   float64
		normL float64
		normR float64
	)
	for i := 0; i < v.Len(); i++ {
		sum += float64(v.Data()[i] * r.Data()[i])
		normL += float64(v.Data()[i] * v.Data()[i])
		normR += float64(r.Data()[i] * r.Data()[i])
	}

	return sum / (math.Sqrt(normL) * math.Sqrt(normR))
}

func (v Vector[T]) EuclideanDistance(r *Vector[T]) float64 {
	if v.Len() != r.Len() {
		return 0
	}

	var sum float64
	for i := 0; i < v.Len(); i++ {
		sum += math.Pow(float64(v.Data()[i]-r.Data()[i]), 2)
	}

	return math.Sqrt(sum)
}

func CreateNormalizedVector[T constraints.Integer | constraints.Float](dim int) []T {
	var ve Vector[T]
	for i := 0; i < dim; i++ {
		ve.Append(T(rand_.Float64()))
	}
	ve.Norm()

	return ve.Data()
}

func CosineDistance[T constraints.Integer | constraints.Float](l, r []T) float64 {
	v1 := NewVector(l...)
	v2 := NewVector(r...)
	return v1.CosineDistance(v2)
}

func EuclideanDistance[T constraints.Integer | constraints.Float](l, r []T) float64 {
	v1 := NewVector(l...)
	v2 := NewVector(r...)
	return v1.EuclideanDistance(v2)
}
