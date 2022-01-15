package context_test

import (
	"context"
	"fmt"
	"testing"
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
