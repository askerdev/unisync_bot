package sqlite

import "database/sql"

const schema = `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id TEXT NOT NULL,
		text TEXT NOT NULL,
		time_at TEXT NOT NULL
  );
`

func ApplyMigrations(db *sql.DB) error {
	_, err := db.Exec(schema)
	return err
}
