package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/unwelcome/iqjtest/internal/config"
)

func Connect(cfg *config.Config, l zerolog.Logger) *sql.DB {
	postgres, err := ConnectToPostgres(cfg)
	if err != nil {
		l.Fatal().Err(err).Msg("Database connection failed")
	}

	l.Trace().Msg("Successfully connected to PostgresSQL!")
	return postgres
}

func ConnectToPostgres(cfg *config.Config) (*sql.DB, error) {
	connStr := cfg.DBConnString()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
