package mysql

import (
	"context"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestTxPipelinedCommitsOnlyAfterCallbackSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectCommit()
	if err := TxPipelined(context.Background(), sqlxDB, func(*sqlx.Tx) error {
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestTxPipelinedRollsBackCallbackError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	want := errors.New("callback failed")

	mock.ExpectBegin()
	mock.ExpectRollback()
	if err := TxPipelined(context.Background(), sqlxDB, func(*sqlx.Tx) error {
		return want
	}); !errors.Is(err, want) {
		t.Fatalf("err = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestTxPipelinedRollsBackPanicAndPropagatesIt(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectBegin()
	mock.ExpectRollback()
	defer func() {
		if recover() == nil {
			t.Fatal("panic was swallowed")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	}()
	_ = TxPipelined(context.Background(), sqlxDB, func(*sqlx.Tx) error {
		panic("callback panicked")
	})
}
