package serverpkg

import (
	"database/sql"
)

// InitDB opens (and creates if necessary) a sqlite database at dbPath.
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	// Keep initialization minimal here; callers may run migrations.
	return db, nil
}
