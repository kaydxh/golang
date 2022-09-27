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
package rand_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	rand_ "github.com/kaydxh/golang/go/math/rand"

	"github.com/stretchr/testify/assert"
)

func TestRand(t *testing.T) {
	s := fmt.Sprintf("%08v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000000))
	fmt.Printf("s: %v\n", s)

	ss := "abc_dvg_123"
	fmt.Printf("ss %v\n", ss)
	ts := strings.TrimPrefix(ss, "abc_dvg")
	fmt.Printf("ts: %v\n", ts)

	ns := "_abc"
	nns := strings.Split(ns, "_")
	fmt.Printf("nns: %v\n", nns)

}

func TestRangeInt(t *testing.T) {
	testCases := []struct {
		min int
		max int
	}{
		{
			min: 10,
			max: 12,
		},
		{
			min: 10000000,
			max: 100000000,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			r, err := rand_.RangeInt(testCase.min, testCase.max)
			if err != nil {
				t.Fatalf("failed to rand int, err: %v", err)
			}
			t.Logf("random: %v", r)

			assert.GreaterOrEqual(t, r, testCase.min)
			assert.LessOrEqual(t, r, testCase.max)

		})
	}
}

func TestRead(t *testing.T) {
	testCases := []struct {
		p []byte
	}{
		{
			p: make([]byte, 10),
		},
		{
			p: make([]byte, 20),
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			n, err := rand_.Read(testCase.p)
			if err != nil {
				t.Fatalf("failed to rand int, err: %v", err)
			}
			t.Logf("read n: %v, p: %v", n, testCase.p)

			assert.Equal(t, len(testCase.p), n)

		})
	}
}

func TestRangeString(t *testing.T) {
	testCases := []struct {
		n int
	}{
		{
			n: 0,
		},
		{
			n: 5,
		},
		{
			n: 8,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			str := rand_.RangeString(testCase.n)
			t.Logf("str: %v", str)
			assert.Equal(t, len(str), testCase.n)

		})
	}

}
