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
package context_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	context_ "github.com/kaydxh/golang/go/context"
)

func withField(ctx context.Context) {
	ctx = context.WithValue(ctx, "abc", "abc")
	fmt.Printf("context: %+v\n", ctx)
}

func TestContext(t *testing.T) {
	ctx := context.Background()
	t.Logf("context: %v", ctx)
	withField(ctx)
	t.Logf("context: %v", ctx)
}

func doA(ctx context.Context) {

	{
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second)
		defer cancel()
	}

	timer := time.NewTimer(3 * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		fmt.Println("doA timeout")

	case <-timer.C:
		fmt.Println("doA finish")
	}
}

func doB(ctx context.Context) {

	{
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	timer := time.NewTimer(3 * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		fmt.Println("doB timeout")

	case <-timer.C:
		fmt.Println("doB finish")
	}
}

func TestContextTimeout(t *testing.T) {
	ctx := context.Background()
	doA(ctx)
	doB(ctx)
}

func TestWithIdleTimeout(t *testing.T) {
	t.Parallel()
	ctx, cancel, touch := context_.WithIdleTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	time.Sleep(60 * time.Millisecond)
	touch()
	select {
	case <-ctx.Done():
		t.Fatal("idle context canceled before the reset timeout")
	case <-time.After(60 * time.Millisecond):
	}
	select {
	case <-ctx.Done():
	case <-time.After(100 * time.Millisecond):
		t.Fatal("idle context was not canceled")
	}
}

func TestExtractIntegerFromContext(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		key      string
		value    string
		expected string
	}{
		{
			key:   "test-1",
			value: "123",
		},
		{
			key:   "test-2",
			value: "test-123",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.key, func(t *testing.T) {
			ctx = context_.SetPairContext(ctx, testCase.key, testCase.value)

			number, err := context_.ExtractIntegerFromContext(ctx, testCase.key)
			if err != nil {
				t.Errorf("expect nil, got %v", err)
				return
			}
			t.Logf("extract value %v by key %v ", number, testCase.key)

		})
	}

}

func TestExtractStringFromContext(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		key      string
		value    string
		expected string
	}{
		{
			key:   "test-1",
			value: "123",
		},
		{
			key:   "test-2",
			value: "test-123",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.key, func(t *testing.T) {
			ctx = context_.SetPairContext(ctx, testCase.key, testCase.value)

			value := context_.ExtractStringFromContext(ctx, testCase.key)
			t.Logf("extract value %v by key %v ", value, testCase.key)

		})
	}

}
