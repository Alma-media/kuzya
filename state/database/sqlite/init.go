package sqlite

import (
	"context"
	"database/sql"
)

var createTable = "CREATE TABLE IF NOT EXISTS \"state\" (" +
	"`device_id` text PRIMARY KEY," +
	"`state` bool NOT NULL," +
	"`created_at` datetime DEFAULT CURRENT_TIMESTAMP," +
	"`updated_at` datetime" +
	");"

func Init(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, createTable)

	return err
}
