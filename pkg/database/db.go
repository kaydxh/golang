package database

import (
	"context"
	"time"
)

var (
	ExecuteTimeout = time.Minute
)

func WithDatabaseExecuteTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {

	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	}
	return ctx, func() {}
}
