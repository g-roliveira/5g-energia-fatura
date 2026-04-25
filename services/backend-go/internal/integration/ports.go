package integration

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store define as operações de persistência do domínio integration (public schema).
type Store interface {
	// Credentials
	InsertCredential(ctx context.Context, c *Credential) error
	GetCredentialByID(ctx context.Context, id string) (*Credential, error)

	// Sessions
	InsertSession(ctx context.Context, s *Session) error
	GetLatestSessionByCredentialID(ctx context.Context, credentialID string) (*Session, error)

	// Consumer Units
	UpsertConsumerUnit(ctx context.Context, u *ConsumerUnit) error
	ListConsumerUnits(ctx context.Context, limit int, status string) ([]ConsumerUnit, error)
	GetConsumerUnitByUC(ctx context.Context, uc string) (*ConsumerUnit, error)

	// Raw Invoices
	UpsertRawInvoice(ctx context.Context, inv *RawInvoice) (*RawInvoice, error)
	GetRawInvoiceByID(ctx context.Context, id uuid.UUID) (*RawInvoice, error)
	GetLatestRawInvoiceByUC(ctx context.Context, uc string) (*RawInvoice, error)
	ListRawInvoicesByUC(ctx context.Context, uc string, limit int) ([]RawInvoice, error)

	// Sync Runs
	InsertSyncRun(ctx context.Context, sr *SyncRun) error
	GetSyncRunByID(ctx context.Context, id string) (*SyncRun, error)
	GetLatestSyncRunByUC(ctx context.Context, uc string) (*SyncRun, error)

	// Jobs (Worker Pool)
	EnqueueJob(ctx context.Context, jobType string, payload map[string]any) (*Job, error)
	ClaimNextJob(ctx context.Context, workerID string) (*Job, error)
	CompleteJob(ctx context.Context, jobID uuid.UUID, result map[string]any) error
	FailJob(ctx context.Context, jobID uuid.UUID, errMsg string) error
}

// pgxStore implementa Store usando pgxpool.
type pgxStore struct {
	pool *pgxpool.Pool
}

// NewStore cria uma nova instância de Store.
func NewStore(pool *pgxpool.Pool) Store {
	return &pgxStore{pool: pool}
}
