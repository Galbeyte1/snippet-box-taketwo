package database

import (
	"database/sql"
	"embed"
	"fmt"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var embeddedMigrations embed.FS

type MigrationsRunner struct {
	db *sql.DB
	fs embed.FS
}

func NewMigrationsRunner(db *sql.DB) *MigrationsRunner {
	return &MigrationsRunner{
		db: db,
		fs: embeddedMigrations,
	}
}

func (m *MigrationsRunner) Run() error {

	initialized, err := m.isDBInit()
	if err != nil {
		return err
	}

	if initialized {
		return nil
	}

	files, err := m.migrationFiles()
	if err != nil {
		return err
	}

	if err := m.applyMigrations(files); err != nil {
		return err
	}

	return nil
}

func (m *MigrationsRunner) isDBInit() (bool, error) {
	var exists bool
	err := m.db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = DATABASE()
			AND table_name = 'snippets'
		)
	`).Scan(&exists)
	if err != nil {
		return exists, fmt.Errorf("failed to check DB initialization: %w", err)
	}
	return exists, nil
}

func (m *MigrationsRunner) migrationFiles() ([]string, error) {
	entries, err := m.fs.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("cannot read migration directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)
	return files, nil
}

func (m *MigrationsRunner) applyMigrations(files []string) error {
	for _, file := range files {
		data, err := m.fs.ReadFile("migrations/" + file)
		if err != nil {
			return fmt.Errorf("cannot read migration file %s: %w", file, err)
		}

		sqlStmt := strings.TrimSpace(string(data))
		if sqlStmt == "" {
			continue
		}

		if _, err := m.db.Exec(sqlStmt); err != nil {
			return fmt.Errorf("error executing migration %s: %w", file, err)
		}
	}
	return nil
}
