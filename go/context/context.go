package context

import (
	"context"
	"time"
)

func WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {

	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	}
	return ctx, func() {}
}
