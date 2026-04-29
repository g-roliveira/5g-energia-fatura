package cycle

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Row representa uma linha do dashboard de ciclo (uma UC).
type Row struct {
	ConsumerUnitID      uuid.UUID `json:"consumer_unit_id"`
	UCCode              string    `json:"uc_code"`
	CustomerName        string    `json:"customer_name"`
	SyncStatus          string    `json:"sync_status"`
	NumeroFatura        *string   `json:"numero_fatura,omitempty"`
	MesReferencia       *string   `json:"mes_referencia,omitempty"`
	CalculationStatus   *string   `json:"calculation_status,omitempty"`
	ValorAziSemDesconto *float64  `json:"valor_azi_sem_desconto,omitempty"`
	ValorAziComDesconto *float64  `json:"valor_azi_com_desconto,omitempty"`
	EconomiaRS          *float64  `json:"economia_rs,omitempty"`
	EconomiaPct         *float64  `json:"economia_pct,omitempty"`
	PDFGenerated        bool      `json:"pdf_generated"`
	NeedsReviewReasons  []string  `json:"needs_review_reasons,omitempty"`
	ErrorMessage        *string   `json:"error_message,omitempty"`
}

// ListRowsRequest filtros para o dashboard.
type ListRowsRequest struct {
	CycleID         uuid.UUID
	Q               string
	SyncStatus      string
	CalcStatus      string
	NeedsReviewOnly bool
	Limit           int
	Offset          int
}

// ListRows retorna o dashboard de um ciclo.
func (s *Service) ListRows(ctx context.Context, req ListRowsRequest) ([]Row, error) {
	if req.Limit <= 0 {
		req.Limit = 100
	}

	query := `
		WITH cycle_ucs AS (
			SELECT
				ccu.consumer_unit_id,
				ccu.status AS sync_status,
				ccu.error_message,
				cu.uc_code,
				COALESCE(c.nome_razao, '') AS customer_name
			FROM public.cycle_consumer_unit ccu
			JOIN public.consumer_unit cu ON cu.id = ccu.consumer_unit_id
			LEFT JOIN public.customer c ON c.id = cu.customer_id
			WHERE ccu.billing_cycle_id = $1
		),
		latest_invoice AS (
			SELECT DISTINCT ON (uir.consumer_unit_id)
				uir.consumer_unit_id,
				uir.numero_fatura,
				uir.mes_referencia,
				uir.billing_record_snapshot
			FROM public.utility_invoice_ref uir
			WHERE uir.billing_cycle_id = $1
			ORDER BY uir.consumer_unit_id, uir.created_at DESC
		),
		latest_calc AS (
			SELECT DISTINCT ON (bc.consumer_unit_id)
				bc.consumer_unit_id,
				bc.status AS calc_status,
				bc.total_sem_desconto,
				bc.total_com_desconto,
				bc.economia_rs,
				bc.economia_pct,
				bc.needs_review_reasons
			FROM public.billing_calculation bc
			WHERE bc.billing_cycle_id = $1
			  AND bc.status != 'superseded'
			ORDER BY bc.consumer_unit_id, bc.version DESC
		),
		pdf_check AS (
			SELECT DISTINCT gd.billing_calculation_id
			FROM public.generated_document gd
			WHERE gd.type = 'customer_invoice_pdf'
		)
		SELECT
			cuc.consumer_unit_id,
			cuc.uc_code,
			cuc.customer_name,
			cuc.sync_status,
			li.numero_fatura,
			li.mes_referencia,
			lc.calc_status,
			COALESCE(lc.total_sem_desconto, 0)::float8,
			COALESCE(lc.total_com_desconto, 0)::float8,
			COALESCE(lc.economia_rs, 0)::float8,
			COALESCE(lc.economia_pct, 0)::float8,
			EXISTS(SELECT 1 FROM pdf_check p
			       JOIN public.billing_calculation bc2 ON bc2.id = p.billing_calculation_id
			       WHERE bc2.consumer_unit_id = cuc.consumer_unit_id
			         AND bc2.billing_cycle_id = $1),
			COALESCE(lc.needs_review_reasons, ARRAY[]::text[]),
			cuc.error_message
		FROM cycle_ucs cuc
		LEFT JOIN latest_invoice li ON li.consumer_unit_id = cuc.consumer_unit_id
		LEFT JOIN latest_calc lc ON lc.consumer_unit_id = cuc.consumer_unit_id
		WHERE 1=1`

	args := []interface{}{req.CycleID}
	argNum := 2

	if req.Q != "" {
		query += fmt.Sprintf(" AND (cuc.uc_code ILIKE $%d OR cuc.customer_name ILIKE $%d)", argNum, argNum)
		args = append(args, "%"+req.Q+"%")
		argNum++
	}
	if req.SyncStatus != "" {
		query += fmt.Sprintf(" AND cuc.sync_status = $%d", argNum)
		args = append(args, req.SyncStatus)
		argNum++
	}
	if req.CalcStatus != "" {
		query += fmt.Sprintf(" AND lc.calc_status = $%d", argNum)
		args = append(args, req.CalcStatus)
		argNum++
	}
	if req.NeedsReviewOnly {
		query += " AND lc.calc_status = 'needs_review'"
	}

	query += fmt.Sprintf(`
		ORDER BY cuc.uc_code
		LIMIT $%d OFFSET $%d`, argNum, argNum+1)
	args = append(args, req.Limit, req.Offset)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Row
	for rows.Next() {
		var r Row
		var semDesc, comDesc, econRS, econPct float64
		var needsReview []string
		err := rows.Scan(
			&r.ConsumerUnitID, &r.UCCode, &r.CustomerName, &r.SyncStatus,
			&r.NumeroFatura, &r.MesReferencia, &r.CalculationStatus,
			&semDesc, &comDesc, &econRS, &econPct,
			&r.PDFGenerated, &needsReview, &r.ErrorMessage,
		)
		if err != nil {
			return nil, err
		}
		r.ValorAziSemDesconto = &semDesc
		r.ValorAziComDesconto = &comDesc
		r.EconomiaRS = &econRS
		r.EconomiaPct = &econPct
		r.NeedsReviewReasons = needsReview
		result = append(result, r)
	}
	return result, rows.Err()
}

// UpdateUCStatus atualiza o status de uma UC no ciclo.
func (s *Service) UpdateUCStatus(ctx context.Context, cycleID, ucID uuid.UUID, status string, calcID *uuid.UUID, errMsg *string) error {
	query := `
		UPDATE public.cycle_consumer_unit
		SET status = $1, calculation_id = $2, error_message = $3, updated_at = NOW()
		WHERE billing_cycle_id = $4 AND consumer_unit_id = $5
	`
	_, err := s.pool.Exec(ctx, query, status, calcID, errMsg, cycleID, ucID)
	return err
}

// GetUCStatus retorna o status de uma UC no ciclo.
func (s *Service) GetUCStatus(ctx context.Context, cycleID, ucID uuid.UUID) (string, error) {
	var status string
	err := s.pool.QueryRow(ctx, `
		SELECT status FROM public.cycle_consumer_unit
		WHERE billing_cycle_id = $1 AND consumer_unit_id = $2
	`, cycleID, ucID).Scan(&status)
	if err == pgx.ErrNoRows {
		return "", fmt.Errorf("not_found")
	}
	return status, err
}
