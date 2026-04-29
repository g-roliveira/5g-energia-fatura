package cycle

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SyncJob representa um job na fila public.sync_job.
type SyncJob struct {
	ID             uuid.UUID
	Type           string
	Payload        map[string]any
	Status         string
	RetryCount     int
	MaxRetries     int
	ErrorMessage   *string
	IdempotencyKey string
	ScheduledFor   time.Time
	StartedAt      *time.Time
	FinishedAt     *time.Time
	CreatedAt      time.Time
}

// SyncJobStore opera na tabela public.sync_job.
type SyncJobStore struct {
	pool *pgxpool.Pool
}

// NewSyncJobStore cria um novo store.
func NewSyncJobStore(pool *pgxpool.Pool) *SyncJobStore {
	return &SyncJobStore{pool: pool}
}

// ClaimNextJob reserva o próximo job pendente usando FOR UPDATE SKIP LOCKED.
func (s *SyncJobStore) ClaimNextJob(ctx context.Context, workerID string) (*SyncJob, error) {
	query := `
		UPDATE public.sync_job
		SET status = 'running', started_at = NOW()
		WHERE id = (
			SELECT id FROM public.sync_job
			WHERE status = 'pending' AND scheduled_for <= NOW()
			ORDER BY scheduled_for ASC, created_at ASC
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		)
		RETURNING id, type, payload_json, status, retry_count, max_retries,
		          error_message, idempotency_key, scheduled_for, started_at,
		          finished_at, created_at
	`
	row := s.pool.QueryRow(ctx, query)
	return scanSyncJob(row)
}

// CompleteJob marca um job como concluído.
func (s *SyncJobStore) CompleteJob(ctx context.Context, jobID uuid.UUID, result map[string]any) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE public.sync_job
		 SET status = 'success', finished_at = NOW()
		 WHERE id = $1`,
		jobID,
	)
	return err
}

// FailJob marca um job como falho e incrementa retry_count.
func (s *SyncJobStore) FailJob(ctx context.Context, jobID uuid.UUID, errMsg string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE public.sync_job
		 SET status = 'failed', error_message = $1, retry_count = retry_count + 1,
		     finished_at = NOW()
		 WHERE id = $2`,
		errMsg, jobID,
	)
	return err
}

// Enqueue insere um job na fila (idempotente).
func (s *SyncJobStore) Enqueue(ctx context.Context, jobType string, payload map[string]any, idempotencyKey string) error {
	payloadJSON, _ := json.Marshal(payload)
	_, err := s.pool.Exec(ctx, `
		INSERT INTO public.sync_job (type, payload_json, status, idempotency_key, scheduled_for)
		VALUES ($1, $2, 'pending', $3, NOW())
		ON CONFLICT (type, idempotency_key) WHERE status IN ('pending','running','retrying','success')
		DO NOTHING
	`, jobType, payloadJSON, idempotencyKey)
	return err
}

// scanSyncJob lê uma linha de sync_job.
func scanSyncJob(row pgx.Row) (*SyncJob, error) {
	var j SyncJob
	var payloadJSON []byte
	var errMsg *string
	err := row.Scan(
		&j.ID, &j.Type, &payloadJSON, &j.Status, &j.RetryCount, &j.MaxRetries,
		&errMsg, &j.IdempotencyKey, &j.ScheduledFor, &j.StartedAt,
		&j.FinishedAt, &j.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanSyncJob: %w", err)
	}
	if errMsg != nil {
		j.ErrorMessage = errMsg
	}
	if len(payloadJSON) > 0 {
		json.Unmarshal(payloadJSON, &j.Payload)
	}
	return &j, nil
}
