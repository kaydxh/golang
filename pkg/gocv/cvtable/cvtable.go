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
package cvtable

import (
	"fmt"
	"sort"
	"strconv"

	io_ "github.com/kaydxh/golang/go/io"
	exp_ "github.com/kaydxh/golang/go/math/exp"
)

const (
	validateRows = 101
	maxScore     = 100
	minScore     = 0
)

type CVTable struct {
	table []float64
}

func NewCVTable(filepath string) (*CVTable, error) {
	table, err := io_.ReadFileLines(filepath)
	if err != nil {
		return nil, err
	}

	c := &CVTable{}
	if table[0] != "0" {
		c.table = append(c.table, 0)
	}

	for _, v := range table {
		sim, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		c.table = append(c.table, sim)
	}
	err = c.Validate()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c CVTable) Validate() error {
	if len(c.table) != validateRows {
		return fmt.Errorf("invalid cv table rows %v, must be %v", len(c.table), validateRows)
	}

	return nil
}

// score 0-100, mapping to sim
func (c CVTable) Sim(score float64) float64 {
	if score <= minScore {
		return 0
	}

	// integerPart integer part of a decimal
	// decimalPart  decimal part of a decimal
	integerPart := int(score)
	decimalPart := score - float64(integerPart)

	// len(c.table) - 1 == 100
	if integerPart >= len(c.table)-1 {
		return 1
	}

	return c.table[integerPart] + (c.table[integerPart+1]-c.table[integerPart])*decimalPart
}

func (c CVTable) Score(sim float64) float64 {
	if sim <= 0 {
		return minScore
	}
	if sim >= 1 {
		return maxScore
	}

	pos := sort.Search(len(c.table), func(i int) bool { return c.table[i] >= sim })
	score := float64(pos)
	if pos > 0 && pos < len(c.table) {
		score += -1.0 + (sim-c.table[pos-1])/(c.table[pos]-c.table[pos-1])
	}

	return exp_.Value(score, minScore, maxScore)
}
