package adjustment

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// Service aplica ajustes manuais em cálculos de faturamento.
// Cada ajuste cria uma nova versão imutável do cálculo.
type Service struct {
	pool *pgxpool.Pool
}

// NewService cria um novo Service.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// ManualAdjustment representa um ajuste manual registrado.
type ManualAdjustment struct {
	ID                 uuid.UUID       `json:"id"`
	BillingCalculationID uuid.UUID   `json:"billing_calculation_id"`
	FieldPath          string          `json:"field_path"`
	OldValue           any             `json:"old_value,omitempty"`
	NewValue           any             `json:"new_value"`
	Reason             string          `json:"reason"`
	CreatedBy          uuid.UUID       `json:"created_by"`
	CreatedAt          time.Time       `json:"created_at"`
}

// ApplyRequest define um ajuste a ser aplicado.
type ApplyRequest struct {
	CalculationID uuid.UUID `json:"calculation_id"`
	FieldPath     string    `json:"field_path"`    // ex: "total_com_desconto", "economia_rs"
	NewValue      any       `json:"new_value"`     // novo valor (number, string, bool)
	Reason        string    `json:"reason"`        // motivo do ajuste
	CreatedBy     uuid.UUID `json:"created_by"`    // quem fez o ajuste
}

// Apply aplica um ajuste manual, criando uma nova versão do cálculo.
func (s *Service) Apply(ctx context.Context, req ApplyRequest) (*ManualAdjustment, error) {
	if req.FieldPath == "" || req.Reason == "" {
		return nil, fmt.Errorf("field_path e reason são obrigatórios")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Buscar cálculo atual
	var current struct {
		ID                   uuid.UUID
		UtilityInvoiceRefID  uuid.UUID
		BillingCycleID       uuid.UUID
		ConsumerUnitID       uuid.UUID
		ContractID           uuid.UUID
		ContractSnapshotJSON []byte
		InputsSnapshotJSON   []byte
		ResultSnapshotJSON   []byte
		TotalSemDesconto     decimal.Decimal
		TotalComDesconto     decimal.Decimal
		EconomiaRS           decimal.Decimal
		EconomiaPct          decimal.Decimal
		Status               string
		Version              int
	}
	err = tx.QueryRow(ctx, `
		SELECT id, utility_invoice_ref_id, billing_cycle_id, consumer_unit_id, contract_id,
		       contract_snapshot_json, inputs_snapshot_json, result_snapshot_json,
		       total_sem_desconto, total_com_desconto, economia_rs, economia_pct,
		       status, version
		FROM public.billing_calculation
		WHERE id = $1
	`, req.CalculationID).Scan(
		&current.ID, &current.UtilityInvoiceRefID, &current.BillingCycleID,
		&current.ConsumerUnitID, &current.ContractID,
		&current.ContractSnapshotJSON, &current.InputsSnapshotJSON, &current.ResultSnapshotJSON,
		&current.TotalSemDesconto, &current.TotalComDesconto, &current.EconomiaRS,
		&current.EconomiaPct, &current.Status, &current.Version,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("cálculo não encontrado")
		}
		return nil, fmt.Errorf("buscar cálculo: %w", err)
	}

	if current.Status == "superseded" {
		return nil, fmt.Errorf("não é possível ajustar uma versão superseded")
	}

	// 2. Extrair old_value do result_snapshot_json
	oldValue, err := extractField(current.ResultSnapshotJSON, req.FieldPath)
	if err != nil {
		return nil, fmt.Errorf("extrair old_value: %w", err)
	}

	// 3. Marcar cálculo atual como superseded
	_, err = tx.Exec(ctx, `
		UPDATE public.billing_calculation
		SET status = 'superseded'
		WHERE id = $1
	`, current.ID)
	if err != nil {
		return nil, fmt.Errorf("superseded calc: %w", err)
	}

	// 4. Calcular novo result_snapshot_json e totais
	newResultJSON, newTotals, err := applyFieldChange(current.ResultSnapshotJSON, req.FieldPath, req.NewValue)
	if err != nil {
		return nil, fmt.Errorf("aplicar mudança: %w", err)
	}

	// 5. Inserir nova versão
	var newCalcID uuid.UUID
	err = tx.QueryRow(ctx, `
		INSERT INTO public.billing_calculation (
			utility_invoice_ref_id, billing_cycle_id, consumer_unit_id, contract_id,
			contract_snapshot_json, inputs_snapshot_json, result_snapshot_json,
			total_sem_desconto, total_com_desconto, economia_rs, economia_pct,
			status, version, calculated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 'needs_review', $12, NOW())
		RETURNING id
	`, current.UtilityInvoiceRefID, current.BillingCycleID, current.ConsumerUnitID, current.ContractID,
		current.ContractSnapshotJSON, current.InputsSnapshotJSON, newResultJSON,
		newTotals.TotalSemDesconto, newTotals.TotalComDesconto, newTotals.EconomiaRS, newTotals.EconomiaPct,
		current.Version+1,
	).Scan(&newCalcID)
	if err != nil {
		return nil, fmt.Errorf("inserir nova versão: %w", err)
	}

	// 6. Registrar ajuste manual
	adj := &ManualAdjustment{
		ID:                 uuid.New(),
		BillingCalculationID: newCalcID,
		FieldPath:          req.FieldPath,
		OldValue:           oldValue,
		NewValue:           req.NewValue,
		Reason:             req.Reason,
		CreatedBy:          req.CreatedBy,
		CreatedAt:          time.Now(),
	}
	oldJSON, _ := json.Marshal(oldValue)
	newJSON, _ := json.Marshal(req.NewValue)
	_, err = tx.Exec(ctx, `
		INSERT INTO public.manual_adjustment (
			id, billing_calculation_id, field_path, old_value, new_value, reason, created_by, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, adj.ID, adj.BillingCalculationID, adj.FieldPath, oldJSON, newJSON, adj.Reason, adj.CreatedBy, adj.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("registrar ajuste: %w", err)
	}

	// 7. Atualizar cycle_consumer_unit com novo calculation_id
	_, err = tx.Exec(ctx, `
		UPDATE public.cycle_consumer_unit
		SET calculation_id = $1, status = 'calculated'
		WHERE billing_cycle_id = $2 AND consumer_unit_id = $3
	`, newCalcID, current.BillingCycleID, current.ConsumerUnitID)
	if err != nil {
		return nil, fmt.Errorf("atualizar ccu: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	return adj, nil
}

// List retorna o histórico de ajustes de um cálculo.
func (s *Service) List(ctx context.Context, calculationID uuid.UUID) ([]ManualAdjustment, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, billing_calculation_id, field_path, old_value, new_value, reason, created_by, created_at
		FROM public.manual_adjustment
		WHERE billing_calculation_id = $1
		ORDER BY created_at DESC
	`, calculationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ManualAdjustment
	for rows.Next() {
		var a ManualAdjustment
		var oldJSON, newJSON []byte
		err := rows.Scan(&a.ID, &a.BillingCalculationID, &a.FieldPath, &oldJSON, &newJSON, &a.Reason, &a.CreatedBy, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(oldJSON, &a.OldValue)
		json.Unmarshal(newJSON, &a.NewValue)
		result = append(result, a)
	}
	return result, rows.Err()
}

// --- helpers ---

func extractField(resultJSON []byte, path string) (any, error) {
	var m map[string]any
	if err := json.Unmarshal(resultJSON, &m); err != nil {
		return nil, err
	}
	v, ok := m[path]
	if !ok {
		return nil, fmt.Errorf("campo %s não encontrado", path)
	}
	return v, nil
}

type totals struct {
	TotalSemDesconto decimal.Decimal
	TotalComDesconto decimal.Decimal
	EconomiaRS       decimal.Decimal
	EconomiaPct      decimal.Decimal
}

func applyFieldChange(resultJSON []byte, path string, newValue any) ([]byte, totals, error) {
	var m map[string]any
	if err := json.Unmarshal(resultJSON, &m); err != nil {
		return nil, totals{}, err
	}
	m[path] = newValue

	// Recalcular economia se total_com_desconto ou total_sem_desconto mudou
	var t totals
	if v, ok := m["total_sem_desconto"]; ok {
		t.TotalSemDesconto = toDecimal(v)
	}
	if v, ok := m["total_com_desconto"]; ok {
		t.TotalComDesconto = toDecimal(v)
	}
	if !t.TotalSemDesconto.IsZero() {
		t.EconomiaRS = t.TotalSemDesconto.Sub(t.TotalComDesconto)
		t.EconomiaPct = t.EconomiaRS.Div(t.TotalSemDesconto)
	}
	m["economia_rs"] = t.EconomiaRS
	m["economia_pct"] = t.EconomiaPct

	newJSON, err := json.Marshal(m)
	return newJSON, t, err
}

func toDecimal(v any) decimal.Decimal {
	switch n := v.(type) {
	case float64:
		return decimal.NewFromFloat(n)
	case string:
		d, _ := decimal.NewFromString(n)
		return d
	case decimal.Decimal:
		return n
	default:
		return decimal.Zero
	}
}
