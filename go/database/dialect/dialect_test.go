package dialect

import (
	"testing"

	"github.com/jmoiron/sqlx"
)

// newNamedDB wraps a nil *sql.DB: only the driver name is read here.
func newNamedDB(driver string) *sqlx.DB {
	return sqlx.NewDb(nil, driver) //nolint
}

func TestOfByDriverName(t *testing.T) {
	// DriverName is read from the sqlx wrapper; the constant strings are
	// the database/sql driver registrations keel uses.
	cases := map[string]Dialect{
		"mysql":  MySQL,
		"sqlite": SQLite,
	}
	for driver, want := range cases {
		db := newNamedDB(driver)
		if got := Of(db); got != want {
			t.Fatalf("Of(%q) = %q, want %q", driver, got, want)
		}
	}
	if got := Of(newNamedDB("postgres")); got != "postgres" {
		t.Fatalf("unknown drivers must be reported as themselves, got %q", got)
	}
}

func TestRowLockPerEngine(t *testing.T) {
	if got := MySQL.RowLock(false); got != " FOR UPDATE" {
		t.Fatalf("mysql lock = %q", got)
	}
	if got := MySQL.RowLock(true); got != " FOR UPDATE SKIP LOCKED" {
		t.Fatalf("mysql skip-locked = %q", got)
	}
	for _, skip := range []bool{false, true} {
		if got := SQLite.RowLock(skip); got != "" {
			t.Fatalf("sqlite has no row locks, got %q", got)
		}
	}
}

func TestNowExprPerEngine(t *testing.T) {
	if got := MySQL.NowExpr(); got != "NOW(3)" {
		t.Fatalf("mysql now = %q", got)
	}
	// SQLite's CURRENT_TIMESTAMP truncates to seconds; the expression must
	// carry sub-second precision.
	if got := SQLite.NowExpr(); got == "CURRENT_TIMESTAMP" {
		t.Fatal("sqlite now must not be second-granular")
	}
}

func TestTranslateRefusesUnknown(t *testing.T) {
	if _, err := Dialect("oracle").Translate("SELECT 1"); err == nil {
		t.Fatal("unknown engines must refuse translation")
	}
	for _, d := range []Dialect{MySQL, SQLite} {
		if _, err := d.Translate("SELECT ?"); err != nil {
			t.Fatalf("%s: %v", d, err)
		}
	}
}
