// Package pgstore provides a pgx connection pool to the backoffice Postgres,
// which hosts the `core` (cadastro) and `billing` (faturamento) schemas.
//
// The backend-go service has two stores:
//
//   - store.SQLiteStore (already existing) for the integration domain:
//     credentials, sessions, consumer_units fetched from the API, sync_runs,
//     invoices from the Coelba.
//
//   - pgstore.Pool (this package) for the business domain: customer
//     registry, contracts, billing cycles, calculations, generated PDFs.
//
// They don't share transactions on purpose. When a billing operation needs
// data that lives in SQLite (e.g. the full PDF of a Coelba invoice), it
// reads from SQLite in a separate call and holds the data in memory for
// the duration of the Postgres transaction.
package pgstore

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config holds the Postgres connection parameters.
type Config struct {
	// URL is the full DSN. If set, other fields are ignored.
	// Example: postgres://user:pass@host:5432/backoffice?sslmode=disable
	URL string

	// Individual fields used when URL is empty.
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string // disable, require, verify-ca, verify-full

	// Pool tuning
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// DefaultConfig returns sensible defaults for local dev.
func DefaultConfig() Config {
	return Config{
		Host:            "127.0.0.1",
		Port:            5432,
		User:            "backoffice",
		Database:        "backoffice",
		SSLMode:         "disable",
		MaxConns:        10,
		MinConns:        2,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}
}

// Open creates a pgx connection pool and verifies connectivity with a ping.
// Returns an error if the pool cannot be created or the database is unreachable.
func Open(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	dsn := cfg.URL
	if dsn == "" {
		dsn = fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode,
		)
	}

	pcfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgstore: parse config: %w", err)
	}

	if cfg.MaxConns > 0 {
		pcfg.MaxConns = cfg.MaxConns
	}
	if cfg.MinConns > 0 {
		pcfg.MinConns = cfg.MinConns
	}
	if cfg.MaxConnLifetime > 0 {
		pcfg.MaxConnLifetime = cfg.MaxConnLifetime
	}
	if cfg.MaxConnIdleTime > 0 {
		pcfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	}

	pool, err := pgxpool.NewWithConfig(ctx, pcfg)
	if err != nil {
		return nil, fmt.Errorf("pgstore: open pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pgstore: ping: %w", err)
	}
	return pool, nil
}
