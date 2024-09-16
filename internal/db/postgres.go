// internal/db/postgres.go
package db

import (
	"avito-tender-service/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.PostgresConn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
