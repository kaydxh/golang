package mysql

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	database_ "github.com/kaydxh/golang/pkg/database"
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
	ctx, cancel := database_.WithDatabaseExecuteTimeout(ctx, timeout)
	defer cancel()

	return h(ctx, db, query)
}
