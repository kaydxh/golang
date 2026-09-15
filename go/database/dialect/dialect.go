// Package dialect adapts SQL written for one engine to run unchanged on
// another where the difference is mechanical: row-locking suffixes and
// "now" expressions. Statements whose semantics differ (UPSERT shapes,
// window functions) are NOT translated here — those belong at the
// statement site with explicit per-engine branches, because a wrong
// mechanical translation of a semantic difference is a data bug.
package dialect

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Dialect names the SQL engine a database speaks.
type Dialect string

const (
	// MySQL is the multi-process server engine: row locks exist, NOW(3)
	// carries milliseconds.
	MySQL Dialect = "mysql"
	// SQLite is the embedded single-writer engine: row locks are
	// meaningless (the writer is serialized), CURRENT_TIMESTAMP is
	// second-granular so now() uses strftime.
	SQLite Dialect = "sqlite"
)

// Of reports which dialect a opened database speaks, by driver name. An
// unknown driver is reported as such: guessing an engine's SQL is how
// data bugs are born.
func Of(db *sqlx.DB) Dialect {
	switch db.DriverName() {
	case "mysql":
		return MySQL
	case "sqlite", "sqlite3":
		return SQLite
	}
	return Dialect(db.DriverName())
}

// RowLock renders the row-locking suffix for a SELECT inside a
// transaction. MySQL takes " FOR UPDATE" (optionally " SKIP LOCKED" so
// concurrent claimers skip rather than block); SQLite has no row locks —
// the engine serializes writers, so the suffix is empty. skipLocked only
// applies on MySQL.
func (d Dialect) RowLock(skipLocked bool) string {
	if d != MySQL {
		return ""
	}
	if skipLocked {
		return " FOR UPDATE SKIP LOCKED"
	}
	return " FOR UPDATE"
}

// NowExpr renders a non-constant "current timestamp" expression with
// sub-second precision: MySQL NOW(3); SQLite strftime to milliseconds
// (CURRENT_TIMESTAMP would truncate to whole seconds).
func (d Dialect) NowExpr() string {
	if d == SQLite {
		return "strftime('%Y-%m-%d %H:%M:%f','now')"
	}
	return "NOW(3)"
}

// Translate rewrites a dialect-portable statement for this engine. The
// only transformation is placeholder normalization for the ?-style
// bindvars both supported engines share — it exists so future mechanical
// differences have one seam instead of spreading through call sites.
func (d Dialect) Translate(query string) (string, error) {
	switch d {
	case MySQL, SQLite:
		return query, nil
	}
	return "", fmt.Errorf("dialect: unknown engine %q", d)
}
