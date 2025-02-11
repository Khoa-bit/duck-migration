package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Migration struct {
	Version int64  `validate:"min=1"`
	Source  string `validate:"required"`
}

func main() {
	// Open an in-memory database.
	conn, err := sqlite.OpenConn(":memory:", sqlite.OpenReadWrite)
	Assert(err == nil, "failed to open database connection", "error", err)
	defer conn.Close()

	Up(context.Background(), embedMigrations, "migrations")
}

func Up(ctx context.Context, migrationsFs fs.FS, dirpath string) {
	// Collect all SQL migration files.
	migrations, err := collectSqlMigrations(embedMigrations, "migrations")
	if err != nil {
		log.Fatalf("failed to collect SQL migrations: %v", err)
	}

	fmt.Println(migrations)

	// Sort the migrations by version number.
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Execute a query.
	err = sqlitex.ExecuteTransient(conn, "SELECT 'hello, world';", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			fmt.Println(stmt.ColumnText(0))
			return nil
		},
	})
	Assert(err == nil, "failed to execute query", "error", err)
}

func collectSqlMigrations(migrationsFs fs.FS, dirpath string) ([]Migration, error) {
	if _, err := fs.Stat(migrationsFs, dirpath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("%s directory does not exist", dirpath)
		}
		return nil, err
	}
	// SQL migration files.
	sqlMigrationFiles, err := fs.Glob(migrationsFs, path.Join(dirpath, "*.sql"))
	if err != nil {
		return nil, err
	}

	sqlMigrations := make([]Migration, 0, len(sqlMigrationFiles))
	for _, file := range sqlMigrationFiles {
		sqlMigrations = append(sqlMigrations, Migration{
			Version: NumericComponent(file),
			Source:  file,
		})
	}

	return sqlMigrations, nil
}

// NumericComponent parses the version from the migration file name.
//
// XXX_descriptivename.ext where XXX specifies the version number and ext specifies the type of
// .sql migration file.
func NumericComponent(filename string) int64 {
	base := filepath.Base(filename)
	Assert(filepath.Ext(base) == ".sql", "migration file does not have .sql file extension", "filename", filename)

	idx := strings.Index(base, "_")
	Assert(idx > 0, "no filename separator '_' found", "filename", filename)

	n, err := strconv.ParseInt(base[:idx], 10, 64)
	Assert(err == nil, "failed to parse version from migration file", "filename", filename, "error", err)
	Assert(n > 0, "migration version must be greater than zero", "filename", filename)

	return n
}
