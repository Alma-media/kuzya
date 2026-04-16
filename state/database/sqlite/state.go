package sqlite

import (
	"database/sql"
	"errors"
	"time"

	db "github.com/Alma-media/kuzya/state/database"
	"github.com/Alma-media/kuzya/state"
	_ "github.com/mattn/go-sqlite3"
)

var (
	selectQuery = `SELECT state FROM state WHERE device_id = ?;`
	insertQuery = `INSERT INTO state (device_id, state) VALUES(?, ?) ON CONFLICT(device_id) DO UPDATE SET state=?, updated_at=CURRENT_TIMESTAMP;`
	listQuery   = `SELECT device_id, state FROM state ORDER BY device_id;`

	errNoDevice = errors.New("device not found")
)

type StateManager struct{}

func (StateManager) Get(tx db.QueryRower, deviceID string) (bool, error) {
	var status bool

	err := tx.QueryRow(selectQuery, deviceID).Scan(&status)
	if err == sql.ErrNoRows {
		return status, errNoDevice
	}

	return status, err
}

func (StateManager) Set(tx db.Execer, deviceID string, s bool) error {
	_, err := tx.Exec(insertQuery, deviceID, s, s, time.Now())

	return err
}

func (StateManager) List(tx db.Querier) ([]state.Device, error) {
	rows, err := tx.Query(listQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var devices []state.Device

	for rows.Next() {
		var (
			deviceID string
			status   bool
		)

		if err := rows.Scan(&deviceID, &status); err != nil {
			return nil, err
		}

		devices = append(devices, state.Device{
			ID:    deviceID,
			State: state.Status(status),
		})
	}

	return devices, rows.Err()
}
