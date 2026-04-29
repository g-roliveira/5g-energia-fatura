package app

import (
	"fmt"
	"strings"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/store"
)

func openIntegrationStore(cfg Config) (store.IntegrationStore, error) {
	dsn := strings.TrimSpace(cfg.IntegrationPGURL)
	if dsn == "" {
		dsn = strings.TrimSpace(cfg.DatabaseURL)
	}
	if dsn == "" {
		dsn = strings.TrimSpace(cfg.BackofficePGURL)
	}

	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL, INTEGRATION_PG_URL ou BACKOFFICE_PG_URL deve estar configurado")
	}
	if !strings.HasPrefix(dsn, "postgres://") && !strings.HasPrefix(dsn, "postgresql://") {
		return nil, fmt.Errorf("integration store DSN deve ser PostgreSQL (got: %s)", dsn)
	}

	pg, err := store.OpenIntegrationPostgres(dsn)
	if err != nil {
		return nil, fmt.Errorf("open integration postgres store: %w", err)
	}
	return pg, nil
}
