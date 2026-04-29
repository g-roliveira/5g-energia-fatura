package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound é retornado quando a query não encontra registro.
var ErrNotFound = errors.New("repo: not found")

// ContractRepo encapsula todas as queries da tabela public.contract.
type ContractRepo struct {
	pool *pgxpool.Pool
}

func NewContractRepo(pool *pgxpool.Pool) *ContractRepo {
	return &ContractRepo{pool: pool}
}

const contractCols = `
    id, customer_id, consumer_unit_id, vigencia_inicio, vigencia_fim,
    fator_repasse_energia, valor_ip_com_desconto, ip_faturamento_mode, ip_faturamento_valor,
    ip_faturamento_percent, bandeira_com_desconto,
    custo_disponibilidade_sempre_cobrado, consumo_minimo_kwh, notes, status,
    created_at, created_by, updated_at
`

// scanContract reads one row into a Contract.
func scanContract(row pgx.Row) (*Contract, error) {
	var c Contract
	err := row.Scan(
		&c.ID, &c.CustomerID, &c.ConsumerUnitID, &c.VigenciaInicio, &c.VigenciaFim,
		&c.FatorRepasseEnergia, &c.ValorIPComDesconto, &c.IPFaturamentoMode, &c.IPFaturamentoValor,
		&c.IPFaturamentoPercent, &c.BandeiraComDesconto,
		&c.CustoDisponibilidadeSempreCobrado, &c.ConsumoMinimoKWh, &c.Notes, &c.Status,
		&c.CreatedAt, &c.CreatedBy, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scanContract: %w", err)
	}
	return &c, nil
}

// GetByID looks up a contract by its primary key.
func (r *ContractRepo) GetByID(ctx context.Context, id uuid.UUID) (*Contract, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+contractCols+` FROM public.contract WHERE id = $1`, id)
	return scanContract(row)
}

// GetActiveByConsumerUnit returns the currently-in-force contract for a UC.
// Uses the partial unique index idx_contract_one_active_per_uc.
func (r *ContractRepo) GetActiveByConsumerUnit(ctx context.Context, ucID uuid.UUID) (*Contract, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+contractCols+`
		   FROM public.contract
		  WHERE consumer_unit_id = $1
		    AND vigencia_fim IS NULL
		    AND status = 'active'`,
		ucID,
	)
	return scanContract(row)
}

// GetByConsumerUnitAtDate returns the contract that was in force for a UC
// at a specific date. Used by the calc engine when recalculating historical
// cycles — we want the contract that was active at the competência, not today.
func (r *ContractRepo) GetByConsumerUnitAtDate(
	ctx context.Context, ucID uuid.UUID, asOf time.Time,
) (*Contract, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+contractCols+`
		   FROM public.contract
		  WHERE consumer_unit_id = $1
		    AND vigencia_inicio <= $2
		    AND (vigencia_fim IS NULL OR vigencia_fim >= $2)
		    AND status IN ('active','ended')
		  ORDER BY vigencia_inicio DESC
		  LIMIT 1`,
		ucID, asOf,
	)
	return scanContract(row)
}

// ListByConsumerUnit returns all versions of contracts for a UC, newest first.
func (r *ContractRepo) ListByConsumerUnit(ctx context.Context, ucID uuid.UUID) ([]*Contract, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+contractCols+`
		   FROM public.contract
		  WHERE consumer_unit_id = $1
		  ORDER BY vigencia_inicio DESC`,
		ucID,
	)
	if err != nil {
		return nil, fmt.Errorf("ListByConsumerUnit: %w", err)
	}
	defer rows.Close()

	var out []*Contract
	for rows.Next() {
		c, err := scanContract(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// Insert creates a new contract row. This is the only mutation — contracts
// are never UPDATE'd in terms of business fields. To "edit" a contract,
// call CloseActive + Insert in the same transaction (see contract.Service).
func (r *ContractRepo) Insert(ctx context.Context, tx pgx.Tx, c *Contract) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.Status == "" {
		c.Status = ContractStatusActive
	}

	_, err := tx.Exec(ctx,
		`INSERT INTO public.contract (
		     id, customer_id, consumer_unit_id, vigencia_inicio, vigencia_fim,
		     fator_repasse_energia, valor_ip_com_desconto, ip_faturamento_mode, ip_faturamento_valor,
		     ip_faturamento_percent, bandeira_com_desconto,
		     custo_disponibilidade_sempre_cobrado, consumo_minimo_kwh, notes, status, created_by
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
		c.ID, c.CustomerID, c.ConsumerUnitID, c.VigenciaInicio, c.VigenciaFim,
			c.FatorRepasseEnergia, c.ValorIPComDesconto, c.IPFaturamentoMode, c.IPFaturamentoValor,
		c.IPFaturamentoPercent, c.BandeiraComDesconto,
		c.CustoDisponibilidadeSempreCobrado, c.ConsumoMinimoKWh, c.Notes, c.Status, c.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("ContractRepo.Insert: %w", err)
	}
	return nil
}

// CloseActive sets vigencia_fim on the currently active contract for a UC.
// Called when a new contract version is being inserted.
// If there's no active contract, returns nil (idempotent).
func (r *ContractRepo) CloseActive(
	ctx context.Context, tx pgx.Tx, ucID uuid.UUID, endDate time.Time,
) error {
	_, err := tx.Exec(ctx,
		`UPDATE public.contract
		    SET vigencia_fim = $2,
		        status = 'ended',
		        updated_at = NOW()
		  WHERE consumer_unit_id = $1
		    AND vigencia_fim IS NULL
		    AND status = 'active'`,
		ucID, endDate,
	)
	if err != nil {
		return fmt.Errorf("ContractRepo.CloseActive: %w", err)
	}
	return nil
}

// BeginTx starts a transaction on the underlying pool. Callers use this to
// wrap CloseActive + Insert atomically.
func (r *ContractRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.BeginTx(ctx, pgx.TxOptions{})
}
