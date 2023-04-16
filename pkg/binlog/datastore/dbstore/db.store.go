/*
 *Copyright (c) 2023, kaydxh
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
package dbstore

import (
	"context"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	ds_ "github.com/kaydxh/golang/pkg/binlog/datastore"
	mysql_ "github.com/kaydxh/golang/pkg/database/mysql"
)

var _ ds_.DataStore = (*DBDataStore)(nil)

type DBDataStore struct {
	*sqlx.DB
}

func NewDBDataStore(db *sqlx.DB) (*DBDataStore, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	s := &DBDataStore{}
	s.DB = db
	return s, nil
}

func (s *DBDataStore) WriteData(ctx context.Context, query string, arg interface{}, key string, p [][]byte) (file *os.File, n int64, err error) {
	n, err = mysql_.ExecContext(ctx, query, mysql_.BuildNamedColumnsValuesBatch(arg), nil, s.DB)
	return nil, n, err
}
