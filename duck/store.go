package duck

const TimeFormat = "2006-01-02T15:04:05.999Z"

type DuckStore interface {
	// CreateTable creates the table if it does not exist.
	DuckCreateTableIfNotExist() error

	// DuckInsertVersion inserts a new version entry.
	DuckInsertVersion(version int64) (Migration, error)

	// DuckListMigrations lists all applied migrations.
	DuckListMigrations() ([]Migration, error)

	// DuckExecuteMigration executes a SQL migration script in a transaction.
	DuckExecuteMigration(sqlMigrationScript string) error
}
