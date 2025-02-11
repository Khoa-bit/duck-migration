package main

import (
	"fmt"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

var _ Store = (*StoreSqliteZombiezen)(nil)

type StoreSqliteZombiezen struct {
	conn *sqlite.Conn
}

// CreateTableIfNotExist implements Store.
func (s *StoreSqliteZombiezen) CreateTableIfNotExist(tableName string) error {
	q := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		version_id INTEGER NOT NULL,
		create_at  TIMESTAMP DEFAULT (datetime('now'))
	)`, tableName)

	return sqlitex.Execute(s.conn, q, nil)
}

// DeleteVersion implements Store.
func (s *StoreSqliteZombiezen) DeleteVersion(tableName string, version int64) error {
	q := fmt.Sprintf("DELETE FROM %s WHERE version_id = ?", tableName)
	return sqlitex.Execute(s.conn, q)
}

// GetLatestVersion implements Store.
func (s *StoreSqliteZombiezen) GetLatestVersion(tableName string) (Migration, error) {
	panic("unimplemented")
}

// GetMigrationByVersion implements Store.
func (s *StoreSqliteZombiezen) GetMigrationByVersion(tableName string, version int64) (Migration, error) {
	panic("unimplemented")
}

// InsertVersion implements Store.
func (s *StoreSqliteZombiezen) InsertVersion(tableName string, version int64) (Migration, error) {
	panic("unimplemented")
}

// ListMigrations implements Store.
func (s *StoreSqliteZombiezen) ListMigrations(tableName string) ([]Migration, error) {
	panic("unimplemented")
}
