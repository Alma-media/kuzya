package database

import (
	"context"
	"database/sql"

	"github.com/Alma-media/kuzya/state"
)

type QueryRower interface {
	QueryRow(query string, args ...interface{}) *sql.Row
}

type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type StateManager interface {
	Get(tx QueryRower, deviceID string) (bool, error)
	Set(tx Execer, deviceID string, state bool) error
}

type Switch struct {
	state StateManager
	db    *sql.DB
}

func NewSwitch(db *sql.DB, state StateManager) *Switch {
	return &Switch{
		db:    db,
		state: state,
	}
}

func (sw *Switch) Switch(deviceID string) (string, error) {
	tx, err := sw.db.BeginTx(context.Background(), nil)
	if err != nil {
		return "", err
	}

	defer tx.Rollback()

	status, err := sw.state.Get(tx, deviceID)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	if err := sw.state.Set(tx, deviceID, !status); err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return state.Status(!status), nil
}

// SetState applies an absolute state ("ON"/"OFF") instead of toggling. An
// unknown device simply starts at the requested state (the upsert below);
// an invalid payload changes nothing.
func (sw *Switch) SetState(deviceID, requested string) (string, error) {
	target, err := state.Parse(requested)
	if err != nil {
		return "", err
	}

	tx, err := sw.db.BeginTx(context.Background(), nil)
	if err != nil {
		return "", err
	}

	defer tx.Rollback()

	if err := sw.state.Set(tx, deviceID, target); err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return state.Status(target), nil
}

func (sw *Switch) Status(deviceID string) (string, error) {
	status, err := sw.state.Get(sw.db, deviceID)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	return state.Status(status), nil
}
