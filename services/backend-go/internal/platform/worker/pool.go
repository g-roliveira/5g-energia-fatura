package worker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/integration"
)

// Handler é a função que processa um job.
type Handler func(ctx context.Context, job *integration.Job) error

// Pool gerencia workers concorrentes que processam jobs do Postgres.
type Pool struct {
	store     integration.Store
	size      int
	pollInterval time.Duration
	handlers  map[string]Handler
	logger    *slog.Logger
	wg        sync.WaitGroup
	stopCh    chan struct{}
	stopped   bool
	mu        sync.Mutex
}

// NewPool cria um novo worker pool.
func NewPool(store integration.Store, size int, pollInterval time.Duration, logger *slog.Logger) *Pool {
	if size <= 0 {
		size = 5
	}
	if pollInterval <= 0 {
		pollInterval = 2 * time.Second
	}
	return &Pool{
		store:        store,
		size:         size,
		pollInterval: pollInterval,
		handlers:     make(map[string]Handler),
		logger:       logger,
		stopCh:       make(chan struct{}),
	}
}

// RegisterHandler registra um handler para um tipo de job.
func (p *Pool) RegisterHandler(jobType string, handler Handler) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers[jobType] = handler
}

// Start inicia o pool de workers.
func (p *Pool) Start(ctx context.Context) {
	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	p.logger.Info("worker_pool_start", "size", p.size, "poll_interval", p.pollInterval.String())
	for i := 0; i < p.size; i++ {
		p.wg.Add(1)
		go p.worker(ctx, fmt.Sprintf("worker-%d", i))
	}
}

// Stop sinaliza o shutdown e aguarda workers terminarem.
func (p *Pool) Stop() {
	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return
	}
	p.stopped = true
	close(p.stopCh)
	p.mu.Unlock()

	p.wg.Wait()
	p.logger.Info("worker_pool_stopped")
}

func (p *Pool) worker(ctx context.Context, workerID string) {
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
			// Sem jobs disponíveis, espera antes de tentar novamente
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
		handler, ok := p.handlers[job.JobType]
		p.mu.Unlock()

		if !ok {
			p.logger.Warn("no_handler_for_job_type", "job_type", job.JobType, "job_id", job.ID)
			p.store.FailJob(ctx, job.ID, fmt.Sprintf("no handler registered for job type: %s", job.JobType))
			continue
		}

		p.logger.Info("job_started", "worker", workerID, "job_id", job.ID, "job_type", job.JobType)
		if err := handler(ctx, job); err != nil {
			p.logger.Error("job_failed", "worker", workerID, "job_id", job.ID, "error", err)
			p.store.FailJob(ctx, job.ID, err.Error())
		} else {
			p.logger.Info("job_completed", "worker", workerID, "job_id", job.ID)
			p.store.CompleteJob(ctx, job.ID, nil)
		}
	}
}
