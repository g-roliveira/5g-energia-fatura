// Package contract implements the business logic for managing commercial
// contracts between a customer and a consumer unit. The fundamental rule:
// a contract is NEVER updated in place. Changing any term (discount %,
// IP formula, bandeira rule) creates a new version. The previous version
// has its vigencia_fim set to the day before the new version starts.
//
// This versioning is what makes billing_calculation.contract_snapshot_json
// meaningful: we can always reconstruct exactly which rules were applied
// for any past competência, even if the customer's contract has changed
// several times since.
package contract

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/billing/repo"
)

// Service holds the contract CRUD logic.
type Service struct {
	repo *repo.ContractRepo
}

func NewService(r *repo.ContractRepo) *Service {
	return &Service{repo: r}
}

// CreateInput is what the HTTP handler / caller passes in to create a
// new contract. Ends up in public.contract.
type CreateInput struct {
	CustomerID                        uuid.UUID
	ConsumerUnitID                    uuid.UUID
	VigenciaInicio                    time.Time
	DescontoPercentual                decimal.Decimal
	IPFaturamentoMode                 repo.IPMode
	IPFaturamentoValor                decimal.Decimal
	IPFaturamentoPercent              decimal.Decimal
	BandeiraComDesconto               bool
	CustoDisponibilidadeSempreCobrado bool
	Notes                             *string
	CreatedBy                         *uuid.UUID
}

// Validate checks business invariants before we touch the database.
func (in *CreateInput) Validate() error {
	if in.CustomerID == uuid.Nil {
		return errors.New("customer_id é obrigatório")
	}
	if in.ConsumerUnitID == uuid.Nil {
		return errors.New("consumer_unit_id é obrigatório")
	}
	if in.VigenciaInicio.IsZero() {
		return errors.New("vigencia_inicio é obrigatória")
	}
	if in.DescontoPercentual.LessThanOrEqual(decimal.Zero) ||
		in.DescontoPercentual.GreaterThan(decimal.NewFromInt(1)) {
		return errors.New("desconto_percentual deve estar no intervalo (0, 1]")
	}
	switch in.IPFaturamentoMode {
	case repo.IPModeFixed, repo.IPModePercent:
		// ok
	case "":
		// default applied later
	default:
		return fmt.Errorf("ip_faturamento_mode inválido: %q", in.IPFaturamentoMode)
	}
	if in.IPFaturamentoMode == repo.IPModePercent &&
		(in.IPFaturamentoPercent.IsZero() || in.IPFaturamentoPercent.GreaterThan(decimal.NewFromInt(1))) {
		return errors.New("ip_faturamento_percent deve estar em (0, 1] quando mode=percent")
	}
	return nil
}

// Create persists a new contract version. If there is an active contract
// for the same UC, it is automatically closed (vigencia_fim = in.VigenciaInicio - 1 day).
// Both operations happen in a single transaction so the invariant
// "exactly one active contract per UC" is never violated.
func (s *Service) Create(ctx context.Context, in CreateInput) (*repo.Contract, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	if in.IPFaturamentoMode == "" {
		in.IPFaturamentoMode = repo.IPModeFixed
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("contract.Create: begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }() // no-op after Commit

	// Close the previous active contract on the day before the new starts.
	// This guarantees no overlap in vigência.
	endOfPrevious := in.VigenciaInicio.AddDate(0, 0, -1)
	if err := s.repo.CloseActive(ctx, tx, in.ConsumerUnitID, endOfPrevious); err != nil {
		return nil, err
	}

	c := &repo.Contract{
		CustomerID:                        in.CustomerID,
		ConsumerUnitID:                    in.ConsumerUnitID,
		VigenciaInicio:                    in.VigenciaInicio,
		VigenciaFim:                       nil, // new contract is open-ended
		DescontoPercentual:                in.DescontoPercentual,
		IPFaturamentoMode:                 in.IPFaturamentoMode,
		IPFaturamentoValor:                in.IPFaturamentoValor,
		IPFaturamentoPercent:              in.IPFaturamentoPercent,
		BandeiraComDesconto:               in.BandeiraComDesconto,
		CustoDisponibilidadeSempreCobrado: in.CustoDisponibilidadeSempreCobrado,
		Notes:                             in.Notes,
		Status:                            repo.ContractStatusActive,
		CreatedBy:                         in.CreatedBy,
	}
	if err := s.repo.Insert(ctx, tx, c); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("contract.Create: commit: %w", err)
	}
	return c, nil
}

// Get returns a contract by ID.
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*repo.Contract, error) {
	return s.repo.GetByID(ctx, id)
}

// GetActiveForUC returns the active contract for a UC, or ErrNotFound.
func (s *Service) GetActiveForUC(ctx context.Context, ucID uuid.UUID) (*repo.Contract, error) {
	return s.repo.GetActiveByConsumerUnit(ctx, ucID)
}

// GetForUCAtDate returns the contract that was in force at a specific date.
// Used by the billing engine when (re)calculating a past competência —
// it must use the contract as it existed *at the time of the competência*,
// not today's contract.
func (s *Service) GetForUCAtDate(
	ctx context.Context, ucID uuid.UUID, asOf time.Time,
) (*repo.Contract, error) {
	return s.repo.GetByConsumerUnitAtDate(ctx, ucID, asOf)
}

// ListForUC returns all versions of contracts for a UC, newest first.
func (s *Service) ListForUC(ctx context.Context, ucID uuid.UUID) ([]*repo.Contract, error) {
	return s.repo.ListByConsumerUnit(ctx, ucID)
}
