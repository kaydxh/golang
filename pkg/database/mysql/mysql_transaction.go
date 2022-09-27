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

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TxDao struct {
	*sqlx.Tx
}

func (d *TxDao) Begin(ctx context.Context, db *sqlx.DB, opts *sql.TxOptions) error {
	if db == nil {
		return fmt.Errorf("unexpected err: db is nil")
	}

	tx, err := db.BeginTxx(ctx, opts)
	if err != nil {
		return err
	}

	d.Tx = tx

	return nil
}

func (d *TxDao) Commit() error {
	if d.Tx == nil {
		return fmt.Errorf("unexpected err: tx is nil")
	}

	err := d.Tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// Rollback ...
func (d *TxDao) Rollback() error {
	if d.Tx == nil {
		return fmt.Errorf("unexpected err: tx is nil")
	}

	err := d.Tx.Rollback()
	if err != nil {
		return err
	}
	return nil
}

func TxPipelined(ctx context.Context, db *sqlx.DB, fn func(*sqlx.Tx) error) (err error) {
	var tx TxDao
	err = tx.Begin(ctx, db, nil)
	if err != nil {
		logrus.WithError(err).Errorf("failed to transaction begin")
		return err
	}

	defer func() {
		if err != nil {
			if txErr := tx.Rollback(); txErr != nil {
				logrus.WithError(err).Errorf("failed to rollback, err: %v", txErr)
				return
			}
			return
		}

		if err = tx.Commit(); err != nil {
			logrus.WithError(err).Errorf("failed to commit")
			return
		}

	}()

	return fn(tx.Tx)
}
