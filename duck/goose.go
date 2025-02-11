package duck

import (
	"duck-migration/tool"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func Up(duckStore DuckStore, migrationsFs fs.FS, dirpath string) {
	collectedMigrations, err := collectSqlMigrations(migrationsFs, dirpath)
	tool.Assert(err == nil, "failed to collect SQL migrations", "error", err, "dirpath", dirpath)

	err = duckStore.DuckCreateTableIfNotExist()
	tool.Assert(err == nil, "failed to create goose_versioning table", "error", err)

	appliedMigrations, err := duckStore.DuckListMigrations()
	var latestAppliedVersion int64
	if len(appliedMigrations) > 0 {
		latestAppliedVersion = appliedMigrations[0].Version
	}

	missingMigrations := findMissingMigrations(collectedMigrations, appliedMigrations, latestAppliedVersion)
	tool.Assert(len(missingMigrations) == 0, "collected missing migrations that is older the latest applied version", "latestAppliedVersion", latestAppliedVersion, "missingMigrations", missingMigrations)

	var migrationsToApply []Migration
	for _, collected := range collectedMigrations {
		if collected.Version > latestAppliedVersion {
			migrationsToApply = append(migrationsToApply, collected)
		}
	}
	if len(migrationsToApply) == 0 {
		fmt.Println("No new migrations to apply.")
	}

	// Sort the migrations by version number.
	sort.Slice(migrationsToApply, func(i, j int) bool {
		return migrationsToApply[i].Version < migrationsToApply[j].Version
	})

	for _, apply := range migrationsToApply {
		err := apply.Up(duckStore, migrationsFs)
		tool.Assert(err == nil, "failed to apply migration", "error", err)
		latestAppliedVersion++
	}

	fmt.Printf("Current version: %d\n", latestAppliedVersion)
}

func collectSqlMigrations(migrationsFs fs.FS, dirpath string) ([]Migration, error) {
	if _, err := fs.Stat(migrationsFs, dirpath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("'%s' directory does not exist: %v", dirpath, err)
		}
		return nil, err
	}
	// SQL migration files.
	sqlMigrationFiles, err := fs.Glob(migrationsFs, path.Join(dirpath, "*.sql"))
	if err != nil {
		return nil, err
	}

	versionSet := make(map[int64]string, len(sqlMigrationFiles))
	sqlMigrations := make([]Migration, 0, len(sqlMigrationFiles))
	for _, file := range sqlMigrationFiles {
		version, err := numericComponent(file)
		if err != nil {
			return nil, err
		}
		if dupFile, ok := versionSet[version]; ok {
			return nil, fmt.Errorf("duplicate migration version between: '%s' and '%s'", dupFile, file)
		}
		versionSet[version] = file

		sqlMigrations = append(sqlMigrations, Migration{
			Version: version,
			Source:  file,
		})
	}

	return sqlMigrations, nil
}

// numericComponent parses the version from the migration file name.
//
// XXX_descriptivename.ext where XXX specifies the version number and ext specifies the type of
// migration, e.g. '.sql'.
func numericComponent(filename string) (int64, error) {
	base := filepath.Base(filename)
	if ext := filepath.Ext(base); ext != ".sql" {
		return 0, errors.New("migration file does not have .sql or .go file extension")
	}
	idx := strings.Index(base, "_")
	if idx < 0 {
		return 0, errors.New("no filename separator '_' found")
	}
	n, err := strconv.ParseInt(base[:idx], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse version from migration file: %s: %w", base, err)
	}
	if n < 1 {
		return 0, errors.New("migration version must be greater than zero")
	}
	return n, nil
}

// findMissingMigrations migrations returns all missing migrations.
// A migrations is considered missing if it has a version less than the
// current known max version.
func findMissingMigrations(collectedMigrations, appliedMigrations []Migration, latestAppliedVersion int64) []Migration {
	existing := make(map[int64]bool)
	for _, known := range appliedMigrations {
		existing[known.Version] = true
	}
	var missing []Migration
	for _, collected := range collectedMigrations {
		if !existing[collected.Version] && collected.Version < latestAppliedVersion {
			missing = append(missing, collected)
		}
	}
	sort.SliceStable(missing, func(i, j int) bool {
		return missing[i].Version < missing[j].Version
	})
	return missing
}
