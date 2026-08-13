package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func InitDB(driver, source string, schemaPath string) (*sql.DB, error) {
	database, err := sql.Open(driver, source)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if driver == "sqlite" || driver == "sqlite3" {
		if _, err := database.Exec("PRAGMA foreign_keys = ON;"); err != nil {
			return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
		}
	}

	if schemaPath != "" {
		schema, err := os.ReadFile(schemaPath)
		if err == nil {
			if _, err := database.Exec(string(schema)); err != nil {
				return nil, fmt.Errorf("failed to execute schema: %w", err)
			}
		}
	}

	return database, nil
}
