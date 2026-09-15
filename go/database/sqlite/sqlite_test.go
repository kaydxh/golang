package sqlite

import (
	"context"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

func openTestDB(t *testing.T) *sqlx.DB {
	t.Helper()
	db, err := Open(filepath.Join(t.TempDir(), "test.db"), Options{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := Ping(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	return db
}

// TestOpenCreatesAndReopens: an opened database is pingable and the file
// persists across close/reopen with its committed data.
func TestOpenCreatesAndReopens(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a.db")
	db, err := Open(path, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`CREATE TABLE t (k TEXT PRIMARY KEY)`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO t (k) VALUES ('committed')`); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	reopened, err := Open(path, Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	var k string
	if err := reopened.Get(&k, `SELECT k FROM t`); err != nil {
		t.Fatalf("reopened read: %v", err)
	}
	if k != "committed" {
		t.Fatalf("k = %q", k)
	}
}

// TestOpenRejectsEmptyPath and creates the parent directory when missing.
func TestOpenPathHandling(t *testing.T) {
	if _, err := Open("", Options{}); err == nil {
		t.Fatal("empty path must be refused")
	}
	nested := filepath.Join(t.TempDir(), "deep", "dir", "b.db")
	db, err := Open(nested, Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
}

// TestLeaseClaimExactlyOnce: the queue pattern — UPDATE ... WHERE
// state='queued' with affected-rows arbitration — claims each job exactly
// once no matter how many workers race.
func TestLeaseClaimExactlyOnce(t *testing.T) {
	db := openTestDB(t)
	if _, err := db.Exec(`CREATE TABLE jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		state TEXT NOT NULL DEFAULT 'queued',
		lease_owner TEXT,
		lease_until DATETIME)`); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if _, err := db.Exec(`INSERT INTO jobs (state) VALUES ('queued')`); err != nil {
			t.Fatal(err)
		}
	}

	claim := func(owner string) int {
		res, err := db.Exec(`UPDATE jobs SET state='running', lease_owner=?, lease_until=? WHERE id=(
			SELECT id FROM jobs WHERE state='queued' ORDER BY id LIMIT 1)`,
			owner, time.Now().Add(time.Minute))
		if err != nil {
			t.Fatal(err)
		}
		n, _ := res.RowsAffected()
		return int(n)
	}
	if claim("w1")+claim("w2") != 2 {
		t.Fatal("each job must be claimed exactly once")
	}
	if claim("w3") != 0 {
		t.Fatal("no job may be claimed twice")
	}
}

// TestTransactionRollback: a failed multi-statement transaction leaves
// nothing behind.
func TestTransactionRollback(t *testing.T) {
	db := openTestDB(t)
	if _, err := db.Exec(`CREATE TABLE t (k TEXT PRIMARY KEY)`); err != nil {
		t.Fatal(err)
	}
	err := func() error {
		tx, err := db.Beginx()
		if err != nil {
			return err
		}
		defer tx.Rollback()
		if _, err := tx.Exec(`INSERT INTO t (k) VALUES ('a')`); err != nil {
			return err
		}
		if _, err := tx.Exec(`INSERT INTO t (k) VALUES ('a')`); err == nil {
			t.Fatal("duplicate insert must fail")
		}
		return tx.Rollback()
	}()
	if err != nil {
		t.Fatal(err)
	}
	var count int
	if err := db.Get(&count, `SELECT COUNT(*) FROM t`); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("rolled-back transaction must leave nothing behind")
	}
}

// TestConcurrentWritersUnderWAL: several writers plus a reader in flight —
// no lost writes, no spurious failures.
func TestConcurrentWritersUnderWAL(t *testing.T) {
	db := openTestDB(t)
	if _, err := db.Exec(`CREATE TABLE w (id INTEGER PRIMARY KEY AUTOINCREMENT, v TEXT)`); err != nil {
		t.Fatal(err)
	}
	const writers, each = 4, 25
	var wg sync.WaitGroup
	for w := range writers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for range each {
				if _, err := db.Exec(`INSERT INTO w (v) VALUES (?)`, id); err != nil {
					t.Errorf("writer %d: %v", id, err)
					return
				}
			}
		}(w)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range each {
			var n int
			if err := db.Get(&n, `SELECT COUNT(*) FROM w`); err != nil {
				t.Errorf("reader: %v", err)
				return
			}
		}
	}()
	wg.Wait()

	var total int
	if err := db.Get(&total, `SELECT COUNT(*) FROM w`); err != nil {
		t.Fatal(err)
	}
	if total != writers*each {
		t.Fatalf("rows = %d, want %d — writes were lost", total, writers*each)
	}
}
