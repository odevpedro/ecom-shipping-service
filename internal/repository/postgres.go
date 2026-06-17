package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}
	return db, nil
}

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS shipping_quotes (
			id UUID PRIMARY KEY,
			from_cep VARCHAR(8) NOT NULL,
			to_cep VARCHAR(8) NOT NULL,
			weight_kg DECIMAL(10,2) NOT NULL,
			price_cents INTEGER NOT NULL,
			estimated_days INTEGER NOT NULL,
			carrier VARCHAR(100) NOT NULL,
			service_name VARCHAR(100) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS tracking_events (
			id UUID PRIMARY KEY,
			order_id VARCHAR(100) NOT NULL,
			location VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			event_date TIMESTAMP WITH TIME ZONE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
	}
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}
