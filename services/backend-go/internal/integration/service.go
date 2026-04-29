package integration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Service expõe as operações do domínio integration (public schema).
type Service struct {
	store Store
}

// NewService cria um novo Service.
func NewService(store Store) *Service {
	return &Service{store: store}
}

// --- Credential ---

func (s *Service) CreateCredential(ctx context.Context, c *Credential) (*Credential, error) {
	if c.Label == "" || c.DocumentoCipher == "" || c.SenhaCipher == "" {
		return nil, fmt.Errorf("label, documento e senha são obrigatórios")
	}
	if err := s.store.InsertCredential(ctx, c); err != nil {
		return nil, fmt.Errorf("insert credential: %w", err)
	}
	return c, nil
}

func (s *Service) GetCredential(ctx context.Context, id string) (*Credential, error) {
	return s.store.GetCredentialByID(ctx, id)
}

// --- Session ---

func (s *Service) CreateSession(ctx context.Context, sess *Session) (*Session, error) {
	if sess.CredentialID == uuid.Nil {
		return nil, fmt.Errorf("credential_id é obrigatório")
	}
	if err := s.store.InsertSession(ctx, sess); err != nil {
		return nil, fmt.Errorf("insert session: %w", err)
	}
	return sess, nil
}

func (s *Service) GetLatestSession(ctx context.Context, credentialID string) (*Session, error) {
	return s.store.GetLatestSessionByCredentialID(ctx, credentialID)
}

// --- Consumer Unit ---

func (s *Service) SyncConsumerUnit(ctx context.Context, u *ConsumerUnit) error {
	if u.UC == "" {
		return fmt.Errorf("uc é obrigatório")
	}
	return s.store.UpsertConsumerUnit(ctx, u)
}

func (s *Service) ListConsumerUnits(ctx context.Context, limit int, status string) ([]ConsumerUnit, error) {
	return s.store.ListConsumerUnits(ctx, limit, status)
}

func (s *Service) GetConsumerUnit(ctx context.Context, uc string) (*ConsumerUnit, error) {
	return s.store.GetConsumerUnitByUC(ctx, uc)
}

// --- Sync Run ---

func (s *Service) RecordSyncRun(ctx context.Context, sr *SyncRun) (*SyncRun, error) {
	if sr.UC == "" || sr.Documento == "" {
		return nil, fmt.Errorf("uc e documento são obrigatórios")
	}
	if err := s.store.InsertSyncRun(ctx, sr); err != nil {
		return nil, fmt.Errorf("insert sync run: %w", err)
	}
	return sr, nil
}

func (s *Service) GetSyncRun(ctx context.Context, id string) (*SyncRun, error) {
	return s.store.GetSyncRunByID(ctx, id)
}

// --- Job Queue ---

func (s *Service) EnqueueJob(ctx context.Context, jobType string, payload map[string]any) (*Job, error) {
	return s.store.EnqueueJob(ctx, jobType, payload)
}

func (s *Service) ClaimNextJob(ctx context.Context, workerID string) (*Job, error) {
	return s.store.ClaimNextJob(ctx, workerID)
}

func (s *Service) CompleteJob(ctx context.Context, jobID uuid.UUID, result map[string]any) error {
	return s.store.CompleteJob(ctx, jobID, result)
}

func (s *Service) FailJob(ctx context.Context, jobID uuid.UUID, errMsg string) error {
	return s.store.FailJob(ctx, jobID, errMsg)
}
