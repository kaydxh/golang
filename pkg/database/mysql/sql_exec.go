package mysql

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	ExecuteTimeout = time.Minute
)

func ExecContext(
	ctx context.Context,
	timeout time.Duration,
	db *sqlx.DB,
	query string,
	h func(ctx context.Context, db *sqlx.DB, query string) error,
) error {
	if h == nil {
		return nil
	}
	ctx, cancel := WithDatabaseExecuteTimeout(ctx, timeout)
	defer cancel()

	return h(ctx, db, query)
}

func WithDatabaseExecuteTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {

	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	}
	return ctx, func() {}
}
