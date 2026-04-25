package app

import (
	"fmt"
	"strings"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/store"
)

func openIntegrationStore(cfg Config) (store.IntegrationStore, error) {
	if dsn := strings.TrimSpace(cfg.IntegrationPGURL); dsn != "" {
		pg, err := store.OpenPostgres(dsn)
		if err != nil {
			return nil, fmt.Errorf("open integration postgres store: %w", err)
		}
		return pg, nil
	}

	dsn := strings.TrimSpace(cfg.DatabaseURL)
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		pg, err := store.OpenPostgres(dsn)
		if err != nil {
			return nil, fmt.Errorf("open postgres store from BACKEND_DATABASE_URL: %w", err)
		}
		return pg, nil
	}

	sqlite, err := store.OpenSQLite(dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite store: %w", err)
	}
	return sqlite, nil
}
