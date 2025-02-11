package main

import "time"

type Migration struct {
	CreatedAt time.Time
	Version   int64  `validate:"min=1"`
	Source    string `validate:"required"`
}

type Store interface {
	// CreateTable returns the SQL query string to create the db version table.
	CreateTableIfNotExist(tableName string) error

	// InsertVersion returns the SQL query string to insert a new version into the db version table.
	InsertVersion(tableName string, version int64) (Migration, error)

	// DeleteVersion returns the SQL query string to delete a version from the db version table.
	DeleteVersion(tableName string, version int64) error

	// GetMigrationByVersion returns the SQL query string to get a single migration by version.
	//
	// The query should return the timestamp and is_applied columns.
	GetMigrationByVersion(tableName string, version int64) (Migration, error)

	// ListMigrations returns the SQL query string to list all migrations in descending order by id.
	//
	// The query should return the version_id and is_applied columns.
	ListMigrations(tableName string) ([]Migration, error)

	// GetLatestVersion returns the SQL query string to get the last version_id from the db version
	// table. Returns a nullable int64 value.
	GetLatestVersion(tableName string) (Migration, error)
}
