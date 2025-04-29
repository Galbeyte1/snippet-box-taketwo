package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

type DBConfig struct {
	Username string
	Password string
	Database string
	Host     string
	Port     string
}

func (cfg DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

func OpenDB(dsn string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	maxAttempts := 10
	for i := 1; i <= maxAttempts; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return db, nil
			}
		}

		slog.Warn("waiting for database...", slog.Int("attempt", i), slog.String("error", err.Error()))
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to database after %d attempts: %w", maxAttempts, err)
}

func VerifyDBMigrations(db *sql.DB) error {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_schema = DATABASE()
			AND table_name = 'snippets'
		)
	`).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check DB schema: %w", err)
	}
	if !exists {
		return fmt.Errorf("database is not initialized")
	}

	return nil
}
