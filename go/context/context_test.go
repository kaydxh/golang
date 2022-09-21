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
