package duck_test

import (
	"duck-migration/duck"
	"duck-migration/tool"
	"log"
	"os"
	"path"
	"testing"

	"zombiezen.com/go/sqlite"
)

func TestSqliteZombiezen(t *testing.T) {
	cwd, err := os.Getwd()
	duckTestDir := path.Join(cwd, "store_sqlite_zombiezen_test")

	testDataDir := "test_data"
	testDataFs := os.DirFS(cwd)

	cleanupFunc := func() {
		err = os.RemoveAll(duckTestDir)
		tool.Assert(err == nil, "Failed to remove test data directory", "err", err)
	}
	cleanupFunc()
	defer cleanupFunc()

	if err = os.MkdirAll(duckTestDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Open an in-memory database.
	conn, err := sqlite.OpenConn(path.Join(duckTestDir, "duck_test.db"))
	tool.Assert(err == nil, "failed to open database connection", "error", err)
	defer conn.Close()

	sqliteDuckStore := duck.StoreSqliteZombiezen{
		Conn: conn,
	}

	duck.Up(&sqliteDuckStore, testDataFs, testDataDir)
}
