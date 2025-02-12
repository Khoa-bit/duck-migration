package duck

import (
	"duck-migration/tool"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

type Migration struct {
	CreatedAt time.Time
	Version   int64
	Source    string
}

func (m *Migration) Up(duckStore Store, migrationsFs fs.FS) error {
	if m.Version <= 0 {
		return fmt.Errorf("ERROR %v: invalid migration version", filepath.Base(m.Source))
	}
	if m.Source == "" {
		return fmt.Errorf("ERROR %v: empty migration source", filepath.Base(m.Source))
	}

	f, err := migrationsFs.Open(m.Source)
	if err != nil {
		return fmt.Errorf("ERROR %v: failed to open SQL migration file: %w", filepath.Base(m.Source), err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			err = fmt.Errorf("ERROR %v: failed to close SQL migration file: %w", filepath.Base(m.Source), err)
		}
	}()

	buffer, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("ERROR %v: failed to read SQL migration file: %w", filepath.Base(m.Source), err)
	}

	if len(buffer) == 0 {
		return fmt.Errorf("ERROR %v: SQL migration file is empty", filepath.Base(m.Source))
	}

	start := time.Now()
	defer func() {
		finish := time.Since(start)
		if err != nil {
			fmt.Printf("%s   %s (%dms)\n", tool.FormatRed("FAIL"), m.Source, finish.Milliseconds())
		} else {
			fmt.Printf("%s   %s (%dms)\n", tool.FormatGreen("OK  "), m.Source, finish.Milliseconds())
		}
	}()
	err = duckStore.DuckExecuteMigration(strings.TrimSpace(string(buffer)))
	if err != nil {
		return fmt.Errorf("ERROR %v: failed to execute SQL migration: %w", filepath.Base(m.Source), err)
	}

	_, err = duckStore.DuckInsertVersion(m.Version)
	if err != nil {
		return fmt.Errorf("ERROR %v: failed to insert migration version: %w", filepath.Base(m.Source), err)
	}

	return nil
}
