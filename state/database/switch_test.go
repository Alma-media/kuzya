package database_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Alma-media/kuzya/state"
	"github.com/Alma-media/kuzya/state/database"
	"github.com/Alma-media/kuzya/state/database/sqlite"
)

func setupSwitch(t *testing.T) *database.Switch {
	t.Helper()

	db, err := sql.Open("sqlite3", "file::memory:")
	if err != nil {
		t.Fatalf("unable to establish database connection: %s", err)
	}

	if err := sqlite.Init(context.Background(), db); err != nil {
		db.Close()
		t.Fatalf("database migration failure: %s", err)
	}

	t.Cleanup(func() { db.Close() })

	return database.NewSwitch(db, new(sqlite.StateManager))
}

func TestSetState(t *testing.T) {
	sw := setupSwitch(t)

	// Sets an unregistered device (upsert) and returns the new state.
	if got, err := sw.SetState("lamp", "ON"); err != nil || got != state.ON {
		t.Fatalf(`SetState("ON") = %q, %v; want %q, nil`, got, err, state.ON)
	}

	// Tolerates case and surrounding whitespace.
	if got, err := sw.SetState("lamp", " off "); err != nil || got != state.OFF {
		t.Fatalf(`SetState(" off ") = %q, %v; want %q, nil`, got, err, state.OFF)
	}

	// Persisted: Status sees the same value.
	if got, err := sw.Status("lamp"); err != nil || got != state.OFF {
		t.Fatalf(`Status = %q, %v; want %q, nil`, got, err, state.OFF)
	}

	// A payload without absolute state is rejected and changes nothing.
	if _, err := sw.SetState("lamp", "TRIG"); err == nil {
		t.Error("TRIG payload must be rejected")
	}

	if got, _ := sw.Status("lamp"); got != state.OFF {
		t.Errorf("state after rejected payload = %q, want %q", got, state.OFF)
	}

	// The whole point: a toggle after a set flips from the set value, so
	// absolute controllers and physical buttons share one consistent phase.
	if got, err := sw.Switch("lamp"); err != nil || got != state.ON {
		t.Fatalf(`Switch after set OFF = %q, %v; want %q, nil`, got, err, state.ON)
	}
}
