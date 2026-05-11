package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	_ "modernc.org/sqlite"

	"mtesense_home/internal/migrations"
)

func Open(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}
	database, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	database.SetMaxOpenConns(1)
	if _, err := database.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}
	return database, nil
}

func Migrate(database *sql.DB) error {
	entries, err := migrations.Files.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		content, err := migrations.Files.ReadFile(name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		if _, err := database.Exec(string(content)); err != nil {
			return fmt.Errorf("run migration %s: %w", name, err)
		}
	}
	return nil
}
