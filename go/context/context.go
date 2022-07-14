package context

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

func WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {

	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	}
	return ctx, func() {}
}

func ExtractStringFromContext(ctx context.Context, key string) string {
	if v, ok := ctx.Value(key).(string); ok {
		return v
	}

	return ""
}

func ExtractIntegerFromContext(ctx context.Context, key string) (int64, error) {
	v, ok := ctx.Value(key).(string)
	if !ok {
		return 0, fmt.Errorf("key[%v] value type is not string", key)
	}

	number, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("key[%v] value type is not number: %v", key, err)
	}

	return number, nil
}

func SetPairContext(ctx context.Context, key, value string) context.Context {
	return context.WithValue(ctx, key, value)
}
