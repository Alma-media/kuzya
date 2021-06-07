package sqlite

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	selectQuery = `SELECT state FROM state WHERE device_id = ?;`
	insertQuery = `INSERT INTO state (device_id, state) VALUES(?, ?) ON CONFLICT(device_id) DO UPDATE SET state=?, updated_at=CURRENT_TIMESTAMP;`
)

type StateManager struct{}

func (StateManager) Get(tx *sql.Tx, deviceID string) (bool, error) {
	var state bool

	return state, tx.QueryRow(selectQuery, deviceID).Scan(&state)
}

func (StateManager) Set(tx *sql.Tx, deviceID string, state bool) error {
	_, err := tx.Exec(insertQuery, deviceID, state, state, time.Now())

	return err
}
