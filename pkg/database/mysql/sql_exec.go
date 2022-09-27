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
package mysql

import (
	"context"
	"database/sql"
	"fmt"

	runtime_ "github.com/kaydxh/golang/go/runtime"
	time_ "github.com/kaydxh/golang/go/time"
	"github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
)

/*
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
*/

func ExecContext(
	ctx context.Context,
	query string,
	arg interface{},
	tx *sqlx.Tx,
	db *sqlx.DB,
) (err error) {
	tc := time_.New(true)
	caller := runtime_.GetShortCaller()
	logger := logrus.WithField("caller", caller)

	clean := func() {
		tc.Tick(caller)
		logger.WithField("cost", tc.String()).Infof("SQL EXECL")
		if err != nil {
			logger.WithError(err).Errorf("sql: %s", query)
		}
	}
	defer clean()

	result, err := NamedExecContext(ctx, query, arg, tx, db)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	logger.Infof("affected rows: %v", rows)

	return nil
}

// pointer struct model for dest
func SelectNamedContext(
	ctx context.Context,
	query string,
	arg interface{},
	dest interface{},
	db *sqlx.DB,
) (err error) {
	tc := time_.New(true)
	caller := runtime_.GetShortCaller()
	logger := logrus.WithField("caller", caller)

	clean := func() {
		tc.Tick(caller)
		logger.WithField("cost", tc.String()).Infof("SQL EXECL")
		if err != nil {
			logger.WithError(err).Errorf("sql: %s", query)
		}
	}
	defer clean()

	// prepare
	ns, err := db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer ns.Close()

	// query
	return ns.SelectContext(ctx, dest, arg)
}

func GetCountContext(ctx context.Context, query string, arg interface{}, db *sqlx.DB) (count uint32, err error) {
	tc := time_.New(true)
	caller := runtime_.GetShortCaller()
	logger := logrus.WithField("caller", caller)

	clean := func() {
		tc.Tick(caller)
		logger.WithField("cost", tc.String()).Infof("SQL EXECL")
		if err != nil {
			logger.WithError(err).Errorf("sql: %s", query)
		}
	}
	defer clean()

	ns, err := db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer ns.Close()

	err = ns.QueryRowContext(ctx, arg).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil

}

func NamedExecContext(
	ctx context.Context,
	query string,
	arg interface{},
	tx *sqlx.Tx,
	db *sqlx.DB,
) (sql.Result, error) {

	if tx != nil {
		return tx.NamedExecContext(ctx, query, arg)
	}

	if db != nil {
		return db.NamedExecContext(ctx, query, arg)
	}

	return nil, fmt.Errorf("db is nil")
}

func PrepareNamedContext(ctx context.Context,
	query string,
	tx *sqlx.Tx,
	db *sqlx.DB,
) (*sqlx.NamedStmt, error) {

	if tx != nil {
		return tx.PrepareNamedContext(ctx, query)
	}

	if db != nil {
		return db.PrepareNamedContext(ctx, query)
	}

	return nil, fmt.Errorf("db is nil")
}
