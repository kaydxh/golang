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
