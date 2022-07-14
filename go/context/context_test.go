package context_test

import (
	"context"
	"fmt"
	"testing"

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
