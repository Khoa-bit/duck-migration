package duck

import (
	"time"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

var _ Store = (*StoreSqliteZombiezen)(nil)

type StoreSqliteZombiezen struct {
	Conn *sqlite.Conn
}

// DuckCreateTableIfNotExist implements Store.
func (s *StoreSqliteZombiezen) DuckCreateTableIfNotExist() error {
	return sqlitex.Execute(s.Conn, `CREATE TABLE IF NOT EXISTS goose_versioning (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at TIMESTAMP DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
		version    INTEGER NOT NULL,
		source     TEXT NOT NULL
	)`, nil)
}

// DuckInsertVersion implements Store.
func (s *StoreSqliteZombiezen) DuckInsertVersion(version int64) (Migration, error) {
	migration := Migration{
		Version: version,
	}
	err := sqlitex.Execute(s.Conn, "INSERT INTO goose_versioning (version, source) VALUES (?, ?)", &sqlitex.ExecOptions{
		Args: []any{version, migration.Source},
	})
	return migration, err
}

// DuckListMigrations implements Store.
func (s *StoreSqliteZombiezen) DuckListMigrations() ([]Migration, error) {
	var migrations []Migration
	err := sqlitex.Execute(s.Conn, "SELECT created_at, version, source FROM goose_versioning ORDER BY version DESC", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			createAt, err := time.Parse(TimeFormat, stmt.ColumnText(0))
			if err != nil {
				return err
			}

			migrations = append(migrations, Migration{
				CreatedAt: createAt,
				Version:   stmt.ColumnInt64(1),
				Source:    stmt.ColumnText(2),
			})
			return nil
		},
	})
	return migrations, err
}

// DuckExecuteMigration implements Store.
func (s *StoreSqliteZombiezen) DuckExecuteMigration(sqlMigration string) error {
	return sqlitex.ExecuteScript(s.Conn, sqlMigration, nil)
}
