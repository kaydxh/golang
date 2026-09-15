// Package sqlite opens pure-Go SQLite databases (modernc.org/sqlite, no
// cgo, so binaries cross-compile) behind the standard database/sql surface
// with the pragmas a concurrent single-writer workload needs.
//
// Concurrency model: journal_mode=WAL lets readers proceed while one
// writer commits; busy_timeout makes a second writer wait instead of
// failing at once. The pool keeps readers parallel — SQLite serializes
// writers itself, and the timeout absorbs the overlap.
//
// The database file must live on a local filesystem: a network filesystem
// violates SQLite's locking assumptions. Callers that choose the path own
// that rule; this package only documents it.
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite" // registers the "sqlite" driver, cgo-free
)

// Options tune Open. The zero value is the recommended default.
type Options struct {
	// MaxOpenConns bounds the pool. Zero uses 4 — enough for parallel
	// readers while writers queue behind busy_timeout. 1 serializes
	// everything and is rarely what a service wants.
	MaxOpenConns int
	// BusyTimeoutMs is how long a blocked writer waits before failing.
	// Zero uses 5000.
	BusyTimeoutMs int
	// ForeignKeys enforces REFERENCES constraints. Default true; a schema
	// whose foreign keys are decorative should say so by passing false.
	ForeignKeys *bool
}

// Open opens (creating if needed) the SQLite database at path.
func Open(path string, opts Options) (*sqlx.DB, error) {
	if path == "" {
		return nil, fmt.Errorf("sqlite: path is empty")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("sqlite: create db dir: %w", err)
	}

	busy := opts.BusyTimeoutMs
	if busy <= 0 {
		busy = 5000
	}
	foreignKeys := true
	if opts.ForeignKeys != nil {
		foreignKeys = *opts.ForeignKeys
	}

	dsn := fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(%d)", path, busy)
	if foreignKeys {
		dsn += "&_pragma=foreign_keys(1)"
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlite: open %s: %w", path, err)
	}
	maxOpen := opts.MaxOpenConns
	if maxOpen <= 0 {
		maxOpen = 4
	}
	db.SetMaxOpenConns(maxOpen)
	return sqlx.NewDb(db, "sqlite"), nil
}

// Ping verifies the database answers.
func Ping(ctx context.Context, db *sqlx.DB) error {
	var one int
	return db.GetContext(ctx, &one, "SELECT 1")
}
