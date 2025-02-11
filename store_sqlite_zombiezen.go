package main

import "zombiezen.com/go/sqlite"

var _ Store = (*StoreSqliteZombiezen)(nil)

type StoreSqliteZombiezen struct {
	conn *sqlite.Conn
}

// CreateTableIfNotExist implements Store.
func (s *StoreSqliteZombiezen) CreateTableIfNotExist(tableName string) error {
	q := `CREATE TABLE %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		version_id INTEGER NOT NULL,
		is_applied INTEGER NOT NULL,
		tstamp TIMESTAMP DEFAULT (datetime('now'))
	)`
	_, err := s.conn.Exec("CREATE TABLE IF NOT EXISTS " + tableName + " (version_id INTEGER PRIMARY KEY, is_applied BOOLEAN NOT NULL DEFAULT FALSE)")
	return err
}

// DeleteVersion implements Store.
func (s *StoreSqliteZombiezen) DeleteVersion(tableName string) string {
	panic("unimplemented")
}

// GetLatestVersion implements Store.
func (s *StoreSqliteZombiezen) GetLatestVersion(tableName string) string {
	panic("unimplemented")
}

// GetMigrationByVersion implements Store.
func (s *StoreSqliteZombiezen) GetMigrationByVersion(tableName string) string {
	panic("unimplemented")
}

// InsertVersion implements Store.
func (s *StoreSqliteZombiezen) InsertVersion(tableName string) string {
	panic("unimplemented")
}

// ListMigrations implements Store.
func (s *StoreSqliteZombiezen) ListMigrations(tableName string) string {
	panic("unimplemented")
}
