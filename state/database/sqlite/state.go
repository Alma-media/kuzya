package sqlite

import (
	"database/sql"
	"errors"
	"time"

	db "github.com/Alma-media/kuzya/state/database"
	_ "github.com/mattn/go-sqlite3"
)

var (
	selectQuery = `SELECT state FROM state WHERE device_id = ?;`
	insertQuery = `INSERT INTO state (device_id, state) VALUES(?, ?) ON CONFLICT(device_id) DO UPDATE SET state=?, updated_at=CURRENT_TIMESTAMP;`

	errNoDevice = errors.New("device not found")
)

type StateManager struct{}

func (StateManager) Get(tx db.QueryRower, deviceID string) (bool, error) {
	var state bool

	err := tx.QueryRow(selectQuery, deviceID).Scan(&state)
	if err == sql.ErrNoRows {
		return state, errNoDevice
	}

	return state, err
}

func (StateManager) Set(tx db.Execer, deviceID string, state bool) error {
	_, err := tx.Exec(insertQuery, deviceID, state, state, time.Now())

	return err
}
