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

	// Initialization means tables are eisting and created
	//	And all the neceessary constraints are applied
	res, err := m.isDBInit()
	if err != nil {
		return err
	}
	// TODO: Instead initalize the uninitialized ie create the tables
	//		This requires knowing which tables need to be initialized
	//		Is DB should return a map of table errors and the tables
	//		mentioned are the objective
	if len(res) == 0 {
		return nil
	} else {
		// Begin working on un-initialized tasks
		// begin creating tables that haven't been created
		// beging adding constraints to tables that are missing/new
		// err = m.initalizeNewTables(res)

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

// func (m *MigrationsRunner) userEmailConstraints() error {
// 	var count int
// 	err := m.db.QueryRow(`
// 		SELECT COUNT(*)
// 		FROM information_schema.TABLE_CONSTRAINTS
// 		WHERE TABLE_NAME = 'users'
// 			AND CONSTRAINT_NAME = 'users_uc_email'
// 			AND CONSTRAINT_TYPE = 'UNIQUE'
// 			AND TABLE_SCHEMA = DATABASE()
// 	`).Scan(&count)
// 	if err != nil {
// 		return fmt.Errorf("failed to check Email Constraints: %w", err)
// 	}

// 	if count > 0 {
// 		// Unique constraint on email already exists
// 		return nil
// 	}

// 	_, err = m.db.Exec(`ALTER TABLE users Add CONSTRAINT users_uc_email UNIQUE (email);`)
// 	if err != nil {
// 		return fmt.Errorf("failed to add unique constraint: %w", err)
// 	}

// 	// Successful Email constraint applied
// 	return nil
// }

/*
TODO:
  - Query Row Snippets
  - Query Row Users
  - Ensure Constraints
  - Return which issues need to be addressed
*/
func (m *MigrationsRunner) tableExists(tableName string) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = DATABASE()
			AND table_name = '%s'
		)
	`, tableName)
	var exists bool
	err := m.db.QueryRow(query).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to query %s: %w", tableName, err)
	}
	return exists, nil
}

func (m *MigrationsRunner) isDBInit() (map[string]error, error) {
	tables := []string{"users", "snippets"}

	tableErrors := make(map[string]error)

	for _, table := range tables {
		exist, err := m.tableExists(table)
		if err != nil {
			tableErrors[table] = fmt.Errorf("error checking %s: %w", table, err)
			continue
		}

		if !exist {
			// tableErrors[table] = fmt.Errorf("%s does not exist", table)

		}
	}

	return tableErrors, nil
}

// func (m *MigrationsRunner) initalizeNewTables(tables map[string]error) error {

// 	sqlStmt :=

// 	if _, err := m.db.Exec(sqlStmt); err != nil {
// 		return fmt.Errorf("error executing sql statement %s: %w", file, err)
// 	}
// 	return nil
// }

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
		name := entry.Name()
		if strings.HasSuffix(name, ".sql") && name != "001_init.sql" {
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
