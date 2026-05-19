package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const checkTimeout = 3 * time.Second

type DB struct {
	pool *pgxpool.Pool
}

type SchemaValidation struct {
	Status           string   `json:"status"`
	MigrationVersion int      `json:"migration_version"`
	Dirty            bool     `json:"dirty"`
	MissingTables    []string `json:"missing_tables"`
}

var expectedPhase2Tables = []string{
	"users",
	"wallets",
	"markets",
	"positions",
	"trades",
	"liquidity",
	"resolutions",
	"settlements",
	"oracle_events",
	"audit_logs",
	"api_keys",
	"webhooks",
	"agent_access",
}

func Connect(ctx context.Context, databaseURL string) (*DB, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	db := &DB{pool: pool}
	if err := db.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return db, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) Ping(parent context.Context) error {
	ctx, cancel := context.WithTimeout(parent, checkTimeout)
	defer cancel()

	if err := db.pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	return nil
}

func (db *DB) ValidateSchema(parent context.Context) (SchemaValidation, error) {
	ctx, cancel := context.WithTimeout(parent, checkTimeout)
	defer cancel()

	result := SchemaValidation{
		Status:        "error",
		MissingTables: []string{},
	}

	var migrationsExists bool
	err := db.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public'
				AND table_name = 'schema_migrations'
		)
	`).Scan(&migrationsExists)
	if err != nil {
		return result, fmt.Errorf("check schema_migrations table: %w", err)
	}
	if !migrationsExists {
		result.MissingTables = append(result.MissingTables, "schema_migrations")
		return result, nil
	}

	if err := db.pool.QueryRow(ctx, `SELECT version, dirty FROM schema_migrations LIMIT 1`).Scan(&result.MigrationVersion, &result.Dirty); err != nil {
		return result, fmt.Errorf("read schema_migrations: %w", err)
	}

	rows, err := db.pool.Query(ctx, `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
			AND table_name = ANY($1)
	`, expectedPhase2Tables)
	if err != nil {
		return result, fmt.Errorf("read phase 2 tables: %w", err)
	}
	defer rows.Close()

	existing := make(map[string]bool, len(expectedPhase2Tables))
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return result, fmt.Errorf("scan phase 2 table: %w", err)
		}
		existing[tableName] = true
	}
	if err := rows.Err(); err != nil {
		return result, fmt.Errorf("iterate phase 2 tables: %w", err)
	}

	for _, tableName := range expectedPhase2Tables {
		if !existing[tableName] {
			result.MissingTables = append(result.MissingTables, tableName)
		}
	}

	if result.MigrationVersion == 13 && !result.Dirty && len(result.MissingTables) == 0 {
		result.Status = "ok"
	}

	return result, nil
}
