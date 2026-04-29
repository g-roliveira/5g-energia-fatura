package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CalculationRepo encapsula queries de public.billing_calculation e
// public.manual_adjustment.
type CalculationRepo struct {
	pool *pgxpool.Pool
}

func NewCalculationRepo(pool *pgxpool.Pool) *CalculationRepo {
	return &CalculationRepo{pool: pool}
}

const calcCols = `
    id, utility_invoice_ref_id, billing_cycle_id, consumer_unit_id,
    contract_id, contract_snapshot_json, inputs_snapshot_json,
    result_snapshot_json, total_sem_desconto, total_com_desconto,
    economia_rs, economia_pct, status, needs_review_reasons, version,
    calculated_at, approved_at, approved_by
`

func scanCalc(row pgx.Row) (*BillingCalculation, error) {
	var c BillingCalculation
	err := row.Scan(
		&c.ID, &c.UtilityInvoiceRefID, &c.BillingCycleID, &c.ConsumerUnitID,
		&c.ContractID, &c.ContractSnapshotJSON, &c.InputsSnapshotJSON,
		&c.ResultSnapshotJSON, &c.TotalSemDesconto, &c.TotalComDesconto,
		&c.EconomiaRS, &c.EconomiaPct, &c.Status, &c.NeedsReviewReasons,
		&c.Version, &c.CalculatedAt, &c.ApprovedAt, &c.ApprovedBy,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scanCalc: %w", err)
	}
	return &c, nil
}

// GetByID retrieves a single calculation by its primary key.
func (r *CalculationRepo) GetByID(ctx context.Context, id uuid.UUID) (*BillingCalculation, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+calcCols+` FROM public.billing_calculation WHERE id = $1`, id)
	return scanCalc(row)
}

// GetCurrentForInvoiceRef returns the latest non-superseded calculation
// for a given utility_invoice_ref. This is what the UI shows as "the
// current calc" for an invoice.
func (r *CalculationRepo) GetCurrentForInvoiceRef(
	ctx context.Context, refID uuid.UUID,
) (*BillingCalculation, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT `+calcCols+`
		   FROM public.billing_calculation
		  WHERE utility_invoice_ref_id = $1
		    AND status != 'superseded'
		  ORDER BY version DESC
		  LIMIT 1`,
		refID,
	)
	return scanCalc(row)
}

// ListByCycle returns all current (non-superseded) calculations for a cycle,
// ordered by consumer unit. Used by the cycle dashboard row endpoint.
func (r *CalculationRepo) ListByCycle(
	ctx context.Context, cycleID uuid.UUID,
) ([]*BillingCalculation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+calcCols+`
		   FROM public.billing_calculation
		  WHERE billing_cycle_id = $1
		    AND status != 'superseded'
		  ORDER BY consumer_unit_id, version DESC`,
		cycleID,
	)
	if err != nil {
		return nil, fmt.Errorf("ListByCycle: %w", err)
	}
	defer rows.Close()

	var out []*BillingCalculation
	for rows.Next() {
		c, err := scanCalc(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// Insert creates a new calculation row. Caller is responsible for marking
// previous versions as superseded (see MarkSuperseded) in the same transaction.
func (r *CalculationRepo) Insert(
	ctx context.Context, tx pgx.Tx, c *BillingCalculation,
) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.Status == "" {
		c.Status = CalcStatusDraft
	}
	if c.Version == 0 {
		c.Version = 1
	}

	_, err := tx.Exec(ctx,
		`INSERT INTO public.billing_calculation (
		     id, utility_invoice_ref_id, billing_cycle_id, consumer_unit_id,
		     contract_id, contract_snapshot_json, inputs_snapshot_json,
		     result_snapshot_json, total_sem_desconto, total_com_desconto,
		     economia_rs, economia_pct, status, needs_review_reasons, version,
		     approved_at, approved_by
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`,
		c.ID, c.UtilityInvoiceRefID, c.BillingCycleID, c.ConsumerUnitID,
		c.ContractID, c.ContractSnapshotJSON, c.InputsSnapshotJSON,
		c.ResultSnapshotJSON, c.TotalSemDesconto, c.TotalComDesconto,
		c.EconomiaRS, c.EconomiaPct, c.Status, c.NeedsReviewReasons, c.Version,
		c.ApprovedAt, c.ApprovedBy,
	)
	if err != nil {
		return fmt.Errorf("CalculationRepo.Insert: %w", err)
	}
	return nil
}

// MarkSuperseded flags all previous calculations for an invoice ref as superseded.
// Called before inserting a new version (from adjustment or recalculate).
func (r *CalculationRepo) MarkSuperseded(
	ctx context.Context, tx pgx.Tx, refID uuid.UUID,
) error {
	_, err := tx.Exec(ctx,
		`UPDATE public.billing_calculation
		    SET status = 'superseded'
		  WHERE utility_invoice_ref_id = $1
		    AND status != 'superseded'`,
		refID,
	)
	if err != nil {
		return fmt.Errorf("MarkSuperseded: %w", err)
	}
	return nil
}

// Approve sets status=approved and records approver + timestamp.
func (r *CalculationRepo) Approve(
	ctx context.Context, id uuid.UUID, approverID uuid.UUID,
) error {
	var query string
	var args []any
	if approverID == uuid.Nil {
		query = `UPDATE public.billing_calculation
			    SET status = 'approved',
			        approved_at = NOW()
			  WHERE id = $1 AND status != 'superseded'`
		args = []any{id}
	} else {
		query = `UPDATE public.billing_calculation
			    SET status = 'approved',
			        approved_at = NOW(),
			        approved_by = $2
			  WHERE id = $1 AND status != 'superseded'`
		args = []any{id, approverID}
	}
	ct, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("Approve: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// NextVersionFor looks up what the next version number should be for a given
// invoice ref. Returns 1 if no prior calc exists.
func (r *CalculationRepo) NextVersionFor(
	ctx context.Context, tx pgx.Tx, refID uuid.UUID,
) (int, error) {
	var maxVersion *int
	err := tx.QueryRow(ctx,
		`SELECT MAX(version) FROM public.billing_calculation
		  WHERE utility_invoice_ref_id = $1`,
		refID,
	).Scan(&maxVersion)
	if err != nil {
		return 0, fmt.Errorf("NextVersionFor: %w", err)
	}
	if maxVersion == nil {
		return 1, nil
	}
	return *maxVersion + 1, nil
}

// BeginTx starts a Postgres transaction for multi-step operations.
func (r *CalculationRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.BeginTx(ctx, pgx.TxOptions{})
}
