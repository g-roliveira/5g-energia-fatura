package cycle

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

// SyncJobHandler processa um SyncJob.
type SyncJobHandler func(ctx context.Context, job *SyncJob) error

// WorkerPool gerencia workers concorrentes que processam jobs do public.sync_job.
type WorkerPool struct {
	store        *SyncJobStore
	size         int
	pollInterval time.Duration
	handlers     map[string]SyncJobHandler
	logger       *slog.Logger
	wg           sync.WaitGroup
	stopCh       chan struct{}
	stopped      bool
	mu           sync.Mutex
}

// NewWorkerPool cria um novo pool de workers para sync_job.
func NewWorkerPool(store *SyncJobStore, size int, pollInterval time.Duration, logger *slog.Logger) *WorkerPool {
	if size <= 0 {
		size = 5
	}
	if pollInterval <= 0 {
		pollInterval = 2 * time.Second
	}
	return &WorkerPool{
		store:        store,
		size:         size,
		pollInterval: pollInterval,
		handlers:     make(map[string]SyncJobHandler),
		logger:       logger,
		stopCh:       make(chan struct{}),
	}
}

// RegisterHandler registra um handler para um tipo de job.
func (p *WorkerPool) RegisterHandler(jobType string, handler SyncJobHandler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers[jobType] = handler
}

// Start inicia o pool de workers.
func (p *WorkerPool) Start(ctx context.Context) {
	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	p.logger.Info("sync_job_worker_pool_start", "size", p.size, "poll_interval", p.pollInterval.String())
	for i := 0; i < p.size; i++ {
		p.wg.Add(1)
		go p.worker(ctx, fmt.Sprintf("sync-worker-%d", i))
	}
}

// Stop sinaliza o shutdown e aguarda workers terminarem.
func (p *WorkerPool) Stop() {
	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return
	}
	p.stopped = true
	close(p.stopCh)
	p.mu.Unlock()

	p.wg.Wait()
	p.logger.Info("sync_job_worker_pool_stopped")
}

func (p *WorkerPool) worker(ctx context.Context, workerID string) {
	defer p.wg.Done()

	for {
		select {
		case <-p.stopCh:
			return
		case <-ctx.Done():
			return
		default:
		}

		job, err := p.store.ClaimNextJob(ctx, workerID)
		if err != nil {
			select {
			case <-p.stopCh:
				return
			case <-time.After(p.pollInterval):
				continue
			}
		}

		if job == nil {
			continue
		}

		p.mu.Lock()
		handler, ok := p.handlers[job.Type]
		p.mu.Unlock()

		if !ok {
			p.logger.Warn("no_handler_for_sync_job_type", "job_type", job.Type, "job_id", job.ID)
			p.store.FailJob(ctx, job.ID, fmt.Sprintf("no handler registered for job type: %s", job.Type))
			continue
		}

		p.logger.Info("sync_job_started", "worker", workerID, "job_id", job.ID, "job_type", job.Type)
		NotifyJobEvent(ctx, p.store.pool, p.logger, job, "running", "Processando "+job.Type)
		if err := handler(ctx, job); err != nil {
			p.logger.Error("sync_job_failed", "worker", workerID, "job_id", job.ID, "error", err)
			NotifyJobEvent(ctx, p.store.pool, p.logger, job, "failed", err.Error())
			p.store.FailJob(ctx, job.ID, err.Error())
		} else {
			p.logger.Info("sync_job_completed", "worker", workerID, "job_id", job.ID)
			NotifyJobEvent(ctx, p.store.pool, p.logger, job, "success", job.Type+" concluído")
			p.store.CompleteJob(ctx, job.ID, nil)
		}
	}
}

// RunOne executa um job específico de forma síncrona (útil para testes).
func (p *WorkerPool) RunOne(ctx context.Context, jobID uuid.UUID) error {
	// Marca o job como running e retorna
	query := `
		UPDATE public.sync_job
		SET status = 'running', started_at = NOW()
		WHERE id = $1 AND status = 'pending'
		RETURNING id, type, payload_json, status, retry_count, max_retries,
		          error_message, idempotency_key, scheduled_for, started_at,
		          finished_at, created_at
	`
	row := p.store.pool.QueryRow(ctx, query, jobID)
	job, err := scanSyncJob(row)
	if err != nil {
		return fmt.Errorf("claim job %s: %w", jobID, err)
	}

	p.mu.Lock()
	handler, ok := p.handlers[job.Type]
	p.mu.Unlock()

	if !ok {
		p.store.FailJob(ctx, job.ID, fmt.Sprintf("no handler registered for job type: %s", job.Type))
		return fmt.Errorf("no handler for job type %s", job.Type)
	}

	NotifyJobEvent(ctx, p.store.pool, p.logger, job, "running", "Processando "+job.Type)
	if err := handler(ctx, job); err != nil {
		NotifyJobEvent(ctx, p.store.pool, p.logger, job, "failed", err.Error())
		p.store.FailJob(ctx, job.ID, err.Error())
		return err
	}
	NotifyJobEvent(ctx, p.store.pool, p.logger, job, "success", job.Type+" concluído")
	p.store.CompleteJob(ctx, job.ID, nil)
	return nil
}
