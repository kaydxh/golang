package mysql

import (
	"sync/atomic"

	"github.com/jmoiron/sqlx"
)

type SQLDB atomic.Value

//check type SQLDB* whether to implement interface of atomic.Value
//var _ atomic.Value = (*SQLDB)(nil)

//check type SQLDB whether to implement interface of atomic.Value
//var _ atomic.Value = SQLDB{}

func (m *SQLDB) Store(value *sqlx.DB) {
	(*atomic.Value)(m).Store(value)
}

func (m *SQLDB) Load() *sqlx.DB {
	value := (*atomic.Value)(m).Load()
	if value == nil {
		return nil
	}
	return value.(*sqlx.DB)
}
