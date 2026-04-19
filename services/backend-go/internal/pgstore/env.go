package pgstore

import (
	"os"
	"strconv"
)

// LoadConfigFromEnv reads Postgres config from env vars. Used by app.LoadConfigFromEnv
// to wire the pool at startup. All env vars are prefixed with BACKOFFICE_PG_ to make
// it clear which Postgres we're talking about (not the integration SQLite).
//
//	BACKOFFICE_PG_URL       — full DSN; wins over individual fields
//	BACKOFFICE_PG_HOST      — default 127.0.0.1
//	BACKOFFICE_PG_PORT      — default 5432
//	BACKOFFICE_PG_USER      — default "backoffice"
//	BACKOFFICE_PG_PASSWORD  — no default
//	BACKOFFICE_PG_DATABASE  — default "backoffice"
//	BACKOFFICE_PG_SSLMODE   — default "disable"
//	BACKOFFICE_PG_MAX_CONNS — default 10
func LoadConfigFromEnv() Config {
	cfg := DefaultConfig()
	if v := os.Getenv("BACKOFFICE_PG_URL"); v != "" {
		cfg.URL = v
		return cfg
	}
	if v := os.Getenv("BACKOFFICE_PG_HOST"); v != "" {
		cfg.Host = v
	}
	if v := os.Getenv("BACKOFFICE_PG_PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Port = n
		}
	}
	if v := os.Getenv("BACKOFFICE_PG_USER"); v != "" {
		cfg.User = v
	}
	if v := os.Getenv("BACKOFFICE_PG_PASSWORD"); v != "" {
		cfg.Password = v
	}
	if v := os.Getenv("BACKOFFICE_PG_DATABASE"); v != "" {
		cfg.Database = v
	}
	if v := os.Getenv("BACKOFFICE_PG_SSLMODE"); v != "" {
		cfg.SSLMode = v
	}
	if v := os.Getenv("BACKOFFICE_PG_MAX_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.MaxConns = int32(n)
		}
	}
	return cfg
}
