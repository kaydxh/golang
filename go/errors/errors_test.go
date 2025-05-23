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
package errors_test

import (
	"errors"
	"fmt"
	"testing"

	errors_ "github.com/kaydxh/golang/go/errors"
)

func TestError(t *testing.T) {
	var errs = []error{fmt.Errorf("error 1"), fmt.Errorf("error 2")}

	var err error
	//Aggregate  implemnet interface error 	Error() string
	err = errors_.NewAggregate(errs)
	//multiErrorStrings := err errors_.NewAggregate(errs).Errors()
	multiErrorStrings := err.Error()
	t.Logf("multiErrorStrings: %v", multiErrorStrings)

}

func TestErrorIs(t *testing.T) {
	var ErrInternal = errors.New("internal error")
	testCases := []struct {
		err1     error
		err2     error
		expected bool
	}{
		{
			err1:     ErrInternal,
			err2:     ErrInternal,
			expected: true,
		},
		{
			// the same error messge is not mean the same error
			err1:     errors.New("internal error"),
			err2:     errors.New("internal error"),
			expected: false,
		},

		{
			err1:     errors.New("internal error1"),
			err2:     errors.New("internal error2"),
			expected: false,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			err := errors_.NewAggregate([]error{testCase.err1})
			if err.Is(testCase.err2) != testCase.expected {
				t.Fatalf("err[%v] is not expected err2[%v]", err, testCase.err2)
			}

		})
	}
}

func TestErrore(t *testing.T) {
	var ErrInternal = errors.New("FailedOperation__Internal")
	testCases := []struct {
		err      error
		code     int32
		expected bool
	}{
		{
			err:      ErrInternal,
			code:     100,
			expected: true,
		},
		{
			err:      errors_.Errore(fmt.Errorf("failed to process,%d", 2), ErrInternal),
			code:     100,
			expected: true,
		},
		{
			err:      fmt.Errorf("failed to process, %v", 2),
			code:     100,
			expected: false,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			//err := errors_.Errore(testCase.code, testCase.err)
			err := errors_.Errore(testCase.err)
			if errors.Is(err, ErrInternal) != testCase.expected {
				t.Fatalf("err[%v], epected[%v] test err[%v]", err, testCase.expected, testCase.err)
			}
			t.Logf("err[%v] expected[%v] test err[%v] ", err, testCase.expected, testCase.err)
		})
	}
}
