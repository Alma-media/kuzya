package sqlite

import (
	"context"
	"database/sql"
	"testing"
)

func setup(t *testing.T) (*sql.DB, func() error) {
	db, err := sql.Open("sqlite3", "file::memory:")
	if err != nil {
		t.Fatalf("unable to establish database connection: %s", err)
	}

	if err := Init(context.Background(), db); err != nil {
		t.Fatalf("database migration failure: %s", err)
	}

	return db, db.Close
}

func TestStateManagerGet(t *testing.T) {
	db, release := setup(t)
	defer release()

	t.Run("get state for registered device", func(t *testing.T) {
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf("unable to start the transaction: %s", err)
		}

		defer tx.Rollback()

		if _, err := tx.Exec(insertQuery, "device-id", true, true); err != nil {
			t.Fatalf("unable to prepare test data: %s", err)
		}

		state, err := new(StateManager).Get(tx, "device-id")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if !state {
			t.Error("state was expected to be true")
		}
	})

	t.Run("get state for unknown/unregistered device", func(t *testing.T) {
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf("unable to start the transaction: %s", err)
		}

		defer tx.Rollback()

		state, err := new(StateManager).Get(tx, "unknown-device")
		if err != errNoDevice {
			t.Errorf("unexpected error: %s", err)
		}

		if state {
			t.Error("state was expected to be false")
		}
	})
}

func TestStateManagerSet(t *testing.T) {
	db, release := setup(t)
	defer release()

	t.Run("set state for unregistered device", func(t *testing.T) {
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf("unable to start the transaction: %s", err)
		}

		defer tx.Rollback()

		if err := new(StateManager).Set(tx, "device-id", true); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		var state bool

		if err := tx.QueryRow(selectQuery, "device-id").Scan(&state); err != nil {
			t.Fatalf("unable to read device state: %s", err)
		}

		if !state {
			t.Error("state was expected to be true")
		}
	})

	t.Run("update state for registered device", func(t *testing.T) {
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf("unable to start the transaction: %s", err)
		}

		defer tx.Rollback()

		if _, err := tx.Exec(insertQuery, "device-id", true, true); err != nil {
			t.Fatalf("unable to prepare test data: %s", err)
		}

		if err := new(StateManager).Set(tx, "device-id", false); err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		var state bool

		if err := tx.QueryRow(selectQuery, "device-id").Scan(&state); err != nil {
			t.Fatalf("unable to read device state: %s", err)
		}

		if state {
			t.Error("state was expected to be false")
		}
	})
}
