package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "" CHECK (LENGTH(date) = 8),
    title VARCHAR(255) NOT NULL DEFAULT "" CHECK (LENGTH(title) <= 255),
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT "" CHECK (LENGTH(repeat) <= 128)
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);
`

var db *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := err != nil

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return fmt.Errorf("error in creating schema: %w", err)
		}
	}

	return nil
}
