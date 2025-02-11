package duck

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type StoreSqliteMattn struct {
	DB *sql.DB
}

// DuckCreateTableIfNotExist creates the table if it does not exist.
func (s *StoreSqliteMattn) DuckCreateTableIfNotExist(_ context.Context) error {
	_, err := s.DB.Exec(`CREATE TABLE IF NOT EXISTS goose_versioning (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		version    INTEGER NOT NULL,
		source     TEXT NOT NULL
	)`)
	return err
}

// DuckInsertVersion inserts a new version entry.
func (s *StoreSqliteMattn) DuckInsertVersion(_ context.Context, version int64) (Migration, error) {
	m := Migration{Version: version}
	_, err := s.DB.Exec("INSERT INTO goose_versioning (version, source) VALUES (?, ?)", version, m.Source)
	if err != nil {
		return m, err
	}
	return m, nil
}

// DuckListMigrations lists all migrations.
func (s *StoreSqliteMattn) DuckListMigrations(_ context.Context) ([]Migration, error) {
	rows, err := s.DB.Query("SELECT created_at, version, source FROM goose_versioning ORDER BY version DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []Migration
	for rows.Next() {
		var m Migration
		var createdAt string
		if err := rows.Scan(&createdAt, &m.Version, &m.Source); err != nil {
			return nil, err
		}
		m.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, m)
	}
	return migrations, nil
}

// DuckExecuteMigration executes an SQL migration.
func (s *StoreSqliteMattn) DuckExecuteMigration(ctx context.Context, sqlMigration string) error {
	_, err := s.DB.Exec(sqlMigration)
	return err
}
