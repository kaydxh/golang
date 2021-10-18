package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TxDao struct {
	*sqlx.Tx
}

func (d *TxDao) Begin(ctx context.Context, db *sqlx.DB, opts *sql.TxOptions) error {
	if db == nil {
		return fmt.Errorf("db is nil")
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
		return fmt.Errorf("tx is nil")
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
		return fmt.Errorf("tx is nil")
	}

	err := d.Tx.Rollback()
	if err != nil {
		return err
	}
	return nil
}
