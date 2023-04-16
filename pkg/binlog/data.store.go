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
package binlog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jmoiron/sqlx"
	mysql_ "github.com/kaydxh/golang/pkg/database/mysql"
	rotate_ "github.com/kaydxh/golang/pkg/file-rotate"
)

type DataStore interface {
	// query arg for db, p for filedata
	// file for filestoreage
	WriteData(ctx context.Context, key string, query string, arg interface{}, p [][]byte) (file *os.File, n int64, err error)
}

type FileDataStore struct {
	rotateFiler    *rotate_.RotateFiler
	rotateFilers   map[string]*rotate_.RotateFiler //message key -> rotateFilter
	rotateFilersMu sync.Mutex
	rootDir        string
	opts           []rotate_.RotateFilerOption
}

func NewFileDataStore(filedir string, options ...rotate_.RotateFilerOption) (*FileDataStore, error) {
	rotate, err := rotate_.NewRotateFiler(filedir, options...)
	if err != nil {
		return nil, err
	}
	s := &FileDataStore{
		rootDir: filedir,
		opts:    options,
	}
	s.rotateFiler = rotate
	return s, nil
}

func (s *FileDataStore) WriteData(ctx context.Context, key string, query string, arg interface{}, p [][]byte) (file *os.File, n int64, err error) {
	rotateFiler := s.getOrCreate(ctx, key)
	file, tn, err := rotateFiler.WriteBytesLine(p)
	return file, int64(tn), err
}

func (s *FileDataStore) getOrCreate(ctx context.Context, key string) *rotate_.RotateFiler {
	if key == "" {
		return s.rotateFiler
	}

	s.rotateFilersMu.Lock()
	defer s.rotateFilersMu.Unlock()
	rotateFiler, ok := s.rotateFilers[key]
	if !ok {
		rotateFiler, _ = rotate_.NewRotateFiler(
			filepath.Join(s.rootDir, key),
			s.opts...,
		)
		s.rotateFilers[key] = rotateFiler
	}

	return rotateFiler
}

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

func (s *DBDataStore) WriteData(ctx context.Context, query string, arg interface{}, p [][]byte) (file *os.File, n int64, err error) {
	n, err = mysql_.ExecContext(ctx, query, mysql_.BuildNamedColumnsValuesBatch(arg), nil, s.DB)
	return nil, n, err
}
