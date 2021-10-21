package mysql

import (
	"context"
	"time"

	context_ "github.com/kaydxh/golang/go/context"

	"github.com/jmoiron/sqlx"
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
	ctx, cancel := context_.WithTimeout(ctx, timeout)
	defer cancel()

	return h(ctx, db, query)
}
